package binance_websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/kryptomind/bidboxapi/bitgetms/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type MarketEvent struct {
	Symbol      string `json:"s"`
	MarketPrice string `json:"p"`
}

type Cache struct {
	Positions []models.Positions
}

const (
	PERCENT_CHANGE = 5
)

func (s *Server) WebsocketTest() {
	var paramsList []string

	coinPair := models.CoinPair{}
	coinPairs, err := coinPair.GetAllCoins(s.DB)
	if err != nil {
		fmt.Println("---- error fetching coins ----", err)
		return
	}

	for _, coinPair := range *coinPairs {
		result := strings.ReplaceAll(coinPair.Coin, "/", "")
		eventString := strings.ToLower(result)
		eventString = fmt.Sprintf("%s%s", eventString, "@markPrice")
		paramsList = append(paramsList, eventString)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Create a cache instance
	cache := &Cache{}

	// Start a goroutine to periodically update the cache
	go func() {
		for {
			position := models.Positions{}
			positions, err := position.GetOpenPositionsByExchange(s.DB, "binance")
			if err != nil {
				log.Println("Error fetching positions:", err)
			} else {
				cache.Positions = *positions
				log.Println("Positions cache updated successfully.")
			}
			time.Sleep(10 * time.Second) // Wait for 10 seconds before updating the cache again
		}
	}()

	for {
		conn, err := connectWebSocket()
		if err != nil {
			log.Println("WebSocket connection error:", err)
			time.Sleep(5 * time.Second) // Wait for 5 seconds before reconnecting
			continue
		}

		err = subscribeToMarketEvents(conn, paramsList)
		if err != nil {
			log.Println("WebSocket subscribe error:", err)
			conn.Close()
			time.Sleep(5 * time.Second) // Wait for 5 seconds before reconnecting
			continue
		}

		go handleWebSocketMessages(conn, s.DB, cache)

		select {
		case <-interrupt:
			log.Println("Received interrupt signal. Closing WebSocket connection...")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("WebSocket close message sending error:", err)
			}
			time.Sleep(1 * time.Second) // Wait for the server to close the connection
			return
		}
	}
}

func connectWebSocket() (*websocket.Conn, error) {
	url := "wss://fstream.binance.com/ws/"
	dialer := &websocket.Dialer{}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func subscribeToMarketEvents(conn *websocket.Conn, paramsList []string) error {
	subscribeRequest := struct {
		Method string   `json:"method"`
		Params []string `json:"params"`
		ID     int      `json:"id"`
	}{
		Method: "SUBSCRIBE",
		Params: paramsList,
		ID:     1,
	}

	err := conn.WriteJSON(subscribeRequest)
	if err != nil {
		return err
	}

	return nil
}

func handleWebSocketMessages(conn *websocket.Conn, db *gorm.DB, cache *Cache) {
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket message receiving error:", err)
			return
		}

		var eventData MarketEvent

		err = json.Unmarshal(message, &eventData)
		if err != nil {
			log.Println("WebSocket message parsing error:", err)
			continue
		}

		go handleMarketUpdate(db, cache, eventData.Symbol, eventData.MarketPrice)
	}
}

func handleMarketUpdate(db *gorm.DB, cache *Cache, coinPair string, marketPrice string) {
	// fmt.Println("🚀 ~ file: binance.go:156 ~ funchandleMarketUpdate ~ coinPair:", coinPair)

	positions := cache.Positions

	var symbolPositions []models.Positions

	for _, pos := range positions {
		if pos.Symbol == coinPair {
			symbolPositions = append(symbolPositions, pos)
		}
	}

	for _, filteredPos := range symbolPositions {
		go handlePositionOnRateUpdate(filteredPos, marketPrice)
	}

}

func handlePositionOnRateUpdate(position models.Positions, marketPrice string) {
	lastUpdatePrice := position.OpenPrice

	if position.LastUpdatePrice != "" {
		lastUpdatePrice = position.LastUpdatePrice
	}

	flEntryPrice, err := strconv.ParseFloat(lastUpdatePrice, 64)
	if err != nil {
		fmt.Println("--- unable to convert entryprice to float ---", err)
	}

	flMarketPriceFloat, err := strconv.ParseFloat(marketPrice, 64)
	if err != nil {
		fmt.Println("--- unable to convert marketPrice to float ---", err)
	}

	percentage_change := (flMarketPriceFloat - flEntryPrice) / flEntryPrice * 100

	switch position.Side {
	case "long":
		if percentage_change > 1 {
			fmt.Println(position.Symbol, "--- long position", "---- percentage change ----", percentage_change)
			take_profit := flEntryPrice * (1 - PERCENT_CHANGE/100)
			fmt.Println("--- take profit for long position at ----", take_profit)

		}

	case "short":
		if percentage_change < -1 {
			// update TP/SL here
			fmt.Println(position.Symbol, "--- short position", "---- percentage change ----", percentage_change)
			take_profit := flEntryPrice * (1 - PERCENT_CHANGE/100)
			fmt.Println("--- take profit for short position at ----", take_profit)

		}

	default:
		fmt.Println("--- defaault case reached ----")
	}

	// fmt.Println("--- last update rate ---", lastUpdatePrice)

}
