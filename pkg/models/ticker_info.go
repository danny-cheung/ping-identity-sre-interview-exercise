package models

type TickerInfo struct {
	Symbol              string                          `json:"symbol"`
	TimeSeries          map[string]TickerInfoDailyEntry `json:"time_series"`
	AverageClosingPrice float64                         `json:"average_closing_price"`
}

type TickerInfoDailyEntry struct {
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}
