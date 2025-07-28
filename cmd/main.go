package main

import (
	"log"
	"rest_service/internal/db"
	"rest_service/internal/handlers"
	"rest_service/internal/subscriptionService"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("could not connect database: %v", err)
	}

	subsRepo := subscriptionService.NewSubscriptionRepository(db)
	subsService := subscriptionService.NewSubscriptionService(subsRepo)
	subsHadlers := handlers.NewSubscriptionHadler(subsService)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/subscriptions", subsHadlers.ListSubscriptions)
	r.POST("/subscriptions", subsHadlers.CreateSubscription)
	r.GET("/subscriptions/:id", subsHadlers.GetSubscriptionByID)
	r.PUT("/subscriptions/:id", subsHadlers.UpdateSubscriptionByID)
	r.DELETE("/subscriptions/:id", subsHadlers.DeleteSubcriptionByID)
	r.GET("/subscriptions/amountSubscriptions", subsHadlers.GetAmountOfsubscriptions)
	r.Run(":8081")

}
