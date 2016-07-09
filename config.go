package main

import (
	"encoding/json"
	"os"
)

// Configuration file for wemolightapi
type Configuration struct {
	Switches []ConfigurationSwitch `json:"switches"`
}

// ConfigurationSwitch switch object
type ConfigurationSwitch struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

func parseConfig(configFile string) (*Configuration, error) {
	file, err := os.Open(configFile)

	if err != nil {
		return nil, err
	}

	configuration := Configuration{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}
