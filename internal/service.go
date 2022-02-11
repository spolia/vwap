package internal

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	vwap "github.com/spolia/vwap/internal/vwapcalculator"
	"github.com/spolia/vwap/internal/websocket"
	"golang.org/x/xerrors"
)

// Service is our main service.
type Service struct {
	webSocketClient websocket.Client
	tradingPairs    []string
	list            *vwap.List
}

// NewService returns a new service.
func NewService(webSocket websocket.Client, tradingPairs []string, list *vwap.List) *Service {
	return &Service{
		webSocketClient: webSocket,
		tradingPairs:    tradingPairs,
		list:            list,
	}
}

// ExecuteEngine execute the vwap task
func (s *Service) ExecuteEngine(ctx context.Context, tradingReceiver chan websocket.Response) error {
	if err := s.webSocketClient.Subscribe(ctx, s.tradingPairs, tradingReceiver); err != nil {
		return xerrors.Errorf("service: %w", err)
	}

	for data := range tradingReceiver {
		if data.Price == "" {
			continue
		}

		price, err := decimal.NewFromString(data.Price)
		if err != nil {
			return fmt.Errorf("service: wrong received price %v", err)
		}

		size, err := decimal.NewFromString(data.Size)
		if err != nil {
			return fmt.Errorf("service: wrong received data size %v", err)
		}

		s.list.Push(vwap.DataPoint{
			Price:     price,
			Volume:    size,
			ProductID: data.ProductID,
		})

		// Just Print
		for k, v := range s.list.VWAP {
			fmt.Println(k, ":", v)
		}
	}

	return nil
}
