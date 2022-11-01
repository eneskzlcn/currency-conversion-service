// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/eneskzlcn/currency-conversion-service/app/wallet (interfaces: WalletService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	wallet "github.com/eneskzlcn/currency-conversion-service/app/wallet"
	gomock "github.com/golang/mock/gomock"
)

// MockWalletService is a mock of WalletService interface.
type MockWalletService struct {
	ctrl     *gomock.Controller
	recorder *MockWalletServiceMockRecorder
}

// MockWalletServiceMockRecorder is the mock recorder for MockWalletService.
type MockWalletServiceMockRecorder struct {
	mock *MockWalletService
}

// NewMockWalletService creates a new mock instance.
func NewMockWalletService(ctrl *gomock.Controller) *MockWalletService {
	mock := &MockWalletService{ctrl: ctrl}
	mock.recorder = &MockWalletServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletService) EXPECT() *MockWalletServiceMockRecorder {
	return m.recorder
}

// GetUserWalletAccounts mocks base method.
func (m *MockWalletService) GetUserWalletAccounts(arg0 context.Context, arg1 int) (wallet.UserWalletAccountsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWalletAccounts", arg0, arg1)
	ret0, _ := ret[0].(wallet.UserWalletAccountsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWalletAccounts indicates an expected call of GetUserWalletAccounts.
func (mr *MockWalletServiceMockRecorder) GetUserWalletAccounts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWalletAccounts", reflect.TypeOf((*MockWalletService)(nil).GetUserWalletAccounts), arg0, arg1)
}
