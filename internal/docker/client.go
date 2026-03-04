package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// Interface to wrap the docker sdk client.
type DockerClient interface {
	ContainerInspect(ctx context.Context, containerID string, options client.ContainerInspectOptions) (client.ContainerInspectResult, error)
	ContainerList(ctx context.Context, options client.ContainerListOptions) (client.ContainerListResult, error)
	ContainerStats(ctx context.Context, containerID string, options client.ContainerStatsOptions) (client.ContainerStatsResult, error)
	Events(ctx context.Context, options client.EventsListOptions) client.EventsResult
	Close() error
}

func NewClient() (*Client, error) {
	client, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Client{
		api: client,
	}, nil
}

// FUNCTION FOR MOCKING DOCKER SDK API.
func NewClientWithMockAPI(mock DockerClient) *Client {
	return &Client{api: mock}
}

// Client that contains docker client or mock client.
type Client struct {
	api DockerClient
}

func (c *Client) InspectContainer(ctx context.Context, containerID string) (client.ContainerInspectResult, error) {
	return c.api.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
}

// TODO: Pass the options as param to ListContainers.
// List all containers in the docker host.
func (c *Client) ListContainers(ctx context.Context) (client.ContainerListResult, error) {
	return c.api.ContainerList(ctx, client.ContainerListOptions{All: true})
}

// Retrieves live resource usage statistics for the specified container.
func (c *Client) ListContainerStats(ctx context.Context, containerID string) (client.ContainerStatsResult, error) {
	return c.api.ContainerStats(ctx, containerID, client.ContainerStatsOptions{Stream: true})
}

// Subscribes to the Events stream.
func (c *Client) Events(ctx context.Context, opts client.EventsListOptions) client.EventsResult {
	return c.api.Events(ctx, opts)
}

func (c *Client) Close() error {
	return c.api.Close()
}
