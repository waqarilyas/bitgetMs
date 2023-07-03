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
	BybitAPIEndpoint = "https://api-testnet.bybit.com"
)

type ErrorResponse struct {
	RetCode    int                    `json:"retCode"`
	RetMsg     string                 `json:"retMsg"`
	Result     map[string]interface{} `json:"result"`
	RetExtInfo map[string]interface{} `json:"retExtInfo"`
	Time       int64                  `json:"time"`
}

type BybitPosition struct {
	Symbol           string `json:"symbol"`
	Leverage         string `json:"leverage"`
	AutoAddMargin    int    `json:"autoAddMargin"`
	AvgPrice         string `json:"avgPrice"`
	LiqPrice         string `json:"liqPrice"`
	RiskLimitValue   string `json:"riskLimitValue"`
	TakeProfit       string `json:"takeProfit"`
	PositionValue    string `json:"positionValue"`
	TpslMode         string `json:"tpslMode"`
	RiskID           int    `json:"riskId"`
	TrailingStop     string `json:"trailingStop"`
	UnrealisedPnl    string `json:"unrealisedPnl"`
	MarkPrice        string `json:"markPrice"`
	AdlRankIndicator int    `json:"adlRankIndicator"`
	CumRealisedPnl   string `json:"cumRealisedPnl"`
	PositionMM       string `json:"positionMM"`
	CreatedTime      string `json:"createdTime"`
	PositionIdx      int    `json:"positionIdx"`
	PositionIM       string `json:"positionIM"`
	UpdatedTime      string `json:"updatedTime"`
	Side             string `json:"side"`
	BustPrice        string `json:"bustPrice"`
	PositionBalance  string `json:"positionBalance"`
	Size             string `json:"size"`
	PositionStatus   string `json:"positionStatus"`
	StopLoss         string `json:"stopLoss"`
	TradeMode        int    `json:"tradeMode"`
}

type PositionsResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		NextPageCursor string          `json:"nextPageCursor"`
		Category       string          `json:"category"`
		List           []BybitPosition `json:"list"`
	} `json:"result"`
	RetExtInfo struct{} `json:"retExtInfo"`
	Time       int64    `json:"time"`
}

func GetBybitAccountPositions(apiKey string, secret string, coin_pair string) (*BybitPosition, error) {
	queryString := "settleCoin=USDT&category=linear"
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	accSignature := GenerateBybitSignature(apiKey, secret, 50000, timestamp, queryString)
	finalURL := BybitAPIEndpoint + "/v5/position/list?" + queryString

	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-BAPI-SIGN-TYPE", "2")
	req.Header.Set("X-BAPI-SIGN", accSignature)
	req.Header.Set("X-BAPI-API-KEY", apiKey)
	req.Header.Set("X-BAPI-TIMESTAMP", strconv.FormatInt(timestamp, 10))
	req.Header.Set("X-BAPI-RECV-WINDOW", "50000")
	req.Header.Set("Content-Type", "application/json")

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
		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse.RetMsg)
	}

	var response PositionsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var requiredPosition BybitPosition

	for _, pos := range response.Result.List {

		if pos.Symbol == coin_pair {
			requiredPosition = pos
			break
		}
	}

	return &requiredPosition, nil
}
