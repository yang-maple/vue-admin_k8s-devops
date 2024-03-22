package utils

import (
	"math/rand"
	"time"
)

func GenRandNum(width int) string {
	numeric := []byte("0123456789")
	r := len(numeric)
	seeds := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, width)
	for i := 0; i < width; i++ {
		result[i] = numeric[seeds.Intn(r)]
	}
	return string(result)
}
