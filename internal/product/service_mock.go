// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -typed -source=service.go -destination=service_mock.go -package=product
//

// Package product is a generated GoMock package.
package product

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockproductDBAccessor is a mock of productDBAccessor interface.
type MockproductDBAccessor struct {
	ctrl     *gomock.Controller
	recorder *MockproductDBAccessorMockRecorder
}

// MockproductDBAccessorMockRecorder is the mock recorder for MockproductDBAccessor.
type MockproductDBAccessorMockRecorder struct {
	mock *MockproductDBAccessor
}

// NewMockproductDBAccessor creates a new mock instance.
func NewMockproductDBAccessor(ctrl *gomock.Controller) *MockproductDBAccessor {
	mock := &MockproductDBAccessor{ctrl: ctrl}
	mock.recorder = &MockproductDBAccessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockproductDBAccessor) EXPECT() *MockproductDBAccessorMockRecorder {
	return m.recorder
}

// GetProductsByVendor mocks base method.
func (m *MockproductDBAccessor) GetProductsByVendor(ctx context.Context, vendorID string, spec GetProductsByVendorSpec) (*AccessorGetProductsByVendorPaginationData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductsByVendor", ctx, vendorID, spec)
	ret0, _ := ret[0].(*AccessorGetProductsByVendorPaginationData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductsByVendor indicates an expected call of GetProductsByVendor.
func (mr *MockproductDBAccessorMockRecorder) GetProductsByVendor(ctx, vendorID, spec any) *MockproductDBAccessorGetProductsByVendorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductsByVendor", reflect.TypeOf((*MockproductDBAccessor)(nil).GetProductsByVendor), ctx, vendorID, spec)
	return &MockproductDBAccessorGetProductsByVendorCall{Call: call}
}

// MockproductDBAccessorGetProductsByVendorCall wrap *gomock.Call
type MockproductDBAccessorGetProductsByVendorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductDBAccessorGetProductsByVendorCall) Return(arg0 *AccessorGetProductsByVendorPaginationData, arg1 error) *MockproductDBAccessorGetProductsByVendorCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductDBAccessorGetProductsByVendorCall) Do(f func(context.Context, string, GetProductsByVendorSpec) (*AccessorGetProductsByVendorPaginationData, error)) *MockproductDBAccessorGetProductsByVendorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductDBAccessorGetProductsByVendorCall) DoAndReturn(f func(context.Context, string, GetProductsByVendorSpec) (*AccessorGetProductsByVendorPaginationData, error)) *MockproductDBAccessorGetProductsByVendorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdatePrice mocks base method.
func (m *MockproductDBAccessor) UpdatePrice(ctx context.Context, price Price) (Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePrice", ctx, price)
	ret0, _ := ret[0].(Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePrice indicates an expected call of UpdatePrice.
func (mr *MockproductDBAccessorMockRecorder) UpdatePrice(ctx, price any) *MockproductDBAccessorUpdatePriceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePrice", reflect.TypeOf((*MockproductDBAccessor)(nil).UpdatePrice), ctx, price)
	return &MockproductDBAccessorUpdatePriceCall{Call: call}
}

// MockproductDBAccessorUpdatePriceCall wrap *gomock.Call
type MockproductDBAccessorUpdatePriceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductDBAccessorUpdatePriceCall) Return(arg0 Price, arg1 error) *MockproductDBAccessorUpdatePriceCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductDBAccessorUpdatePriceCall) Do(f func(context.Context, Price) (Price, error)) *MockproductDBAccessorUpdatePriceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductDBAccessorUpdatePriceCall) DoAndReturn(f func(context.Context, Price) (Price, error)) *MockproductDBAccessorUpdatePriceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateProduct mocks base method.
func (m *MockproductDBAccessor) UpdateProduct(ctx context.Context, payload Product) (Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", ctx, payload)
	ret0, _ := ret[0].(Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProduct indicates an expected call of UpdateProduct.
func (mr *MockproductDBAccessorMockRecorder) UpdateProduct(ctx, payload any) *MockproductDBAccessorUpdateProductCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockproductDBAccessor)(nil).UpdateProduct), ctx, payload)
	return &MockproductDBAccessorUpdateProductCall{Call: call}
}

// MockproductDBAccessorUpdateProductCall wrap *gomock.Call
type MockproductDBAccessorUpdateProductCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductDBAccessorUpdateProductCall) Return(arg0 Product, arg1 error) *MockproductDBAccessorUpdateProductCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductDBAccessorUpdateProductCall) Do(f func(context.Context, Product) (Product, error)) *MockproductDBAccessorUpdateProductCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductDBAccessorUpdateProductCall) DoAndReturn(f func(context.Context, Product) (Product, error)) *MockproductDBAccessorUpdateProductCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
