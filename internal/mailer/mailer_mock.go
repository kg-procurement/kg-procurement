// Code generated by MockGen. DO NOT EDIT.
// Source: mailer.go
//
// Generated by this command:
//
//	mockgen -typed -source=mailer.go -destination=mailer_mock.go -package=mailer
//

// Package mailer is a generated GoMock package.
package mailer

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockEmailProvider is a mock of EmailProvider interface.
type MockEmailProvider struct {
	ctrl     *gomock.Controller
	recorder *MockEmailProviderMockRecorder
}

// MockEmailProviderMockRecorder is the mock recorder for MockEmailProvider.
type MockEmailProviderMockRecorder struct {
	mock *MockEmailProvider
}

// NewMockEmailProvider creates a new mock instance.
func NewMockEmailProvider(ctrl *gomock.Controller) *MockEmailProvider {
	mock := &MockEmailProvider{ctrl: ctrl}
	mock.recorder = &MockEmailProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailProvider) EXPECT() *MockEmailProviderMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockEmailProvider) SendEmail(email Email) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", email)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockEmailProviderMockRecorder) SendEmail(email any) *MockEmailProviderSendEmailCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockEmailProvider)(nil).SendEmail), email)
	return &MockEmailProviderSendEmailCall{Call: call}
}

// MockEmailProviderSendEmailCall wrap *gomock.Call
type MockEmailProviderSendEmailCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockEmailProviderSendEmailCall) Return(arg0 error) *MockEmailProviderSendEmailCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockEmailProviderSendEmailCall) Do(f func(Email) error) *MockEmailProviderSendEmailCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockEmailProviderSendEmailCall) DoAndReturn(f func(Email) error) *MockEmailProviderSendEmailCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
