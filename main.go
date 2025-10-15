// main.go
package main

import (
	"log"
	"os"
	"todo-server/pkg/db"
	"todo-server/pkg/server"
)

func main() {
	// Получаем путь к БД из TODO_DBFILE или используем по умолчанию
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализируем БД
	if err := db.Init(dbFile); err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}

	log.Println("Запуск сервера...")
	server.Start()
}
