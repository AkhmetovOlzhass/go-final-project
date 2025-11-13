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
	TopicHandler *handler.TopicHandler
}

func NewContainer(jwtSecret string) *Container {
	dbConn := db.Connect()

	userRepo := repository.NewUserRepository(dbConn)
	tokenRepo := repository.NewTokenRepository(dbConn)
	topicRepo := repository.NewTopicRepository(dbConn)

	s3Service, err := service.NewS3Service()
	if err != nil {
		log.Fatalf("failed to init S3 service: %v", err)
	}

	authService := service.NewAuthService(userRepo, tokenRepo, jwtSecret)
	userService := service.NewUserService(userRepo, s3Service)
	topicService := service.NewTopicService(topicRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, s3Service)
	topicHandler := handler.NewTopicHandler(topicService)


	return &Container{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		TopicHandler: topicHandler,
	}
}
