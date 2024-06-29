package testhelpers

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	return RandomStringN(10)
}

func RandomStringN(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandomInt32() int32 {
	return int32(RandomInt())
}

func RandomInt() int {
	return RandomIntRange(1, 10000)
}

func RandomIntRange(min, max int) int {
	return min + rand.Intn(max-min)
}
