package app

import (
	"log"

	"learning-platform/internal/db"
	"learning-platform/internal/handler"
	"learning-platform/internal/repository"
	"learning-platform/internal/service"
)

type Container struct {
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
	TaskHandler *handler.TaskHandler
}

func NewContainer(jwtSecret string) *Container {
	dbConn := db.Connect()

	userRepo := repository.NewUserRepository(dbConn)
	tokenRepo := repository.NewTokenRepository(dbConn)

	taskRepo := repository.NewTaskRepository(dbConn)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	s3Service, err := service.NewS3Service()
	if err != nil {
		log.Fatalf("failed to init S3 service: %v", err)
	}

	authService := service.NewAuthService(userRepo, tokenRepo, jwtSecret)
	userService := service.NewUserService(userRepo, s3Service)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, s3Service)

	return &Container{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		TaskHandler: taskHandler,
	}
}
