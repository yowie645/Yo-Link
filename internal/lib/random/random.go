package random

import (
	"math/rand"
	"sync"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var src struct {
	sync.Mutex
	src rand.Source
}

func init() {
	src.src = rand.NewSource(time.Now().UnixNano())
}

func NewRandomString(n int) string {
	b := make([]byte, n)

	src.Lock()

	for i, cache, remain := n-1, src.src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	src.Unlock()

	return string(b)
}
