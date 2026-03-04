package engine

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/testutils"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
)

// TODO: Create a setup function.

func TestEngine_StartSuccess(t *testing.T) {

	mockSum := &container.Summary{
		ID:    "a1b2c3",
		Names: []string{"mock-container"},
		State: container.StateExited,
		Image: "mock-image",
	}

	mockContainer := NewContainerFromListContainers(*mockSum)

	mockAPIClient := &testutils.MockDockerClient{
		MockContainers: client.ContainerListResult{
			Items: []container.Summary{*mockSum},
		},
		MockContainerStats: client.ContainerStatsResult{
			Body: io.NopCloser(strings.NewReader("")),
		},
		MockEventResults: client.EventsResult{
			Messages: make(chan events.Message),
			Err:      make(chan error),
		},
		Err: nil,
	}
	mockDockerClient := docker.NewClientWithMockAPI(mockAPIClient)

	mockEngine := NewEngine(*mockDockerClient)

	ctx, cancel := context.WithCancel(t.Context())
	err := mockEngine.Start(ctx)
	cancel()

	assert.NoError(t, err)
	assert.Equal(
		t,
		mockEngine.Containers["a1b2c3"],
		mockContainer,
		"Should set the containers field.",
	)
}

// TODO: Add test for loadContainers and event subscription.
