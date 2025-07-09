package model

import "net/http"

type Response struct {
	Id      string `json:"id"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SuccessResponse() *Response {
	return &Response{Message: "success", Status: http.StatusOK}
}

func SuccessResponseWithId(id string) *Response {
	return &Response{Message: "success", Status: http.StatusOK, Id: id}
}

func ErrorResponse() *Response {
	return &Response{Message: "unexpected internal error", Status: http.StatusInternalServerError}
}
