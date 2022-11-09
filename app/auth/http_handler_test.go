package auth_test

import (
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/auth"
	"github.com/eneskzlcn/currency-conversion-service/app/common/testutil"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http/httptest"
	"testing"
)

type HttpHandler interface {
	RegisterRoutes(app *fiber.App)
	Login(ctx *fiber.Ctx) error
}

func TestHandler_Login(t *testing.T) {
	httpHandler, mockAuthService := createHandlerAndMockAuthService(t)
	app := fiber.New()
	app.Post("/login", httpHandler.Login)
	t.Run("given not valid login request then it should return status bad request", func(t *testing.T) {
		loginRequestData := "asdf"
		request := testutil.MakeTestRequestWithBody(fiber.MethodPost, "/login", loginRequestData)
		resp, err := app.Test(request)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
	t.Run("given valid login request and error occurred on service then it should return internal server error", func(t *testing.T) {
		loginRequestData := auth.LoginRequest{
			Username: "test",
			Password: "test",
		}
		mockAuthService.EXPECT().Tokenize(gomock.Any(), loginRequestData).
			Return(auth.TokenResponse{}, errors.New("error occurred"))
		request := testutil.MakeTestRequestWithBody(fiber.MethodPost, "/login", loginRequestData)
		resp, err := app.Test(request)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("given valid request and token successfully generated then it should return token", func(t *testing.T) {
		loginRequestData := auth.LoginRequest{
			Username: "test",
			Password: "test",
		}
		expectedResponse := auth.TokenResponse{AccessToken: "someaccesstoken"}
		mockAuthService.EXPECT().Tokenize(gomock.Any(), loginRequestData).
			Return(expectedResponse, nil)
		request := testutil.MakeTestRequestWithBody(fiber.MethodPost, "/login", loginRequestData)
		resp, err := app.Test(request)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		testutil.AssertBodyEqual(t, resp.Body, expectedResponse)
	})
}
func TestRegisterRoutesSuccessfullyRegistersTheEndpointsToTheApp(t *testing.T) {
	app := fiber.New()
	ctrl := gomock.NewController(t)
	mockAuthService := mocks.NewMockService(ctrl)
	httpHandler := auth.NewHttpHandler(mockAuthService, zap.S())
	httpHandler.RegisterRoutes(app)
	resp, err := app.Test(httptest.NewRequest(fiber.MethodPost, "/auth/login", nil))
	assert.Nil(t, err)
	assert.NotEqual(t, fiber.StatusNotFound, resp.StatusCode)
}
func createHandlerAndMockAuthService(t *testing.T) (HttpHandler, *mocks.MockService) {
	ctrl := gomock.NewController(t)
	mockAuthService := mocks.NewMockService(ctrl)
	httpHandler := auth.NewHttpHandler(mockAuthService, zap.S())
	return httpHandler, mockAuthService
}
