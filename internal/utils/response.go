package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type Json struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Error   interface{} `json:"error"`
}

func (j *Json) Response(success bool, result interface{}, message string, status int, error interface{}) *Json {
	return &Json{success, result, message, status, error}
}

func (j *Json) NewResponse(success bool, w http.ResponseWriter, result interface{}, message string, status int, error interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(j.Response(success, result, message, status, error))
	log.Println(j.Response(success, result, message, status, error))
}
