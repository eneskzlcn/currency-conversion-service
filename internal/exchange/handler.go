package exchange

import "github.com/gofiber/fiber/v2"

type ExchangeService interface {
}
type Handler struct {
	exchangeService ExchangeService
}

func NewHandler(service ExchangeService) *Handler {
	if service == nil {
		return nil
	}
	return &Handler{exchangeService: service}
}
func (h *Handler) GetExchangeRate(ctx *fiber.Ctx) error {
	return nil
}
func (h *Handler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/exchange")
	appGroup.Get("/rate", h.GetExchangeRate)

}
