package main

import (
	"log"
	"notification-service/config"
	"notification-service/internal/routes"
)

func main() {
	serverConfig, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	defer func() {
		sqlDB, _ := serverConfig.DB.DB()
		err := sqlDB.Close()
		if err != nil {
			return
		}
		log.Println("✅ Database connection closed")
	}()

	if err := serverConfig.Start(); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}

	engine := serverConfig.Gin

	routes.RegisterRoutes(engine, serverConfig.Controller.NotificationController)
	// Run server
	log.Println("Starting server on :8083")
	err = engine.Run(":" + serverConfig.Config.AppPort)
	if err != nil {
		return
	}
}
