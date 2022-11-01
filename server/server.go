package server

import (
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/eneskzlcn/currency-conversion-service/config"
	_ "github.com/eneskzlcn/currency-conversion-service/docs"
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
	s.app.Get("/swagger/*", swagger.HandlerDefault)
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description Get the status of server.
// @Tags Health
// @Accept */*
// @Success 200
// @Failure 404
// @Router /health [get]
func (s *Server) healthCheck(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

// @title Currency Conversion Service
// @version 2.0
// @description This is a currency conversion service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost
// @BasePath /
// @schemes http
func (s *Server) Start() error {
	address := fmt.Sprintf(":%s", s.config.Port)
	shutDownChan := make(chan os.Signal, 1)
	signal.Notify(shutDownChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutDownChan
		if err := s.app.Shutdown(); err != nil {
			s.logger.Error(err)
		}
	}()
	return s.app.Listen(address)
}
