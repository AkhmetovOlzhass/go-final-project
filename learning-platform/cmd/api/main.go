// @title Learning Platform API
// @version 1.0
// @description API for Learning Platform — users, auth, tasks, topics
// @host localhost:8080
// @BasePath /api/v1

package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	
	"learning-platform/internal/app"
	"learning-platform/internal/otelinit"
)

func main() {
	shutdown := otelinit.Init("learning-platform")
    defer shutdown()

	_ = godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("missing JWT_SECRET")
	}

	container := app.NewContainer(jwtSecret)
	router := app.SetupRouter(container)

	log.Printf("Server running on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
