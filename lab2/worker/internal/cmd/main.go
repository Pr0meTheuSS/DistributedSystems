package main

import (
	"log"
	"worker/internal/app"
	"worker/internal/di"
)

func main() {
	container, err := di.InitContainer()
	if err != nil {
		log.Fatalf("[❌] Failed to init container: %v", err)
	}

	app := app.NewApp(container)

	if err := app.Run(); err != nil {
		log.Fatalf("[❌] Application error: %v", err)
	}
}
