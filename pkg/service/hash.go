package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"yuudidi.com/pkg/utils"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

var functionMap = map[string]func([]byte, []byte) []byte{
	"Argon2": argon2di,
	"bcrypt": bcryptFunc,
	"scrypt": scryptFunc,
	"PBKDF2": pbkdf2Func,
}

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// generateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func argon2di(password, salt []byte) []byte {
	return argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
}

func bcryptFunc(password, salt []byte) []byte {
	//bcrypt keep salt inline
	result, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		utils.Log.Warnf("bcrypt hash generation failed")
		panic(err)
	}
	return result
}

func scryptFunc(password, salt []byte) []byte {
	result, err := scrypt.Key(password, salt, 32768, 8, 1, 32)
	if err != nil {
		utils.Log.Warnf("scrypt hash generation failed")
		panic(err)
	}
	return result
}

func pbkdf2Func(password, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 4096, 32, sha1.New)
}
