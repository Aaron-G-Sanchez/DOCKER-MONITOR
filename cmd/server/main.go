package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aaron-g-sanchez/PROJECTS/DOCKER-MONITOR/internal/docker"
)

func main() {

	ctx := context.Background()

	client, err := docker.NewClient(ctx)
	if err != nil {
		log.Fatalf("Error initializing Docker client: %v\n", err)
	}
	defer client.Close()

	// TODO: Add options as a param.
	containers, err := client.ListContainers()
	if err != nil {
		log.Fatalf("Error retrieving containers: %v\n", err)
	}

	// TODO: Remove after testing.
	fmt.Println("ACTIVE CONTAINERS: ")
	fmt.Printf("%+v\n", containers.Items)

}
