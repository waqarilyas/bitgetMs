package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

func GenerateBitgetSignature(apiSecret string, method string, uri string, timestamp string) string {
	message := fmt.Sprintf("%s%s%s", timestamp, method, uri)
	hmac := hmac.New(sha256.New, []byte(apiSecret))
	hmac.Write([]byte(message))
	signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
	return signature
}

func GenerateBinanceSignature(params map[string]string, secretKey string) string {
	var queryString string
	for key, value := range params {
		queryString += key + "=" + url.QueryEscape(value) + "&"
	}

	queryString = strings.TrimSuffix(queryString, "&")
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(queryString))
	signature := hex.EncodeToString(mac.Sum(nil))
	return signature
}

func GenerateBybitSignature(apiKey, apiSecret string, recvWindow, timestamp int64, queryString string) string {
	dataToSign := fmt.Sprintf("%d%s%d%s", timestamp, apiKey, recvWindow, queryString)
	hmacKey := []byte(apiSecret)
	hmacHash := hmac.New(sha256.New, hmacKey)
	hmacHash.Write([]byte(dataToSign))
	signature := hex.EncodeToString(hmacHash.Sum(nil))
	return signature
}
