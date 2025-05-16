package main

import (
	"github.com/danny-cheung/ping-identity-sre-interview-exercise/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", handlers.Ticker)

	r.Run() // listen and serve on 0.0.0.0:8080
}
