package testutil

import (
	"bytes"
	"encoding/json"
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

func AssertBodyEqual(t *testing.T, responseBody io.Reader, expectedValue interface{}) {
	var actualBody interface{}
	_ = json.NewDecoder(responseBody).Decode(&actualBody)

	expectedBodyAsJSON, _ := json.Marshal(expectedValue)

	var expectedBody interface{}
	_ = json.Unmarshal(expectedBodyAsJSON, &expectedBody)
	assert.Equal(t, expectedBody, actualBody)
}
