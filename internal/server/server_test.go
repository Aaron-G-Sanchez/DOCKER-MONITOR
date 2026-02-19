package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/engine"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockContainers := &client.ContainerListResult{
		Items: []container.Summary{
			{
				ID:    "mock-container",
				Names: []string{"mock-container"},
				Image: "mock-image:latest",
			},
			{
				ID:    "mock-container-two",
				Names: []string{"mock-container-two"},
				Image: "mock-image",
			},
		},
	}

	mockServer := setup(mockContainers, t)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)

	mockServer.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func setup(containers *client.ContainerListResult, t *testing.T) *Server {

	mockAPIClient := &testutils.MockAPIClient{
		MockContainers: *containers,
		MockContainerStats: client.ContainerStatsResult{
			Body: io.NopCloser(strings.NewReader("")),
		},
	}

	mockDockerClient := docker.NewClientWithMockAPI(mockAPIClient)
	mockEngine := engine.NewEngine(*mockDockerClient)
	defer mockEngine.Client.Close()

	mockEngine.Start(t.Context())

	mockServer := NewServer(mockEngine)

	return mockServer
}
