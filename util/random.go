package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Create a orivate RNG instance with its own seed
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomInt generates a random integer in the range [min, max]
func RandomInt(min, max int64) int64 {
	if max <= min {
		return min
	}
	return min + rng.Int63n(max-min+1)
}

// RandomString generatess a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()

}
