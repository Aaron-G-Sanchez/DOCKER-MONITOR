package engine

import (
	"github.com/moby/moby/api/types/container"
)

// TODO: Add check for cgroups (v1 & v2)
func CalculateMemUsage(stat container.MemoryStats) float64 {
	return float64(stat.Usage - stat.Stats["inactive_file"])
}

func bytesToMB(num float64) float64 {
	return num / (1024 * 1024)
}

func CalculateMemUsagePerc(usedMem float64, stat container.MemoryStats) float64 {
	return usedMem / float64(stat.Limit) * 100
}

func CalculateCPUPerc(stats *container.StatsResponse) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemCPUDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	cpuCount := stats.CPUStats.OnlineCPUs

	return (cpuDelta / systemCPUDelta) * float64(cpuCount) * 100
}
