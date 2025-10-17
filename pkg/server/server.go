package server

import (
	"log"
	"net/http"
	"os"
	"todo-server/pkg/api"
)

const defaultPort = "7540"
const webDir = "./web"

// Start запускает HTTP-сервер и возвращает указатель на него для graceful shutdown
func Start() *http.Server {
	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	srv := &http.Server{
		Addr: ":" + port,
	}

	log.Printf("Сервер запущен на порту %s", port)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Сервер завершил работу с ошибкой: %v", err)
		}
	}()

	return srv
}
