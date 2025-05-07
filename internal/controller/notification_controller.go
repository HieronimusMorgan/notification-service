package controller

import (
	"net/http"
	"notification-service/internal/models"
	"notification-service/internal/services"

	"github.com/gin-gonic/gin"
)

type NotificationController interface {
	Send(c *gin.Context)
}

type notificationController struct {
	service services.NotificationService
}

func NewNotificationController(service services.NotificationService) NotificationController {
	return &notificationController{service: service}
}

func (ctrl *notificationController) Send(c *gin.Context) {
	var notif models.Notification
	if err := c.ShouldBindJSON(&notif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification enqueued"})
}
