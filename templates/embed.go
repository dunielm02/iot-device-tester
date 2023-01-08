package templates

import (
	"embed"
	"net/http"
)

//go:embed *
var f embed.FS

func Serve() {
	http.Handle("/", http.FileServer(http.FS(f)))
}
