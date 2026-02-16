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
	ctx            context.Context
	Client         docker.DockerClient
	Containers     *client.ContainerListResult
	ContainerStats map[string]*container.StatsResponse
}

func CreateEngine(ctx context.Context, client docker.DockerClient) *MonitorEngine {
	return &MonitorEngine{
		ctx:    ctx,
		Client: client,
	}
}

func (eng *MonitorEngine) Start() error {
	if err := eng.refreshContainers(); err != nil {
		return err
	}

	eng.ContainerStats = make(map[string]*container.StatsResponse)

	for _, container := range eng.Containers.Items {
		eng.getContainerStats(container.ID)
	}

	return nil
}

// TODO: Refactor to only set engine.containers with active containers only.
func (eng *MonitorEngine) refreshContainers() error {
	result, err := eng.Client.ListContainers(eng.ctx)
	if err != nil {
		return err
	}

	eng.Containers = &result
	return nil
}

func (eng *MonitorEngine) getContainerStats(id string) {
	stats, err := eng.Client.ListContainerStats(eng.ctx, id)
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
