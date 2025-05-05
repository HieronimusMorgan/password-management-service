package config

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"password-management-service/internal/controller"
	"password-management-service/internal/middleware"
	"password-management-service/internal/repository"
	"password-management-service/internal/services"
	"password-management-service/internal/utils/encryption"
	"password-management-service/internal/utils/jwt"
	"password-management-service/internal/utils/redis"
	"syscall"
)

func NewServerConfig() (*ServerConfig, error) {
	cfg := LoadConfig()
	redisClient := InitRedis(cfg)
	redisService := redis.NewRedisService(*redisClient)
	db := InitDatabase(cfg)
	engine := InitGin()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
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
		JWTService: jwt.NewJWTService(cfg.JWTSecret),
	}

	server.initEncryption()
	server.initNats()
	server.initRepository()
	server.initServices()
	server.initController()
	server.initMiddleware()
	server.initCron()
	return server, nil
}

// InitGin initializes the Gin engine with appropriate configurations
func InitGin() *gin.Engine {
	// Set Gin mode based on environment
	if ginMode := gin.Mode(); ginMode != gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
		logrus.Warn("⚠ Running in DEBUG mode. Use `GIN_MODE=release` in production.")
	} else {
		logrus.Info("✅ Running in RELEASE mode.")
	}

	// Create a new Gin router
	engine := gin.New()

	// Middleware
	engine.Use(gin.Recovery()) // Handles panics and prevents crashes
	engine.Use(gin.Logger())   // Logs HTTP requests

	// Security Headers (Prevents Clickjacking & XSS Attacks)
	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Next()
	})

	logrus.Info("🚀 Gin HTTP server initialized successfully")
	return engine
}

// Start initializes everything and returns an error if something fails
func (s *ServerConfig) Start() error {
	log.Println("✅ Server configuration initialized successfully!")
	return nil
}

func (s *ServerConfig) initEncryption() {
	s.Encryption = Encryption{
		EncryptionService: encryption.NewEncryption(),
	}
}

// initNats initializes the application services
func (s *ServerConfig) initNats() {
	s.Nats = Nats{}
}

func (s *ServerConfig) initRepository() {
	s.Repository = Repository{
		UserRepository:              repository.NewUserRepository(*s.DB),
		UserKeysRepository:          repository.NewUserKeysRepository(*s.DB),
		PasswordEntryRepository:     repository.NewPasswordEntryRepository(*s.DB),
		PasswordEntryKeysRepository: repository.NewPasswordEntryKeysRepository(*s.DB),
		PasswordGroupRepository:     repository.NewPasswordGroupRepository(*s.DB),
		PasswordHistoryRepository:   repository.NewPasswordHistoryRepository(*s.DB),
	}
}

func (s *ServerConfig) initServices() {
	s.Services = Services{
		PasswordEntryService: services.NewPasswordEntryService(
			s.Repository.UserRepository,
			s.Repository.UserKeysRepository,
			s.Repository.PasswordEntryRepository,
			s.Repository.PasswordEntryKeysRepository,
			s.Repository.PasswordGroupRepository,
			s.Encryption.EncryptionService,
			s.Redis),
		PasswordGroupService: services.NewPasswordGroupService(
			s.Repository.UserRepository,
			s.Repository.PasswordGroupRepository,
			s.Repository.PasswordEntryRepository,
			s.Redis),
	}
}

func (s *ServerConfig) initController() {
	s.Controller = Controller{
		PasswordEntryController: controller.NewPasswordEntryController(s.Services.PasswordEntryService, s.JWTService),
		PasswordGroupController: controller.NewPasswordGroupController(s.Services.PasswordGroupService, s.JWTService),
	}
}

func (s *ServerConfig) initMiddleware() {
	s.Middleware = Middleware{
		PasswordMiddleware: middleware.NewPasswordMiddleware(s.JWTService),
		AdminMiddleware:    middleware.NewAdminMiddleware(s.JWTService),
	}
}
func (s *ServerConfig) initCron() {
	s.Cron = Cron{}
	//s.Cron.CronService.Start()
}
