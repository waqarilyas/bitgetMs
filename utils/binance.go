package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	BinanceAPIEndpoint = "https://testnet.binancefuture.com"
)

type Position struct {
	Symbol                 string `json:"symbol"`
	InitialMargin          string `json:"initialMargin"`
	MaintMargin            string `json:"maintMargin"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	Leverage               string `json:"leverage"`
	Isolated               bool   `json:"isolated"`
	EntryPrice             string `json:"entryPrice"`
	MaxNotional            string `json:"maxNotional"`
	PositionSide           string `json:"positionSide"`
	PositionAmt            string `json:"positionAmt"`
	Notional               string `json:"notional"`
	IsolatedWallet         string `json:"isolatedWallet"`
	UpdateTime             int    `json:"updateTime"`
	BidNotional            string `json:"bidNotional"`
	AskNotional            string `json:"askNotional"`
	LiquidationPrice       string `json:"liquidationPrice"`
	MarkPrice              string `json:"markPrice"`
	MarginType             string `json:"marginType"`
}

type BinanceErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type BinancePositionMode struct {
	DualSidePosition bool `json:"dualSidePosition"`
}

func GetBinanceAccountOpenPositions(apiKey string, secret string, coinPair string) (*Position, error) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	params := map[string]string{
		"timestamp": strconv.FormatInt(timestamp, 10),
	}

	accSignature := GenerateBinanceSignature(params, secret)
	finalURL := BinanceAPIEndpoint + "/fapi/v2/positionRisk?timestamp=" + strconv.FormatInt(timestamp, 10) + "&signature=" + accSignature

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MBX-APIKEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse BinanceErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse.Msg)
	}

	var positions []Position

	err = json.Unmarshal(body, &positions)
	if err != nil {
		return nil, err
	}

	var requiredPosition Position
	for _, pos := range positions {

		if pos.Symbol == coinPair {
			requiredPosition = pos
			break
		}
	}

	return &requiredPosition, nil
}

func GetBinanceAccountPositionMode(apiKey string, secret string) (*BinancePositionMode, error) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	params := map[string]string{
		"timestamp": strconv.FormatInt(timestamp, 10),
	}

	accSignature := GenerateBinanceSignature(params, secret)
	finalURL := BinanceAPIEndpoint + "/dapi/v1/positionSide/dual?timestamp=" + strconv.FormatInt(timestamp, 10) + "&signature=" + accSignature

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MBX-APIKEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse BinanceErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse.Msg)
	}

	var positionMode BinancePositionMode

	err = json.Unmarshal(body, &positionMode)
	if err != nil {
		return nil, err
	}

	// var requiredPosition Position
	// for _, pos := range positions {

	// 	if pos.Symbol == coinPair {
	// 		requiredPosition = pos
	// 		break
	// 	}
	// }

	return &positionMode, nil
}
