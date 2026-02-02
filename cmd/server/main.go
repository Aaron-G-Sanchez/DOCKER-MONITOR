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

		// Time
		cur := time.Now()
		hr, min, sec := cur.Clock()
		y, m, d := cur.Date()
		var monthNum int
		switch m {
		case time.January:
			monthNum = 1
		case time.February:
			monthNum = 2
		case time.March:
			monthNum = 3
		case time.April:
			monthNum = 4
		case time.May:
			monthNum = 5
		case time.June:
			monthNum = 6
		case time.July:
			monthNum = 7
		case time.August:
			monthNum = 8
		case time.September:
			monthNum = 9
		case time.October:
			monthNum = 10
		case time.November:
			monthNum = 11
		case time.December:
			monthNum = 12
		default:
			panic("Invalid month")
		}

		dateTime := fmt.Sprintf("%v-%v-%v %v:%v:%v", y, monthNum, d, hr, min, sec)

		// MEM%
		usedMem := s.MemoryStats.Usage - s.MemoryStats.Stats["inactive_file"]
		memPercent := (float64(usedMem) / float64(s.MemoryStats.Limit)) * 100
		sprint := fmt.Sprintf("%.2f", memPercent)

		usedMemInMb := usedMem / 1048576

		// TODO: Need to modify the cpu percent to display in decimal point.
		// CPU Usage
		cpuD := s.CPUStats.CPUUsage.TotalUsage
		systemCpuD := s.CPUStats.SystemUsage - s.PreCPUStats.SystemUsage
		numCpu := s.CPUStats.OnlineCPUs
		cpuPercent := (cpuD / systemCpuD) * uint64(numCpu) * 100

		fmt.Printf("Time: %v, Memory Percent : %v%%, Mem Usage: %dmb, cpu usage: %d\n",
			dateTime,
			sprint,
			usedMemInMb,
			cpuPercent)
	}

}
