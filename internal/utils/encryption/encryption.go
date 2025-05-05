package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type Encryption interface {
	EncryptPasswordEntry(password, notes string, pubKey *rsa.PublicKey) (string, string, string, error)
	DecryptPasswordEntry(encPassword, encNotes, wrappedAESKey string, privKey *rsa.PrivateKey) (string, string, error)
}

type encryption struct {
}

func NewEncryption() Encryption {
	return &encryption{}
}

func (e *encryption) EncryptPasswordEntry(password, notes string, pubKey *rsa.PublicKey) (string, string, string, error) {
	aesKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		return "", "", "", err
	}

	// AES encryption for password
	encPassword, err := encryptWithAES([]byte(password), aesKey)
	if err != nil {
		return "", "", "", err
	}

	// AES encryption for notes (optional)
	encNotes := ""
	if notes != "" {
		encNotes, err = encryptWithAES([]byte(notes), aesKey)
		if err != nil {
			return "", "", "", err
		}
	}

	// Wrap AES key with RSA public key
	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, aesKey, nil)
	if err != nil {
		return "", "", "", err
	}

	return encPassword, encNotes, base64.StdEncoding.EncodeToString(encryptedAESKey), nil
}

func (e *encryption) DecryptPasswordEntry(encPassword, encNotes, wrappedAESKey string, privKey *rsa.PrivateKey) (string, string, error) {
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, decode(wrappedAESKey), nil)
	if err != nil {
		return "", "", err
	}

	decPass, err := decryptAES(encPassword, aesKey)
	if err != nil {
		return "", "", err
	}

	decNotes := ""
	if encNotes != "" {
		decNotes, _ = decryptAES(encNotes, aesKey)
	}
	return decPass, decNotes, nil
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

// decryptAES is already defined in canvas
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
