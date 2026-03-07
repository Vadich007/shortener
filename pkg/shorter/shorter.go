package shorter

import (
	"crypto/md5"
	"encoding/hex"
)

func Shorten(originalLink string) string {
	hash := md5.Sum([]byte(originalLink))
	shortCode := hex.EncodeToString(hash[:])
	return shortCode
}
