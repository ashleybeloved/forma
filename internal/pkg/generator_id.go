package pkg

import (
	"math/rand/v2"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func GenerateShortID(length int) string {
	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		randomIndex := rand.IntN(len(alphabet))
		sb.WriteByte(alphabet[randomIndex])
	}

	return sb.String()
}
