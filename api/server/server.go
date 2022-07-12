package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kajikaji0725/Examination/api/db"
	"github.com/kajikaji0725/Examination/api/model"
)

const (
	timeFormat = "2006-01-02"
)

type Client struct {
	client     *http.Client
	controller *db.Controller
}

func NewClient(config *db.Config) (*Client, error) {

	controller, err := db.NewController(config)
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	return &Client{
		client:     &httpClient,
		controller: controller,
	}, nil
}

func (c *Client) NewRouter() *gin.Engine {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8080",
		},
		AllowMethods: []string{
			"GET",
		},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))

	router.GET("/address", c.handleFetchAddress)
	router.GET("/address/access_logs", c.handleFetchRequestCount)

	return router
}

func (c *Client) fetchAddress(address string) (*model.ResponseJson, error) {

	const endpoint = "https://geoapi.heartrails.com/api/json"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("method", "searchByPostal")
	params.Add("postal", address)
	req.URL.RawQuery = params.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var addressJson model.ResponseJson

	err = json.Unmarshal(respByte, &addressJson)
	if err != nil {
		return nil, err
	}

	if addressJson.Response.Error != "" {
		return nil, errors.New(addressJson.Response.Error)
	}

	return &addressJson, nil
}

func (c *Client) handleFetchAddress(gc *gin.Context) {

	address := gc.Query("postal_code")

	resp, err := c.fetchAddress(address)
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"ErrorMessage": err.Error()})
		return
	}

	resp.Response.Location[0].Date, err = time.Parse("2006-01-02", time.Now().Format(timeFormat))
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"ErrorMessage": err.Error()})
		return
	}

	err = c.controller.SetDBPostalCode(&resp.Response.Location[0])
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"ErrorMessage": err.Error()})
		return
	}

	gc.JSON(http.StatusOK, resp.Response.Location)
}

func (c *Client) handleFetchRequestCount(gc *gin.Context) {

	resp, err := c.controller.FetchDBRequestCount(gc.Query("from"), gc.Query("to"), timeFormat)
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"ErrorMessage": err.Error()})
		return
	}

	postalRequestCount := map[string]int{}
	for _, count := range resp {
		postalRequestCount[count.Postal]++
	}

	postalRequest := model.AccessLogs{}
	for requestPostal, requestCount := range postalRequestCount {
		postalRequest.AccessLog = append(postalRequest.AccessLog,
			model.AccessLog{
				Postal:       requestPostal,
				RequestCount: requestCount,
			},
		)
	}

	gc.JSON(http.StatusOK, postalRequest)
}
