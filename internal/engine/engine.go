package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type MonitorEngine struct {
	mu             sync.Mutex
	Client         docker.DockerClient
	Containers     *client.ContainerListResult
	ContainerStats map[string]*container.StatsResponse
}

func CreateEngine(client docker.DockerClient) *MonitorEngine {
	return &MonitorEngine{
		Client: client,
	}
}

func (eng *MonitorEngine) Start(ctx context.Context) error {
	if err := eng.refreshContainers(ctx); err != nil {
		return err
	}

	eng.ContainerStats = make(map[string]*container.StatsResponse)

	for _, container := range eng.Containers.Items {
		eng.getContainerStats(ctx, container.ID)
	}

	return nil
}

func (eng *MonitorEngine) refreshContainers(ctx context.Context) error {
	result, err := eng.Client.ListContainers(ctx)
	if err != nil {
		return err
	}

	eng.Containers = &result
	return nil
}

func (eng *MonitorEngine) getContainerStats(ctx context.Context, id string) {
	stats, err := eng.Client.ListContainerStats(ctx, id)
	if err != nil {
		log.Fatalf("Error Reading stats: %v\n", err)
	}
	defer stats.Body.Close()

	decoder := json.NewDecoder(stats.Body)

	for {

		var statResult *container.StatsResponse

		if err := decoder.Decode(&statResult); err != nil {
			return
		}

		eng.mu.Lock()
		eng.ContainerStats[id] = statResult
		eng.mu.Unlock()

		fmt.Println(eng.ContainerStats[id])

	}
}
