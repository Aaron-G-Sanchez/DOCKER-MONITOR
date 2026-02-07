package docker

import (
	"context"

	"github.com/moby/moby/client"
)

type DockerClient struct {
	client *client.Client
}

func NewClient() (*DockerClient, error) {
	client, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &DockerClient{
		client: client,
	}, nil
}

func (dc *DockerClient) ListContainers(ctx context.Context) (*client.ContainerListResult, error) {
	containers, err := dc.client.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}
	return &containers, nil
}
