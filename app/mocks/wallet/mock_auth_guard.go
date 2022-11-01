// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/eneskzlcn/currency-conversion-service/app/wallet (interfaces: AuthGuard)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	fiber "github.com/gofiber/fiber/v2"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthGuard is a mock of AuthGuard interface.
type MockAuthGuard struct {
	ctrl     *gomock.Controller
	recorder *MockAuthGuardMockRecorder
}

// MockAuthGuardMockRecorder is the mock recorder for MockAuthGuard.
type MockAuthGuardMockRecorder struct {
	mock *MockAuthGuard
}

// NewMockAuthGuard creates a new mock instance.
func NewMockAuthGuard(ctrl *gomock.Controller) *MockAuthGuard {
	mock := &MockAuthGuard{ctrl: ctrl}
	mock.recorder = &MockAuthGuardMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthGuard) EXPECT() *MockAuthGuardMockRecorder {
	return m.recorder
}

// ProtectWithJWT mocks base method.
func (m *MockAuthGuard) ProtectWithJWT(arg0 func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProtectWithJWT", arg0)
	ret0, _ := ret[0].(func(*fiber.Ctx) error)
	return ret0
}

// ProtectWithJWT indicates an expected call of ProtectWithJWT.
func (mr *MockAuthGuardMockRecorder) ProtectWithJWT(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProtectWithJWT", reflect.TypeOf((*MockAuthGuard)(nil).ProtectWithJWT), arg0)
}
