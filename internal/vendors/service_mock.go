// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
//

// Package vendors is a generated GoMock package.
package vendors

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockvendorDBAccessor is a mock of vendorDBAccessor interface.
type MockvendorDBAccessor struct {
	ctrl     *gomock.Controller
	recorder *MockvendorDBAccessorMockRecorder
}

// MockvendorDBAccessorMockRecorder is the mock recorder for MockvendorDBAccessor.
type MockvendorDBAccessorMockRecorder struct {
	mock *MockvendorDBAccessor
}

// NewMockvendorDBAccessor creates a new mock instance.
func NewMockvendorDBAccessor(ctrl *gomock.Controller) *MockvendorDBAccessor {
	mock := &MockvendorDBAccessor{ctrl: ctrl}
	mock.recorder = &MockvendorDBAccessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockvendorDBAccessor) EXPECT() *MockvendorDBAccessorMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockvendorDBAccessor) GetAll(ctx context.Context, spec GetAllVendorSpec) (*AccessorGetAllPaginationData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, spec)
	ret0, _ := ret[0].(*AccessorGetAllPaginationData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockvendorDBAccessorMockRecorder) GetAll(ctx, spec any) *MockvendorDBAccessorGetAllCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockvendorDBAccessor)(nil).GetAll), ctx, spec)
	return &MockvendorDBAccessorGetAllCall{Call: call}
}

// MockvendorDBAccessorGetAllCall wrap *gomock.Call
type MockvendorDBAccessorGetAllCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockvendorDBAccessorGetAllCall) Return(arg0 *AccessorGetAllPaginationData, arg1 error) *MockvendorDBAccessorGetAllCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockvendorDBAccessorGetAllCall) Do(f func(context.Context, GetAllVendorSpec) (*AccessorGetAllPaginationData, error)) *MockvendorDBAccessorGetAllCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockvendorDBAccessorGetAllCall) DoAndReturn(f func(context.Context, GetAllVendorSpec) (*AccessorGetAllPaginationData, error)) *MockvendorDBAccessorGetAllCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetSomeStuff mocks base method.
func (m *MockvendorDBAccessor) GetSomeStuff(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSomeStuff", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSomeStuff indicates an expected call of GetSomeStuff.
func (mr *MockvendorDBAccessorMockRecorder) GetSomeStuff(ctx any) *MockvendorDBAccessorGetSomeStuffCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSomeStuff", reflect.TypeOf((*MockvendorDBAccessor)(nil).GetSomeStuff), ctx)
	return &MockvendorDBAccessorGetSomeStuffCall{Call: call}
}

// MockvendorDBAccessorGetSomeStuffCall wrap *gomock.Call
type MockvendorDBAccessorGetSomeStuffCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockvendorDBAccessorGetSomeStuffCall) Return(arg0 []string, arg1 error) *MockvendorDBAccessorGetSomeStuffCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockvendorDBAccessorGetSomeStuffCall) Do(f func(context.Context) ([]string, error)) *MockvendorDBAccessorGetSomeStuffCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockvendorDBAccessorGetSomeStuffCall) DoAndReturn(f func(context.Context) ([]string, error)) *MockvendorDBAccessorGetSomeStuffCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
