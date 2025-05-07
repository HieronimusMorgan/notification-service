package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"notification-service/internal/controller"
	"notification-service/internal/repository"
	"notification-service/internal/services"
	"notification-service/internal/utils"
	controllercron "notification-service/internal/utils/cron/controller"
	repositorycron "notification-service/internal/utils/cron/repository"
	servicescron "notification-service/internal/utils/cron/service"
	nt "notification-service/internal/utils/nats"
)

// ServerConfig holds all initialized components
type ServerConfig struct {
	Gin        *gin.Engine
	Config     *Config
	DB         *gorm.DB
	Redis      utils.RedisService
	JWTService utils.JWTService
	Controller Controller
	Services   Services
	Repository Repository
	Cron       Cron
	Nats       Nats
}

// Services holds all service dependencies
type Services struct {
	NotificationService services.NotificationService
}

// Repository contains repository (database access objects)
type Repository struct {
	NotificationRepository repository.NotificationRepository
}

type Controller struct {
	NotificationController controller.NotificationController
}

type Cron struct {
	CronService    servicescron.CronService
	CronRepository repositorycron.CronRepository
	CronController controllercron.CronJobController
}

type Nats struct {
	NatsService nt.Service
}
