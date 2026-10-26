package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Не удалось преобразовать в JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func Err(w http.ResponseWriter, code int, msg string) {
	log.Printf("Ответ с ошибкой %d: %s", code, msg)

	JSON(w, code, map[string]string{
		"error": msg,
	})
}
