package main

import (
	"net/http"
)

func Setup() {
	http.HandleFunc("/api/login", login)

	http.HandleFunc("/api/v1/reading/", receiveData)

	http.HandleFunc("/ws", webSocket)
}
