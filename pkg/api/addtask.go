// pkg/api/addtask.go
package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-server/pkg/db"
)

// addTaskHandler обрабатывает POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "invalid JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err := CheckDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "database error"})
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}
