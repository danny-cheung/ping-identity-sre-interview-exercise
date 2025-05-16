package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/danny-cheung/ping-identity-sre-interview-exercise/pkg/service"
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
	match, err := regexp.MatchString("^[A-Z0-9.]+$", tickerSymbol)
	if err != nil {
		log.Panic("Unable to validate SYMBOL")
	}
	if !match {
		log.Panic("Env var 'SYMBOL' is not valid.")
	}

	tickerDayString, found := os.LookupEnv("NDAYS")
	if !found {
		log.Panic("Env var 'NDAYS' not found.")
	}
	tickerDays, err = strconv.Atoi(tickerDayString)
	if err != nil {
		log.Panicf("Unable to convert '%s' to int.", tickerDayString)
	}
	if tickerDays < 1 {
		log.Panic("Env var NDAYS cannot be less than 1.")
	}
}

func NewTicker(tckr service.Ticker) gin.HandlerFunc {
	return func(c *gin.Context) {
		out, err := tckr.Ticker(tickerSymbol, tickerDays)
		if err != nil {
			c.Error(fmt.Errorf("unable to get ticker info: %w", err))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK, out)
	}
}
