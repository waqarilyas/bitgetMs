package bitget

type AccountData struct {
	Code        string          `json:"code"`
	Msg         string          `json:"msg"`
	RequestTime int64           `json:"requestTime"`
	Data        []AccountDetail `json:"data"`
}

type AccountDetail struct {
	MarginCoin        string `json:"marginCoin"`
	Locked            string `json:"locked"`
	Available         string `json:"available"`
	CrossMaxAvailable string `json:"crossMaxAvailable"`
	FixedMaxAvailable string `json:"fixedMaxAvailable"`
	MaxTransferOut    string `json:"maxTransferOut"`
	Equity            string `json:"equity"`
	USDTEquity        string `json:"usdtEquity"`
	BTCEquity         string `json:"btcEquity"`
	CrossRiskRate     string `json:"crossRiskRate"`
	UnrealizedPL      string `json:"unrealizedPL"`
	Bonus             string `json:"bonus"`
}

type MarginData struct {
	MarginCoin        string `json:"marginCoin"`
	Symbol            string `json:"symbol"`
	HoldSide          string `json:"holdSide"`
	OpenDelegateCount string `json:"openDelegateCount"`
	Margin            string `json:"margin"`
	Available         string `json:"available"`
	Locked            string `json:"locked"`
	Total             string `json:"total"`
	Leverage          int    `json:"leverage"`
	AchievedProfits   string `json:"achievedProfits"`
	AverageOpenPrice  string `json:"averageOpenPrice"`
	MarginMode        string `json:"marginMode"`
	HoldMode          string `json:"holdMode"`
	UnrealizedPL      string `json:"unrealizedPL"`
	LiquidationPrice  string `json:"liquidationPrice"`
	KeepMarginRate    string `json:"keepMarginRate"`
	MarketPrice       string `json:"marketPrice"`
	CTime             string `json:"cTime"`
}

type MarginDataResponse struct {
	Code        string       `json:"code"`
	Msg         string       `json:"msg"`
	RequestTime int64        `json:"requestTime"`
	Data        []MarginData `json:"data"`
}

type ApiKeyResponse struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
	Data        User   `json:"data"`
}

type User struct {
	UserId           string   `json:"user_id"`
	InviterId        string   `json:"inviter_id"`
	AgentInviterCode string   `json:"agent_inviter_code"`
	Channel          string   `json:"channel"`
	Ips              string   `json:"ips"`
	Authorities      []string `json:"authorities"`
	ParentId         int64    `json:"parentId"`
	Trader           bool     `json:"trader"`
	IsSpotTrader     bool     `json:"isSpotTrader"`
}
