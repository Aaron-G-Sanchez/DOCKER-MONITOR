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
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/client"
)

func NewEngine(client docker.Client) *MonitorEngine {
	return &MonitorEngine{
		Client: client,
	}
}

// TODO: Update ContainerStats field.
type MonitorEngine struct {
	Mu             sync.Mutex
	Client         docker.Client
	Containers     *client.ContainerListResult
	ContainerStats map[string]*container.StatsResponse
}

func (eng *MonitorEngine) Start(ctx context.Context) error {
	eventChan := make(chan events.Message)

	// TODO: Add subscription to docker events.
	// Subscribe to the client event stream and handle
	// container start and stop events.
	go eng.handleEvents(ctx, eventChan)
	go eng.monitorEvents(ctx, eventChan)

	if err := eng.refreshContainers(ctx); err != nil {
		return err
	}

	eng.ContainerStats = make(map[string]*container.StatsResponse)

	for _, container := range eng.Containers.Items {
		go eng.getContainerStats(ctx, container.ID)
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

func (eng *MonitorEngine) refreshContainers(ctx context.Context) error {
	result, err := eng.Client.ListContainers(ctx)
	if err != nil {
		return err
	}

	eng.Containers = &result
	return nil
}

func (eng *MonitorEngine) getContainerStats(ctx context.Context, id string) {
	stats, err := eng.Client.ListContainerStats(ctx, id)
	if err != nil {
		log.Printf("Error Reading stats: %v\n", err)
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
		var statResult *container.StatsResponse

		if err := decoder.Decode(&statResult); err != nil {
			if err == io.EOF || err == context.Canceled {
				log.Printf("Stopping monitoring for %s (context canceled or stream ended)\n", id)
			} else {
				log.Printf("Error decoding stats for %s: %v\n", id, err)
			}
			return
		}

		eng.Mu.Lock()
		eng.ContainerStats[id] = statResult
		eng.Mu.Unlock()

	}
}

func (eng *MonitorEngine) handleEvents(ctx context.Context, eventChan <-chan events.Message) {
	addContainer := func(ctx context.Context, id string) {
		go eng.getContainerStats(ctx, id)
	}

	// TODO: Add event handling for die events.
	for e := range eventChan {
		// fmt.Println(e.Action)
		if e.Action == events.ActionStart {
			addContainer(ctx, e.Actor.ID)
		}
	}

}
