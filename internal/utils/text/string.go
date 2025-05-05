package text

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"password-management-service/internal/utils"
	"password-management-service/internal/utils/redis"
	"regexp"
	"strings"
)

func ValidationTrimSpace(s string) string {
	trim := strings.TrimSpace(s)
	trim = strings.Join(strings.Fields(trim), " ") // Remove extra spaces
	return trim
}

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsername.MatchString(username) {
		return errors.New("username can only contain alphanumeric characters and underscores")
	}

	return nil
}

func GenerateInviteToken() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 32-character hex string
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func NullableStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func GetOptionalString(context *gin.Context, field string) *string {
	val := context.PostForm(field)
	if val == "" {
		return nil
	}
	return &val
}

func CheckCredentialKey(redis redis.RedisService, credential, clientID string) error {
	credentialKeyMap := struct {
		CredentialKey string `json:"credential_key"`
	}{}
	err := redis.GetData(utils.CredentialKey, clientID, &credentialKeyMap)

	log.Info().Str("checkCredential", credentialKeyMap.CredentialKey).Msg("credential key retrieved from Redis")
	log.Info().Str("credential", credential).Msg("credential key retrieved from Redis")
	if err != nil {
		log.Error().Str("credential key", credential).Err(err).Msg("Failed to retrieve credential key from Redis")
		return err
	}

	if credentialKeyMap.CredentialKey != credential {
		return errors.New("credential key not matched")
	}
	_ = redis.DeleteData(utils.CredentialKey, clientID)

	return nil
}

func DerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
