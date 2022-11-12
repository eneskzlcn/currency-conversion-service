// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/eneskzlcn/currency-conversion-service/app/wallet (interfaces: Service)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	message "github.com/eneskzlcn/currency-conversion-service/app/message"
	wallet "github.com/eneskzlcn/currency-conversion-service/app/wallet"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetUserWalletAccounts mocks base method.
func (m *MockService) GetUserWalletAccounts(arg0 context.Context, arg1 int) (wallet.UserWalletAccountsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWalletAccounts", arg0, arg1)
	ret0, _ := ret[0].(wallet.UserWalletAccountsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWalletAccounts indicates an expected call of GetUserWalletAccounts.
func (mr *MockServiceMockRecorder) GetUserWalletAccounts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWalletAccounts", reflect.TypeOf((*MockService)(nil).GetUserWalletAccounts), arg0, arg1)
}

// TransferBalancesBetweenUserWallets mocks base method.
func (m *MockService) TransferBalancesBetweenUserWallets(arg0 context.Context, arg1 message.CurrencyConvertedMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferBalancesBetweenUserWallets", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransferBalancesBetweenUserWallets indicates an expected call of TransferBalancesBetweenUserWallets.
func (mr *MockServiceMockRecorder) TransferBalancesBetweenUserWallets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferBalancesBetweenUserWallets", reflect.TypeOf((*MockService)(nil).TransferBalancesBetweenUserWallets), arg0, arg1)
}
