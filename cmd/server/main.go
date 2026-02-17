package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/engine"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/server"
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
	server := server.NewServer(engine)

	fmt.Println("Starting server on :6060")
	if err := server.Start(":6060"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v\n", err)
	}

}
