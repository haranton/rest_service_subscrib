package main

import (
	"log"
	"os"
	"path/filepath"
	"rest_service/internal/db"
	"rest_service/internal/handlers"
	"rest_service/internal/subscriptionService"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "rest_service/docs" // <-- импорт сгенерированных swagger docs
)

// @title           Subscription API
// @version         1.0
// @description     REST API для управления подписками.
// @termsOfService  http://example.com/terms/
// @contact.name    Поддержка API
// @contact.email   support@example.com
// @host            localhost:8081
// @BasePath        /
func main() {

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal("Ошибка получения текущей директории:", err)
	}

	// Поднимаемся на уровень выше, если мы в cmd/
	if filepath.Base(projectRoot) == "cmd" {
		projectRoot = filepath.Dir(projectRoot)
	}

	// Загружаем .env из корня проекта
	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Не удалось загрузить .env файл:", err)
		// Не прерываем выполнение, возможно переменные заданы в системе
	}

	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("could not connect database: %v", err)
	}

	subsRepo := subscriptionService.NewSubscriptionRepository(db)
	subsService := subscriptionService.NewSubscriptionService(subsRepo)
	subsHadlers := handlers.NewSubscriptionHadler(subsService)

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.GET("/subscriptions", subsHadlers.ListSubscriptions)
	r.POST("/subscriptions", subsHadlers.CreateSubscription)
	r.GET("/subscriptions/:id", subsHadlers.GetSubscriptionByID)
	r.PUT("/subscriptions/:id", subsHadlers.UpdateSubscriptionByID)
	r.DELETE("/subscriptions/:id", subsHadlers.DeleteSubcriptionByID)
	r.GET("/subscriptions/amountSubscriptions", subsHadlers.GetAmountOfsubscriptions)

	r.Run(":8081")
}
