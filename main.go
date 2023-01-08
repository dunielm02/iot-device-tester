package main

import (
	"log"
	"net/http"
	"websocket/templates"

	"github.com/olahol/melody"
)

var m *melody.Melody

func main() {
	m = melody.New()
	templates.Serve()
	Setup()

	log.Println("listening")

	http.ListenAndServeTLS(":8000", "certificate.pem", "key.pem", nil)
}
