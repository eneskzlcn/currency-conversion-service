package exchange

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/app/common"
	"github.com/eneskzlcn/currency-conversion-service/app/common/httperror"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ExchangeService interface {
	PrepareExchangeRateOffer(ctx context.Context, userID int, request ExchangeRateRequest) (ExchangeRateResponse, error)
}
type AuthGuard interface {
	ProtectWithJWT(handler fiber.Handler) fiber.Handler
}
type HttpHandler struct {
	exchangeService ExchangeService
	authGuard       AuthGuard
	logger          *zap.SugaredLogger
}

func NewHttpHandler(service ExchangeService, guard AuthGuard, logger *zap.SugaredLogger) *HttpHandler {
	return &HttpHandler{exchangeService: service, authGuard: guard, logger: logger}
}

//GetExchangeRate godoc
//@Summary Create an exchange rate offer
//@Description creates an exchange rate offer for given currencies
//@Param exchangeRateRequest body ExchangeRateRequest true "body params"
//@Param accessToken header string true "header params"
// @Tags Exchange
//@Accept  json
//@Produce  json
//@Success 200 {object} ExchangeRateResponse
//@Failure 400 {object} httperror.HttpError
//@Failure 401 {string} string "Unauthorized"
//@Failure 404
//@Failure 500 {object} httperror.HttpError
//@Router /exchange/rate [get]
func (h *HttpHandler) GetExchangeRate(ctx *fiber.Ctx) error {
	userID, exists := ctx.Locals(common.USER_ID_CTX_KEY).(int)
	h.logger.Info("Exchange Rate Offer Request Arrived. User ID: %d", userID)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httperror.NewBadRequestError(common.UserNotInContext))
	}
	var request ExchangeRateRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httperror.NewBadRequestError(err.Error()))
	}
	exchangeRate, err := h.exchangeService.PrepareExchangeRateOffer(ctx.Context(), userID, request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperror.NewInternalServerError(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(exchangeRate)
}
func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/exchange")
	appGroup.Get("/rate", h.authGuard.ProtectWithJWT(h.GetExchangeRate))
}
