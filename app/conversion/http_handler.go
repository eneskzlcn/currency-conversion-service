package conversion

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/app/common"
	"github.com/eneskzlcn/currency-conversion-service/app/common/httperror"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Service interface {
	ConvertCurrencies(ctx context.Context, userID int, request CurrencyConversionOfferRequest) (bool, error)
}
type AuthGuard interface {
	ProtectWithJWT(handler fiber.Handler) fiber.Handler
}

type httpHandler struct {
	service   Service
	authGuard AuthGuard
	logger    *zap.SugaredLogger
}

func NewHttpHandler(service Service, guard AuthGuard, logger *zap.SugaredLogger) *httpHandler {
	return &httpHandler{service: service, authGuard: guard, logger: logger}
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
func (h *httpHandler) CurrencyConversionOffer(ctx *fiber.Ctx) error {
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
	success, err := h.service.
		ConvertCurrencies(ctx.Context(), userID, offerRequest)
	if err != nil || !success {
		h.logger.Debug(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperror.NewInternalServerError(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).SendString(SuccessfulCurrencyConversionMessage)
}
func (h *httpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/conversion")
	appGroup.Post("/offer", h.authGuard.ProtectWithJWT(h.CurrencyConversionOffer))
}
