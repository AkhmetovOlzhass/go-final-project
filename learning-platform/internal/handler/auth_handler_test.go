package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"learning-platform/internal/service"
)

func setupAuthHandlerForTest() (*gin.Engine, *AuthHandler) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	var authService *service.AuthService = nil

	h := NewAuthHandler(authService)

	return r, h
}

func TestAuthHandler_GetMe_ReturnsUserIDAndMessage(t *testing.T) {
	router, handler := setupAuthHandlerForTest()

	router.GET("/me", func(c *gin.Context) {
		c.Set("user_id", "12345")
		handler.GetMe(c)
	})

	req, err := http.NewRequest(http.MethodGet, "/me", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)

	assert.Equal(t, "12345", body["user_id"])
	assert.Equal(t, "authorized", body["message"])
}

func TestAuthHandler_GetMe_EmptyUserID(t *testing.T) {
	router, handler := setupAuthHandlerForTest()

	router.GET("/me-no-user", func(c *gin.Context) {
		handler.GetMe(c)
	})

	req, err := http.NewRequest(http.MethodGet, "/me-no-user", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)

	assert.Equal(t, "", body["user_id"])
	assert.Equal(t, "authorized", body["message"])
}

func TestAuthHandler_Register_InvalidBody_Returns400(t *testing.T) {
	router, handler := setupAuthHandlerForTest()
	router.POST("/register", handler.Register)

	body := bytes.NewBufferString(`{}`)

	req, err := http.NewRequest(http.MethodPost, "/register", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
