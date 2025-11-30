package app

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	_ "learning-platform/docs"
	"learning-platform/internal/middleware"
)

func SetupRouter(c *Container) *gin.Engine {

	router := gin.Default()

	p := ginprometheus.NewPrometheus("learning_platform")
	p.Use(router)

	router.Use(otelgin.Middleware("learning-platform"))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", c.AuthHandler.Register)
		auth.POST("/verify", c.AuthHandler.Verify)
		auth.POST("/login", c.AuthHandler.Login)
		auth.POST("/refresh", c.AuthHandler.Refresh)
	}

	user := api.Group("/user", middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))
	{
		user.GET("/profile", c.UserHandler.GetProfile)
		user.PUT("/profile", c.UserHandler.UpdateProfile)
		user.GET("/all", c.UserHandler.GetAllUsers)
	}

	topic := api.Group("/topics", middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))
	{
		topic.GET("", c.TopicHandler.GetAll)
		topic.GET("/:id", c.TopicHandler.GetByID)

		protectedTopic := topic.Group("")
		protectedTopic.Use(middleware.RoleMiddleware("Teacher", "Admin", "Student"))
		{
			protectedTopic.POST("", c.TopicHandler.Create)
			protectedTopic.PUT("/:id", c.TopicHandler.Update)
			protectedTopic.DELETE("/:id", c.TopicHandler.Delete)
		}
	}

	tasks := api.Group("/tasks", middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))
	{
		tasks.GET("", c.TaskHandler.GetAllTasks)
		tasks.GET("/drafts", c.TaskHandler.GetDraftTasks)
		tasks.GET("/topic/:topicId", c.TaskHandler.GetTasksByTopic)
		tasks.GET("/:id", c.TaskHandler.GetTask)
		tasks.GET("/my/tasks", c.TaskHandler.GetMyTasks)

		protectedTasks := tasks.Group("")
		protectedTasks.Use(middleware.RoleMiddleware("Teacher", "Admin"))
		{
			protectedTasks.POST("/:id/publish", c.TaskHandler.PublishTask)
			protectedTasks.POST("", c.TaskHandler.CreateTask)
			protectedTasks.PUT("/:id", c.TaskHandler.UpdateTask)
			protectedTasks.DELETE("/:id", c.TaskHandler.DeleteTask)
		}
	}

	return router
}
