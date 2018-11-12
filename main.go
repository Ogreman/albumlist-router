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
	teamId := c.PostForm("team_id")
	if teamId == "" {
		c.String(http.StatusOK, "Failed")
	}
	newClient := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("https://albumlistbot.herokuapp.com/api/mapping/%s", teamId)
	botRequest, _ := http.NewRequest("GET", url, nil)
	botRequest.Header.Set("Content-Type", "application/json")

	var appUrl string
	getJSON(newClient, botRequest, &appUrl)

	uri := c.Query("uri")
	fullUrl := fmt.Sprintf("%sslack/%s", appUrl, uri)
	log.Printf("Routing to: %s", fullUrl)
	appRequest, _ := http.NewRequest("POST", fullUrl, nil)

	var appResponse struct{}
	getJSON(newClient, appRequest, &appResponse)

	c.JSON(http.StatusOK, appResponse)
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
