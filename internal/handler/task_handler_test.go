package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
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

type fakeTaskRepo struct {
	tasks []models.Task
}

func (r *fakeTaskRepo) GetAll(ctx context.Context) ([]models.Task, error) {
	return r.tasks, nil
}

func (r *fakeTaskRepo) GetDraft(ctx context.Context) ([]models.Task, error) {
	return nil, nil
}

func (r *fakeTaskRepo) Create(ctx context.Context, t *models.Task) error {
	t.ID = "task-123"
	r.tasks = append(r.tasks, *t)
	return nil
}

func (r *fakeTaskRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	return nil, nil
}

func (r *fakeTaskRepo) UpdateStatus(ctx context.Context, id string, s models.TaskStatus) error {
	return nil
}

func (r *fakeTaskRepo) Update(ctx context.Context, t *models.Task) error { return nil }
func (r *fakeTaskRepo) Delete(ctx context.Context, id string) error      { return nil }

func (r *fakeTaskRepo) GetByTopic(ctx context.Context, topicID string) ([]models.Task, error) {
	return nil, nil
}

func (r *fakeTaskRepo) GetByAuthor(ctx context.Context, authorID string) ([]models.Task, error) {
	return nil, nil
}

func setupTaskRouter(t *testing.T) (*gin.Engine, *fakeTaskRepo) {
	gin.SetMode(gin.TestMode)

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	repo := &fakeTaskRepo{}
	taskService := service.NewTaskService(repo, rdb)

	s3 := &service.S3Service{}

	h := NewTaskHandler(taskService, s3)

	router := gin.Default()
	router.POST("/tasks", h.CreateTask)
	router.GET("/tasks", h.GetAllTasks)

	return router, repo
}

func TestTaskHandler_CreateTask(t *testing.T) {
	router, repo := setupTaskRouter(t)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	writer.WriteField("title", "Test Task")
	writer.WriteField("bodyMd", "Markdown text")
	writer.WriteField("difficulty", "EASY")
	writer.WriteField("status", "DRAFT")
	writer.WriteField("topicId", "topic-1")
	writer.WriteField("officialSolution", "solution")
	writer.WriteField("correctAnswer", "42")
	writer.WriteField("answerType", "TEXT")
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/tasks", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	require.Len(t, repo.tasks, 1)
	assert.Equal(t, "Test Task", repo.tasks[0].Title)
}

func TestTaskHandler_GetAllTasks(t *testing.T) {
	router, repo := setupTaskRouter(t)

	repo.tasks = []models.Task{
		{ID: "1", Title: "T1"},
		{ID: "2", Title: "T2"},
	}

	req := httptest.NewRequest("GET", "/tasks", nil)
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
