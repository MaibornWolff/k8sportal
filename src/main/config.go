package main

var config struct {
	Mongodb struct {
		Host       string
		Database   string
		Collection string
	}
	LogLevel string `envconfig:"default=info"`
}
