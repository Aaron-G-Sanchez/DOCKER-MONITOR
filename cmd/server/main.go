package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
)

func main() {

	ctx := context.Background()

	client, err := docker.NewClient()
	if err != nil {
		log.Fatalf("Error initializing Docker client: %v\n", err)
	}

	// TODO: Add options as a param.
	containers, err := client.ListContainers(ctx)
	if err != nil {
		log.Fatalf("Error retrieving containers: %v\n", err)
	}

	// TODO: Remove after testing.
	fmt.Println("ACTIVE CONTAINERS: ")
	fmt.Printf("%+v\n", containers.Items)

}
