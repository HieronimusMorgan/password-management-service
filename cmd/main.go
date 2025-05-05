package main

import (
	"log"
	"password-management-service/config"
	"password-management-service/internal/routes"
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
	// Initialize routes
	routes.PasswordEntryRoutes(engine, serverConfig.Middleware, serverConfig.Controller.PasswordEntryController)
	routes.PasswordGroupRoutes(engine, serverConfig.Middleware, serverConfig.Controller.PasswordGroupController)

	// Run server
	log.Println("Starting server on :8082")
	err = engine.Run(":8082")
	if err != nil {
		return
	}
}
