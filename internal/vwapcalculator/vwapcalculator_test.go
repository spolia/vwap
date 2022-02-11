package vwapcalculator_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
	vwap "github.com/spolia/vwap/internal/vwapcalculator"
	"github.com/stretchr/testify/require"
)

func TestPush(t *testing.T) {
	t.Parallel()
	// Given the list with maxSize = 1
	list, err := vwap.NewList([]vwap.DataPoint{}, 1)
	require.NoError(t, err)

	first := vwap.DataPoint{Price: decimal.NewFromInt(1), Volume: decimal.NewFromInt(1)}
	second := vwap.DataPoint{Price: decimal.NewFromInt(2), Volume: decimal.NewFromInt(2)}

	// When
	list.Push(first)
	// then
	require.Equal(t, 1, list.Len())
	require.Equal(t, first, list.DataPoints[0])

	// When
	list.Push(second)
	// then
	require.Equal(t, 1, list.Len())
	require.Equal(t, second, list.DataPoints[0])
}

func TestConcurrentPush(t *testing.T) {
	t.Parallel()
	// Given the list with maxSize = 2
	list, err := vwap.NewList([]vwap.DataPoint{}, 2)
	require.NoError(t, err)

	first := vwap.DataPoint{Price: decimal.NewFromInt(1), Volume: decimal.NewFromInt(1)}
	second := vwap.DataPoint{Price: decimal.NewFromInt(2), Volume: decimal.NewFromInt(2)}
	third := vwap.DataPoint{Price: decimal.NewFromInt(3), Volume: decimal.NewFromInt(3)}

	// When there are concurrent pushes
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		list.Push(first)
		wg.Done()
	}()

	go func() {
		list.Push(second)
		wg.Done()
	}()

	go func() {
		list.Push(third)
		wg.Done()
	}()
	wg.Wait()

	// Then
	require.Len(t, list.DataPoints, 2)
}

func TestVWAP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		dataPoints   []vwap.DataPoint
		expectedVwap map[string]decimal.Decimal
		maxSize      uint
	}{
		{
			name:       "When no datapoints",
			dataPoints: []vwap.DataPoint{},
			expectedVwap: map[string]decimal.Decimal{
				"BTC-USD": decimal.Zero,
				"ETH-USD": decimal.Zero,
			},
		},
		{
			name: "When datapoints are greater than maxSize",
			dataPoints: []vwap.DataPoint{
				{Price: decimal.NewFromInt(10), Volume: decimal.NewFromInt(10), ProductID: "BTC-USD"},
				{Price: decimal.NewFromInt(10), Volume: decimal.NewFromInt(10), ProductID: "BTC-USD"},
				{Price: decimal.NewFromInt(31), Volume: decimal.NewFromInt(30), ProductID: "ETH-USD"},
				{Price: decimal.NewFromInt(21), Volume: decimal.NewFromInt(20), ProductID: "BTC-USD"},
				{Price: decimal.NewFromInt(41), Volume: decimal.NewFromInt(33), ProductID: "ETH-USD"},
			},
			maxSize: 4,
			expectedVwap: map[string]decimal.Decimal{
				"BTC-USD": decimal.RequireFromString("17.3333333333333333"),
				"ETH-USD": decimal.RequireFromString("36.2380952380952381"),
			},
		},
		{
			name: "When datapoints are less than maxSize",
			dataPoints: []vwap.DataPoint{
				{Price: decimal.NewFromInt(10), Volume: decimal.RequireFromString("10.1"), ProductID: "BTC-USD"},
				{Price: decimal.NewFromInt(10), Volume: decimal.RequireFromString("10.1"), ProductID: "BTC-USD"},
			},
			expectedVwap: map[string]decimal.Decimal{
				"BTC-USD": decimal.RequireFromString("10"),
			},
			maxSize: 4,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			list, err := vwap.NewList([]vwap.DataPoint{}, tt.maxSize)
			require.NoError(t, err)

			for _, d := range tt.dataPoints {
				list.Push(d)
			}

			for k := range tt.expectedVwap {
				require.Equal(t, tt.expectedVwap[k].String(), list.VWAP[k].String(),
					fmt.Sprintf("Test: %s does not meet the expected VWap: %s", tt.name, k))
			}
		})
	}
}
