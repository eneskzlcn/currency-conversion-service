package auth

import (
	"github.com/eneskzlcn/currency-conversion-service/app/common/httperror"
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

// Login godoc
// @Summary Authenticate user
// @Description authenticates given user by giving an access token.
// @Param loginCredentials body LoginRequest true "body params"
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Success 200 {object} TokenResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /auth/login [post]
func (h *Handler) Login(ctx *fiber.Ctx) error {
	var request LoginRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httperror.NewBadRequestError(err.Error()))
	}
	tokenResponse, err := h.authService.Tokenize(ctx.Context(), request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	h.logger.Debugf("User logged in. Username: %s", request.Username)

	return ctx.Status(fiber.StatusCreated).JSON(tokenResponse)
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/auth")
	appGroup.Post("/login", h.Login)
}
