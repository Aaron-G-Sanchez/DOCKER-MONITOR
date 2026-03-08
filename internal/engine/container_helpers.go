package engine

import (
	"fmt"

	"github.com/moby/moby/api/types/container"
)

// TODO: Add check for cgroups (v1 & v2)
func calculateMemUsage(stat container.MemoryStats) float64 {
	return float64(stat.Usage - stat.Stats["inactive_file"])
}

func calculateMemUsagePerc(usedMem float64, stat container.MemoryStats) float64 {
	return usedMem / float64(stat.Limit) * 100
}

func format(in float64) string {
	return fmt.Sprintf("%.2f", in)
}
