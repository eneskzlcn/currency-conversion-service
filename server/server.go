package server

import (
	"fmt"
	"github.com/eneskzlcn/currency-conversion-service/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Handler interface {
	RegisterRoutes(app *fiber.App)
}

type Server struct {
	app    *fiber.App
	config config.Server
	logger *zap.SugaredLogger
}

func New(handlers []Handler, config config.Server, logger *zap.SugaredLogger) *Server {
	app := fiber.New()
	app.Use(cors.New(cors.ConfigDefault))
	for _, handler := range handlers {
		handler.RegisterRoutes(app)
	}
	server := &Server{
		app:    app,
		config: config,
		logger: logger,
	}
	server.AddRoutes()
	return server
}
func (s *Server) AddRoutes() {
	s.app.Get("/health", s.healthCheck)
}
func (s *Server) healthCheck(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

func (s *Server) Start() error {
	address := fmt.Sprintf(":%s", s.config.Port)
	shutDownChan := make(chan os.Signal, 1)
	signal.Notify(shutDownChan, os.Interrupt, syscall.SIGTERM, os.Kill)
	go func() {
		<-shutDownChan
		if err := s.app.Shutdown(); err != nil {
			s.logger.Error(err)
		}
	}()
	return s.app.Listen(address)
}
