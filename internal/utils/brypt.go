package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	specialChars = "-@$!%*?&"
)

func GeneratePassword(length int) (string, error) {
	if length < 4 {
		return "", fmt.Errorf("password length must be at least 4")
	}

	allChars := lowerChars + upperChars + digitChars + specialChars

	password := make([]byte, 0, length)

	requiredSets := []string{
		lowerChars,
		upperChars,
		digitChars,
		specialChars,
	}

	for _, set := range requiredSets {
		ch, err := randomChar(set)
		if err != nil {
			return "", err
		}
		password = append(password, ch)
	}

	for len(password) < length {
		ch, err := randomChar(allChars)
		if err != nil {
			return "", err
		}
		password = append(password, ch)
	}

	if err := shuffle(password); err != nil {
		return "", err
	}

	return string(password), nil
}

func randomChar(charset string) (byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		return 0, err
	}
	return charset[n.Int64()], nil
}

func shuffle(data []byte) error {
	for i := len(data) - 1; i > 0; i-- {
		jv, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return err
		}
		j := int(jv.Int64())
		data[i], data[j] = data[j], data[i]
	}
	return nil
}
