package utils

import (
	"math/rand"
	"time"
)

var letters = func() []rune {
	alphabet := []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	for i := 'A'; i <= 'Z'; i++ {
		alphabet = append(alphabet, i)
	}
	for i := 'a'; i <= 'z'; i++ {
		alphabet = append(alphabet, i)
	}
	return alphabet
}()

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().Unix() + rand.Int63())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
