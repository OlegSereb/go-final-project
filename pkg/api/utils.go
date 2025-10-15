// pkg/api/utils.go
package api

import (
	"encoding/json"
	"net/http"
	"time"
)

const DateFormat = "20060102"

// afterNow проверяет, что date > now (дата в будущем)
func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	if y1 > y2 {
		return true
	}
	if y1 < y2 {
		return false
	}
	if m1 > m2 {
		return true
	}
	if m1 < m2 {
		return false
	}
	return d1 > d2
}

// writeJSON отправляет данные в формате JSON
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
