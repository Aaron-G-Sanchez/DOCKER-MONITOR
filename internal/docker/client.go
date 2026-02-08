package docker

import (
	"context"

	"github.com/moby/moby/client"
)

type APIClient interface {
	ContainerList(ctx context.Context, options client.ContainerListOptions) (client.ContainerListResult, error)
	Close() error
}

type DockerClient struct {
	api APIClient
}

func NewClient() (*DockerClient, error) {
	client, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &DockerClient{
		api: client,
	}, nil
}

func NewClientWithAPI(api APIClient) *DockerClient {
	return &DockerClient{api: api}
}

func (dc *DockerClient) ListContainers(ctx context.Context) (client.ContainerListResult, error) {
	return dc.api.ContainerList(ctx, client.ContainerListOptions{All: true})
}

func (dc *DockerClient) Close() error {
	return dc.api.Close()
}
