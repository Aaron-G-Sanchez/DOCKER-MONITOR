package engine

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/moby/moby/api/types/container"
)

// TODO: Add Windows support.
type Stat struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	OSType           string  `json:"os_type"`
	CPUPercentage    float64 `json:"cpu_percentage"`
	Memory           float64 `json:"memory"`
	MemoryPercentage float64 `json:"memory_percentage"`
	NetworkRx        float64 `json:"network_rx"`
	NetworkTx        float64 `json:"network_tx"`
}

type ContainerDTO struct {
	ID    string                   `json:"id"`
	Names []string                 `json:"names"`
	State container.ContainerState `json:"state"`
	Stats *Stat                    `json:"stats,omitempty"`
}

type Container struct {
	mu         sync.Mutex
	id         string
	names      []string
	state      container.ContainerState
	stats      *Stat
	cancelFunc context.CancelFunc
}

// Creates a new container from ListContainer response.
func NewContainerFromListContainers(s container.Summary) *Container {
	return &Container{
		id:    s.ID,
		names: s.Names,
		state: s.State,
	}
}

// Creates a new container from InspectContainer response.
func NewContainerFromInspectContainer(s container.InspectResponse) *Container {
	return &Container{
		id:    s.ID,
		names: []string{s.Name},
		state: s.State.Status,
	}
}

func (c *Container) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state == "running"
}

func (c *Container) CollectStats(ctx context.Context, client *docker.Client) {
	stats, err := client.ListContainerStats(ctx, c.id)
	if err != nil {
		log.Printf("Error reading stats: %v\n", err)
		return
	}
	defer stats.Body.Close()

	decoder := json.NewDecoder(stats.Body)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var rawStatResult *container.StatsResponse

		if err := decoder.Decode(&rawStatResult); err != nil {
			if err == io.EOF || err == context.Canceled {
				log.Printf("Stopping monitoring for %s (context canceled or stream ended)\n", c.id)
			} else {
				log.Printf("Error decoding stats for %s: %v\n", c.id, err)
			}
			return
		}

		usedMem := CalculateMemUsage(rawStatResult.MemoryStats)
		netRx, netTx := CalculateNetworkIO(rawStatResult.Networks)
		c.SetStats(&Stat{
			ID:               rawStatResult.ID,
			Name:             rawStatResult.Name,
			OSType:           rawStatResult.OSType,
			CPUPercentage:    CalculateCPUPerc(rawStatResult),
			Memory:           bytesToMB(usedMem),
			MemoryPercentage: CalculateMemUsagePerc(usedMem, rawStatResult.MemoryStats),
			NetworkRx:        netRx,
			NetworkTx:        netTx,
		})
	}
}

func (c *Container) SetStats(stats *Stat) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats = stats
}

func (c *Container) ToDTO() ContainerDTO {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: Create function to parse out single name and remove backslash.
	return ContainerDTO{
		ID:    c.id,
		Names: c.names,
		State: c.state,
		Stats: c.stats,
	}
}
