package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
	"learning-platform/internal/service" 
    "learning-platform/internal/response"
    "learning-platform/internal/mapper"
	"learning-platform/internal/dto"
)

type UserHandler struct {
	userService *service.UserService 
	s3    *service.S3Service
}

func NewUserHandler(userService *service.UserService, s3 *service.S3Service) *UserHandler {
    return &UserHandler{userService: userService, s3: s3}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
    users, err := h.userService.GetAllUsers()
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to fetch users")
        return
    }

    responses := mapper.ToUserList(users)
    response.Success(c, responses)
}


func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userId")
	user, err := h.userService.FindByID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, mapper.ToUserResponse(user))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
    userID := c.GetString("userId")

    var req dto.UpdateProfileRequest
    if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    var avatarURL *string
    file, _ := c.FormFile("avatar")
    if file != nil {
        url, uploadErr := h.s3.UploadFile(file)
        if uploadErr != nil {
            response.Error(c, http.StatusInternalServerError, "Failed to upload avatar")
            return
        }
        avatarURL = &url
    }

    updatedProfile, err := h.userService.Update(userID, req.Email, req.DisplayName, avatarURL)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(c, mapper.ToUserResponse(updatedProfile))
}
