package service

import "github.com/danny-cheung/ping-identity-sre-interview-exercise/pkg/models"

// Retrieves data from its configured source and return it in a standardised struct.
type Ticker interface {
	Ticker(string, int) (*models.TickerInfo, error)
}
