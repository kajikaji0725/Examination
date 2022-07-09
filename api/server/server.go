package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kajikaji0725/Examination/api/model"
)

func NewRequest(address string) ([]byte, error) {

	url := "https://geoapi.heartrails.com/api/json?method=searchByPostal&postal="

	resp, err := http.Get(url + address)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}

func fetchAddress(c *gin.Context) {

	var addressJson model.ResponseJson
	var errorJson model.ResponseErrorJson

	address := c.Query("postal_code")

	resp, err := NewRequest(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		c.Abort()
		return
	}

	json.Unmarshal(resp, &addressJson)
	json.Unmarshal(resp, &errorJson)

	if errorJson.Response.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorJson.Response.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"postal_code": addressJson.Responses.Location[0].Postal, "address": addressJson.Responses.Location[0].Prefecture + addressJson.Responses.Location[0].City})
	}
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
		MaxAge: 24 * time.Hour,
	}))

	router.GET("/address", fetchAddress)

	return router

}
