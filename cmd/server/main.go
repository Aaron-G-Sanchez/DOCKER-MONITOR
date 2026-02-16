package main

import (
	"context"
	"log"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/engine"
)

func main() {

	ctx := context.Background()

	run(ctx)
}

func run(ctx context.Context) {
	client, err := docker.NewClient()
	if err != nil {
		log.Fatalf("Error starting docker client: %v\n", err)
	}

	// Create and start the monitor engine.
	engine := engine.CreateEngine(ctx, *client)
	defer engine.Client.Close()

	if err := engine.Start(); err != nil {
		log.Fatalf("Error running engine: %v\n", err)
	}

	// Create and run the server.
	// server := api.NewServer(engine)
	// if err := server.Start(":9876"); err != nil {
	// 	log.Fatalf("Error starting server: %v\n", err)
	// }

}
