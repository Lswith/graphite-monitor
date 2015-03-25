package main

import (
	"encoding/json"
	"io"
	"log"
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

func ReadConfig(r io.Reader) Config {
	decoder := json.NewDecoder(r)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	return configuration
}
