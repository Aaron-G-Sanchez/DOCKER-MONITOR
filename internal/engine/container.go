package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/moby/moby/api/types/container"
)

// TODO: Add Windows support.
// TODO: Create calculation function.
type Stat struct {
	ID               string
	Name             string
	OSType           string
	CPUPercentage    float64 // Calc %
	Memory           float64 // Calc
	MemoryPercentage float64 // Calc %
	NetworkRx        float64 // Calc network io
	NetworkTx        float64 // Calc network io
}

// TODO: Calculate network traffic, and mem/cpu usage.
func NewStat(entry container.StatsResponse) *Stat {
	return &Stat{
		ID:   entry.ID,
		Name: entry.Name,
	}
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

		statEntry := NewStat(*rawStatResult)

		c.mu.Lock()
		c.stats = statEntry
		fmt.Printf("CONTAINER: %s\n", statEntry.Name)
		c.mu.Unlock()

	}
}
