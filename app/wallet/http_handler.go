package wallet

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/app/common"
	"github.com/eneskzlcn/currency-conversion-service/app/common/httperror"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WalletService interface {
	GetUserWalletAccounts(ctx context.Context, userID int) (UserWalletAccountsResponse, error)
}
type AuthGuard interface {
	ProtectWithJWT(handler fiber.Handler) fiber.Handler
}
type HttpHandler struct {
	walletService WalletService
	authGuard     AuthGuard
	logger        *zap.SugaredLogger
}

func NewHttpHandler(service WalletService, guard AuthGuard, logger *zap.SugaredLogger) *HttpHandler {
	return &HttpHandler{
		walletService: service,
		authGuard:     guard,
		logger:        logger,
	}
}

//GetUserWalletAccounts godoc
//@Summary Shows user wallet accounts
//@Description shows user wallet accounts for all existing currency
// @Tags Wallet
//@Accept  json
//@Produce  json
//@Param accessToken header string true "header params"
//@Success 200 {object} UserWalletAccountsResponse
//@Failure 400 {object} httperror.HttpError
//@Failure 404
//@Failure 401 {string} string "Unauthorized"
//@Failure 500 {object} httperror.HttpError
//@Router /wallets [get]
func (h *HttpHandler) GetUserWalletAccounts(ctx *fiber.Ctx) error {
	userID, exists := ctx.Locals(common.USER_ID_CTX_KEY).(int)
	h.logger.Infof("List wallet accounts requesta arrived. User ID:%d", userID)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httperror.NewBadRequestError(common.UserNotInContext))
	}
	userWalletAccounts, err := h.walletService.GetUserWalletAccounts(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperror.NewInternalServerError(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(userWalletAccounts)
}
func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/wallets", h.authGuard.ProtectWithJWT(h.GetUserWalletAccounts))
}
