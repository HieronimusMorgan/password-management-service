package repository

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
	"password-management-service/internal/models/user"
	"password-management-service/internal/utils"
)

type UserKeysRepository interface {
	AddUserKey(key *user.UserKey) error
	GetUserKeys(userID uint) (*user.UserKey, error)
	GetPublicKeyByUserID(userID uint) (*rsa.PublicKey, error)
	GetPrivateKeyByUserID(userID uint, clientID string) (*rsa.PrivateKey, error)
}

type userKeysRepository struct {
	db gorm.DB
}

func NewUserKeysRepository(db gorm.DB) UserKeysRepository {
	return &userKeysRepository{
		db: db,
	}
}

func (r *userKeysRepository) AddUserKey(key *user.UserKey) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableUserKeyName).Create(&key).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *userKeysRepository) GetUserKeys(userID uint) (*user.UserKey, error) {
	var userKey user.UserKey
	if err := r.db.Table(utils.TableUserKeyName).Where("user_id = ?", userID).First(&userKey).Error; err != nil {
		return &userKey, err
	}
	return &userKey, nil
}

func (r *userKeysRepository) GetPublicKeyByUserID(userID uint) (*rsa.PublicKey, error) {
	var userKey user.UserKey
	if err := r.db.Table(utils.TableUserKeyName).Where("user_id = ?", userID).First(&userKey).Error; err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(userKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PublicKey(decoded)
}

func (r *userKeysRepository) GetPrivateKeyByUserID(userID uint, clientID string) (*rsa.PrivateKey, error) {
	var userKey user.UserKey
	if err := r.db.Table(utils.TableUserKeyName).Where("user_id = ?", userID).First(&userKey).Error; err != nil {
		return nil, err
	}

	salt, _ := base64.StdEncoding.DecodeString(userKey.Salt)
	encPrivKey, _ := base64.StdEncoding.DecodeString(userKey.EncryptedPrivateKey)

	aesKey := argon2.IDKey(
		[]byte(clientID), // Input
		salt,             // Salt
		5,                // Time (iterations)
		128*1024,         // Memory (128 MB)
		8,                // Threads (parallelism)
		32,               // Output key size
	)
	block, _ := aes.NewCipher(aesKey)
	gcm, _ := cipher.NewGCM(block)
	nonce := encPrivKey[:gcm.NonceSize()]
	cipherText := encPrivKey[gcm.NonceSize():]
	privDER, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(privDER)
}
