package engine

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/testutils"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
)

// TODO: Create a setup function.

func TestEngine_StartSuccess(t *testing.T) {
	mockAPIClient := &testutils.MockAPIClient{
		MockContainers: client.ContainerListResult{
			Items: []container.Summary{
				{
					ID:    "1234",
					Names: []string{"mock-container"},
					Image: "mock-image",
				},
			},
		},
		MockContainerStats: client.ContainerStatsResult{
			Body: io.NopCloser(strings.NewReader("")),
		},
		Err: nil,
	}
	mockDockerClient := docker.NewClientWithMockAPI(mockAPIClient)

	mockEngine := CreateEngine(*mockDockerClient)

	ctx, cancel := context.WithCancel(t.Context())
	err := mockEngine.Start(ctx)
	cancel()

	assert.NoError(t, err)
	assert.Equal(
		t,
		mockAPIClient.MockContainers,
		*mockEngine.Containers,
		"Should set the containers field.",
	)

}

func TestEngine_getContainerStats(t *testing.T) {
	mockStatsResponse := container.StatsResponse{
		ID:     "1234",
		Name:   "mock-container",
		OSType: "ubuntu",
	}
	mockStatsResponseJson, err := json.Marshal(mockStatsResponse)
	assert.NoError(t, err)

	mockStatsString := string(mockStatsResponseJson) + "\n"

	mockBody := io.NopCloser(strings.NewReader(mockStatsString))

	mockAPIClient := &testutils.MockAPIClient{
		MockContainers: client.ContainerListResult{
			Items: []container.Summary{
				{
					ID:    "1234",
					Names: []string{"mock-container"},
					Image: "mock-image",
				},
			},
		},
		MockContainerStats: client.ContainerStatsResult{
			Body: mockBody,
		},
		Err: nil,
	}

	mockDockerClient := docker.NewClientWithMockAPI(mockAPIClient)

	mockEngine := CreateEngine(*mockDockerClient)
	mockEngine.ContainerStats = make(map[string]*container.StatsResponse)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	go mockEngine.getContainerStats(ctx, "1234")

	assert.Eventually(t, func() bool {
		mockEngine.Mu.Lock()
		defer mockEngine.Mu.Unlock()
		return mockEngine.ContainerStats["1234"] != nil

	}, time.Second, 100*time.Millisecond, "Should collect container stats")

	mockEngine.Mu.Lock()
	assert.Equal(t, "1234", mockEngine.ContainerStats["1234"].ID)
	mockEngine.Mu.Unlock()
}
