package server

import (
	"context"
	"net/http"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/engine"
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
	s.router.GET("/", s.handleDemo())
}

// TODO: Move and replace handler function.
func (s *Server) handleDemo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"containers": *s.monitorEngine.Containers,
		})
	}
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
