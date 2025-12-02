package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"learning-platform/internal/kafka"
	"learning-platform/internal/models"
	"learning-platform/internal/repository"
	"learning-platform/internal/service"
)

type fakeUserRepoForHandler struct {
	users map[string]*models.User
}

func newFakeUserRepoForHandler() *fakeUserRepoForHandler {
	return &fakeUserRepoForHandler{
		users: make(map[string]*models.User),
	}
}

func (f *fakeUserRepoForHandler) Create(ctx context.Context, u *models.User) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	f.users[u.ID.String()] = u
	return nil
}

func (f *fakeUserRepoForHandler) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (f *fakeUserRepoForHandler) FindByID(ctx context.Context, id string) (*models.User, error) {
	if u, ok := f.users[id]; ok {
		return u, nil
	}
	return nil, nil
}

func (f *fakeUserRepoForHandler) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	u, ok := f.users[id]
	if !ok || u == nil {
		return nil
	}

	if email, ok := updates["email"].(string); ok {
		u.Email = email
	}
	if name, ok := updates["displayName"].(string); ok {
		u.DisplayName = name
	}
	if avatar, ok := updates["avatarUrl"].(string); ok {
		u.AvatarURL = &avatar
	}

	return nil
}

func (f *fakeUserRepoForHandler) GetAll(ctx context.Context) ([]models.User, error) {
	out := make([]models.User, 0, len(f.users))
	for _, u := range f.users {
		if u != nil {
			out = append(out, *u)
		}
	}
	return out, nil
}

var _ repository.IUserRepository = (*fakeUserRepoForHandler)(nil)

func newTestRouterWithAuthHandler(authSvc *service.AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewAuthHandler(authSvc)

	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.Refresh)
	r.POST("/auth/verify", h.Verify)
	r.GET("/auth/me", func(c *gin.Context) {
		h.GetMe(c)
	})

	return r
}

func TestAuthHandler_GetMe_Success(t *testing.T) {
	userRepo := newFakeUserRepoForHandler()
	ctx := context.Background()

	user := &models.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		DisplayName: "Test User",
	}
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	var producer *kafka.EmailProducer = nil
	authSvc := service.NewAuthService(userRepo, nil, nil, producer, "secret")

	router := gin.Default()
	router.GET("/auth/me", func(c *gin.Context) {
		c.Set("userId", user.ID.String())
		NewAuthHandler(authSvc).GetMe(c)
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/auth/me", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	router := newTestRouterWithAuthHandler(nil)

	body := bytes.NewBufferString(`{}`) 
	req, err := http.NewRequest(http.MethodPost, "/auth/register", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	router := newTestRouterWithAuthHandler(nil)

	body := bytes.NewBufferString(`{}`)
	req, err := http.NewRequest(http.MethodPost, "/auth/login", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Refresh_InvalidBody(t *testing.T) {
	router := newTestRouterWithAuthHandler(nil)

	body := bytes.NewBufferString(`{}`)
	req, err := http.NewRequest(http.MethodPost, "/auth/refresh", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}


func TestAuthHandler_Verify_InvalidBody(t *testing.T) {
	router := newTestRouterWithAuthHandler(nil)

	body := bytes.NewBufferString(`{}`)
	req, err := http.NewRequest(http.MethodPost, "/auth/verify", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
