package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	rSuccess = `success`
	rError   = `error`
)

type response struct {
	Meta meta        `json:"meta"`
	Data interface{} `json:"data"`
}
type meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(w http.ResponseWriter, message string, code int, status string, data interface{}) response {
	meta := meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	response := response{
		Meta: meta,
		Data: data,
	}
	w.WriteHeader(code)
	bytes, _ := json.Marshal(response)
	json := string(bytes[:])
	fmt.Println(json)
	fmt.Fprint(w, json)
	return response
}
