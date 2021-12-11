package main

import (
	"fmt"
	"net/http"

	"github.com/caarlos0/env"
)

type serviceConfig struct {
	ServerPort      int    `env:"PORT" envDefault:"8080"`
	ServerMode      string `env:"MODE" envDefault:"release"`
	HealthCheckPort int    `env:"HEALTH_CHECK_PORT" envDefault:"8081"`
}

func getConfig() (*serviceConfig, error) {
	config := &serviceConfig{}
	err := env.Parse(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func live(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "live\n")
}

func ready(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ready\n")
}

func main() {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}

	go func() {
		http.HandleFunc("/health/live", live)
		http.HandleFunc("/health/ready", ready)

		http.ListenAndServe(fmt.Sprintf(":%d", config.HealthCheckPort), nil)
	}()

	http.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), nil)
}
