package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Subscription struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ServiceName string         `gorm:"not null" json:"service_name"`
	Price       int            `gorm:"not null" json:"price"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     *time.Time     `json:"end_date,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type RequestBody struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"` // формат "MM-YYYY"
	EndDate     *string   `json:"end_date,omitempty"`            // тоже "MM-YYYY"
}

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

func main() {

	initDB()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/subscriptions", listSubscriptions)
	r.POST("/subscriptions", createSubscriptions)
	r.GET("/subscriptions/:id", getSubscriptionByID)
	r.PUT("/subscriptions/:id", updateSubcriptionByID)
	r.DELETE("/subscriptions/:id", deleteSubcriptionByID)
	r.GET("/subscriptions/amountSubscriptions", getAmountOfsubscriptions)
	r.Run(":8081")

}

func listSubscriptions(c *gin.Context) {

	var subscriptions []Subscription

	if err := db.Find(&subscriptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid request"})
	}

	c.JSON(http.StatusOK, subscriptions)

}

func createSubscriptions(c *gin.Context) {

	var req RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	start, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неправильный формат start_date (ожидается MM-YYYY)"})
		return
	}

	// Парсим дату окончания (если есть)
	var end *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsedEnd, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неправильный формат end_date (ожидается MM-YYYY)"})
			return
		}
		end = &parsedEnd
	}

	sub := Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   start,
		EndDate:     end,
	}

	if err := db.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	c.JSON(http.StatusOK, sub)
}

func getSubscriptionByID(c *gin.Context) {

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var sub Subscription
	if err := db.First(&sub, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	c.JSON(http.StatusOK, sub)

}

func updateSubcriptionByID(c *gin.Context) {

	var subscriptionReuest RequestBody
	if err := c.ShouldBindJSON(&subscriptionReuest); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var existingSub Subscription
	if err := db.First(&existingSub, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Парсим даты
	start, err := time.Parse("01-2006", subscriptionReuest.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неправильный формат start_date"})
		return
	}

	var end *time.Time
	if subscriptionReuest.EndDate != nil && *subscriptionReuest.EndDate != "" {
		parsedEnd, err := time.Parse("01-2006", *subscriptionReuest.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неправильный формат end_date"})
			return
		}
		end = &parsedEnd
	}

	existingSub.ServiceName = subscriptionReuest.ServiceName
	existingSub.Price = subscriptionReuest.Price
	existingSub.UserID = subscriptionReuest.UserID
	existingSub.StartDate = start
	existingSub.EndDate = end

	if err := db.Save(&existingSub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	c.JSON(http.StatusOK, existingSub)

}

func deleteSubcriptionByID(c *gin.Context) {

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := db.Delete(&Subscription{}, "id =?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	c.JSON(http.StatusNoContent, "")
}

func getAmountOfsubscriptions(c *gin.Context) {

	userIDStr := c.Query("user_id")
	serviceName := c.Query("name_service")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("01-2006", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date должен быть в формате MM-YYYY"})
		return
	}
	endDate, err := time.Parse("01-2006", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date должен быть в формате MM-YYYY"})
		return
	}

	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date не может быть раньше start_date"})
		return
	}

	query := db.Model(&Subscription{})

	// Фильтрация по дате: пересечение диапазонов
	query = query.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", endDate, startDate)

	// Фильтрация по user_id
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
			return
		}
		query = query.Where("user_id = ?", userID)
	}

	// Фильтрация по названию сервиса
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	var subscriptions []Subscription
	if err := query.Find(&subscriptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка запроса к базе"})
		return
	}

	total := 0
	for _, s := range subscriptions {
		// Рассчитываем количество месяцев пересечения
		subStart := s.StartDate
		subEnd := time.Now()
		if s.EndDate != nil {
			subEnd = *s.EndDate
		}

		// Ограничиваем период рамками фильтра
		if subStart.Before(startDate) {
			subStart = startDate
		}
		if subEnd.After(endDate) {
			subEnd = endDate
		}

		months := monthsBetween(subStart, subEnd)
		total += months * s.Price
	}

	c.JSON(http.StatusOK, gin.H{
		"total_price": total,
	})

}

func monthsBetween(start, end time.Time) int {
	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())
	return yearDiff*12 + monthDiff + 1 // +1, чтобы включить начальный месяц
}
