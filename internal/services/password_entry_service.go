package services

import (
	"errors"
	"github.com/rs/zerolog/log"
	"password-management-service/internal/dto/in"
	"password-management-service/internal/dto/out"
	"password-management-service/internal/models/password"
	"password-management-service/internal/repository"
	"password-management-service/internal/utils"
	"password-management-service/internal/utils/encryption"
	"password-management-service/internal/utils/redis"
	"password-management-service/internal/utils/text"
)

type PasswordEntryService interface {
	AddPasswordEntry(passwordEntryRequest *in.PasswordEntryRequest, clientID string, requestID string) error
	UpdatePasswordEntry(passwordEntryID uint, passwordEntryRequest *in.PasswordEntryRequest, clientID string) error
	AddGroupPasswordEntry(req struct {
		GroupID uint `json:"group_id"`
		EntryID uint `json:"entry_id"`
	}, clientID string) error
	GetPasswordEntryByID(passwordEntryID uint, clientID string) (interface{}, error)
	GetListPasswordEntries(clientID string, tags string, index int, size int) (interface{}, int64, error)
	DeletePasswordEntry(passwordEntryID uint, clientID string) error
}

type passwordEntryService struct {
	UserRepository             repository.UserRepository
	UserKeyRepository          repository.UserKeysRepository
	PasswordEntryRepository    repository.PasswordEntryRepository
	PasswordEntryKeyRepository repository.PasswordEntryKeysRepository
	PasswordTagRepository      repository.PasswordTagRepository
	PasswordGroupRepository    repository.PasswordGroupRepository
	EncryptionService          encryption.Encryption
	Redis                      redis.RedisService
}

func NewPasswordEntryService(
	userRepository repository.UserRepository,
	userKeyRepository repository.UserKeysRepository,
	passwordEntryRepository repository.PasswordEntryRepository,
	passwordEntryKeysRepository repository.PasswordEntryKeysRepository,
	PasswordTagRepository repository.PasswordTagRepository,
	PasswordGroupRepository repository.PasswordGroupRepository,
	encryptionService encryption.Encryption,
	redis redis.RedisService) PasswordEntryService {
	return &passwordEntryService{
		UserRepository:             userRepository,
		UserKeyRepository:          userKeyRepository,
		PasswordEntryRepository:    passwordEntryRepository,
		PasswordEntryKeyRepository: passwordEntryKeysRepository,
		PasswordTagRepository:      PasswordTagRepository,
		PasswordGroupRepository:    PasswordGroupRepository,
		EncryptionService:          encryptionService,
		Redis:                      redis,
	}
}

func (s *passwordEntryService) AddPasswordEntry(passwordEntryRequest *in.PasswordEntryRequest, clientID string, verifyCode string) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve data from Redis")
		return err
	}

	// Verify the code
	var verify *out.VerifyPinCodeResponse
	if err := s.Redis.GetData(utils.PinVerify, data.ClientID, &verify); err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to verify code")
		return err
	}

	if verify.RequestID != verifyCode {
		log.Error().Str("clientID", clientID).Msg("Invalid verification code")
		return errors.New("invalid verification code")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user by client ID")
		return err
	}

	key, err := s.UserKeyRepository.GetUserKeys(user.UserID)
	if key == nil && err != nil {
		userKey, err := s.EncryptionService.GenerateUserKey(user)
		if err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("Failed to generate user key pair")
			return err
		}
		if err := s.UserKeyRepository.AddUserKey(userKey); err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("Failed to add user key")
			return err
		}
	}

	publicKey, err := s.UserKeyRepository.GetPublicKeyByUserID(user.UserID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user keys")
		return err
	}

	encryptedUsername, encryptPassword, notes, wrappedKey, err := s.EncryptionService.EncryptPasswordEntry(passwordEntryRequest.Username, passwordEntryRequest.Password, text.DerefString(passwordEntryRequest.Notes), publicKey)
	if err != nil {
		return err
	}

	passwordEntry := password.PasswordEntry{
		Title:             passwordEntryRequest.Title,
		UserID:            user.UserID,
		Username:          encryptedUsername,
		EncryptedPassword: encryptPassword,
		EncryptedNotes:    &notes,
		URL:               passwordEntryRequest.URL,
		CreatedBy:         &clientID,
		UpdatedBy:         &clientID,
	}

	passwordEntryKey := password.PasswordEntryKey{
		EncryptedSymmetricKey: wrappedKey,
	}

	if err := s.PasswordEntryRepository.AddPasswordEntry(&passwordEntry, &passwordEntryKey, *passwordEntryRequest.Tags, user.UserID); err != nil {
		return err
	}

	return nil
}

