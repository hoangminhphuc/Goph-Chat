package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/hoangminhphuc/goph-chat/common"
	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt(length int) (string, error) {

	switch {
	case length <= 0:
		length = common.DefaultSaltLength
	case length > common.MaxSaltLength:
		length = common.MaxSaltLength
	}

	saltBytes := make([]byte, length)
	// fills the slice with random unpredictable bytes by crypto/rand
	if _, err := rand.Read(saltBytes); err != nil {
		return "", err
	}

	// Encode using URL-safe base64 without padding
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(saltBytes), nil
}

func GenerateFixedSalt(seed string, length int) (string, error) {
	if length <= 0 {
		length = common.DefaultSaltLength
	} else if length > common.MaxSaltLength {
		length = common.MaxSaltLength
	}

	hash := sha256.Sum256([]byte(seed))
	salt := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])

	if len(salt) > length {
		salt = salt[:length]
	} else if len(salt) < length {
		padding := make([]byte, length-len(salt))
		salt += base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(padding)[:length-len(salt)]
	}

	return salt, nil
}


type Hasher interface {
	Hash(plainPassword string) string 
	Compare(hashedValue, plainText string) bool
}

type BcryptHash struct{}

func NewBcryptHash() *BcryptHash {
	return &BcryptHash{}
}

func (*BcryptHash) Hash(plainPassword string) string {
	// Set cost = 12. Higher cost = more secure, but slower.

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), common.BCRYPT_COST)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(hashedBytes)
}

func (*BcryptHash) Compare(hashedValue, plainText string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(plainText))
	return err == nil
}