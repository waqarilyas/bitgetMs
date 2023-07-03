package controllers

import (
	"net/http"
	"sync"

	"github.com/kryptomind/bidboxapi/bitgetms/response"
)

var mutex sync.Mutex

type Account struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
	Data        []struct {
		MarginCoin        string `json:"marginCoin"`
		Locked            string `json:"locked"`
		Available         string `json:"available"`
		CrossMaxAvailable string `json:"crossMaxAvailable"`
		FixedMaxAvailable string `json:"fixedMaxAvailable"`
		MaxTransferOut    string `json:"maxTransferOut"`
		Equity            string `json:"equity"`
		UsdtEquity        string `json:"usdtEquity"`
		BtcEquity         string `json:"btcEquity"`
		CrossRiskRate     string `json:"crossRiskRate"`
		UnrealizedPL      string `json:"unrealizedPL"`
		Bonus             string `json:"bonus"`
	} `json:"data"`
}

type User struct {
	Capital      int
	Trade_amount float64
	First_order  float64
}

type TradeRequest struct {
	CoinPair string `json:"coin_pair"`
	Long     int    `json:"long"`
	Exchange string `json:"exchange"`
}

const (
	TAKE_PROFIT_PERCENTAGE = 20.0
	STOP_LOSS_PERCENTAGE   = 20.0
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "Trade Service")
}
