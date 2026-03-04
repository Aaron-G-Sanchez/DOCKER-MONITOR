package engine

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/client"
)

func NewEngine(client docker.Client) *MonitorEngine {
	return &MonitorEngine{
		Client:     client,
		Containers: make(map[string]*Container),
	}
}

// TODO: #33 Update ContainerStats field.
type MonitorEngine struct {
	Mu         sync.RWMutex
	Client     docker.Client
	Containers map[string]*Container
}

func (eng *MonitorEngine) Start(ctx context.Context) error {
	eventChan := make(chan events.Message)

	// TODO: Add subscription to docker events.
	// Subscribe to the client event stream and handle
	// container start and stop events.
	go eng.handleEvents(ctx, eventChan)
	go eng.monitorEvents(ctx, eventChan)

	if err := eng.loadContainers(ctx); err != nil {
		return err
	}

	return nil
}

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

// Loads containers from Docker host and starts collecting stats for
// any running containers.
func (eng *MonitorEngine) loadContainers(ctx context.Context) error {
	result, err := eng.Client.ListContainers(ctx)
	if err != nil {
		return err
	}

	for _, item := range result.Items {
		container := NewContainer(item)
		eng.Containers[container.id] = container

		if container.IsRunning() {
			childCtx, cancel := context.WithCancel(ctx)
			container.cancelFunc = cancel
			go container.CollectStats(childCtx, &eng.Client)
		}
	}

	return nil
}

func (eng *MonitorEngine) handleEvents(ctx context.Context, eventChan <-chan events.Message) {
	// Create cancellable context and start stat collection.
	collectStats := func(ctx context.Context, container *Container) {
		container.mu.Lock()
		childCtx, cancel := context.WithCancel(ctx)
		container.cancelFunc = cancel
		container.mu.Unlock()

		go container.CollectStats(childCtx, &eng.Client)
	}

	// TODO: Create new container and start stat collection. <----- [STOPPED HERE]
	// TODO: Add event handling for die events.
	for e := range eventChan {

		if e.Actor.ID == "" {
			continue
		}

		switch e.Action {
		case events.ActionStart:
			eng.Mu.Lock()
			container, exists := eng.Containers[e.Actor.ID]
			if !exists {
				fmt.Printf("CREATE NEW CONTAINER INSTANCE: %v", e.Actor.ID)
				eng.Mu.Unlock()
				continue
			}

			eng.Mu.Unlock()
			collectStats(ctx, container)
		default:
			continue
		}
	}
}
