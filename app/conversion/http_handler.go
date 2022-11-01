package conversion

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/app/common"
	"github.com/eneskzlcn/currency-conversion-service/app/common/httperror"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ConversionService interface {
	ConvertCurrencies(ctx context.Context, userID int, request CurrencyConversionOfferRequest) (bool, error)
}
type AuthGuard interface {
	ProtectWithJWT(handler fiber.Handler) fiber.Handler
}

type HttpHandler struct {
	conversionService ConversionService
	authGuard         AuthGuard
	logger            *zap.SugaredLogger
}

func NewHttpHandler(service ConversionService, guard AuthGuard, logger *zap.SugaredLogger) *HttpHandler {
	return &HttpHandler{conversionService: service, authGuard: guard, logger: logger}
}

//CurrencyConversionOffer godoc
//@Summary Accepts currency conversion offer
//@Description convert currencies by given conversion offer
//@Param conversionOffer body CurrencyConversionOfferRequest true "body params"
//@Param accessToken header string true "header params"
// @Tags Conversion
//@Accept  json
//@Produce  json
//@Success 200 {string} string SuccessfulCurrencyConversionMessage
//@Failure 401 {string} string "Unauthorized"
//@Failure 400 {object} httperror.HttpError
//@Failure 404
//@Failure 500 {object} httperror.HttpError
//@Router /conversion/offer [post]
func (h *HttpHandler) CurrencyConversionOffer(ctx *fiber.Ctx) error {
	userID, exists := ctx.Locals(common.USER_ID_CTX_KEY).(int)
	h.logger.Infof("Currency Conversion Offer Request Arrived. User ID: %d", userID)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httperror.NewBadRequestError(common.UserNotInContext))
	}
	var offerRequest CurrencyConversionOfferRequest
	if err := ctx.BodyParser(&offerRequest); err != nil {
		h.logger.Debug("Can not parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(httperror.NewBadRequestError(err.Error()))
	}
	success, err := h.conversionService.
		ConvertCurrencies(ctx.Context(), userID, offerRequest)
	if err != nil || !success {
		h.logger.Debug(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperror.NewInternalServerError(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).SendString(SuccessfulCurrencyConversionMessage)
}
func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/conversion")
	appGroup.Post("/offer", h.authGuard.ProtectWithJWT(h.CurrencyConversionOffer))
}
