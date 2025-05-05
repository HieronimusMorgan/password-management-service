package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"password-management-service/internal/controller"
	"password-management-service/internal/middleware"
	"password-management-service/internal/repository"
	"password-management-service/internal/services"
	"password-management-service/internal/utils/encryption"
	"password-management-service/internal/utils/jwt"
	"password-management-service/internal/utils/redis"
)

// ServerConfig holds all initialized components
type ServerConfig struct {
	Gin        *gin.Engine
	Config     *Config
	DB         *gorm.DB
	Redis      redis.RedisService
	JWTService jwt.Service
	Cron       Cron
	Nats       Nats
	Encryption Encryption
	Controller Controller
	Services   Services
	Repository Repository
	Middleware Middleware
}

// Services holds all service dependencies
type Services struct {
	PasswordEntryService services.PasswordEntryService
}

// Repository contains repository (database access objects)
type Repository struct {
	UserRepository              repository.UserRepository
	UserKeysRepository          repository.UserKeysRepository
	PasswordEntryRepository     repository.PasswordEntryRepository
	PasswordEntryKeysRepository repository.PasswordEntryKeysRepository
	PasswordHistoryRepository   repository.PasswordHistoryRepository
}

type Controller struct {
	PasswordEntryController controller.PasswordEntryController
}

type Middleware struct {
	PasswordMiddleware middleware.PasswordMiddleware
	AdminMiddleware    middleware.AdminMiddleware
}

type Cron struct {
}

type Nats struct {
}

type Encryption struct {
	EncryptionService encryption.Encryption
}
