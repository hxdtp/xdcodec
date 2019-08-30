package xdcodec

import (
	crand "crypto/rand"
)

func RandString(n int) string {
	p := make([]byte, n, n)
	_, _ = crand.Read(p)
	return string(p)
}
