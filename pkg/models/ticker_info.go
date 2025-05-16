package models

type TickerInfo struct {
	Symbol              string                          `json:"symbol"`
	TimeSeries          map[string]TickerInfoDailyEntry `json:"time_series"`
	AverageClosingPrice float64                         `json:"average_closing_price"`
}

type TickerInfoDailyEntry struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}
