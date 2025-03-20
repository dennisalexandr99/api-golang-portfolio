package utility

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func StringToHMACSHA256(message, secret string) string {
	hmacHash := hmac.New(sha256.New, []byte(secret))
	hmacHash.Write([]byte(message))
	hash := hmacHash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash)
}
