package testutils

import (
	"context"

	"github.com/moby/moby/client"
)

type MockDockerClient struct {
	MockContainerInspect client.ContainerInspectResult
	MockContainers       client.ContainerListResult
	MockContainerStats   client.ContainerStatsResult
	MockEventResults     client.EventsResult
	Err                  error
}

func (mock *MockDockerClient) ContainerInspect(
	ctx context.Context,
	_ string,
	_ client.ContainerInspectOptions,
) (client.ContainerInspectResult, error) {
	return mock.MockContainerInspect, mock.Err
}

func (mock *MockDockerClient) ContainerList(
	ctx context.Context,
	_ client.ContainerListOptions,
) (client.ContainerListResult, error) {
	return mock.MockContainers, mock.Err
}

func (mock *MockDockerClient) ContainerStats(
	ctx context.Context,
	containerID string,
	_ client.ContainerStatsOptions,
) (client.ContainerStatsResult, error) {
	return mock.MockContainerStats, mock.Err
}

func (mock *MockDockerClient) Events(
	ctx context.Context,
	_ client.EventsListOptions,
) client.EventsResult {
	return mock.MockEventResults
}

func (mock *MockDockerClient) Close() error {
	return mock.Err
}
