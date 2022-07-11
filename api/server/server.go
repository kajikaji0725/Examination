package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kajikaji0725/Examination/api/model"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	return &Client{&httpClient}
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

	gc.JSON(http.StatusOK, gin.H{"postal_code": resp.Response.Location[0].Postal, "address": resp.Response.Location[0].Prefecture + resp.Response.Location[0].City})
}

func(c *Client) NewRouter() *gin.Engine {

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

	return router
}
