package engine

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/client"
)

type MonitorEngine struct {
	Mu         sync.RWMutex
	Client     docker.Client
	Containers map[string]*Container
}

func NewEngine(client docker.Client) *MonitorEngine {
	return &MonitorEngine{
		Client:     client,
		Containers: make(map[string]*Container),
	}
}

// Launches event subscription and container monitoring process.
func (eng *MonitorEngine) Start(ctx context.Context) error {
	eventChan := make(chan events.Message)

	// Subscribe to the client event stream and handle
	// container start and stop events.
	go eng.handleEvents(ctx, eventChan)
	go eng.monitorEvents(ctx, eventChan)

	if err := eng.loadContainers(ctx); err != nil {
		return err
	}

	return nil
}

// Captures a snapshot of the current Containers state.
func (eng *MonitorEngine) ContainerSnapshot() []ContainerDTO {
	eng.Mu.RLock()
	containers := eng.Containers
	eng.Mu.RUnlock()

	snap := make([]ContainerDTO, 0, len(containers))
	for _, c := range containers {
		snap = append(snap, c.ToDTO())
	}

	return snap
}

// Loads containers from Docker host and starts collecting stats for
// any running containers.
func (eng *MonitorEngine) loadContainers(ctx context.Context) error {
	result, err := eng.Client.ListContainers(ctx)
	if err != nil {
		return err
	}

	for _, item := range result.Items {
		container := NewContainerFromListContainers(item)
		eng.Containers[container.id] = container

		if container.IsRunning() {
			childCtx, cancel := context.WithCancel(ctx)
			container.cancelFunc = cancel
			go container.CollectStats(childCtx, &eng.Client)
		}
	}

	return nil
}

// Subscribes and outputs docker daemon events.
func (eng *MonitorEngine) monitorEvents(
	ctx context.Context,
	output chan<- events.Message,
) {
	defer close(output)
	res := eng.Client.Events(ctx, client.EventsListOptions{})

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Event monitor closed\n")
			return
		case event := <-res.Messages:
			output <- event
		case err := <-res.Err:
			log.Printf("Error retrieving events: %v", err)
			return
		}
	}
}

// Checks incoming events.
func (eng *MonitorEngine) handleEvents(
	ctx context.Context,
	eventChan <-chan events.Message,
) {
	for e := range eventChan {
		if e.Actor.ID == "" {
			continue
		}

		switch e.Action {
		case events.ActionStart:
			container, err := eng.getOrCreateContainer(ctx, e.Actor.ID)
			if err != nil {
				log.Printf("Error getting container: %v\n", err)
				continue
			}
			eng.collectStats(ctx, container)
		case events.ActionDie:
			eng.stopStatCollection(e.Actor.ID)
		default:
			continue
		}
	}
}

// Creates new cancellable context and starts stat collection for the provided
// container.
func (eng *MonitorEngine) collectStats(
	ctx context.Context,
	container *Container,
) {
	container.mu.Lock()
	childCtx, cancel := context.WithCancel(ctx)
	container.cancelFunc = cancel
	container.mu.Unlock()

	go container.CollectStats(childCtx, &eng.Client)
}

// Checks for a valid container on the MonitorEngine or creates a new one.
func (eng *MonitorEngine) getOrCreateContainer(
	ctx context.Context,
	id string,
) (*Container, error) {
	eng.Mu.RLock()
	con, exists := eng.Containers[id]
	eng.Mu.RUnlock()

	if exists {
		con.mu.Lock()
		con.state = container.StateRunning
		con.mu.Unlock()
		return con, nil
	}

	info, err := eng.Client.InspectContainer(ctx, id)
	if err != nil {
		log.Printf("Error inspecting container %s: %v", id, err)
		return nil, err
	}

	con = NewContainerFromInspectContainer(info.Container)

	eng.Mu.Lock()
	if con, exist := eng.Containers[id]; exist {
		eng.Mu.Unlock()
		con.mu.Lock()
		con.state = container.StateRunning
		con.mu.Unlock()
		return con, nil
	}

	eng.Containers[id] = con
	eng.Mu.Unlock()

	return con, nil
}

// Stops stat collection for the provided container.
func (eng *MonitorEngine) stopStatCollection(id string) {
	eng.Mu.Lock()
	// Find the container in eng.Containers.
	con, exist := eng.Containers[id]
	if !exist {
		log.Printf("No container with id: %v\n", id)
		eng.Mu.Unlock()
	}
	eng.Mu.Unlock()

	// Call the container.cancel function.
	con.mu.Lock()
	con.cancelFunc()
	con.state = container.StateExited
	con.mu.Unlock()
}
