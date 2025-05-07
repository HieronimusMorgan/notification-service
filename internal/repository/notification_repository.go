package repository

import (
	"gorm.io/gorm"
	"notification-service/internal/models"
)

type NotificationRepository interface {
	Save(notification *models.Notification) error
	Update(notification *models.Notification) error
	MarkAsSent(id uint) error
	GetPendingNotifications() ([]models.Notification, error)
}

type notificationRepository struct {
	db gorm.DB
}

func NewNotificationRepository(db gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Save(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) Update(notification *models.Notification) error {
	return r.db.Model(notification).Updates(notification).Error
}

func (r *notificationRepository) MarkAsSent(id uint) error {
	return r.db.Model(&models.Notification{}).Where("id = ?", id).Update("status", "sent").Error
}

func (r *notificationRepository) GetPendingNotifications() ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("status = ?", "pending").Find(&notifications).Error
	return notifications, err
}
