package engine

import (
	"context"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/moby/moby/client"
)

type MonitorEngine struct {
	ctx        context.Context
	Client     docker.DockerClient
	Containers *client.ContainerListResult
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

	return nil
}

func (eng *MonitorEngine) refreshContainers() error {
	result, err := eng.Client.ListContainers(eng.ctx)
	if err != nil {
		return err
	}

	eng.Containers = &result
	return nil
}
