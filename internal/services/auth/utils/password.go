package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(inputPassword string, pepper string) (string, error) {
	
	// Generate random salt

	// make is used to create a byte slice of the specified length
	salt := make([]byte, saltLen)
	// rand.Read fills the slice with random bytes
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	if pepper == "" {
		return "", errors.New("password pepper is not configured")
	}
	
	pepperedPassword := inputPassword + pepper
	// hash the password with Argon2id
	// make a byte slice to hold the hash
	hash := argon2.IDKey([]byte(pepperedPassword), salt, timeCost, memoryKB, threads, keyLen)

	// Encode salt+hash into a single string we can store in DB
	// base64.RawStdEncoding is used to encode the byte slices to base64 strings
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Store all parameters inside the string so you can verify later even if params change
	// Format: argon2id$v=19$t=2$m=65536$p=1$<salt>$<hash>
	passwordHash := fmt.Sprintf(
		"argon2id$v=%d$t=%d$m=%d$p=%d$%s$%s",
		argon2.Version, timeCost, memoryKB, threads, b64Salt, b64Hash,
	)
	return passwordHash, nil

}


func VerifyPassword(userPassword string, pepper string, hashedPassword string) (bool, error) {
	fmt.Println(hashedPassword)
	// expected format:
	// argon2id$v=19$t=2$m=65536$p=1$<salt>$<hash>
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 7 {
		return false, errors.New("invalid password hash format")
	}

	var (
		timeCost uint32
		memoryKB uint32
		threads  uint8
	)

	// paste params from parts to separate variables
	// fmt.Sscanf("id=10 name=alex", "id=%d name=%s", &a, &b)
	// a == 10
	// b == "alex"
	if _, err := fmt.Sscanf(parts[2], "t=%d", &timeCost); err != nil {
		return false, err
	}
	if _, err := fmt.Sscanf(parts[3], "m=%d", &memoryKB); err != nil {
		return false, err
	}
	if _, err := fmt.Sscanf(parts[4], "p=%d", &threads); err != nil {
		return false, err
	}

	// decode base64 salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[6])
	if err != nil {
		return false, err
	}

	// hash input password + pepper
	hash := argon2.IDKey(
		[]byte(userPassword+pepper),
		salt,
		timeCost,
		memoryKB,
		threads,
		uint32(len(expectedHash)),
	)

	// constant-time compare
	// to prevent timing attacks


	if subtle.ConstantTimeCompare(hash, expectedHash) == 1 {
		return true, nil
	}

	return false, nil
}

