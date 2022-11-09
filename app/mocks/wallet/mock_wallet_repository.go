// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/eneskzlcn/currency-conversion-service/app/wallet (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/eneskzlcn/currency-conversion-service/app/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AdjustUserBalanceOnGivenCurrency mocks base method.
func (m *MockRepository) AdjustUserBalanceOnGivenCurrency(arg0 context.Context, arg1 int, arg2 string, arg3 float32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AdjustUserBalanceOnGivenCurrency", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AdjustUserBalanceOnGivenCurrency indicates an expected call of AdjustUserBalanceOnGivenCurrency.
func (mr *MockRepositoryMockRecorder) AdjustUserBalanceOnGivenCurrency(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdjustUserBalanceOnGivenCurrency", reflect.TypeOf((*MockRepository)(nil).AdjustUserBalanceOnGivenCurrency), arg0, arg1, arg2, arg3)
}

// GetUserBalanceOnGivenCurrency mocks base method.
func (m *MockRepository) GetUserBalanceOnGivenCurrency(arg0 context.Context, arg1 int, arg2 string) (float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalanceOnGivenCurrency", arg0, arg1, arg2)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalanceOnGivenCurrency indicates an expected call of GetUserBalanceOnGivenCurrency.
func (mr *MockRepositoryMockRecorder) GetUserBalanceOnGivenCurrency(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalanceOnGivenCurrency", reflect.TypeOf((*MockRepository)(nil).GetUserBalanceOnGivenCurrency), arg0, arg1, arg2)
}

// GetUserWalletAccounts mocks base method.
func (m *MockRepository) GetUserWalletAccounts(arg0 context.Context, arg1 int) ([]entity.UserWallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWalletAccounts", arg0, arg1)
	ret0, _ := ret[0].([]entity.UserWallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWalletAccounts indicates an expected call of GetUserWalletAccounts.
func (mr *MockRepositoryMockRecorder) GetUserWalletAccounts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWalletAccounts", reflect.TypeOf((*MockRepository)(nil).GetUserWalletAccounts), arg0, arg1)
}

// IsUserWithUserIDExists mocks base method.
func (m *MockRepository) IsUserWithUserIDExists(arg0 context.Context, arg1 int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserWithUserIDExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUserWithUserIDExists indicates an expected call of IsUserWithUserIDExists.
func (mr *MockRepositoryMockRecorder) IsUserWithUserIDExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserWithUserIDExists", reflect.TypeOf((*MockRepository)(nil).IsUserWithUserIDExists), arg0, arg1)
}
