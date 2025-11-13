package app

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"learning-platform/internal/middleware"
)

func SetupRouter(c *Container) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", c.AuthHandler.Register)
		auth.POST("/login", c.AuthHandler.Login)
		auth.POST("/refresh", c.AuthHandler.Refresh)
	}

	user := api.Group("/user")
	user.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))
	{
		user.GET("/profile", c.UserHandler.GetProfile)
		user.PUT("/profile", c.UserHandler.UpdateProfile)
	}

	tasks := api.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))
	{
		tasks.GET("/topic/:topic_id", c.TaskHandler.GetTasksByTopic)
		tasks.GET("/:id", c.TaskHandler.GetTask)
		tasks.POST("", c.TaskHandler.CreateTask)
		tasks.PUT("/:id", c.TaskHandler.UpdateTask)
		tasks.DELETE("/:id", c.TaskHandler.DeleteTask)
		tasks.GET("/my/tasks", c.TaskHandler.GetMyTasks)
	}


	return router
}
