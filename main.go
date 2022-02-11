package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	ws "github.com/gorilla/websocket"
	"github.com/spolia/vwap/internal"
	"github.com/spolia/vwap/internal/vwapcalculator"
	"github.com/spolia/vwap/internal/websocket"
	"github.com/spolia/vwap/internal/websocket/coinbase"
)

const (
	defaultTradingPairs = "BTC-USD,ETH-USD,ETH-BTC"
	defaultWindowSize   = 200
)

func main() {
	ctx := context.Background()

	fmt.Println("trading pair: ", defaultTradingPairs, "subscribe to coinbase websocket url:", coinbase.DefaultURL, "with a window size: ", defaultWindowSize)
	tradingPairsArr := strings.Split(defaultTradingPairs, ",")

	conn, _, err := ws.DefaultDialer.Dial(coinbase.DefaultURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	list, err := vwapcalculator.NewList([]vwapcalculator.DataPoint{}, defaultWindowSize)
	if err != nil {
		log.Fatal(err)
	}

	service := internal.NewService(coinbase.NewClient(conn), tradingPairsArr, &list)
	tradingReceiver := make(chan websocket.Response, 1)

	if err = service.GetTrading(ctx, tradingReceiver); err != nil {
		log.Fatal(err)
	}
}
