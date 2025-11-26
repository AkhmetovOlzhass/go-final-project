package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"

    "learning-platform/internal/service"
    "learning-platform/internal/mapper"
    "learning-platform/internal/dto"
    "learning-platform/internal/response"
)

type AuthHandler struct {
    auth *service.AuthService
}

func NewAuthHandler(a *service.AuthService) *AuthHandler {
    return &AuthHandler{auth: a}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    if err := h.auth.Register(req.Email, req.Password, req.DisplayName); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    response.SuccessWithStatus(c, http.StatusCreated, dto.RegisterResponse{
        Message: "registered",
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    access, refresh, err := h.auth.Login(req.Email, req.Password)
    if err != nil {
        response.Error(c, http.StatusUnauthorized, err.Error())
        return
    }

    response.Success(c, dto.AuthTokensResponse{
        AccessToken:  access,
        RefreshToken: refresh,
    })
}

func (h *AuthHandler) Refresh(c *gin.Context) {
    var req dto.RefreshRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    access, refresh, err := h.auth.Refresh(req.RefreshToken)
    if err != nil {
        response.Error(c, http.StatusUnauthorized, err.Error())
        return
    }

    response.Success(c, dto.AuthTokensResponse{
        AccessToken: access,
        RefreshToken: refresh,
    })
}

func (h *AuthHandler) GetMe(c *gin.Context) {
    userID := c.GetString("userId")
    user, err := h.auth.GetUserByID(userID)
    if err != nil {
        response.Error(c, http.StatusNotFound, "user not found")
        return
    }

    response.Success(c, mapper.ToMeResponse(user))
}
