package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func routeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Yo")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/route", routeHandler)

	router.Run(":" + port)
}
