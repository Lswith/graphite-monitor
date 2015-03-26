package main

import (
	"encoding/json"
	"io"
)

type Config struct {
	Endpoint      string
	Interval      string
	Target        string
	Threshold     float64
	Frequency     string
	Rule          string
	EmailServer   string
	EmailTo       string
	EmailFrom     string
	EmailUser     string
	EmailPassword string
	EmailPort     string
	EmailSubject  string
}

func ReadConfig(r io.Reader) (Config, error) {
	decoder := json.NewDecoder(r)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		return Config{}, err
	}
	return configuration, nil
}
