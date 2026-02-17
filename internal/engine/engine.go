package engine

import (
	"context"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/moby/moby/client"
)

type MonitorEngine struct {
	Client     docker.DockerClient
	Containers *client.ContainerListResult
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
