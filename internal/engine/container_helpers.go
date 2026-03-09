package engine

import (
	"github.com/moby/moby/api/types/container"
)

// TODO: Add check for cgroups (v1 & v2)
func calculateMemUsage(stat container.MemoryStats) float64 {
	return float64(stat.Usage - stat.Stats["inactive_file"])
}

func bytesToMB(num float64) float64 {
	return num / (1024 * 1024)
}

func calculateMemUsagePerc(usedMem float64, stat container.MemoryStats) float64 {
	return usedMem / float64(stat.Limit) * 100
}
