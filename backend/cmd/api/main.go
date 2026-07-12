package main

import (
	"log"

	"retrofit/backend/internal/app"
	"retrofit/backend/internal/config"
)

func main() {
	cfg := config.Load()

	server, err := app.New(cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
