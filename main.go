package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func getJSON(client *http.Client, request *http.Request, target interface{}) error {
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(target)
}

func routeHandler(c *gin.Context) {
	team_id := c.PostForm("team_id")
	newClient := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("https://albumlistbot.herokuapp.com/api/mapping/%s", team_id)
	bot_req, _ := http.NewRequest("GET", url, nil)
	bot_req.Header.Set("Content-Type", "application/json")

	var app_url string
	getJSON(newClient, bot_req, &app_url)

	uri := c.Query("uri")
	full_url := fmt.Sprintf("%sslack/%s", app_url, uri)
	log.Printf("Routing to: %s", full_url)
	app_req, _ := http.NewRequest("POST", full_url, nil)

	var app_response struct{}
	getJSON(newClient, app_req, &app_response)

	c.JSON(http.StatusOK, app_response)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.POST("/route", routeHandler)

	router.Run(":" + port)
}
