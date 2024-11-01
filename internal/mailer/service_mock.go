// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -typed -source=service.go -destination=service_mock.go -package=mailer
//

// Package mailer is a generated GoMock package.
package mailer

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockemailStatusDBAccessor is a mock of emailStatusDBAccessor interface.
type MockemailStatusDBAccessor struct {
	ctrl     *gomock.Controller
	recorder *MockemailStatusDBAccessorMockRecorder
}

// MockemailStatusDBAccessorMockRecorder is the mock recorder for MockemailStatusDBAccessor.
type MockemailStatusDBAccessorMockRecorder struct {
	mock *MockemailStatusDBAccessor
}

// NewMockemailStatusDBAccessor creates a new mock instance.
func NewMockemailStatusDBAccessor(ctrl *gomock.Controller) *MockemailStatusDBAccessor {
	mock := &MockemailStatusDBAccessor{ctrl: ctrl}
	mock.recorder = &MockemailStatusDBAccessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockemailStatusDBAccessor) EXPECT() *MockemailStatusDBAccessorMockRecorder {
	return m.recorder
}

// WriteEmailStatus mocks base method.
func (m *MockemailStatusDBAccessor) WriteEmailStatus(ctx context.Context, payload EmailStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteEmailStatus", ctx, payload)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteEmailStatus indicates an expected call of WriteEmailStatus.
func (mr *MockemailStatusDBAccessorMockRecorder) WriteEmailStatus(ctx, payload any) *MockemailStatusDBAccessorWriteEmailStatusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteEmailStatus", reflect.TypeOf((*MockemailStatusDBAccessor)(nil).WriteEmailStatus), ctx, payload)
	return &MockemailStatusDBAccessorWriteEmailStatusCall{Call: call}
}

// MockemailStatusDBAccessorWriteEmailStatusCall wrap *gomock.Call
type MockemailStatusDBAccessorWriteEmailStatusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockemailStatusDBAccessorWriteEmailStatusCall) Return(arg0 error) *MockemailStatusDBAccessorWriteEmailStatusCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockemailStatusDBAccessorWriteEmailStatusCall) Do(f func(context.Context, EmailStatus) error) *MockemailStatusDBAccessorWriteEmailStatusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockemailStatusDBAccessorWriteEmailStatusCall) DoAndReturn(f func(context.Context, EmailStatus) error) *MockemailStatusDBAccessorWriteEmailStatusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}