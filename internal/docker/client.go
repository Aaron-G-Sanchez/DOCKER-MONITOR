package docker

import (
	"context"

	"github.com/moby/moby/client"
)

type APIClient interface {
	ContainerList(ctx context.Context, options client.ContainerListOptions) (client.ContainerListResult, error)
	ContainerStats(ctx context.Context, containerID string, options client.ContainerStatsOptions) (client.ContainerStatsResult, error)
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

// FUNCTION FOR MOCKING DOCKER SDK API.
func NewClientWithMockAPI(api APIClient) *DockerClient {
	return &DockerClient{api: api}
}

// List all containers in the docker host.
func (dc *DockerClient) ListContainers(ctx context.Context) (client.ContainerListResult, error) {
	return dc.api.ContainerList(ctx, client.ContainerListOptions{All: true})
}

// Retrieves live resource usage statistics for the specified container.
func (dc *DockerClient) ListContainerStats(ctx context.Context, containerID string) (client.ContainerStatsResult, error) {
	return dc.api.ContainerStats(ctx, containerID, client.ContainerStatsOptions{})
}

func (dc *DockerClient) Close() error {
	return dc.api.Close()
}
