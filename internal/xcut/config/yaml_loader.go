package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

// loadFromYAML loads configuration from YAML file
func loadFromYAML(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	
	return &config, nil
}