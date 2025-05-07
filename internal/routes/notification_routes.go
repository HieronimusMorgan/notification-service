package routes

import (
	"github.com/gin-gonic/gin"
	"notification-service/internal/controller"
)

func RegisterRoutes(r *gin.Engine, ctrl controller.NotificationController) {
	r.POST("/notify", ctrl.Send)
}
