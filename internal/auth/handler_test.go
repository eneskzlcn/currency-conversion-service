package auth_test

import (
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/auth"
	"github.com/eneskzlcn/currency-conversion-service/internal/common/testutil"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestNewHandler(t *testing.T) {
	t.Run("given auth service then it should return new Handler when NewHandler called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := auth.NewHandler(mockAuthService)
		assert.NotNil(t, handler)
	})
	t.Run("given nil auth service then it should return nil when NewHandler called", func(t *testing.T) {
		handler := auth.NewHandler(nil)
		assert.Nil(t, handler)
	})
}
func TestHandler_Login(t *testing.T) {
	handler, mockAuthService := createHandlerAndMockAuthService(t)
	app := fiber.New()
	app.Post("/login", handler.Login)
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
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := auth.NewHandler(mockAuthService)
	handler.RegisterRoutes(app)
	resp, err := app.Test(httptest.NewRequest(fiber.MethodPost, "/login", nil))
	assert.Nil(t, err)
	assert.NotEqual(t, fiber.StatusNotFound, resp.StatusCode)
}
func createHandlerAndMockAuthService(t *testing.T) (*auth.Handler, *mocks.MockAuthService) {
	ctrl := gomock.NewController(t)
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := auth.NewHandler(mockAuthService)
	return handler, mockAuthService
}
