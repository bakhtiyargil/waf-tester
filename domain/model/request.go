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

func (testRequest *TestRequest) GetApi() string {
	return testRequest.Host + testRequest.Path + testRequest.Method
}
