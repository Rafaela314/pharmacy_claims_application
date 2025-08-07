package util

import (
	mathrand "math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

const numbers = "0123456789"
const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Create a private RNG instance with its own seed
var rng = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

// RandomInt generates a random integer in the range [min, max]
func RandomInt(min, max int64) int64 {
	if max <= min {
		return min
	}
	return min + rng.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[mathrand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomString generates a random string of length n
func RandomNumericString(n int) string {
	var sb strings.Builder
	k := len(numbers)

	for i := 0; i < n; i++ {
		c := numbers[mathrand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomUUID generates a random UUID using crypto/rand
func RandomUUID() uuid.UUID {
	return uuid.New()
}

// RandomMoney generates a random amount of money as float64
func RandomMoney() float64 {
	// Generate random amount between 0.01 and 999.99
	dollars := RandomInt(0, 999)
	cents := RandomInt(0, 99)

	// Convert to float64 with 2 decimal places
	return float64(dollars) + float64(cents)/100.0
}

// RandomFloat65 generates a random float64 in the range [min, max]
func RandomFloat65(min, max float64) float64 {
	if max <= min {
		return min
	}

	// Generate random value between min and max
	rangeValue := max - min
	randomValue := float64(RandomInt(0, int64(rangeValue*100))) / 100.0

	return min + randomValue
}
