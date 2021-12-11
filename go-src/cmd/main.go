package main

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

type serviceConfig struct {
	ServerPort int    `env:"PORT" envDefault:"8080"`
	ServerMode string `env:"MODE" envDefault:"release"`
}

func getConfig() (*serviceConfig, error) {
	config := &serviceConfig{}
	err := env.Parse(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("server is running with config %+v\n", config)
	time.Sleep(1800 * time.Minute) // 30 minutes
}
