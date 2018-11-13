package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	botUrl := os.Getenv("BOT_URL")
	if botUrl == "" {
		botUrl = "https://albumlistbot.herokuapp.com"
	}

	c.Request.ParseForm()

	teamId := c.DefaultPostForm("team_id", "")
	if teamId == "" {
		c.String(http.StatusOK, "Failed")
		return
	}

	newClient := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("%s/api/mapping/%s", botUrl, teamId)
	botRequest, _ := http.NewRequest("GET", url, nil)
	botRequest.Header.Set("Content-Type", "application/json")
	var appUrl string
		getJSON(newClient, botRequest, &appUrl)

	uri := c.Query("uri")
	fullUrl := fmt.Sprintf("%sslack/%s", appUrl, uri)
	log.Printf("Routing to: %s", fullUrl)

	listRequest, _ := http.NewRequest("POST", fullUrl, strings.NewReader(c.Request.Form.Encode()))
	listRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	listRequest.Header.Add("Content-Length", strconv.Itoa(len(c.Request.Form.Encode())))
	var appResponse map[string]interface{}
		getJSON(newClient, listRequest, &appResponse)

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
