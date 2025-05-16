package service

import (
	"time"

	"github.com/danny-cheung/ping-identity-sre-interview-exercise/pkg/models"
)

type Mock struct{}

func NewMock() *Mock {
	return new(Mock)
}

func (m *Mock) Ticker(tickerSymbol string, tickerDays int) (*models.TickerInfo, error) {
	out := &models.TickerInfo{
		Symbol:     tickerSymbol,
		TimeSeries: map[string]models.TickerInfoDailyEntry{},
	}

	for i := range tickerDays {
		entryDate := time.Now().Add(time.Hour * time.Duration(24*i*-1)).Format(time.DateOnly)

		out.TimeSeries[entryDate] = models.TickerInfoDailyEntry{
			Open:   float64(i),
			Close:  float64(i + 1),
			Low:    float64(i / 2),
			High:   float64(2 * i),
			Volume: float64(i * i),
		}
	}
	out.AverageClosingPrice = float64(1+tickerDays) / float64(tickerDays)

	return out, nil
}
