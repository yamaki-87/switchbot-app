package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

func NewNonce() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

func CreateSignature(token, secret, nonce string, timestamp int64) string {
	data := fmt.Sprintf("%s%d%s", token, timestamp, nonce)

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(data))

	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	sign = strings.ToUpper(sign)

	return sign
}
