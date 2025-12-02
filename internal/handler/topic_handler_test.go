package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-platform/internal/models"
	"learning-platform/internal/service"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

type fakeTopicRepo struct {
	topics []models.Topic
}

func (r *fakeTopicRepo) Create(ctx context.Context, t *models.Topic) error {
	if t.ID == "" {
		t.ID = "topic-1"
	}
	r.topics = append(r.topics, *t)
	return nil
}

func (r *fakeTopicRepo) FindAll(ctx context.Context) ([]models.Topic, error) {
	return r.topics, nil
}

func (r *fakeTopicRepo) FindByID(ctx context.Context, id string) (*models.Topic, error) {
	for _, t := range r.topics {
		if t.ID == id {
			cp := t
			return &cp, nil
		}
	}
	return nil, assert.AnError
}

func (r *fakeTopicRepo) Update(ctx context.Context, t *models.Topic) error {
	for i, topic := range r.topics {
		if topic.ID == t.ID {
			r.topics[i] = *t
			return nil
		}
	}
	return nil
}

func (r *fakeTopicRepo) Delete(ctx context.Context, id string) error {
	for i, topic := range r.topics {
		if topic.ID == id {
			r.topics = append(r.topics[:i], r.topics[i+1:]...)
			return nil
		}
	}
	return nil
}

func setupTopicRouter(t *testing.T) (*gin.Engine, *fakeTopicRepo) {
	gin.SetMode(gin.TestMode)

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	repo := &fakeTopicRepo{}
	svc := service.NewTopicService(repo, rdb)
	h := NewTopicHandler(svc)

	r := gin.Default()
	r.POST("/topics", h.Create)
	r.GET("/topics", h.GetAll)
	r.GET("/topics/:id", h.GetByID)
	r.PUT("/topics/:id", h.Update)
	r.DELETE("/topics/:id", h.Delete)

	return r, repo
}

func TestTopicHandler_Create(t *testing.T) {
	router, repo := setupTopicRouter(t)

	body := `{
		"title": "Algebra",
		"slug": "algebra",
		"schoolClass": "GRADE_7"
	}`

	req := httptest.NewRequest("POST", "/topics", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	require.Len(t, repo.topics, 1)
	assert.Equal(t, "Algebra", repo.topics[0].Title)
}

func TestTopicHandler_GetAll(t *testing.T) {
	router, repo := setupTopicRouter(t)

	repo.topics = []models.Topic{
		{ID: "1", Title: "T1", Slug: "t1", SchoolClass: "GRADE_5"},
		{ID: "2", Title: "T2", Slug: "t2", SchoolClass: "GRADE_6"},
	}

	req := httptest.NewRequest("GET", "/topics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	data, ok := resp["data"].([]interface{})
	require.True(t, ok)
	assert.Len(t, data, 2)
}

func TestTopicHandler_GetByID(t *testing.T) {
	router, repo := setupTopicRouter(t)

	repo.topics = []models.Topic{
		{ID: "id-1", Title: "Geometry", Slug: "geometry", SchoolClass: "GRADE_8"},
	}

	req := httptest.NewRequest("GET", "/topics/id-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestTopicHandler_Update(t *testing.T) {
	router, repo := setupTopicRouter(t)

	repo.topics = []models.Topic{
		{ID: "id-1", Title: "Old", Slug: "old", SchoolClass: "GRADE_7"},
	}

	body := `{
		"title": "New Title",
		"slug": "new-slug",
		"schoolClass": "GRADE_7"
	}`

	req := httptest.NewRequest("PUT", "/topics/id-1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	require.Len(t, repo.topics, 1)
	assert.Equal(t, "New Title", repo.topics[0].Title)
	assert.Equal(t, "new-slug", repo.topics[0].Slug)
}

func TestTopicHandler_Delete(t *testing.T) {
	router, repo := setupTopicRouter(t)

	repo.topics = []models.Topic{
		{ID: "id-1", Title: "DeleteMe", Slug: "del", SchoolClass: "GRADE_7"},
	}

	req := httptest.NewRequest("DELETE", "/topics/id-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Len(t, repo.topics, 0)
}
