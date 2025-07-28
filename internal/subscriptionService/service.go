package subscriptionService

import (
	"errors"
	"time"

	"github.com/google/uuid"
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

type ParametersСalculatingSum struct {
	StartDate   time.Time
	EndDate     time.Time
	UserID      uuid.UUID
	ServiceName string
}

type RequestParametersСalculatingSum struct {
	StartDate   string
	EndDate     string
	UserID      string
	ServiceName string
}

type SubscriptionService interface {
	ListSubscriptions() ([]Subscription, error)
	CreateSubscriptions(r RequestBody) (Subscription, error)
	GetSubscriptionByID(id string) (Subscription, error)
	UpdateSubcriptionByID(r RequestBody, id string) (Subscription, error)
	DeleteSubcriptionByID(id string) error
	GetAmountOfsubscriptions(RequestParametersСalculatingSum) (int, error)
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

func (sub *subService) DeleteSubcriptionByID(id string) error {
	return sub.repo.deleteSubcriptionByID(id)
}

func (subService *subService) GetAmountOfsubscriptions(params RequestParametersСalculatingSum) (int, error) {

	startDate, err := time.Parse("01-2006", params.StartDate)
	if err != nil {
		return -1, errors.New("start_date должен быть в формате MM-YYYY")

	}
	endDate, err := time.Parse("01-2006", params.EndDate)
	if err != nil {
		return -1, errors.New("start_date должен быть в формате MM-YYYY")
	}

	if endDate.Before(startDate) {
		return -1, errors.New("end_date не может быть раньше start_date")
	}

	userID := uuid.Nil
	if params.UserID != "" {
		var err error
		userID, err = uuid.Parse(params.UserID)
		if err != nil {
			return -1, errors.New("невалидный UUID")
		}
	}

	validParams := ParametersСalculatingSum{
		StartDate:   startDate,
		EndDate:     endDate,
		UserID:      userID,
		ServiceName: params.ServiceName,
	}

	subs, err := subService.repo.getAmountOfSubscriptions(validParams)

	total := 0
	for _, s := range subs {
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

	return total, nil

}

func monthsBetween(start, end time.Time) int {
	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())
	return yearDiff*12 + monthDiff + 1 // +1, чтобы включить начальный месяц
}
