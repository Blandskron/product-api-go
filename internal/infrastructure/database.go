package infrastructure

import (
	"log"
	"os"
	"product-api-go/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("FATAL: DATABASE_DSN environment variable is not set.")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Database connection successful")
	db.AutoMigrate(&domain.Product{}, &domain.Sale{}, &domain.Purchase{})
	log.Println("Database migration completed")
	return db
}
