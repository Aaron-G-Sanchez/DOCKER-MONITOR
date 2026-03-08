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

func TestMemUsage(t *testing.T) {
	s := &container.MemoryStats{
		Usage: 100600,
		Stats: map[string]uint64{
			"inactive_file": 600,
		},
	}

	got := calculateMemUsage(*s)

	assert.Equal(t, float64(100000), got)

}
