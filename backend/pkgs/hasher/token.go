package hasher

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"math/big"
	"strconv"
	"strings"
)

const prefix = "pc_"

type Token struct {
	Raw  string
	Hash []byte
}

func generateToken(bits int) Token {
	randomBytes := make([]byte, bits)
	_, _ = rand.Read(randomBytes)

	plainText := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	plainText = prefix + strings.ToLower(plainText)

	return Token{
		Raw:  plainText,
		Hash: HashToken(plainText),
	}
}

// GenerateCode generates a 6 digit number, generally used for user verification or some similar
// two-factor type verification.
//
// This function will panic if [rand.Int] returns an error - in practice this is basically impossible.
func GenerateCode() Token {
	max := big.NewInt(900000) // Upper bound (exclusive) for random number
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic("integrity error, failed to generate random number: " + err.Error())
	}

	v := strconv.Itoa(int(n.Int64()) + 100000) // Ensure 6-digit number by adding 100000

	return Token{
		Raw:  v,
		Hash: HashToken(v),
	}
}

func GenerateShortToken() Token { return generateToken(16) }

// GenerateToken generates a new token prefixed with set prefix
func GenerateToken() Token { return generateToken(64) }

// HashToken hashes a token using SHA256
func HashToken(plainTextToken string) []byte {
	hash := sha256.Sum256([]byte(plainTextToken))
	return hash[:]
}
