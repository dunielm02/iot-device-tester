package main

import (
	"iotTester/database"
	httpServer "iotTester/http"
	"iotTester/mqtt"
	"iotTester/templates"
	"log"
	"net/http"
)

func main() {
	templates.Serve()
	httpServer.Setup()

	err := database.ConnectToSQLite()

	if err != nil {
		log.Fatal(err)
	}

	mqtt.Connect()

	log.Println("Http listening")

	http.ListenAndServe(":8000", nil)
}
