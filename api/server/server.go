package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kajikaji0725/Examination/api/model"
)

func fetchAddress(address string) (*model.ResponseJson, error) {

	var addressJson model.ResponseJson

	endpoint := "https://geoapi.heartrails.com/api/json"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("method", "searchByPostal")
	params.Add("postal", address)
	req.URL.RawQuery = params.Encode()

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respByte, &addressJson)
	if err != nil {
		return nil, err
	}

	if addressJson.Response.Error != nil {
		return nil, addressJson.Response.Error
	}

	return &addressJson, nil
}

func handleFetchAddress(c *gin.Context) {

	address := c.Query("postal_code")

	resp, err := fetchAddress(address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"ErrorMessage": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"postal_code": resp.Response.Location[0].Postal, "address": resp.Response.Location[0].Prefecture + resp.Response.Location[0].City})
}

func NewRouter() *gin.Engine {

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

	router.GET("/address", handleFetchAddress)

	return router
}
