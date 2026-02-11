package main

import (
	"context"
	"fmt"
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

	engine := engine.CreateEngine(ctx, *client)
	defer engine.Client.Close()

	if err := engine.Start(); err != nil {
		log.Fatalf("Error running engine: %v\n", err)
	}

	fmt.Printf("CONTAINERS: \n")
	fmt.Printf("%+v\n", *engine.Containers)

}
