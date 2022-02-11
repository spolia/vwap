package internal_test

import (
	"context"
	"testing"

	"github.com/spolia/vwap/internal"
	"github.com/spolia/vwap/internal/vwapcalculator"
	"github.com/spolia/vwap/internal/websocket"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_When_ReceivedMessage_IsOk_Then_Returns_Ok(t *testing.T) {
	// Given
	var clientMock clientMock
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// When
	list, err := vwapcalculator.NewList([]vwapcalculator.DataPoint{}, 3)
	service := internal.NewService(&clientMock, []string{"BTC-USD"}, &list)
	tradingReceiver := make(chan websocket.Response, 1)
	err = service.ExecuteEngine(ctx, tradingReceiver)

	// Then
	require.NoError(t, err)
	require.True(t, list.Len() > 0)
	require.Equal(t, "BTC-USD", list.DataPoints[0].ProductID)
}

func TestService_When_ReceivedMessage_IsWrong_Then_Returns_Error(t *testing.T) {
	// Given
	var clientMock clientMockFail
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// When
	list, err := vwapcalculator.NewList([]vwapcalculator.DataPoint{}, 3)
	service := internal.NewService(&clientMock, []string{"BTC-USD"}, &list)
	tradingReceiver := make(chan websocket.Response, 1)
	err = service.ExecuteEngine(ctx, tradingReceiver)

	// Then
	require.Error(t, err)
	require.True(t, list.Len() == 0)
}

type clientMock struct {
	mock.Mock
}

func (c *clientMock) Subscribe(ctx context.Context, tradingPairs []string, tradingGenerator chan websocket.Response) error {
	go func() {
		select {
		default:
			tradingGenerator <- websocket.Response{
				Type:      "match",
				Size:      "20",
				Price:     "0.0713052736957975",
				ProductID: "BTC-USD",
			}

			close(tradingGenerator)

		case <-ctx.Done():
			return
		}
	}()

	return nil
}

type clientMockFail struct {
	mock.Mock
}

func (c *clientMockFail) Subscribe(ctx context.Context, tradingPairs []string, tradingGenerator chan websocket.Response) error {
	go func() {
		select {
		default:
			tradingGenerator <- websocket.Response{
				Type:      "match",
				Size:      "20",
				Price:     "wrong-price",
				ProductID: "BTC-USD",
			}

			close(tradingGenerator)

		case <-ctx.Done():
			return
		}
	}()

	return nil
}
