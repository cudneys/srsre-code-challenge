package tools

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

// Gets the SHA256 Sum
func getSha256Sum(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

// Gets the SHA512 Sum
func getSha512Sum(password string) string {
	h := sha512.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func GetSum(password, sumType string) (string, error) {
	if sumType == "256" {
		return getSha256Sum(password), nil
	} else if sumType == "512" {
		return getSha512Sum(password), nil
	} else {
		return "", errors.New(fmt.Sprintf("Invalid SHA Sum Type: %s", sumType))
	}
}

func GetEnvValue(name, defaultValue string) string {
	v, e := os.LookupEnv(name)
	if !e {
		return defaultValue
	}
	return v
}
