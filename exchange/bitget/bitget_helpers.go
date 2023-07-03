package bitget

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	// "strconv"

	// "github.com/kryptomind/bidboxapi/KeyService/internal/shared"
)

func GenerateBitgetSignature(apiSecret string, method string, uri string, timestamp string) string {
	message := fmt.Sprintf("%s%s%s", timestamp, method, uri)
	hmac := hmac.New(sha256.New, []byte(apiSecret))
	hmac.Write([]byte(message))
	signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
	return signature
}
