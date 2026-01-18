package security

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SaltSize   = 16
	HashSize   = 32
	Iterations = 500000
)

type PasswordHasher struct{}

func NewHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (p *PasswordHasher) Hash(password string) (string, error) {
	salt := make([]byte, SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := pbkdf2.Key([]byte(password), salt, Iterations, HashSize, sha512.New)
	hashHex := hex.EncodeToString(hash)
	saltHex := hex.EncodeToString(salt)

	return fmt.Sprintf("%s-%s", hashHex, saltHex), nil
}

func (p *PasswordHasher) Verify(password string, passwordHashed string) bool {
	parts := strings.Split(passwordHashed, "-")
	if len(parts) != 2 {
		return false
	}

	hash, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}

	salt, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	inputHash := pbkdf2.Key([]byte(password), salt, Iterations, HashSize, sha512.New)

	return subtle.ConstantTimeCompare(hash, inputHash) == 1
}
