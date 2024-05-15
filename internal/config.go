package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Load(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
