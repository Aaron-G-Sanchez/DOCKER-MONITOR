package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
)

// CLIENT STRUCT FOR MOCKING.
type MockAPIClient struct {
	containers client.ContainerListResult
	err        error
}

func (mock *MockAPIClient) ContainerList(ctx context.Context, _ client.ContainerListOptions) (client.ContainerListResult, error) {
	return mock.containers, mock.err
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
					Names: []string{"demo-container"},
					Image: "demo-image",
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

			mockClient := NewClientWithAPI(mockAPI)
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
