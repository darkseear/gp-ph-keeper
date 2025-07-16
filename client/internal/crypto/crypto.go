package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"log"

	"golang.org/x/crypto/pbkdf2"
)

// CryptoService - структура с ключем криптографии.
type CryptoService struct {
	key []byte
}

const (
	SALT_BYTE_SIZE    = 16
	HASH_BYTE_SIZE    = 32
	PBKDF2_ITERATIONS = 10000
)

// New -  экземпляр с созданием ключа.
func New(masterPassword string) *CryptoService {
	salt := make([]byte, SALT_BYTE_SIZE)
	if _, err := rand.Read(salt); err != nil {
		log.Fatalf("failed to generate salt: %v\n", err)
	}
	key := pbkdf2.Key([]byte(masterPassword), salt, PBKDF2_ITERATIONS, HASH_BYTE_SIZE, sha256.New)
	return &CryptoService{key: key}
}

// Encrypt - шифруем.
func (c *CryptoService) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt - расшифровываем.
func (c *CryptoService) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		data[:gcm.NonceSize()],
		data[gcm.NonceSize():],
		nil,
	)
}
