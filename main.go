package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscriptions struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       float64    `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

var db *gorm.DB

func initDB() {

	dsn := "host=localhost user=postgers password=123 dname=db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Subscriptions{}); err != nil {
		log.Fatalf("could not migrate: %v", err)
	}
}

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/subscriptions", listSubscriptions)
	r.POST("/subscriptions", createSubscriptions)
	r.GET("/subscriptions", getS)
	r.Run(":8081")

}

func listSubscriptions(c *gin.Context) {

	var subscriptions []Subscriptions

	if err := db.Find(&subscriptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid request"})
	}

	c.JSON(http.StatusOK, subscriptions)

}

func createSubscriptions(c *gin.Context) {
	var subscription Subscriptions
	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := db.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	c.JSON(http.StatusOK, subscription)
}
