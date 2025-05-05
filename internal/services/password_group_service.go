package services

import (
	"errors"
	"github.com/rs/zerolog/log"
	"password-management-service/internal/dto/in"
	"password-management-service/internal/models/password"
	"password-management-service/internal/repository"
	"password-management-service/internal/utils"
	"password-management-service/internal/utils/redis"
)

type PasswordGroupService interface {
	AddPasswordGroup(group *in.PasswordGroupRequest, clientID string) (interface{}, error)
	UpdatePasswordGroup(groupID uint, req struct {
		Name string `json:"name" binding:"required"`
	}, clientID string) (interface{}, error)
	GetListPasswordGroup(clientID string) (interface{}, error)
	GetPasswordGroupByID(groupID uint, clientID string) (interface{}, error)
	DeletePasswordGroupByID(groupID uint, clientID string) error
}

type passwordGroupService struct {
	UserRepository          repository.UserRepository
	PasswordGroupRepository repository.PasswordGroupRepository
	PasswordEntryRepository repository.PasswordEntryRepository
	Redis                   redis.RedisService
}

func NewPasswordGroupService(
	userRepository repository.UserRepository,
	passwordGroupRepository repository.PasswordGroupRepository,
	PasswordEntryRepository repository.PasswordEntryRepository,
	redis redis.RedisService) PasswordGroupService {
	return &passwordGroupService{
		UserRepository:          userRepository,
		PasswordGroupRepository: passwordGroupRepository,
		PasswordEntryRepository: PasswordEntryRepository,
		Redis:                   redis,
	}
}

func (s *passwordGroupService) AddPasswordGroup(group *in.PasswordGroupRequest, clientID string) (interface{}, error) {
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

	passwordGroup := password.PasswordGroup{
		UserID:    user.UserID,
		Name:      group.Name,
		CreatedBy: &user.ClientID,
		UpdatedBy: &user.ClientID,
	}

	if err := s.PasswordGroupRepository.AddPasswordGroup(&passwordGroup); err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to add password group")
		return nil, err
	}

	return passwordGroup, nil
}

func (s *passwordGroupService) UpdatePasswordGroup(groupID uint, req struct {
	Name string `json:"name" binding:"required"`
}, clientID string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve data from Redis")
		return err, nil
	}
	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user by client ID")
		return err, nil
	}

	passwordGroup, err := s.PasswordGroupRepository.GetPasswordGroupByUserIDAndGroupID(user.UserID, groupID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password group by user ID and group ID")
		return err, nil
	}

	passwordGroup.Name = req.Name
	passwordGroup.UpdatedBy = &user.ClientID

	if err := s.PasswordGroupRepository.UpdatePasswordGroup(passwordGroup); err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to update password group")
		return err, nil
	}

	return passwordGroup, nil
}

func (s *passwordGroupService) GetListPasswordGroup(clientID string) (interface{}, error) {
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

	passwordGroups, err := s.PasswordGroupRepository.GetPasswordGroupByUserID(user.UserID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password groups by user ID")
		return nil, err
	}

	return passwordGroups, nil
}

func (s *passwordGroupService) GetPasswordGroupByID(groupID uint, clientID string) (interface{}, error) {
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

	passwordGroup, err := s.PasswordGroupRepository.GetPasswordGroupByUserIDAndGroupID(user.UserID, groupID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password group by user ID and group ID")
		return nil, err
	}

	return passwordGroup, nil
}

func (s *passwordGroupService) DeletePasswordGroupByID(groupID uint, clientID string) error {
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

	passwordGroup, err := s.PasswordGroupRepository.GetPasswordGroupByUserIDAndGroupID(user.UserID, groupID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password group by user ID and group ID")
		return err
	}
	entry, _ := s.PasswordEntryRepository.GetPasswordEntryByGroupID(passwordGroup.GroupID)
	if len(entry) > 0 {
		return errors.New("cannot delete password group with entries")
	}

	if err := s.PasswordGroupRepository.DeletePasswordGroupByID(groupID, user.ClientID); err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to delete password group by ID")
		return err
	}

	return nil
}
