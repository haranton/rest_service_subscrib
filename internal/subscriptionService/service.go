package subscriptionService

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionService interface {
	ListSubscriptions() ([]Subscription, error)
	CreateSubscriptions(r RequestBody) (Subscription, error)
	GetSubscriptionByID(id string) (Subscription, error)
	UpdateSubcriptionByID(r RequestBody, id string) (Subscription, error)
	DeleteSubcriptionByID(id string) error
	//getAmountOfsubscriptions)
}

type subService struct {
	repo SubscriptionRepository
}

func NewSubscriptionService(r SubscriptionRepository) SubscriptionService {
	return &subService{repo: r}
}

func (sub *subService) ListSubscriptions() ([]Subscription, error) {
	return sub.repo.ListSubscriptions()
}

func (sub *subService) CreateSubscriptions(req RequestBody) (Subscription, error) {

	start, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return Subscription{}, errors.New("неправильный формат start_date (ожидается MM-YYYY)")
	}

	// Парсим дату окончания (если есть)
	var end *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsedEnd, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return Subscription{}, errors.New("неправильный формат end_date (ожидается MM-YYYY)")
		}
		end = &parsedEnd
	}

	subNew := Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   start,
		EndDate:     end,
	}

	return subNew, err

}

func (sub *subService) GetSubscriptionByID(id string) (Subscription, error) {
	return sub.repo.getSubscriptionByID(id)
}

func (sub *subService) UpdateSubcriptionByID(req RequestBody, id string) (Subscription, error) {

	existingSub, err := sub.repo.getSubscriptionByID(id)
	if err != nil {
		return Subscription{}, err
	}
	// Парсим даты
	start, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return Subscription{}, err
	}

	var end *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsedEnd, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return Subscription{}, err
		}
		end = &parsedEnd
	}

	existingSub.ServiceName = req.ServiceName
	existingSub.Price = req.Price
	existingSub.UserID = req.UserID
	existingSub.StartDate = start
	existingSub.EndDate = end

	return existingSub, nil
}

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

func (sub *subService) DeleteSubcriptionByID(id string) error {
	return sub.repo.deleteSubcriptionByID(id)
}

type RequestBody struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"` // формат "MM-YYYY"
	EndDate     *string   `json:"end_date,omitempty"`            // тоже "MM-YYYY"
}
