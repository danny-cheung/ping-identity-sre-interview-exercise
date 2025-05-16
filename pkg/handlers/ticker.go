package handlers

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
	"github.com/gin-gonic/gin"
)

var (
	tickerSymbol string
	tickerDays   int
)

func init() {
	var found bool
	var err error

	tickerSymbol, found = os.LookupEnv("SYMBOL")
	if !found {
		log.Panic("Env var 'SYMBOL' not found.")
	}
	if tickerSymbol == "" {
		log.Panic("Env var 'SYMBOL' is empty")
	}

	tickerDayString, found := os.LookupEnv("NDAYS")
	if !found {
		log.Panic("Env var 'NDAYS' not found.")
	}
	tickerDays, err = strconv.Atoi(tickerDayString)
	if err != nil {
		log.Panicf("Unable to convert '%s' to int.", tickerDayString)
	}
}

func Ticker(c *gin.Context) {
	apiKey := "C227WD9W3LUVKVV9"
	url := fmt.Sprintf("https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s", apiKey, tickerSymbol)

	httpResp, httpErr := http.Get(url)
	if httpErr != nil {
		c.Error(fmt.Errorf("error connecting to AlphaVantage: %w", httpErr))
		c.Status(http.StatusInternalServerError)
		return
	}

	avRespBodyBuf := new(bytes.Buffer)
	_, err := avRespBodyBuf.ReadFrom(httpResp.Body)
	if err != nil {
		c.Error(fmt.Errorf("error reading AlphaVantage data: %w", httpErr))
		c.Status(http.StatusInternalServerError)
		return
	}

	println(avRespBodyBuf.String())

	bodyObj := new(models.AlphaVantageResponse)

	err = json.Unmarshal(avRespBodyBuf.Bytes(), bodyObj)
	if err != nil {
		c.Error(fmt.Errorf("error unmarshalling AlphaVantage data: %w", err))
		c.Status(http.StatusInternalServerError)
		return
	}

	seriesKeys := make([]string, 0)
	for key := range bodyObj.TimeSeries {
		seriesKeys = append(seriesKeys, key)
	}

	// Sort in ascending order (oldest->newest)
	slices.Sort(seriesKeys)
	// Sort in descending order (newest->oldest)
	slices.Reverse(seriesKeys)

	out := models.TickerInfo{
		Symbol:     tickerSymbol,
		TimeSeries: map[string]models.TickerInfoDailyEntry{},
	}

	closingSum := 0.0
	closingCount := 0.0
	for _, key := range seriesKeys[:tickerDays] {
		out.TimeSeries[key] = models.TickerInfoDailyEntry{
			Open:   bodyObj.TimeSeries[key].Open,
			Close:  bodyObj.TimeSeries[key].Close,
			Low:    bodyObj.TimeSeries[key].Low,
			High:   bodyObj.TimeSeries[key].High,
			Volume: bodyObj.TimeSeries[key].Volume,
		}

		close, err := strconv.ParseFloat(bodyObj.TimeSeries[key].Close, 64)
		if err != nil {
			c.Error(fmt.Errorf("unable to parse closing price of %s on %s: %w", tickerSymbol, key, err))
			c.Status(http.StatusInternalServerError)
		}

		closingCount++
		closingSum += close
	}

	out.AverageClosingPrice = closingSum / closingCount

	c.IndentedJSON(http.StatusOK, out)
}
