// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/eneskzlcn/currency-conversion-service/app/exchange (interfaces: ExchangeRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/eneskzlcn/currency-conversion-service/app/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockExchangeRepository is a mock of ExchangeRepository interface.
type MockExchangeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockExchangeRepositoryMockRecorder
}

// MockExchangeRepositoryMockRecorder is the mock recorder for MockExchangeRepository.
type MockExchangeRepositoryMockRecorder struct {
	mock *MockExchangeRepository
}

// NewMockExchangeRepository creates a new mock instance.
func NewMockExchangeRepository(ctrl *gomock.Controller) *MockExchangeRepository {
	mock := &MockExchangeRepository{ctrl: ctrl}
	mock.recorder = &MockExchangeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExchangeRepository) EXPECT() *MockExchangeRepositoryMockRecorder {
	return m.recorder
}

// GetExchangeValuesForGivenCurrencies mocks base method.
func (m *MockExchangeRepository) GetExchangeValuesForGivenCurrencies(arg0 context.Context, arg1, arg2 string) (entity.Exchange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExchangeValuesForGivenCurrencies", arg0, arg1, arg2)
	ret0, _ := ret[0].(entity.Exchange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExchangeValuesForGivenCurrencies indicates an expected call of GetExchangeValuesForGivenCurrencies.
func (mr *MockExchangeRepositoryMockRecorder) GetExchangeValuesForGivenCurrencies(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExchangeValuesForGivenCurrencies", reflect.TypeOf((*MockExchangeRepository)(nil).GetExchangeValuesForGivenCurrencies), arg0, arg1, arg2)
}

// IsCurrencyExists mocks base method.
func (m *MockExchangeRepository) IsCurrencyExists(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCurrencyExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsCurrencyExists indicates an expected call of IsCurrencyExists.
func (mr *MockExchangeRepositoryMockRecorder) IsCurrencyExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCurrencyExists", reflect.TypeOf((*MockExchangeRepository)(nil).IsCurrencyExists), arg0, arg1)
}
