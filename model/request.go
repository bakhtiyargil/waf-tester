package model

type TestRequest struct {
	Host     string   `json:"host"`
	Path     string   `json:"path"`
	Method   string   `json:"method"`
	Criteria Criteria `json:"criteria"`
}

type Criteria struct {
	TextToSearch string `json:"textToSearch"`
	HttpStatus   string `json:"httpStatus"`
}
