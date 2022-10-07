package main

import (
	"log"
	"ykjam/new_tx/config"
	"ykjam/new_tx/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}
	app.Run(cfg)
}
