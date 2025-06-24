package model

import "net/http"

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func SuccessResponse() *Response {
	return &Response{Message: "success", Status: http.StatusOK}
}
