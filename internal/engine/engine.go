package engine

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func NewEngine(client docker.DockerClient) *MonitorEngine {
	return &MonitorEngine{
		Client: client,
	}
}

type MonitorEngine struct {
	Mu             sync.Mutex
	Client         docker.DockerClient
	Containers     *client.ContainerListResult
	ContainerStats map[string]*container.StatsResponse
}

func (eng *MonitorEngine) Start(ctx context.Context) error {
	if err := eng.refreshContainers(ctx); err != nil {
		return err
	}

	eng.ContainerStats = make(map[string]*container.StatsResponse)

	for _, container := range eng.Containers.Items {
		go eng.getContainerStats(ctx, container.ID)
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
		log.Printf("Error Reading stats: %v\n", err)
		return
	}
	defer stats.Body.Close()

	decoder := json.NewDecoder(stats.Body)

	for {

		select {
		case <-ctx.Done():
			return
		default:
		}
		var statResult *container.StatsResponse

		if err := decoder.Decode(&statResult); err != nil {
			if err == io.EOF || err == context.Canceled {
				log.Printf("Stopping monitoring for %s (context canceled or stream ended)", id)
			} else {
				log.Printf("Error decoding stats for %s: %v", id, err)
			}
			return
		}

		eng.Mu.Lock()
		eng.ContainerStats[id] = statResult
		eng.Mu.Unlock()

	}
}
