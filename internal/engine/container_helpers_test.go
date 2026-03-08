package engine

import (
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			"rounds down",
			.2222222,
			"0.22",
		},
		{
			"rounds up",
			.247,
			"0.25",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := format(tc.input)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestCalculateMemUsage(t *testing.T) {
	s := &container.MemoryStats{
		Usage: 200,
		Stats: map[string]uint64{
			"inactive_file": 100,
		},
	}

	got := calculateMemUsage(*s)

	assert.Equal(t, float64(100), got)
}

func TestCalculateMemUsagePerc(t *testing.T) {
	s := &container.MemoryStats{
		Usage: 100,
		Stats: map[string]uint64{
			"inactive_file": 20,
		},
		Limit: 1280,
	}

	uMem := calculateMemUsage(*s)

	got := calculateMemUsagePerc(uMem, *s)

	assert.Equal(t, float64(6.25), got)
}
