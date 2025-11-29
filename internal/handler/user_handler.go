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
	s3          *service.S3Service
}

func NewUserHandler(userService *service.UserService, s3 *service.S3Service) *UserHandler {
	return &UserHandler{userService: userService, s3: s3}
}

// GetAllUsers godoc
// @Summary Get all users
// @Tags users
// @Description Returns list of all users
// @Produce json
// @Success 200 {object} response.SuccessWrapper{data=[]dto.UserResponse}
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /user/all [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	responses := mapper.ToUserList(users)
	response.Success(c, responses)
}

// GetProfile godoc
// @Summary Get current user profile
// @Tags users
// @Description Returns profile of the currently authenticated user
// @Produce json
// @Success 200 {object} response.SuccessWrapper{data=dto.UserResponse}
// @Failure 404 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString("userId")
	user, err := h.userService.FindByID(ctx, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, mapper.ToUserResponse(user))
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags users
// @Description Updates profile fields and optionally uploads avatar
// @Accept multipart/form-data
// @Produce json
// @Param email formData string false "Updated email"
// @Param displayName formData string false "Updated display name"
// @Param avatar formData file false "Avatar image"
// @Success 200 {object} response.SuccessWrapper{data=dto.UserResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString("userId")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var avatarURL *string
	file, _ := c.FormFile("avatar")
	if file != nil {
		url, uploadErr := h.s3.UploadFile(ctx, file)
		if uploadErr != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to upload avatar")
			return
		}
		avatarURL = &url
	}

	updatedProfile, err := h.userService.Update(ctx, userID, req.Email, req.DisplayName, avatarURL)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, mapper.ToUserResponse(updatedProfile))
}
