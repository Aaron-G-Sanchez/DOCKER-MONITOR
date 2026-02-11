package engine

import (
	"testing"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/testutils"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
)

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
		Err: nil,
	}
	mockDockerClient := docker.NewClientWithMockAPI(mockAPIClient)

	mockEngine := CreateEngine(t.Context(), *mockDockerClient)

	err := mockEngine.Start()

	assert.NoError(t, err)
	assert.Equal(
		t,
		mockAPIClient.MockContainers,
		*mockEngine.Containers,
		"Should set the containers field.",
	)

}
