package helper

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateId(length int) string {
	letter := []rune("qwertyuiopasdfghjklzxcvbnm")
	result := make([]rune, length)

	for i := range result {
		result[i] = letter[rand.Intn(len(letter))]
	}

	return string(result)
}
