// pkg/api/done.go
package api

import (
	"net/http"
	"time"
	"todo-server/pkg/db"
)

// doneTaskHandler обрабатывает GET /api/task/done?id=123
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Если задача одноразовая — удаляем
	if task.Repeat == "" || task.Repeat == "none" {
		if err := db.DeleteTask(id); err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		// Повторяющаяся задача — переносим на следующую дату
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Не удалось вычислить следующую дату"}`, http.StatusInternalServerError)
			return
		}

		if err := db.UpdateDate(nextDate, id); err != nil {
			http.Error(w, `{"error":"Ошибка обновления даты"}`, http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]interface{}{})
}
