package random

import (
	"math/rand"
	"time"
)

func String(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var byteStr []byte
	for range n {
		byteStr = append(byteStr, charset[r.Intn(len(charset))])
	}

	return string(byteStr)
}
