package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-server/pkg/db"
)

// TaskHandler обрабатывает все методы для /api/task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	// Преобразуем в APITask (id как строка)
	apiTask := APITask{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	writeJSON(w, http.StatusOK, apiTask)
}

// updateTaskHandler обрабатывает PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if input.ID == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Преобразуем ID в число
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Неверный идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Создаём задачу для обновления
	task := &db.Task{
		ID:      id,
		Date:    input.Date,
		Title:   input.Title,
		Comment: input.Comment,
		Repeat:  input.Repeat,
	}

	// Проверяем дату (как при добавлении)
	if err := CheckDate(task); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Обновляем в БД
	if err := db.UpdateTask(task); err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON с кодом 200
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// deleteTaskHandler обрабатывает DELETE /api/task?id=123
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
