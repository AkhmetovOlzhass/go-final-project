package main

import (
	"log"
	"os"

	"learning-platform/internal/app"
	"github.com/joho/godotenv"
)

func main() {
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
