package middleware

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "learning-platform/internal/models"
  "learning-platform/internal/service"
)

func BanMiddleware(userService *service.UserService) gin.HandlerFunc {
  return func(c *gin.Context) {

    // user_id кладётся AuthMiddleware
    userIDRaw, exists := c.Get("user_id")
    if !exists {
      c.Next()
      return
    }

    userID, ok := userIDRaw.(string)
    if !ok {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
        "error": "invalid user id",
      })
      return
    }

    user, err := userService.FindByID(c.Request.Context(), userID)
    if err != nil {
      c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
        "error": "internal error",
      })
      return
    }

    if user == nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
        "error": "user not found",
      })
      return
    }

    // 🔒 ЖЁСТКАЯ ПРОВЕРКА БАНА
    if user.IsBanned == models.IsBannedBanned {
      c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
        "error":       "user is banned",
        "bannedUntil": user.BannedUntil,
        "banReason":   user.BanReason,
      })
      return
    }

    c.Next()
  }
}