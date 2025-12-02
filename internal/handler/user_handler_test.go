package handler

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-platform/internal/models"
	"learning-platform/internal/service"
)


func setupUserRouter() (*gin.Engine, *fakeUserRepoForHandler, *service.UserService) {
	gin.SetMode(gin.TestMode)

	repo := newFakeUserRepoForHandler()
	userSvc := service.NewUserService(repo)
	s3 := &service.S3Service{} 

	h := NewUserHandler(userSvc, s3)

	r := gin.Default()
	r.GET("/user/all", h.GetAllUsers)

	return r, repo, userSvc
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	router, repo, _ := setupUserRouter()

	u1 := &models.User{
		ID:          uuid.New(),
		Email:       "a@test.com",
		DisplayName: "User A",
	}
	u2 := &models.User{
		ID:          uuid.New(),
		Email:       "b@test.com",
		DisplayName: "User B",
	}
	require.NoError(t, repo.Create(nil, u1))
	require.NoError(t, repo.Create(nil, u2))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/user/all", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestUserHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newFakeUserRepoForHandler()
	userSvc := service.NewUserService(repo)
	s3 := &service.S3Service{}

	h := NewUserHandler(userSvc, s3)

	r := gin.Default()

	user := &models.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		DisplayName: "Test User",
	}
	require.NoError(t, repo.Create(nil, user))

	r.GET("/user/profile", func(c *gin.Context) {
		c.Set("userId", user.ID.String())
		h.GetProfile(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/user/profile", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newFakeUserRepoForHandler()
	userSvc := service.NewUserService(repo)
	s3 := &service.S3Service{}

	h := NewUserHandler(userSvc, s3)

	r := gin.Default()

	user := &models.User{
		ID:          uuid.New(),
		Email:       "old@example.com",
		DisplayName: "Old Name",
	}
	require.NoError(t, repo.Create(nil, user))

	r.PUT("/user/profile", func(c *gin.Context) {
		c.Set("userId", user.ID.String())
		h.UpdateProfile(c)
	})

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("email", "new@example.com")
	_ = writer.WriteField("displayName", "New Name")
	_ = writer.Close()

	req := httptest.NewRequest("PUT", "/user/profile", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	updated, err := repo.FindByID(nil, user.ID.String())
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "new@example.com", updated.Email)
	assert.Equal(t, "New Name", updated.DisplayName)
}
