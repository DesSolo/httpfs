package handlers

import (
	"encoding/json"
	"net/http"
)

// respondJson возвращает ответ в JSON представлении
func respondJson(w http.ResponseWriter, statusCode int, v any) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// errorMessage сообщение об ошибке
type errorMessage struct {
	Message string `json:"message"`
}
