package main

import (
	"context"
	"log"
	"net/http"
	"rest_service/internal/db/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Subscriptions struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       float64    `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

var db *pgxpool.Pool

func main() {

	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer database.CloseDB()

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	subscriptions := r.Group("/subscriptions")
	{
		subscriptions.GET("", listSubscriptions)
	}

	r.Run(":8081")

	// r.GET("/users", getAllUsers)
	// r.GET("/users/:id", getUserByIDHandler)
	// r.POST("/users", createUser)
	// r.PUT("users/:id", updateUser)
	// r.DELETE("users/:id", deleteUser)
	// r.GET("/users/search", searchUsers)

}

func listSubscriptions(c *gin.Context) {
	rows, err := db.Query(context.Background(), "Select * from subscriptions")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()

	var subscriptions []Subscriptions
	for rows.Next() {
		var sub Subscriptions
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		subscriptions = append(subscriptions, sub)
	}

	c.JSON(http.StatusOK, subscriptions)

}
