package engine

import (
	"context"
	"sync"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
)

func NewEngine(client docker.Client) *MonitorEngine {
	return &MonitorEngine{
		Client: client,
	}
}

// TODO: #33 Update ContainerStats field.
type MonitorEngine struct {
	Mu         sync.Mutex
	Client     docker.Client
	Containers map[string]*Container
}

func (eng *MonitorEngine) Start(ctx context.Context) error {
	// eventChan := make(chan events.Message)
	eng.Containers = make(map[string]*Container)

	// TODO: Add subscription to docker events.
	// Subscribe to the client event stream and handle
	// container start and stop events.
	// go eng.handleEvents(ctx, eventChan)
	// go eng.monitorEvents(ctx, eventChan)

	if err := eng.loadContainers(ctx); err != nil {
		return err
	}

	// eng.ContainerStats = make(map[string]*container.StatsResponse)

	// for _, container := range eng.Containers.Items {
	// 	go eng.getContainerStats(ctx, container.ID)
	// }

	return nil
}

// func (eng *MonitorEngine) monitorEvents(
// 	ctx context.Context,
// 	output chan<- events.Message,
// ) {
// 	defer close(output)
// 	res := eng.Client.Events(ctx, client.EventsListOptions{})

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Printf("Event monitor closed\n")
// 			return
// 		case event := <-res.Messages:
// 			output <- event
// 		case err := <-res.Err:
// 			log.Printf("Error retrieving events: %v", err)
// 			return
// 		}
// 	}
// }

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
			// TODO: Implement.
			go container.CollectStats(childCtx, &eng.Client)
		}
	}

	return nil
}

// func (eng *MonitorEngine) getContainerStats(ctx context.Context, id string) {
// 	stats, err := eng.Client.ListContainerStats(ctx, id)
// 	if err != nil {
// 		log.Printf("Error Reading stats: %v\n", err)
// 		return
// 	}
// 	defer stats.Body.Close()

// 	decoder := json.NewDecoder(stats.Body)

// 	for {

// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 		}
// 		var statResult *container.StatsResponse

// 		if err := decoder.Decode(&statResult); err != nil {
// 			if err == io.EOF || err == context.Canceled {
// 				log.Printf("Stopping monitoring for %s (context canceled or stream ended)\n", id)
// 			} else {
// 				log.Printf("Error decoding stats for %s: %v\n", id, err)
// 			}
// 			return
// 		}

// 		// TODO: Create a Stats entry with details from statResult

// 		eng.Mu.Lock()
// 		eng.ContainerStats[id] = statResult
// 		eng.Mu.Unlock()

// 	}
// }

// func (eng *MonitorEngine) handleEvents(ctx context.Context, eventChan <-chan events.Message) {
// 	collectStats := func(ctx context.Context, id string) {
// 		go eng.getContainerStats(ctx, id)
// 	}

// 	// TODO: Add event handling for die events.
// 	for e := range eventChan {
// 		// TODO: Add check to ensure container id is present.
// 		if e.Actor.ID == "" {
// 			continue
// 		}

// 		switch e.Action {
// 		case events.ActionStart:
// 			// TODO: Add check to make sure container is not already having stats collected.
// 			if _, present := eng.ContainerStats[e.Actor.ID]; !present {
// 				collectStats(ctx, e.Actor.ID)
// 			}
// 		default:
// 			continue
// 		}
// 	}
// }
