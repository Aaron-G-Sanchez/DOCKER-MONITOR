package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/engine"
	"github.com/aaron-g-sanchez/DOCKER-MONITOR/web/templates"
	"github.com/gin-gonic/gin"
)

type Server struct {
	monitorEngine *engine.MonitorEngine
	router        *gin.Engine
	http          *http.Server
}

// Initialize a new custom server instance.
func NewServer(monitor *engine.MonitorEngine) *Server {
	server := &Server{
		monitorEngine: monitor,
		router:        gin.Default(),
	}
	server.CreateRoutes()
	return server
}

// Assign routes and handlers to the router.
func (s *Server) CreateRoutes() {
	s.router.Static("/static", "./web/static")

	s.router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	s.router.GET("/", func(ctx *gin.Context) {
		render(ctx, http.StatusOK, templates.Home())
	})

	s.router.GET("/demo", s.handleDemo())
	s.router.GET("/containers", s.handleContainerData())
}

// Create http.Server instance and launch the server.
func (s *Server) Start(addr string) error {
	s.http = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	return s.http.ListenAndServe()
}

// Shutdown the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func (s *Server) handleContainerData() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		fmt.Println("New Client")

		clientGone := ctx.Request.Context().Done()

		rc := http.NewResponseController(ctx.Writer)
		t := time.NewTicker(time.Second)
		defer t.Stop()

		for {
			select {
			case <-clientGone:
				fmt.Println("Client disconnected")
				return
			case <-t.C:
				containers := s.monitorEngine.ContainerSnapshot()

				json, err := json.Marshal(containers)
				if err != nil {
					fmt.Println("Error marshalling container data", err)
				}

				// TODO: Get container data and pass to the data field.
				if _, err := fmt.Fprintf(ctx.Writer, "data: %s\n\n", json); err != nil {
					return
				}
				err = rc.Flush()
				if err != nil {
					fmt.Println("Error writing to response writer: ", err)
					return
				}
			}
		}

	}
}

// TODO: Move and replace handler function.
func (s *Server) handleDemo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s.monitorEngine.Mu.RLock()
		containers := s.monitorEngine.Containers
		s.monitorEngine.Mu.RUnlock()

		ctx.JSON(http.StatusOK, gin.H{
			"containers": containers,
		})
	}
}
