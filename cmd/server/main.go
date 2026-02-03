package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func main() {

	// Sample code from the Docker SDK docs
	ctx := context.Background()
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	testContainerID := "49283dece168"

	stats, err := apiClient.ContainerStats(ctx, testContainerID, client.ContainerStatsOptions{Stream: true})
	if err != nil {
		panic(err)
	}
	defer stats.Body.Close()

	decoder := json.NewDecoder(stats.Body)

	for {
		var s container.StatsResponse

		if err := decoder.Decode(&s); err != nil {
			log.Fatal(err)
		}

		// Get / Format Time
		cur := time.Now()
		hr, min, sec := cur.Clock()
		y, _, d := cur.Date()
		m := int(cur.Month())

		dateTime := fmt.Sprintf("%v-%v-%v %v:%v:%v", y, m, d, hr, min, sec)

		// Calc. Memory Stats
		usedMem := s.MemoryStats.Usage - s.MemoryStats.Stats["inactive_file"]
		memPercent := (float64(usedMem) / float64(s.MemoryStats.Limit)) * 100
		memResult := fmt.Sprintf("%.2f", memPercent)

		usedMemInMb := usedMem / 1048576

		// Calc. CPU Usage
		cpuD := float64(s.CPUStats.CPUUsage.TotalUsage) - float64(s.PreCPUStats.CPUUsage.TotalUsage)
		systemCpuD := float64(s.CPUStats.SystemUsage) - float64(s.PreCPUStats.SystemUsage)
		numCpu := s.CPUStats.OnlineCPUs
		cpuPercent := (cpuD / systemCpuD) * float64(numCpu) * 100

		cpuResult := fmt.Sprintf("%.2f", cpuPercent)

		fmt.Printf("Time: %v, Memory Percent : %v%%, Mem Usage: %dmb, cpu usage: %v\n",
			dateTime,
			memResult,
			usedMemInMb,
			cpuResult)
	}

}
