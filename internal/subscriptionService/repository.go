package subscriptionService

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	ListSubscriptions(page, limit int) ([]Subscription, int64, int, error)
	createSubscriptions(sub Subscription) (Subscription, error)
	getSubscriptionByID(id string) (Subscription, error)
	updateSubcriptionByID(sub Subscription) error
	deleteSubcriptionByID(id string) error
	getAmountOfSubscriptions(params ParametersСalculatingSum) ([]Subscription, error)
}

type subRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subRepository{db: db}
}

func (r *subRepository) createSubscriptions(sub Subscription) (Subscription, error) {

	r.db = r.db.Debug()

	if err := r.db.Create(&sub).Error; err != nil {
		log.Printf("Ошибка создания подписки: %v", err)
		return Subscription{}, err
	}
	log.Printf("Сохранено: %+v", sub)
	return sub, nil
}

func (r *subRepository) ListSubscriptions(page, limit int) ([]Subscription, int64, int, error) {
	var subs []Subscription

	offset := (page - 1) * limit

	err := r.db.Offset(offset).Limit(limit).Find(&subs).Error

	var totalItems int64
	if err := r.db.Model(&Subscription{}).Count(&totalItems).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return subs, totalItems, totalPages, err
}

func (r *subRepository) getSubscriptionByID(id string) (Subscription, error) {
	var sub Subscription
	err := r.db.First(&sub, "id = ?", id).Error
	return sub, err
}

func (r *subRepository) updateSubcriptionByID(sub Subscription) error {
	var existingSub Subscription
	result := r.db.First(&existingSub, "id = ?", sub.ID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("подписка с ID %d не найдена", sub.ID)
		}
		return result.Error
	}

	err := r.db.Save(&sub).Error
	return err
}

func (r *subRepository) deleteSubcriptionByID(id string) error {
	var sub Subscription
	result := r.db.First(&sub, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("подписка с ID %s не найдена", id)
		}
		return result.Error
	}

	// Удаляем только если подписка существует
	result = r.db.Delete(&Subscription{}, "id = ?", id)
	if result.Error != nil {
		log.Printf("Ошибка удаления подписки ID %s: %v", id, result.Error)
		return result.Error
	}

	return nil
}

func (r *subRepository) getAmountOfSubscriptions(params ParametersСalculatingSum) ([]Subscription, error) {

	query := r.db.Model(&Subscription{})

	// Фильтрации
	query = query.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", params.EndDate, params.StartDate)

	if params.UserID != uuid.Nil {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.ServiceName != "" {
		query = query.Where("service_name = ?", params.ServiceName)
	}

	var subscriptions []Subscription
	if err := query.Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}
