package main

import (
	"log"

	"github.com/kajikaji0725/Examination/api/db"
	"github.com/kajikaji0725/Examination/api/server"
)

func main() {

	config := db.Config{
		Host:     "postal-count-db",
		Username: "root",
		Password: "root",
		DBname:   "root",
		Port:     "5432",
	}

	c, err := server.NewClient(&config)
	if err != nil {
		log.Fatal(err)
	}

	router := c.NewRouter()

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