func (s *passwordEntryService) UpdatePasswordEntry(passwordEntryID uint, passwordEntryRequest *in.PasswordEntryRequest, clientID string) error {
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

	key, err := s.UserKeyRepository.GetUserKeys(user.UserID)
	if key == nil && err != nil {
		userKey, err := s.EncryptionService.GenerateUserKey(user)
		if err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("Failed to generate user key pair")
			return err
		}
		if err := s.UserKeyRepository.AddUserKey(userKey); err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("Failed to add user key")
			return err
		}
	}

	entry, err := s.PasswordEntryRepository.GetPasswordEntryByEntryIDAndUserID(passwordEntryID, user.UserID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password entry")
		return err
	}

	publicKey, err := s.UserKeyRepository.GetPublicKeyByUserID(user.UserID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user keys")
		return err
	}

	encryptedUsername, encryptPassword, notes, wrappedKey, err := s.EncryptionService.EncryptPasswordEntry(passwordEntryRequest.Username, passwordEntryRequest.Password, text.DerefString(passwordEntryRequest.Notes), publicKey)
	if err != nil {
		return err
	}

	passwordTags, err := s.PasswordTagRepository.GetPasswordTagsByEntryID(entry.EntryID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password tags")
	}

	passwordEntry := password.PasswordEntry{
		EntryID:           entry.EntryID,
		Title:             passwordEntryRequest.Title,
		UserID:            user.UserID,
		Username:          encryptedUsername,
		EncryptedPassword: encryptPassword,
		EncryptedNotes:    &notes,
		URL:               passwordEntryRequest.URL,
		Tags:              passwordTags,
		UpdatedBy:         &clientID,
	}

	passwordEntryKey := password.PasswordEntryKey{
		EntryID:               entry.EntryID,
		EncryptedSymmetricKey: wrappedKey,
	}

	if err := s.PasswordEntryRepository.UpdatePasswordEntryAndEntryKey(passwordEntry, passwordEntryKey); err != nil {
		return err
	}
	return nil
}

func (s *passwordEntryService) AddGroupPasswordEntry(req struct {
	GroupID uint `json:"group_id"`
	EntryID uint `json:"entry_id"`
}, clientID string) error {
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

	entry, err := s.PasswordEntryRepository.GetPasswordEntryByEntryIDAndUserID(req.EntryID, user.UserID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password entry")
		return err
	}
	if entry == nil {
		log.Error().Str("clientID", clientID).Msg("Password entry not found")
		return errors.New("password entry not found")
	}

	group, err := s.PasswordGroupRepository.GetPasswordGroupByUserIDAndGroupID(user.UserID, req.GroupID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password group")
		return err
	}

	if group == nil {
		log.Error().Str("clientID", clientID).Msg("Password group not found")
		return errors.New("password group not found")
	}

	entry.GroupID = &group.GroupID
	if err := s.PasswordEntryRepository.UpdatePasswordEntry(entry); err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to update password entry")
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

	passwordEntry, err := s.PasswordEntryRepository.GetPasswordEntryByEntryIDAndUserID(passwordEntryID, user.UserID)
	if err != nil {
		return nil, err
	}
	if passwordEntry == nil {
		log.Error().Str("clientID", clientID).Msg("Password entry not found")
		return nil, errors.New("password entry not found")
	}

	privateKey, err := s.UserKeyRepository.GetPrivateKeyByUserID(user.UserID, user.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user keys")
		return nil, err
	}
	if privateKey == nil {
		log.Error().Str("clientID", clientID).Msg("User public key not found")
		return nil, errors.New("user public key not found")
	}

	passwordEntryKey, err := s.PasswordEntryKeyRepository.GetPasswordEntryKeyByEntryID(passwordEntry.EntryID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve password entry key")
		return nil, err
	}
	if passwordEntryKey == nil {
		log.Error().Str("clientID", clientID).Msg("Password entry key not found")
		return nil, errors.New("password entry key not found")
	}

	encUsername, encPass, encNotes, err := s.EncryptionService.DecryptPasswordEntry(passwordEntry.Username, passwordEntry.EncryptedPassword, *passwordEntry.EncryptedNotes, passwordEntryKey.EncryptedSymmetricKey, privateKey)
	passwordEntry.Username = encUsername
	passwordEntry.EncryptedPassword = encPass
	passwordEntry.EncryptedNotes = &encNotes
	return passwordEntry, nil
}

func (s *passwordEntryService) GetListPasswordEntries(clientID string, tags string, index int, size int) (interface{}, int64, error) {
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
	if user == nil {
		log.Error().Str("clientID", clientID).Msg("User not found")
		return nil, 0, errors.New("user not found")
	}

	totalPasswordEntries, err := s.PasswordEntryRepository.GetCountPasswordEntriesByUserID(user.UserID)

	passwordEntries, err := s.PasswordEntryRepository.GetListPasswordEntryResponse(user.UserID, tags, index, size)
	if err != nil {
		return nil, 0, err
	}

	return passwordEntries, totalPasswordEntries, nil
}

func (s *passwordEntryService) DeletePasswordEntry(passwordEntryID uint, clientID string) error {
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
	if err := s.PasswordEntryRepository.DeletePasswordEntry(passwordEntryID); err != nil {
		return err
	}

	return nil
}
