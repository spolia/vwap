package coinbase

import (
	"context"
	"fmt"
	"log"

	ws "github.com/gorilla/websocket"
	"github.com/spolia/vwap/internal/websocket"
	"golang.org/x/xerrors"
)

type client struct {
	conn *ws.Conn
}

// NewClient returns a new websocket client.
func NewClient(conn *ws.Conn) websocket.Client {
	return &client{
		conn: conn,
	}
}

// Subscribe subscribes to the coinbase websocket.
func (c *client) Subscribe(ctx context.Context, tradingPairs []string, tradingSender chan websocket.Response) error {
	if len(tradingPairs) == 0 {
		return xerrors.New("client: there is no trading pairs")
	}

	var subscription = &Request{
		Type:       RequestTypeSubscribe,
		ProductIDs: tradingPairs,
		Channels:   []Channel{{Name: "matches"}},
	}

	if err := c.conn.WriteJSON(subscription); err != nil {
		return fmt.Errorf("failed writing the websocket %w", err)
	}

	var response Response
	if err := c.conn.ReadJSON(&response); err != nil {
		return fmt.Errorf("failed reading the subscription response: %w", err)
	}

	if response.Type == "error" {
		return fmt.Errorf("subscription error %s", response.Message)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				if err := c.conn.Close(); err != nil {
					log.Printf("failed closing coinbase websocket connection: %s", err)
				}
				return

			default:
				var message Response
				if err := c.conn.ReadJSON(&message); err != nil {
					log.Printf("failed reading messages: %s", err)
					continue
				}
				tradingSender <- websocket.Response{
					Type:      message.Type,
					Size:      message.Size,
					Price:     message.Price,
					ProductID: message.ProductID,
				}
			}
		}
	}()

	return nil
}
