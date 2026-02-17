package server

import (
	"net/http"

	"github.com/aaron-g-sanchez/DOCKER-MONITOR/internal/engine"
	"github.com/gin-gonic/gin"
)

// TODO: Refactor to use http.Server
type Server struct {
	monitorEngine *engine.MonitorEngine
	router        *gin.Engine
	http          *http.Server
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

// TODO: Move and replace handler function.
func (s *Server) handleDemo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"containers": *s.monitorEngine.Containers,
		})
	}
}

func (s *Server) Start(addr string) error {
	s.http = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	return s.http.ListenAndServe()
}
