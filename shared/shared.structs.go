package shared

type AccountData struct {
	Available     float64 `json:"available"`
	Equity        float64 `json:"equity"`
	MarginBalance float64 `json:"marginBalance"`
	UnrealizedPL  float64 `json:"unrealizedPl"`
}

type DecryptedKeys struct {
	ApiKey     string `json:"api_key"`
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase"`
}

type PositionsData struct {
	Symbol           string `json:"symbol"`
	MarginCoin       string `json:"marginCoin"`
	HoldSide         string `json:"holdSide"`
	Margin           string `json:"margin"`
	Available        string `json:"available"`
	Total            string `json:"total"`
	MarginMode       string `json:"marginMode"`
	HoldMode         string `json:"holdMode"`
	LiquidationPrice string `json:"liquidationPrice"`
	MarketPrice      string `json:"marketPrice"`
	EntryPrice       string `json:"entryPrice"`
	CreationTime     string `json:"creationTime"`
	UnrealizedPL     string `json:"unrealizedPl"`
	Leverage         string `json:"leverage"`
}
