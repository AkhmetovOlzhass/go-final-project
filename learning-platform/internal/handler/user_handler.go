package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"learning-platform/internal/service" 
)

type UserHandler struct {
	users *service.UserService 
	s3    *service.S3Service
}

func NewUserHandler(u *service.UserService, s3 *service.S3Service) *UserHandler {
    return &UserHandler{users: u, s3: s3}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.users.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"email":       user.Email,
		"displayName": user.DisplayName,
		"role":        user.Role,
		"avatarUrl":        user.AvatarURL,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	email := c.PostForm("email")
	displayName := c.PostForm("displayName")

	var avatarURL *string
	file, err := c.FormFile("avatar")
	if err == nil && file != nil {
		url, uploadErr := h.s3.UploadFile(file)
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload avatar"})
			return
		}
		avatarURL = &url
	}

	if err := h.users.Update(userID, &email, &displayName, avatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
