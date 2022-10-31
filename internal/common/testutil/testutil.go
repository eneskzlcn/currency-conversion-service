package testutil

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MakeTestRequestWithBody(method string, route string, body interface{}) *http.Request {
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(method, route, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	return req
}
func MakeTestRequestWithoutBody(method string, route string) *http.Request {
	req := httptest.NewRequest(method, route, nil)
	return req
}
func AssertBodyEqual(t *testing.T, responseBody io.Reader, expectedValue interface{}) {
	var actualBody interface{}
	_ = json.NewDecoder(responseBody).Decode(&actualBody)

	expectedBodyAsJSON, _ := json.Marshal(expectedValue)

	var expectedBody interface{}
	_ = json.Unmarshal(expectedBodyAsJSON, &expectedBody)
	assert.Equal(t, expectedBody, actualBody)
}
func AssertRouteRegistered(t *testing.T, app *fiber.App, method, route string) {
	resp, err := app.Test(httptest.NewRequest(method, route, nil))
	assert.Nil(t, err)
	assert.NotEqual(t, fiber.StatusNotFound, resp.StatusCode)
}
func MakeTestRequestWithoutBodyToProtectedEndpoint(method string, route string, token string) *http.Request {
	req := httptest.NewRequest(method, route, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	return req
}
