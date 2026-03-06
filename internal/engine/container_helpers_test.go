package engine

import (
	"testing"

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
