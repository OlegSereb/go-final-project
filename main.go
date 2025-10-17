package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-server/pkg/db"
	"todo-server/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	// Гарантируем закрытие БД при выходе из main
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Ошибка при закрытии БД:", err)
		}
	}()

	log.Println("Запуск сервера...")
	srv := server.Start()

	// Ждём сигнал завершения (Ctrl+C или SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Получен сигнал завершения, выключение сервера...")

	// Даём 5 секунд на graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Ошибка при остановке сервера:", err)
	}

	log.Println("Сервер успешно остановлен.")
}
