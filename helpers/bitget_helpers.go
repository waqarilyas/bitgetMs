package helpers

import (
	"encoding/json"
	// "errors"
	"crypto/aes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"crypto/cipher"
	"errors"
	"encoding/base64"
)
type bitgetServerTimeStampResponse struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int    `json:"requestTime"`
	Data        string `json:"data"`
}
func DecryptStrings(encodedCiphertext string) (string, error) {
	keyString := os.Getenv("ENCRYPTION_PASS")

	key := []byte(keyString)
	data, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(data) < 12 {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:12], data[12:]

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
func GetBitgetServerTimeStamp() string {

	url := "https://api.bitget.com/api/spot/v1/public/time"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Print("-----error in request---", err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print("---read error---", err.Error())
	}

	var response bitgetServerTimeStampResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Print("---json unmarshal error---", err.Error())
	}

	return response.Data

}