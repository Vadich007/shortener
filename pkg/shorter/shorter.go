package shorter

import (
	"crypto/md5"
	"encoding/hex"
)

func Shorten(originalURL string) string {
	hash := md5.Sum([]byte(originalURL))
	shortCode := hex.EncodeToString(hash[:])
	return shortCode
}
