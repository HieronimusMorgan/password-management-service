package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"password-management-service/internal/models/user"
	"time"
)

type Encryption interface {
	GenerateUserKey(user *user.Users) (*user.UserKey, error)
	EncryptPasswordEntry(username, password, notes string, pubKey *rsa.PublicKey) (string, string, string, string, error)
	DecryptPasswordEntry(encUsername, encPassword, encNotes, wrappedAESKey string, privateKey *rsa.PrivateKey) (string, string, string, error)
}

type encryption struct {
}

func NewEncryption() Encryption {
	return &encryption{}
}

func (e *encryption) GenerateUserKey(data *user.Users) (*user.UserKey, error) {
	salt := make([]byte, 32)
	n, err := rand.Read(salt)
	if err != nil || n != 32 {
		return nil, fmt.Errorf("failed to generate secure salt: %w", err)
	}

	aesKey := argon2.IDKey(
		[]byte(data.ClientID),
		salt,
		5,
		128*1024,
		8,
		32,
	)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, privateKeyBytes, nil)

	return &user.UserKey{
		UserID:              data.UserID,
		PublicKey:           base64.StdEncoding.EncodeToString(publicKeyBytes),
		EncryptedPrivateKey: base64.StdEncoding.EncodeToString(ciphertext),
		EncryptionAlgorithm: "RSA-2048 + AES-GCM + Argon2id",
		Salt:                base64.StdEncoding.EncodeToString(salt),
		CreatedAt:           time.Now(),
		CreatedBy:           &data.ClientID,
		UpdatedAt:           time.Now(),
		UpdatedBy:           &data.UpdatedBy,
	}, nil
}

func (e *encryption) EncryptPasswordEntry(username, password, notes string, pubKey *rsa.PublicKey) (string, string, string, string, error) {
	aesKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		return "", "", "", "", err
	}

	encUsername, err := encryptWithAES([]byte(username), aesKey)
	if err != nil {
		return "", "", "", "", err
	}

	encPassword, err := encryptWithAES([]byte(password), aesKey)
	if err != nil {
		return "", "", "", "", err
	}

	encNotes := ""
	if notes != "" {
		encNotes, err = encryptWithAES([]byte(notes), aesKey)
		if err != nil {
			return "", "", "", "", err
		}
	}

	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, aesKey, nil)
	if err != nil {
		return "", "", "", "", err
	}

	return encUsername, encPassword, encNotes, base64.StdEncoding.EncodeToString(encryptedAESKey), nil
}

func (e *encryption) DecryptPasswordEntry(encUsername, encPassword, encNotes, wrappedAESKey string, privateKey *rsa.PrivateKey) (string, string, string, error) {
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, decode(wrappedAESKey), nil)
	if err != nil {
		return "", "", "", err
	}

	decUsername, err := decryptAES(encUsername, aesKey)
	if err != nil {
		return "", "", "", err
	}

	decPass, err := decryptAES(encPassword, aesKey)
	if err != nil {
		return "", "", "", err
	}

	decNotes := ""
	if encNotes != "" {
		decNotes, err = decryptAES(encNotes, aesKey)
		if err != nil {
			return "", "", "", err
		}
	}
	return decUsername, decPass, decNotes, nil
}

func encryptWithAES(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decode(b64 string) []byte {
	decoded, _ := base64.StdEncoding.DecodeString(b64)
	return decoded
}

func decryptAES(b64cipher string, key []byte) (string, error) {
	ciphertext := decode(b64cipher)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	cipherData := ciphertext[nonceSize:]
	plain, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
