package engine

import (
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/assert"
)

func TestToDTO(t *testing.T) {
	stats := &Stat{
		ID:            "stat-id",
		CPUPercentage: 25.2,
		Memory:        512.00,
	}

	c := &Container{
		id:    "a1b2c3",
		names: []string{"mock-con"},
		state: container.StateExited,
		stats: stats,
	}

	got := c.ToDTO()

	assert.Equal(t, got.Stats, stats)
}
