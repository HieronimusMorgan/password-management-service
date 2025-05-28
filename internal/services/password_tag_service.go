package services

import (
	"github.com/rs/zerolog/log"
	"password-management-service/internal/models/password"
	"password-management-service/internal/repository"
	"password-management-service/internal/utils"
	"password-management-service/internal/utils/redis"
)

type PasswordTagService interface {
	AddPasswordTag(req struct {
		Name string `json:"name" binding:"required"`
	}, clientID string) (interface{}, error)
	UpdatePasswordTag(tagID uint, req struct {
		Name string `json:"name" binding:"required"`
	}, clientID string) (interface{}, error)
	GetListPasswordTag(clientID string, index, size int) (interface{}, int64, error)
	DeletePasswordTagByID(tagID uint, clientID string) error
}

type passwordTagService struct {
	UserRepository        repository.UserRepository
	PasswordTagRepository repository.PasswordTagRepository
	Redis                 redis.RedisService
}

func NewPasswordTagService(
	userRepository repository.UserRepository,
	passwordTagRepository repository.PasswordTagRepository,
	redis redis.RedisService) PasswordTagService {
	return &passwordTagService{
		UserRepository:        userRepository,
		PasswordTagRepository: passwordTagRepository,
		Redis:                 redis,
	}
}

func (s *passwordTagService) AddPasswordTag(req struct {
	Name string `json:"name" binding:"required"`
}, clientID string) (interface{}, error) {

	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve data from Redis")
		return nil, err
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user by client ID")
		return nil, err
	}

	passwordTag := &password.PasswordTag{
		UserID:    user.UserID,
		Name:      req.Name,
		CreatedBy: &user.ClientID,
		UpdatedBy: &user.ClientID,
	}

	err = s.PasswordTagRepository.AddPasswordTag(passwordTag)
	if err != nil {
		return nil, err
	}

	return passwordTag, nil
}

func (s *passwordTagService) UpdatePasswordTag(tagID uint, req struct {
	Name string `json:"name" binding:"required"`
}, clientID string) (interface{}, error) {

	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve data from Redis")
		return nil, err
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user by client ID")
		return nil, err
	}

	passwordTag, err := s.PasswordTagRepository.GetPasswordTagByIDAndUserID(tagID, user.UserID)
	if err != nil {
		return nil, err
	}

	passwordTag.Name = req.Name
	passwordTag.UpdatedBy = &user.ClientID

	err = s.PasswordTagRepository.UpdatePasswordTag(passwordTag)
	if err != nil {
		return nil, err
	}

	return passwordTag, nil
}

func (s *passwordTagService) GetListPasswordTag(clientID string, index, size int) (interface{}, int64, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve data from Redis")
		return nil, 0, err
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user by client ID")
		return nil, 0, err
	}

	count, err := s.PasswordTagRepository.GetCountPasswordTag(user.UserID)

	passwordTags, err := s.PasswordTagRepository.GetListPasswordTag(user.UserID, index, size)
	if err != nil {
		return nil, count, err
	}

	return passwordTags, count, nil
}

func (s *passwordTagService) DeletePasswordTagByID(tagID uint, clientID string) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve data from Redis")
		return err
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user by client ID")
		return err
	}

	passwordTag, err := s.PasswordTagRepository.GetPasswordTagByIDAndUserID(tagID, user.UserID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password tag by ID")
		return err
	}
	passwordTag.DeletedBy = &user.ClientID

	err = s.PasswordTagRepository.DeletePasswordTag(passwordTag)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to delete password tag by ID")
		return err
	}
	return nil
}
