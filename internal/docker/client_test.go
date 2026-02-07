package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_Success(t *testing.T) {

	client, err := NewClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
