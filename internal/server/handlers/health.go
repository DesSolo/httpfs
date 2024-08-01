package handlers

import "net/http"

// healthResponse ответ на запрос состояния
type healthResponse struct {
	// состояние сервера
	Healthy bool `json:"health"`
}

// HandleHealth обработчик проверки состояния сервера
func HandleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondJson(w, http.StatusOK, healthResponse{
			Healthy: true,
		})
	}
}
