package engine

import (
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/assert"
)

func TestCalculateMemUsage(t *testing.T) {
	s := &container.MemoryStats{
		Usage: 200,
		Stats: map[string]uint64{
			"inactive_file": 100,
		},
	}

	got := CalculateMemUsage(*s)

	assert.Equal(t, float64(100), got)
}

func TestBytesToMB(t *testing.T) {
	got := bytesToMB(float64(1024))

	assert.Equal(t, float64(1024), got)
}

func TestCalculateMemUsagePerc(t *testing.T) {
	s := &container.MemoryStats{
		Usage: 100,
		Stats: map[string]uint64{
			"inactive_file": 20,
		},
		Limit: 1280,
	}

	uMem := CalculateMemUsage(*s)

	got := CalculateMemUsagePerc(uMem, *s)

	assert.Equal(t, float64(6.25), got)
}

func TestCalculateCPUPerc(t *testing.T) {
	s := &container.StatsResponse{
		CPUStats: container.CPUStats{
			CPUUsage: container.CPUUsage{
				TotalUsage: 200,
			},
			SystemUsage: 2000,
			OnlineCPUs:  4,
		},
		PreCPUStats: container.CPUStats{
			CPUUsage: container.CPUUsage{
				TotalUsage: 100,
			},
			SystemUsage: 1000,
		},
	}

	got := CalculateCPUPerc(s)

	assert.Equal(t, float64(40), got)
}
