package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {

	dsn := "host=localhost user=postgres password=123 dbname=db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Subscription{}); err != nil {
		log.Fatalf("could not migrate: %v", err)
	}
}
