package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/docker"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/engine"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/server"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	run(ctx)
}

func run(ctx context.Context) {
	client, err := docker.NewClient()
	if err != nil {
		log.Fatalf("Error starting docker client: %v\n", err)
	}

	// CREATE AND RUN THE DOCKER ENGINE.
	engine := engine.NewEngine(*client)
	defer engine.Client.Close()

	if err := engine.Start(ctx); err != nil {
		log.Fatalf("Error running engine: %v\n", err)
	}

	// CREATE AND RUN THE SERVER INSTANCE.
	server := server.NewServer(engine)

	go func() {
		fmt.Println("Starting server on :6060")
		if err := server.Start(":6060"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutdown initiated")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	fmt.Print("Server exiting.")
}
