package app

import (
	"log"

	"learning-platform/internal/db"
	"learning-platform/internal/handler"
	"learning-platform/internal/kafka"
	"learning-platform/internal/repository"
	"learning-platform/internal/service"
)

type Container struct {
	AuthHandler  *handler.AuthHandler
	UserHandler  *handler.UserHandler
	TaskHandler  *handler.TaskHandler
	TopicHandler *handler.TopicHandler
}

func NewContainer(jwtSecret string) *Container {
	dbConn := db.Connect()

	s3Service, err := service.NewS3Service()
	if err != nil {
		log.Fatalf("failed to init S3 service: %v", err)
	}

	emailProducer := kafka.NewEmailProducer()

	userRepo := repository.NewUserRepository(dbConn)
	verifyRepo := repository.NewVerificationRepository(dbConn)
	tokenRepo := repository.NewTokenRepository(dbConn)
	topicRepo := repository.NewTopicRepository(dbConn)
	taskRepo := repository.NewTaskRepository(dbConn)

	authService := service.NewAuthService(userRepo, verifyRepo, tokenRepo, emailProducer, jwtSecret)
	userService := service.NewUserService(userRepo)
	topicService := service.NewTopicService(topicRepo)
	taskService := service.NewTaskService(taskRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, s3Service)
	topicHandler := handler.NewTopicHandler(topicService)
	taskHandler := handler.NewTaskHandler(taskService, s3Service)

	return &Container{
		AuthHandler:  authHandler,
		UserHandler:  userHandler,
		TaskHandler:  taskHandler,
		TopicHandler: topicHandler,
	}
}
