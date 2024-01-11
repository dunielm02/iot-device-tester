package httpServer

import (
	"net/http"

	"github.com/olahol/melody"
)

var webSocket *melody.Melody

func Setup() {
	webSocket = melody.New()

	http.HandleFunc("/api/", receiveData)

	http.HandleFunc("/ws", webSocketReq)
}
