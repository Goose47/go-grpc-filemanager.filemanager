package random

import (
	"math/rand"
)

func String(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var byteStr []byte
	for range n {
		byteStr = append(byteStr, charset[rand.Intn(len(charset))])
	}

	return string(byteStr)
}
