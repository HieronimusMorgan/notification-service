package config

import (
	"log"
	"notification-service/internal/controller"
	"notification-service/internal/repository"
	"notification-service/internal/services"
	"notification-service/internal/utils"
	controllercron "notification-service/internal/utils/cron/controller"
	repositorycron "notification-service/internal/utils/cron/repository"
	"notification-service/internal/utils/cron/service"
	nt "notification-service/internal/utils/nats"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewServerConfig() (*ServerConfig, error) {
	cfg := LoadConfig()
	redisClient := InitRedis(cfg)
	redisService := utils.NewRedisService(*redisClient)
	db := InitDatabase(cfg)
	engine := InitGin()

	// Graceful Shutdown Handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("🛑 Shutting down gracefully...")

		// Close database and Redis before exiting
		CloseDatabase(db)
		CloseRedis(redisClient)

		os.Exit(0)
	}()

	server := &ServerConfig{
		Gin:        engine,
		Config:     cfg,
		DB:         db,
		Redis:      redisService,
		JWTService: utils.NewJWTService(cfg.JWTSecret),
	}

	server.initRepository()
	server.initServices()
	server.initController()
	server.initCron()
	server.initNats()
	return server, nil
}

// initRepository initializes database access objects (Repository)
func (s *ServerConfig) initRepository() {
	s.Repository = Repository{
		NotificationRepository: repository.NewNotificationRepository(*s.DB),
	}
}

// initServices initializes the application services
func (s *ServerConfig) initServices() {
	s.Services = Services{
		NotificationService: services.NewNotificationService(s.Repository.NotificationRepository,
			s.Config.FCMFilePath,
			s.Config.FCMProjectID,
			s.Config.SMTPHost,
			s.Config.SMTPPort,
			s.Config.SMTPEmail,
			s.Config.SMTPPassword),
	}

}

// Start initializes everything and returns an error if something fails
func (s *ServerConfig) Start() error {
	log.Println("✅ Server configuration initialized successfully!")
	return nil
}

func (s *ServerConfig) initController() {
	s.Controller = Controller{
		NotificationController: controller.NewNotificationController(s.Services.NotificationService),
	}
}

func (s *ServerConfig) initCron() {
	s.Cron = Cron{
		CronRepository: repositorycron.NewCronRepository(*s.DB),
		CronService:    service.NewCronService(*s.DB, repositorycron.NewCronRepository(*s.DB)),
		CronController: controllercron.NewCronJobController(service.NewCronService(*s.DB, repositorycron.NewCronRepository(*s.DB))),
	}
	s.Cron.CronService.Start()
}

func (s *ServerConfig) initNats() {
	s.Nats = Nats{
		NatsService: nt.NewNatsService(s.Config.NatsUrl, s.Services.NotificationService),
	}

	go func() {
		for {
			s.Nats.NatsService.RetryPending("notifications.send")
			time.Sleep(2 * time.Minute)
		}
	}()

}
