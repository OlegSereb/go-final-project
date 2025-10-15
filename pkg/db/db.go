// pkg/db/db.go
package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

// Init инициализирует базу данных
func Init(dbFile string) error {
	// Проверяем, существует ли файл
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем БД
	var errOpen error
	db, errOpen = sql.Open("sqlite", dbFile)
	if errOpen != nil {
		return errOpen
	}

	// Если файл не существовал — создаём схему
	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddTask добавляет задачу в БД
func AddTask(task *Task) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
