package db

import (
	"log"

	subscriptionService "rest_service/internal/subscriptionService"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() (*gorm.DB, error) {

	dsn := "host=localhost user=postgres password=123 dbname=db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&subscriptionService.Subscription{}); err != nil {
		log.Fatalf("could not migrate: %v", err)
	}

	return db, nil

}
