package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"

    "learning-platform/internal/service"
    "learning-platform/internal/dto"
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
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.auth.Register(req.Email, req.Password, req.DisplayName); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    access, refresh, err := h.auth.Login(req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token":  access,
        "refresh_token": refresh,
    })
}

func (h *AuthHandler) Refresh(c *gin.Context) {
    var req dto.RefreshRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    access, refresh, err := h.auth.Refresh(req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token":  access,
        "refresh_token": refresh,
    })
}

func (h *AuthHandler) GetMe(c *gin.Context) {
    userID := c.GetString("user_id")
    c.JSON(http.StatusOK, gin.H{
        "user_id": userID,
        "message": "authorized",
    })
}
