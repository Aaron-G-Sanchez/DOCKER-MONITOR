package docker

import (
	"context"

	"github.com/moby/moby/client"
)

type DockerClient struct {
	ctx    context.Context
	client *client.Client
}

func NewClient(ctx context.Context) (*DockerClient, error) {
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &DockerClient{
		ctx:    ctx,
		client: apiClient,
	}, nil
}

func (dc *DockerClient) ListContainers() (*client.ContainerListResult, error) {
	containers, err := dc.client.ContainerList(dc.ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return &containers, nil
}

func (dc *DockerClient) Close() error {
	return dc.client.Close()
}
