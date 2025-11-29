package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"learning-platform/internal/dto"
	"learning-platform/internal/mapper"
	"learning-platform/internal/response"
	"learning-platform/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(a *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: a}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Description Creates a new user and sends email verification code
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register payload"
// @Success 201 {object} response.SuccessWrapper{data=dto.RegisterResponse}
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.auth.Register(ctx, req.Email, req.Password, req.DisplayName); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithStatus(c, http.StatusCreated, dto.RegisterResponse{
		Message: "Verification code sent to email",
	})
}

// Verify godoc
// @Summary Verify email
// @Tags auth
// @Description Verifies user email via code
// @Accept json
// @Produce json
// @Param request body dto.VerifyEmailRequest true "Verify email payload"
// @Success 201 {object} response.SuccessWrapper{data=dto.VerifyResponse}
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/verify [post]
func (h *AuthHandler) Verify(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	if err := h.auth.VerifyEmail(ctx, req.Email, req.Code); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithStatus(c, http.StatusCreated, dto.VerifyResponse{
		Message: "Email verified successfully",
	})
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Description Returns access + refresh tokens
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.SuccessWrapper{data=dto.AuthTokensResponse}
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	access, refresh, err := h.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, dto.AuthTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// Refresh godoc
// @Summary Refresh JWT tokens
// @Tags auth
// @Description Returns new access + refresh tokens
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh token payload"
// @Success 200 {object} response.SuccessWrapper{data=dto.AuthTokensResponse}
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	access, refresh, err := h.auth.Refresh(ctx, req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, dto.AuthTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// GetMe godoc
// @Summary Get current user profile
// @Tags auth
// @Description Returns user info for current token
// @Produce json
// @Success 200 {object} response.SuccessWrapper{data=dto.MeResponse}
// @Failure 404 {object} response.ErrorResponse
// @Router /auth/me [get]
// @Security BearerAuth
func (h *AuthHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString("userId")
	user, err := h.auth.GetUserByID(ctx, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "user not found")
		return
	}

	response.Success(c, mapper.ToMeResponse(user))
}
