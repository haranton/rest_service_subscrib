package subscriptionService

import "gorm.io/gorm"

type SubscriptionRepository interface {
	ListSubscriptions() ([]Subscription, error)
	createSubscriptions(sub Subscription) error
	getSubscriptionByID(id string) (Subscription, error)
	updateSubcriptionByID(sub Subscription) error
	deleteSubcriptionByID(id string) error
	//getAmountOfsubscriptions)

}

type subRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subRepository{db: db}
}

func (r *subRepository) createSubscriptions(sub Subscription) error {
	return r.db.Create(&sub).Error
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
