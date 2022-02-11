package coinbase_test

import (
	"context"
	"testing"

	ws "github.com/gorilla/websocket"
	"github.com/spolia/vwap/internal/websocket"
	"github.com/spolia/vwap/internal/websocket/coinbase"
	"github.com/stretchr/testify/require"
)

func TestSubscribe_WebsocketConection_Ok(t *testing.T) {
	// Given
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tradingReceiver := make(chan websocket.Response)
	conn, _, err := ws.DefaultDialer.Dial(coinbase.DefaultURL, nil)
	require.NoError(t, err)
	// When
	err = coinbase.NewClient(conn).Subscribe(ctx, []string{"BTC-USD"}, tradingReceiver)
	// Then
	require.NoError(t, err)
	for d := range tradingReceiver {
		if d.Type == "last_match" {
			require.Equal(t, "BTC-USD", d.ProductID)
			break
		}
	}

}

func TestSubscribe_WebsocketConection_Fail(t *testing.T) {
	// Given
	tradingReceiver := make(chan websocket.Response)
	conn, _, err := ws.DefaultDialer.Dial(coinbase.DefaultURL, nil)
	require.NoError(t, err)
	// When
	err = coinbase.NewClient(conn).Subscribe(context.Background(), []string{"wrong"}, tradingReceiver)
	// Then
	require.Error(t, err)
	require.True(t, len(tradingReceiver) == 0)
}
