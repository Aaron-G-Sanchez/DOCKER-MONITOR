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

func TestContainerSnapshot(t *testing.T) {
	mockContainer := &Container{
		id:    "abc",
		names: []string{"mock-container"},
		state: container.StateExited,
	}

	expected := []ContainerDTO{{
		ID:    mockContainer.id,
		Names: mockContainer.names,
		State: mockContainer.state,
	},
	}

	mockAPI := &testutils.MockDockerClient{}
	mockClient := docker.NewClientWithMockAPI(mockAPI)

	eng := NewEngine(*mockClient)

	eng.Containers["abc"] = mockContainer

	got := eng.ContainerSnapshot()

	assert.Equal(t, got, expected)

}

func TestGetOrCreateContainer_ExistingContainer(t *testing.T) {
	mockApi := &testutils.MockDockerClient{}
	mockClient := docker.NewClientWithMockAPI(mockApi)

	eng := NewEngine(*mockClient)

	mockContainer := &Container{
		id:    "a1b2c3",
		names: []string{"mock-container"},
	}
	eng.Containers["a1b2ce"] = mockContainer

	got, err := eng.getOrCreateContainer(t.Context(), "a1b2ce")
	assert.NoError(t, err, "Should not throw error")

	assert.Equal(t, mockContainer, got)
}

func TestGetOrCreateContainer_NewContainer(t *testing.T) {
	expectedContainer := &Container{
		id:    "123",
		names: []string{"mock-container-inspect"},
		state: container.StateExited,
	}

	mockContainerInspect := container.InspectResponse{
		ID:   "123",
		Name: "mock-container-inspect",
		State: &container.State{
			Status: container.StateExited,
		},
	}

	mockApi := &testutils.MockDockerClient{
		MockContainerInspect: client.ContainerInspectResult{
			Container: mockContainerInspect,
		},
		Err: nil,
	}

	mockClient := docker.NewClientWithMockAPI(mockApi)

	eng := NewEngine(*mockClient)

	got, err := eng.getOrCreateContainer(t.Context(), "123")

	assert.NoError(t, err, "Should not throw error")
	assert.Equal(t, expectedContainer, got, "Should create new container")
}
