package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := "postgres://postgres:password@localhost:5433/eventmanagement?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	log.Println("Connected and migrated PostgreSQL")
	return db
}
