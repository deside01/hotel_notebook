package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Неудалось преобразовать в JSON: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func Err(w http.ResponseWriter, code int, msg string) {
	log.Printf("Ответ с ошибкой %d: %s", code, msg)

	JSON(w, code, map[string]string{
		"error": msg,
	})
}
