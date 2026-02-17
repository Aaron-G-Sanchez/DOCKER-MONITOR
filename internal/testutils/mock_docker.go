package testutils

import (
	"context"

	"github.com/moby/moby/client"
)

type MockAPIClient struct {
	MockContainers     client.ContainerListResult
	MockContainerStats client.ContainerStatsResult
	Err                error
}

func (mock *MockAPIClient) ContainerList(
	ctx context.Context,
	_ client.ContainerListOptions,
) (client.ContainerListResult, error) {
	return mock.MockContainers, mock.Err
}

func (mock *MockAPIClient) ContainerStats(
	ctx context.Context,
	containerID string,
	_ client.ContainerStatsOptions,
) (client.ContainerStatsResult, error) {
	return mock.MockContainerStats, mock.Err
}

func (mock *MockAPIClient) Close() error {
	return mock.Err
}
