// pkg/server/server.go
package server

import (
	"log"
	"net/http"
	"os"
	"todo-server/pkg/api"
)

const defaultPort = "7540"
const webDir = "./web"

func Start() {
	// Инициализируем API до файлового сервера!
	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
