package db

import (
	"fmt"
	"log"
	"os"

	subscriptionService "rest_service/internal/subscriptionService"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() (*gorm.DB, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env: %v", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&subscriptionService.Subscription{}); err != nil {
		log.Fatalf("could not migrate: %v", err)
	}

	return db, nil

}
