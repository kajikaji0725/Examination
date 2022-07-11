package main

import (
	"log"

	"github.com/kajikaji0725/Examination/api/server"
)

func main() {

	c := server.NewClient()

	router := c.NewRouter()

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
