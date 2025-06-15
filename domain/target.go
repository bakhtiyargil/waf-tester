package domain

import (
	"strings"
)

type Target struct {
	Host   string
	Path   string
	Method string
}

func (t *Target) GetUrl() string {
	if strings.HasPrefix(t.Host, "http://") || strings.HasPrefix(t.Host, "https://") {
		return t.Host + t.Path
	}
	return "http://" + t.Host + t.Path
}

// add singleton
func GetTestTargetInstance() *Target {
	return &Target{
		Host:   "http://waffy.xyz",
		Path:   "/DVWA",
		Method: "GET",
	}
}
