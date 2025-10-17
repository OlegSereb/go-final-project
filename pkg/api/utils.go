// pkg/api/utils.go
package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON отправляет данные в формате JSON
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
