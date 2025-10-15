// pkg/db/task.go
package db

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      int64  `json:"id" db:"id"`
	Date    string `json:"date" db:"date"` // YYYYMMDD
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"` // d N, y, и т.д.
}

// Tasks возвращает список задач, отсортированных по дате (максимум limit штук)
func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Важно: возвращаем пустой слайс, а не nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, rows.Err()
}

// SearchTasks ищет задачи по поисковому запросу
func SearchTasks(query string, limit int) ([]*Task, error) {
	// Проверяем, является ли запрос датой в формате DD.MM.YYYY
	dateRegex := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
	if dateRegex.MatchString(query) {
		// Преобразуем DD.MM.YYYY → YYYYMMDD
		t, err := time.Parse("02.01.2006", query)
		if err != nil {
			// Если не удалось распарсить дату — ищем как текст
			return searchByText(query, limit)
		}
		targetDate := t.Format("20060102")

		rows, err := db.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = ?
			ORDER BY date
			LIMIT ?
		`, targetDate, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		return scanTasks(rows)
	}

	// Поиск по тексту
	return searchByText(query, limit)
}

// searchByText ищет подстроку в title или comment (регистронезависимо)
func searchByText(query string, limit int) ([]*Task, error) {
	searchPattern := "%" + query + "%"
	rows, err := db.Query(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE title LIKE ? OR comment LIKE ?
		ORDER BY date
		LIMIT ?
	`, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

// scanTasks — вспомогательная функция для сканирования задач
func scanTasks(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, rows.Err()
}

// GetTask возвращает задачу по ID (id как строка)
func GetTask(idStr string) (*Task, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID")
	}

	task := &Task{}
	err = db.QueryRow(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	return task, nil
}

// UpdateTask обновляет задачу в БД
func UpdateTask(task *Task) error {
	res, err := db.Exec(`
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`, task.Date, task.Title, task.Comment, task.Repeat, task.ID)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// DeleteTask удаляет задачу по ID (id как строка)
func DeleteTask(idStr string) error {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID")
	}

	res, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// UpdateDate обновляет только дату задачи
func UpdateDate(nextDate, idStr string) error {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID")
	}

	res, err := db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
