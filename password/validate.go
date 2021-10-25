package password

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"

	passwordvalidator "github.com/wagslane/go-password-validator"
)

func Validate(password string) (bool, error) {
	// e, err := getEntropy()
	// if err != nil {
	// 	return false, err
	// }

	//entropy := passwordvalidator.GetEntropy(e)

	const minEntropyBits = 60

	err := passwordvalidator.Validate(password, minEntropyBits)

	if err != nil {
		return false, err
	}
	return true, nil
}

func getEntropy() (string, error) {
	assertAvailablePRNG()
	token, err := GenerateRandomStringURLSafe(32)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		return "", err
	}
	fmt.Println(token)
	token, err = GenerateRandomString(256)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		return "", err
	}
	return token, nil
}

func assertAvailablePRNG() {
	// Assert that a cryptographically secure PRNG is available.
	// Panic otherwise.
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

// GenerateRandomStringURLSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}
