package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/danny-cheung/ping-identity-sre-interview-exercise/pkg/models"
)

type AlphaVantageService struct{}

func NewAlphaVantageService() *AlphaVantageService {
	return new(AlphaVantageService)
}

func (avs *AlphaVantageService) Ticker(tickerSymbol string, tickerDays int) (*models.TickerInfo, error) {
	alphavantageApiKey, found := os.LookupEnv("APIKEY")
	if !found {
		log.Panic("Env var 'APIKEY' not found")
	}

	url := fmt.Sprintf("https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s", alphavantageApiKey, tickerSymbol)

	// Retrieve data from Alpha Vantage API
	httpResp, httpErr := http.Get(url)
	if httpErr != nil {
		return nil, fmt.Errorf("error connecting to AlphaVantage: %w", httpErr)
	}

	avRespBodyBuf := new(bytes.Buffer)
	_, err := avRespBodyBuf.ReadFrom(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading AlphaVantage data: %w", httpErr)
	}

	// println(avRespBodyBuf.String())

	bodyObj := new(models.AlphaVantageResponse)

	err = json.Unmarshal(avRespBodyBuf.Bytes(), bodyObj)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling AlphaVantage data: %w", err)
	}

	// Look for the last ${NDAYS} number of records
	seriesKeys := make([]string, 0)
	for key := range bodyObj.TimeSeries {
		seriesKeys = append(seriesKeys, key)
	}

	// Sort in ascending order (oldest->newest)
	slices.Sort(seriesKeys)
	// Sort in descending order (newest->oldest)
	slices.Reverse(seriesKeys)

	out := &models.TickerInfo{
		Symbol:     tickerSymbol,
		TimeSeries: map[string]models.TickerInfoDailyEntry{},
	}

	// Process the API results and populate our response
	closingSum := 0.0
	closingCount := 0.0
	for _, key := range seriesKeys[:tickerDays] {
		open, err := strconv.ParseFloat(bodyObj.TimeSeries[key].Open, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse opening price of %s on %s: %w", tickerSymbol, key, err)
		}

		close, err := strconv.ParseFloat(bodyObj.TimeSeries[key].Close, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse closing price of %s on %s: %w", tickerSymbol, key, err)
		}

		low, err := strconv.ParseFloat(bodyObj.TimeSeries[key].Low, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse low price of %s on %s: %w", tickerSymbol, key, err)
		}

		high, err := strconv.ParseFloat(bodyObj.TimeSeries[key].High, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse high price of %s on %s: %w", tickerSymbol, key, err)
		}

		volume, err := strconv.ParseFloat(bodyObj.TimeSeries[key].Volume, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse volume of %s on %s: %w", tickerSymbol, key, err)
		}

		out.TimeSeries[key] = models.TickerInfoDailyEntry{
			Open:   open,
			Close:  close,
			Low:    low,
			High:   high,
			Volume: volume,
		}

		closingCount++
		closingSum += close
	}
	// Avoid Divide By Zero errors
	if closingCount != 0 {
		out.AverageClosingPrice = closingSum / closingCount
	}

	return out, nil
}
