package object

import (
	"fmt"
	"crypto/sha1"
)

func Hash(data []byte) string {
	hash := sha1.Sum(data)
	return fmt.Sprintf("%x", hash)
}