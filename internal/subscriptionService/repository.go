package subscriptionService

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	ListSubscriptions() ([]Subscription, error)
	createSubscriptions(sub Subscription) error
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

func (r *subRepository) createSubscriptions(sub Subscription) error {

	tx := r.db.Begin()
	if err := tx.Create(&sub).Error; err != nil {
		tx.Rollback()
		log.Printf("Ошибка: %v", err)
		return err
	}
	tx.Commit()
	log.Printf("Сохранено: %+v", sub)
	return nil

}

func (r *subRepository) ListSubscriptions() ([]Subscription, error) {
	var subs []Subscription
	err := r.db.Find(&subs).Error
	return subs, err
}

func (r *subRepository) getSubscriptionByID(id string) (Subscription, error) {
	var sub Subscription
	err := r.db.First(&sub, "id = ?", id).Error
	return sub, err
}

func (r *subRepository) updateSubcriptionByID(sub Subscription) error {
	err := r.db.Save(&sub).Error
	return err
}

func (r *subRepository) deleteSubcriptionByID(id string) error {
	return r.db.Delete(&Subscription{}, "id = ?", id).Error
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
