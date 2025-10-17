// pkg/api/validation.go
package api

import (
	"errors"
	"time"
	"todo-server/pkg/db"
)

// CheckDate корректирует дату задачи по правилам ТЗ
func CheckDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(DateFormat)

	if task.Date == "" {
		task.Date = today
	}

	// Проверяем формат даты
	if _, err := time.Parse(DateFormat, task.Date); err != nil {
		return errors.New("invalid date format, expected YYYYMMDD")
	}

	// Если дата в прошлом (строго меньше сегодня)
	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	return nil
}
