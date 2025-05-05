package services

import (
	"errors"
	"github.com/rs/zerolog/log"
	"password-management-service/internal/dto/in"
	"password-management-service/internal/models/password"
	"password-management-service/internal/repository"
	"password-management-service/internal/utils"
	"password-management-service/internal/utils/encryption"
	"password-management-service/internal/utils/redis"
	"password-management-service/internal/utils/text"
)

type PasswordEntryService interface {
	AddPasswordEntry(passwordEntryRequest *in.PasswordEntryRequest, clientID string) error
	UpdatePasswordEntry(passwordEntryID uint, passwordEntryRequest *in.PasswordEntryRequest, clientID string) error
	GetPasswordEntryByID(passwordEntryID uint, clientID string) (interface{}, error)
	GetListPasswordEntries(clientID string) (interface{}, error)
	DeletePasswordEntry(passwordEntryID uint, clientID string) error
}

type passwordEntryService struct {
	UserRepository          repository.UserRepository
	UserKeyRepository       repository.UserKeysRepository
	PasswordEntryRepository repository.PasswordEntryRepository
	PasswordEntryKeys       repository.PasswordEntryKeysRepository
	EncryptionService       encryption.Encryption
	Redis                   redis.RedisService
}

func NewPasswordEntryService(
	userRepository repository.UserRepository,
	userKeyRepository repository.UserKeysRepository,
	passwordEntryRepository repository.PasswordEntryRepository,
	passwordEntryKeysRepository repository.PasswordEntryKeysRepository,
	encryptionService encryption.Encryption,
	redis redis.RedisService) PasswordEntryService {
	return &passwordEntryService{
		UserRepository:          userRepository,
		UserKeyRepository:       userKeyRepository,
		PasswordEntryRepository: passwordEntryRepository,
		PasswordEntryKeys:       passwordEntryKeysRepository,
		EncryptionService:       encryptionService,
		Redis:                   redis,
	}
}

func (s *passwordEntryService) AddPasswordEntry(passwordEntryRequest *in.PasswordEntryRequest, clientID string) error {
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
	if user == nil {
		log.Error().Str("clientID", clientID).Msg("User not found")
		return errors.New("user not found")
	}

	pubKey, err := s.UserKeyRepository.GetPublicKeyByUserID(user.UserID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user keys")
		return err
	}
	encPass, encNotes, wrappedKey, err := s.EncryptionService.EncryptPasswordEntry(passwordEntryRequest.Password, text.DerefString(passwordEntryRequest.Notes), pubKey)
	if err != nil {
		return err
	}

	passwordEntry := password.PasswordEntry{
		Title:             passwordEntryRequest.Title,
		UserID:            user.UserID,
		Username:          passwordEntryRequest.Username,
		EncryptedPassword: encPass,
		EncryptedNotes:    &encNotes,
		URL:               passwordEntryRequest.URL,
		Tags:              passwordEntryRequest.Tags,
	}

	passwordEntryKey := password.PasswordEntryKey{
		EncryptedSymmetricKey: wrappedKey,
	}

	if err := s.PasswordEntryRepository.AddPasswordEntry(&passwordEntry, &passwordEntryKey); err != nil {
		return err
	}

	return nil
}

func (s *passwordEntryService) UpdatePasswordEntry(passwordEntryID uint, passwordEntryRequest *in.PasswordEntryRequest, clientID string) error {
	passwordEntry := password.PasswordEntry{
		Title:             passwordEntryRequest.Title,
		Username:          passwordEntryRequest.Username,
		EncryptedPassword: passwordEntryRequest.Password,
		EncryptedNotes:    passwordEntryRequest.Notes,
		URL:               passwordEntryRequest.URL,
		Tags:              passwordEntryRequest.Tags,
	}

	if err := s.PasswordEntryRepository.UpdatePasswordEntry(passwordEntry); err != nil {
		return err
	}

	return nil
}

func (s *passwordEntryService) GetPasswordEntryByID(passwordEntryID uint, clientID string) (interface{}, error) {
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
	if user == nil {
		log.Error().Str("clientID", clientID).Msg("User not found")
		return nil, errors.New("user not found")
	}

	privKey, err := s.UserKeyRepository.GetPrivateKeyByUserID(user.UserID, user.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user keys")
		return nil, err
	}
	if privKey == nil {
		log.Error().Str("clientID", clientID).Msg("User public key not found")
		return nil, errors.New("user public key not found")
	}

	passwordEntry, err := s.PasswordEntryRepository.GetPasswordEntryByEntryID(passwordEntryID)
	if err != nil {
		return nil, err
	}
	if passwordEntry == nil {
		log.Error().Str("clientID", clientID).Msg("Password entry not found")
		return nil, errors.New("password entry not found")
	}
	passwordEntryKey, err := s.PasswordEntryKeys.GetPasswordEntryKeyByEntryID(passwordEntry.EntryID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password entry key")
		return nil, err
	}
	if passwordEntryKey == nil {
		log.Error().Str("clientID", clientID).Msg("Password entry key not found")
		return nil, errors.New("password entry key not found")
	}

	encPass, encNotes, err := s.EncryptionService.DecryptPasswordEntry(passwordEntry.EncryptedPassword, passwordEntryKey.EncryptedSymmetricKey, passwordEntryKey.EncryptedSymmetricKey, privKey)
	passwordEntry.EncryptedPassword = encPass
	passwordEntry.EncryptedNotes = &encNotes
	return passwordEntry, nil
}

func (s *passwordEntryService) GetListPasswordEntries(clientID string) (interface{}, error) {
	passwordEntries, err := s.PasswordEntryRepository.GetPasswordEntryByGroupID(1)
	if err != nil {
		return nil, err
	}

	return passwordEntries, nil
}

func (s *passwordEntryService) DeletePasswordEntry(passwordEntryID uint, clientID string) error {
	if err := s.PasswordEntryRepository.DeletePasswordEntry(passwordEntryID); err != nil {
		return err
	}

	return nil
}

func logError(method, clientID string, err error, message string) (interface{}, error) {
	// Log the error with the method name, client ID, and message
	log.Error().Str("method", method).Str("clientID", clientID).Err(err).Msg(message)
	return nil, err
}

func logErrorWithNoReturn(method, clientID string, err error, message string) error {
	// Log the error with the method name, client ID, and message
	log.Error().Str("method", method).Str("clientID", clientID).Err(err).Msg(message)
	return err
}
