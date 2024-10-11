// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -typed -source=service.go -destination=service_mock.go -package=account
//

// Package account is a generated GoMock package.
package account

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockaccountDBAccessor is a mock of accountDBAccessor interface.
type MockaccountDBAccessor struct {
	ctrl     *gomock.Controller
	recorder *MockaccountDBAccessorMockRecorder
}

// MockaccountDBAccessorMockRecorder is the mock recorder for MockaccountDBAccessor.
type MockaccountDBAccessorMockRecorder struct {
	mock *MockaccountDBAccessor
}

// NewMockaccountDBAccessor creates a new mock instance.
func NewMockaccountDBAccessor(ctrl *gomock.Controller) *MockaccountDBAccessor {
	mock := &MockaccountDBAccessor{ctrl: ctrl}
	mock.recorder = &MockaccountDBAccessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockaccountDBAccessor) EXPECT() *MockaccountDBAccessorMockRecorder {
	return m.recorder
}

// RegisterAccount mocks base method.
func (m *MockaccountDBAccessor) RegisterAccount(ctx context.Context, account Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterAccount", ctx, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterAccount indicates an expected call of RegisterAccount.
func (mr *MockaccountDBAccessorMockRecorder) RegisterAccount(ctx, account any) *MockaccountDBAccessorRegisterAccountCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterAccount", reflect.TypeOf((*MockaccountDBAccessor)(nil).RegisterAccount), ctx, account)
	return &MockaccountDBAccessorRegisterAccountCall{Call: call}
}

// MockaccountDBAccessorRegisterAccountCall wrap *gomock.Call
type MockaccountDBAccessorRegisterAccountCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockaccountDBAccessorRegisterAccountCall) Return(arg0 error) *MockaccountDBAccessorRegisterAccountCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockaccountDBAccessorRegisterAccountCall) Do(f func(context.Context, Account) error) *MockaccountDBAccessorRegisterAccountCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockaccountDBAccessorRegisterAccountCall) DoAndReturn(f func(context.Context, Account) error) *MockaccountDBAccessorRegisterAccountCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}