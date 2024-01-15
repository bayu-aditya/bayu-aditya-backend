package util

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}
