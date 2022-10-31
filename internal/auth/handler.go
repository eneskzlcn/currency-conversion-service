package auth

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	authService AuthService
	logger      *zap.SugaredLogger
}

func NewHandler(service AuthService, logger *zap.SugaredLogger) *Handler {
	return &Handler{authService: service, logger: logger}
}
func (h *Handler) Login(ctx *fiber.Ctx) error {
	var request LoginRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	tokenResponse, err := h.authService.Tokenize(ctx.Context(), request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.Status(fiber.StatusCreated).JSON(tokenResponse)
}
func (h *Handler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/auth")
	appGroup.Post("/login", h.Login)
}
