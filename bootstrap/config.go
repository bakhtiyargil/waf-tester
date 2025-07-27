package bootstrap

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"waf-tester/config"
)

func InitConfig(file string) (c *config.Config) {
	fileType := filepath.Ext(file)
	if fileType != ".yaml" && fileType != ".yml" {
		log.Fatalf("unsupported file type %s, file type must be .yml or .yaml", fileType)
	}

	configFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error while reading file %s, %v", file, err)
	}

	tmpl, err := template.New("config").Option("missingkey=error").Parse(string(configFile))
	if err != nil {
		log.Fatalf("error while parsing template: %v", err)
	}

	var processed bytes.Buffer
	err = tmpl.Execute(&processed, getEnvMap())
	if err != nil {
		log.Fatalf("error while executing template: %v", err)
	}

	err = yaml.Unmarshal(processed.Bytes(), &c)
	if err != nil {
		log.Fatalf("error while parsing file %s, %v", file, err)
	}
	return c
}

func getEnvMap() map[string]string {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		parts := bytes.SplitN([]byte(e), []byte("="), 2)
		if len(parts) == 2 {
			envs[string(parts[0])] = string(parts[1])
		}
	}
	return envs
}
