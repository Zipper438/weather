package main

import (
	"fmt"
	"log"

	"example.com/m/internal/api"
	"example.com/m/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	temp, err := api.NewWeatherClient(cfg.APIKey).GetWeather("Moscow")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Температура в Москве: %.1f градусов", temp)
}
