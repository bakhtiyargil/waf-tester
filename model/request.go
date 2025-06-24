package model

type TestRequest struct {
	Host   string `json:"Host"`
	Path   string `json:"Path"`
	Method string `json:"Method"`
}
