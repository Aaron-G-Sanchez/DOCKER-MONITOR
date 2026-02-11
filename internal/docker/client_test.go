package docker

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
)

// TODO: Replace with the mock in the test utils.
// CLIENT STRUCT FOR MOCKING.
type MockAPIClient struct {
	containers     client.ContainerListResult
	containerStats client.ContainerStatsResult
	err            error
}

func (mock *MockAPIClient) ContainerList(ctx context.Context, _ client.ContainerListOptions) (client.ContainerListResult, error) {
	return mock.containers, mock.err
}

func (mock *MockAPIClient) ContainerStats(ctx context.Context, containerID string, _ client.ContainerStatsOptions) (client.ContainerStatsResult, error) {
	return mock.containerStats, mock.err
}

func (mock *MockAPIClient) Close() error {
	return mock.err
}

func TestNewClient_Success(t *testing.T) {

	client, err := NewClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestListContainers(t *testing.T) {
	tests := []struct {
		name       string
		containers client.ContainerListResult
		err        error
		len        int
		wantErr    bool
	}{
		{

			name: "success with container",
			containers: client.ContainerListResult{
				Items: []container.Summary{{
					ID:    "1234",
					Names: []string{"mock-container"},
					Image: "mock-image",
				}},
			},
			err:     nil,
			len:     1,
			wantErr: false,
		}, {
			name:       "error fetching containers",
			containers: client.ContainerListResult{},
			err:        errors.New("error from API"),
			len:        0,
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI := &MockAPIClient{
				containers: tc.containers,
				err:        tc.err,
			}

			mockClient := NewClientWithMockAPI(mockAPI)
			result, err := mockClient.ListContainers(t.Context())

			if tc.wantErr {
				assert.Error(t, tc.err, "Should throw error")
			} else {
				assert.NoError(t, err, "Should not throw error")
				assert.Len(t, result.Items, tc.len, "Should have 1 container")
				assert.Equal(t, tc.containers.Items, result.Items, "Should be equal")
			}
		})
	}
}

func TestListContainerStats(t *testing.T) {
	mockStatsResponse := container.StatsResponse{
		ID:     "1234",
		Name:   "mock-container",
		OSType: "ubuntu",
	}
	mockStatsResponseJson, err := json.Marshal(mockStatsResponse)
	if err != nil {
		t.Errorf("Error marshalling response data: %v", err)
	}
	mockStatsString := string(mockStatsResponseJson)

	tests := []struct {
		name           string
		containerStats client.ContainerStatsResult
		containerID    string
		err            error
		wantErr        bool
	}{
		{
			name: "success with container stats",
			containerStats: client.ContainerStatsResult{
				Body: io.NopCloser(strings.NewReader(mockStatsString)),
			},
			containerID: "1234",
			err:         nil,
			wantErr:     false,
		}, {
			name:           "error fetching container stats",
			containerStats: client.ContainerStatsResult{},
			containerID:    "1234",
			err:            errors.New("error from API"),
			wantErr:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockApi := &MockAPIClient{
				containerStats: tc.containerStats,
				err:            tc.err,
			}

			mockClient := NewClientWithMockAPI(mockApi)
			result, err := mockClient.ListContainerStats(t.Context(), tc.containerID)

			if tc.wantErr {
				assert.Error(t, err, "Should throw error")
				return
			}

			assert.NoError(t, err, "Should not throw error")
			defer result.Body.Close()

			decoder := json.NewDecoder(result.Body)

			var got container.StatsResponse

			if err := decoder.Decode(&got); err != nil {
				t.Errorf("Error decoding body: %v\n", err)
			}

			assert.Equal(t, mockStatsResponse, got, "Should return the proper stats")
		})
	}
}
