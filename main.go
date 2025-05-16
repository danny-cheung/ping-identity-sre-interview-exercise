package main

import (
	"github.com/danny-cheung/ping-identity-sre-interview-exercise/pkg/handlers"
	"github.com/danny-cheung/ping-identity-sre-interview-exercise/pkg/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", handlers.NewTicker(service.NewAlphaVantageService()))
	r.GET("/health", handlers.Health)

	r.Run() // listen and serve on 0.0.0.0:8080
}
