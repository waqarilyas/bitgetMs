package utils

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	// "io/ioutil"

	helpers "github.com/WAQAR5/bitget-helpers"
)

func initCoinPiars() ([]string, error) {

	var coin_pairs []string

	// Open the CSV file
	file, err := os.Open("test.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}, err
	}

	// Iterate over the records and print the values
	for _, record := range records {
		for i, value := range record {
			if i == 1 {
				coin_pairs = append(coin_pairs, value)
			}
		}
	}
	return coin_pairs, nil
}

func AiStub(orders int) ([]string, []string, error) {
	rand.Seed(time.Now().UnixNano()) // initialize the random number generator with the current time

	// define an array of integers
	coins := []string{"SETHSUSDT_SUMCBL"}
	// create a slice to hold the selected elements
	list := make([]string, orders)

	// generate the random indices and select the corresponding elements
	for i := 0; i < orders; i++ {
		index := rand.Intn(len(coins)) // generate a random index within the range of the array
		list[i] = coins[index]         // select the element at the random index and add it to the selected slice
	}

	middle := orders / 2
	return list[middle:], list[:middle], nil
}

type IndexPriceReq struct {
	Code string `json:"code"`
	Data struct {
		Symbol    string `json:"symbol"`
		Index     string `json:"index"`
		Timestamp string `json:"timestamp"`
	} `json:"data"`
	Msg         string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
}

type IndexPriceReqBybit struct {
	RetCode    int         `json:"retCode"`
	RetMsg     string      `json:"retMsg"`
	Result     ResultData  `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type ResultData struct {
	List     []SymbolDetails `json:"list"`
}
type SymbolDetails struct {
	MarkPrice               string `json:"markPrice"`
}

type IndexPriceReqBinance struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}

func GetSize(symbol string, first_order float64) (float64, float64, error) {
	url := "https://api.bitget.com/api/mix/v1/market/index?symbol=" + symbol

	client := http.Client{}

	res, err := client.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer res.Body.Close()

	price := IndexPriceReq{}
	if err := json.NewDecoder(res.Body).Decode(&price); err != nil {
		return 0, 0, err
	}

	if price.Code != "00000" {
		return 0, 0, errors.New(price.Msg)
	}

	fprice, err := strconv.ParseFloat(price.Data.Index, 64)
	if err != nil {
		return 0, 0, err
	}

	tradeAmount := first_order / fprice

	decimalMultiplier := math.Pow(10, float64(3))
	fixedAmount := math.Round(tradeAmount*decimalMultiplier) / decimalMultiplier

	return fixedAmount, fprice, nil
}

func GetSizeBybit(symbol string, first_order float64) (float64, float64, error) {
	url := "https://api-testnet.bybit.com/v5/market/tickers?category=inverse&symbol=" + symbol

	client := http.Client{}

	res, err := client.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer res.Body.Close()

	price := IndexPriceReqBybit{}
	if err := json.NewDecoder(res.Body).Decode(&price); err != nil {
		return 0, 0, err
	}

	if price.RetCode != 0 {
		return 0, 0, errors.New(price.RetMsg)
	}

	fprice, err := strconv.ParseFloat(price.Result.List[0].MarkPrice, 64)
	if err != nil {
		return 0, 0, err
	}

	return first_order / fprice, fprice, nil
}

func BinanceRequest(symbol string) (string, error) {
	url := "https://fapi.binance.com/fapi/v1/ticker/price?symbol=" + symbol
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()
	price := IndexPriceReqBinance{}
	if err := json.NewDecoder(res.Body).Decode(&price); err != nil {
		return "", err
	}
	return price.Price, nil
}
func GetBinanceSize(symbol string, first_order float64) (float64, float64, error) {
	binancePrice, err := BinanceRequest(symbol)
	if err != nil {
		return 0, 0, err
	}
	fprice, err := strconv.ParseFloat(binancePrice, 64)
	if err != nil {
		return 0, 0, err
	}
	return first_order / fprice, fprice, nil
}

func CheckBalance(val int) error {
	if val < 200 {
		return errors.New("Balance should be at least 200 USDT")
	}
	return nil
}

func DecryptKeys(api_key string, secret_key string, passphrase string, service string) (string, string, string, error) {
	api_key, err := helpers.DecryptStrings(api_key)
	fmt.Println("in decrypt keys function")
	if err != nil {
		log.Fatal(err)
		return "", "", "", err
	}
	secret_key, err = helpers.DecryptStrings(secret_key)
	if err != nil {
		log.Fatal(err)
		return "", "", "", err
	}
	if service == "bitget" || service == "okx" {
		passphrase, err = helpers.DecryptStrings(passphrase)
		if err != nil {
			log.Fatal(err)
			return "", "", "", err
		}
	} else if service == "binance" {
		passphrase = ""
	}

	return api_key, secret_key, passphrase, nil
}

func ConvertStrToFloat64(val string) float64 {
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0.0
	}

	return num
}

func ConvertStrToFloat32(val string) float32 {
	num64, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0.0
	}

	num32 := float32(num64)

	if !math.IsInf(float64(num32), 0) {
		return num32
	}

	return 0.0
}
