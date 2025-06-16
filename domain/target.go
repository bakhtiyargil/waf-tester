package domain

import (
	"strings"
)

type Target struct {
	Host   string
	Path   string
	Method string
}

var testTarget *Target

func (t *Target) GetUrl() string {
	if strings.HasPrefix(t.Host, "http://") || strings.HasPrefix(t.Host, "https://") {
		return t.Host + t.Path
	}
	return "http://" + t.Host + t.Path
}

// make concurrent safe
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
