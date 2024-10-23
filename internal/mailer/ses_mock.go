// Code generated by MockGen. DO NOT EDIT.
// Source: ses.go
//
// Generated by this command:
//
//	mockgen -typed -source=ses.go -destination=ses_mock.go -package=mailer
//

// Package mailer is a generated GoMock package.
package mailer

import (
	context "context"
	reflect "reflect"

	ses "github.com/aws/aws-sdk-go-v2/service/ses"
	gomock "go.uber.org/mock/gomock"
)

// MockSESSendEmailAPI is a mock of SESSendEmailAPI interface.
type MockSESSendEmailAPI struct {
	ctrl     *gomock.Controller
	recorder *MockSESSendEmailAPIMockRecorder
}

// MockSESSendEmailAPIMockRecorder is the mock recorder for MockSESSendEmailAPI.
type MockSESSendEmailAPIMockRecorder struct {
	mock *MockSESSendEmailAPI
}

// NewMockSESSendEmailAPI creates a new mock instance.
func NewMockSESSendEmailAPI(ctrl *gomock.Controller) *MockSESSendEmailAPI {
	mock := &MockSESSendEmailAPI{ctrl: ctrl}
	mock.recorder = &MockSESSendEmailAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSESSendEmailAPI) EXPECT() *MockSESSendEmailAPIMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockSESSendEmailAPI) SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendEmail", varargs...)
	ret0, _ := ret[0].(*ses.SendEmailOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockSESSendEmailAPIMockRecorder) SendEmail(ctx, params any, optFns ...any) *MockSESSendEmailAPISendEmailCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockSESSendEmailAPI)(nil).SendEmail), varargs...)
	return &MockSESSendEmailAPISendEmailCall{Call: call}
}

// MockSESSendEmailAPISendEmailCall wrap *gomock.Call
type MockSESSendEmailAPISendEmailCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSESSendEmailAPISendEmailCall) Return(arg0 *ses.SendEmailOutput, arg1 error) *MockSESSendEmailAPISendEmailCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSESSendEmailAPISendEmailCall) Do(f func(context.Context, *ses.SendEmailInput, ...func(*ses.Options)) (*ses.SendEmailOutput, error)) *MockSESSendEmailAPISendEmailCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSESSendEmailAPISendEmailCall) DoAndReturn(f func(context.Context, *ses.SendEmailInput, ...func(*ses.Options)) (*ses.SendEmailOutput, error)) *MockSESSendEmailAPISendEmailCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}