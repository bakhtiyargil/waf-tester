package model

import "strings"

const (
	TextToSearch = iota
	HttpStatus
)

type Target struct {
	Host     string
	Path     string
	Method   string
	Payload  string
	Criteria map[int8]string
}

var testTarget *Target

func (t *Target) GetUrl() string {
	if strings.HasPrefix(t.Host, "http://") || strings.HasPrefix(t.Host, "https://") {
		return t.Host + t.Path
	}
	return "http://" + t.Host + t.Path
}

func GetTestTargetInstance() *Target {
	if testTarget != nil {
		return testTarget
	} else {
		return &Target{
			Host:   "http://waffy.xyz",
			Path:   "/DVWA",
			Method: "GET",
		}
	}
}

func FromRequest(request *TestRequest) *Target {
	target := Target{
		Host:   request.Host,
		Path:   request.Path,
		Method: request.Method,
		Criteria: map[int8]string{
			TextToSearch: request.Criteria.TextToSearch,
			HttpStatus:   request.Criteria.HttpStatus,
		},
	}
	return &target
}
