package api

import (
	"net/http"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/engine"
	"github.com/gin-gonic/gin"
)

type Server struct {
	monitorEngine *engine.MonitorEngine
	router        *gin.Engine
}

func NewServer(monitor *engine.MonitorEngine) *Server {
	server := &Server{
		monitorEngine: monitor,
		router:        gin.Default(),
	}
	server.CreateRoutes()
	return server
}

func (s *Server) CreateRoutes() {
	s.router.GET("/", s.handleDemo())
}

// TODO: Replace handler function.
func (s *Server) handleDemo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"containers": *s.monitorEngine.Containers,
		})
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
