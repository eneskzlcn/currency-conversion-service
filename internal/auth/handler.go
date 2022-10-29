package auth

import "github.com/gofiber/fiber/v2"

type Handler struct {
	authService AuthService
}

func NewHandler(service AuthService) *Handler {
	if service == nil {
		return nil
	}
	return &Handler{authService: service}
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
	app.Post("/login", h.Login)
}
