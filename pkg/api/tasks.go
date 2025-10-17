// pkg/api/tasks.go
package api

import (
	"net/http"
	"strconv"
	"todo-server/pkg/db"
)

// APITask — задача для JSON-ответа (id как строка)
type APITask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// TasksResp — ответ для списка задач
type TasksResp struct {
	Tasks []APITask `json:"tasks"`
}

// tasksHandler обрабатывает GET /api/tasks[?search=...]
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	limit := defaultTaskLimit

	var dbTasks []*db.Task
	var err error

	if searchQuery != "" {
		dbTasks, err = db.SearchTasks(searchQuery, limit)
	} else {
		dbTasks, err = db.Tasks(limit)
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "database error"})
		return
	}

	// Преобразуем []*db.Task → []APITask
	apiTasks := make([]APITask, len(dbTasks))
	for i, t := range dbTasks {
		apiTasks[i] = APITask{
			ID:      strconv.FormatInt(t.ID, 10),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		}
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: apiTasks})
}
