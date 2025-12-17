package middleware

import (
	"net/http"

	"learning-platform/internal/models"
	"learning-platform/internal/service"

	"time"

	"github.com/gin-gonic/gin"
)

func BanMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDRaw, exists := c.Get("userId")
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

		if user.IsBanned == models.IsBannedBanned {

			now := time.Now().UTC()

			if user.BannedUntil != nil && user.BannedUntil.Before(now) {
				if _, err := userService.UnbanProfile(c.Request.Context(), user.ID.String()); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "failed to unban user",
					})
					return
				}
				c.Next()
				return
			}

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
