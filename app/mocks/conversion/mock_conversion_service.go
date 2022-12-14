// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/eneskzlcn/currency-conversion-service/app/conversion (interfaces: ConversionService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	conversion "github.com/eneskzlcn/currency-conversion-service/app/conversion"
	gomock "github.com/golang/mock/gomock"
)

// MockConversionService is a mock of ConversionService interface.
type MockConversionService struct {
	ctrl     *gomock.Controller
	recorder *MockConversionServiceMockRecorder
}

// MockConversionServiceMockRecorder is the mock recorder for MockConversionService.
type MockConversionServiceMockRecorder struct {
	mock *MockConversionService
}

// NewMockConversionService creates a new mock instance.
func NewMockConversionService(ctrl *gomock.Controller) *MockConversionService {
	mock := &MockConversionService{ctrl: ctrl}
	mock.recorder = &MockConversionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConversionService) EXPECT() *MockConversionServiceMockRecorder {
	return m.recorder
}

// ConvertCurrencies mocks base method.
func (m *MockConversionService) ConvertCurrencies(arg0 context.Context, arg1 int, arg2 conversion.CurrencyConversionOfferRequest) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertCurrencies", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConvertCurrencies indicates an expected call of ConvertCurrencies.
func (mr *MockConversionServiceMockRecorder) ConvertCurrencies(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertCurrencies", reflect.TypeOf((*MockConversionService)(nil).ConvertCurrencies), arg0, arg1, arg2)
}
