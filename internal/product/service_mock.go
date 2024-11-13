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

// GetAllProductVendors mocks base method.
func (m *MockproductDBAccessor) GetAllProductVendors(ctx context.Context, spec GetProductVendorsSpec) (*AccessorGetProductVendorsPaginationData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllProductVendors", ctx, spec)
	ret0, _ := ret[0].(*AccessorGetProductVendorsPaginationData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllProductVendors indicates an expected call of GetAllProductVendors.
func (mr *MockproductDBAccessorMockRecorder) GetAllProductVendors(ctx, spec any) *MockproductDBAccessorGetAllProductVendorsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllProductVendors", reflect.TypeOf((*MockproductDBAccessor)(nil).GetAllProductVendors), ctx, spec)
	return &MockproductDBAccessorGetAllProductVendorsCall{Call: call}
}

// MockproductDBAccessorGetAllProductVendorsCall wrap *gomock.Call
type MockproductDBAccessorGetAllProductVendorsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductDBAccessorGetAllProductVendorsCall) Return(arg0 *AccessorGetProductVendorsPaginationData, arg1 error) *MockproductDBAccessorGetAllProductVendorsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductDBAccessorGetAllProductVendorsCall) Do(f func(context.Context, GetProductVendorsSpec) (*AccessorGetProductVendorsPaginationData, error)) *MockproductDBAccessorGetAllProductVendorsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductDBAccessorGetAllProductVendorsCall) DoAndReturn(f func(context.Context, GetProductVendorsSpec) (*AccessorGetProductVendorsPaginationData, error)) *MockproductDBAccessorGetAllProductVendorsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetProductVendorsByVendor mocks base method.
func (m *MockproductDBAccessor) GetProductVendorsByVendor(ctx context.Context, vendorID string, spec GetProductVendorByVendorSpec) (*AccessorGetProductVendorsPaginationData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductVendorsByVendor", ctx, vendorID, spec)
	ret0, _ := ret[0].(*AccessorGetProductVendorsPaginationData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductVendorsByVendor indicates an expected call of GetProductVendorsByVendor.
func (mr *MockproductDBAccessorMockRecorder) GetProductVendorsByVendor(ctx, vendorID, spec any) *MockproductDBAccessorGetProductVendorsByVendorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductVendorsByVendor", reflect.TypeOf((*MockproductDBAccessor)(nil).GetProductVendorsByVendor), ctx, vendorID, spec)
	return &MockproductDBAccessorGetProductVendorsByVendorCall{Call: call}
}

// MockproductDBAccessorGetProductVendorsByVendorCall wrap *gomock.Call
type MockproductDBAccessorGetProductVendorsByVendorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductDBAccessorGetProductVendorsByVendorCall) Return(arg0 *AccessorGetProductVendorsPaginationData, arg1 error) *MockproductDBAccessorGetProductVendorsByVendorCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductDBAccessorGetProductVendorsByVendorCall) Do(f func(context.Context, string, GetProductVendorByVendorSpec) (*AccessorGetProductVendorsPaginationData, error)) *MockproductDBAccessorGetProductVendorsByVendorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductDBAccessorGetProductVendorsByVendorCall) DoAndReturn(f func(context.Context, string, GetProductVendorByVendorSpec) (*AccessorGetProductVendorsPaginationData, error)) *MockproductDBAccessorGetProductVendorsByVendorCall {
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

// getPriceByPVID mocks base method.
func (m *MockproductDBAccessor) getPriceByPVID(ctx context.Context, pvID string) (*Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getPriceByPVID", ctx, pvID)
	ret0, _ := ret[0].(*Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getPriceByPVID indicates an expected call of getPriceByPVID.
func (mr *MockproductDBAccessorMockRecorder) getPriceByPVID(ctx, pvID any) *MockproductDBAccessorgetPriceByPVIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getPriceByPVID", reflect.TypeOf((*MockproductDBAccessor)(nil).getPriceByPVID), ctx, pvID)
	return &MockproductDBAccessorgetPriceByPVIDCall{Call: call}
}

// MockproductDBAccessorgetPriceByPVIDCall wrap *gomock.Call
type MockproductDBAccessorgetPriceByPVIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductDBAccessorgetPriceByPVIDCall) Return(arg0 *Price, arg1 error) *MockproductDBAccessorgetPriceByPVIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductDBAccessorgetPriceByPVIDCall) Do(f func(context.Context, string) (*Price, error)) *MockproductDBAccessorgetPriceByPVIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductDBAccessorgetPriceByPVIDCall) DoAndReturn(f func(context.Context, string) (*Price, error)) *MockproductDBAccessorgetPriceByPVIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// getProductByID mocks base method.
func (m *MockproductDBAccessor) getProductByID(ctx context.Context, productID string) (*Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getProductByID", ctx, productID)
	ret0, _ := ret[0].(*Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getProductByID indicates an expected call of getProductByID.
func (mr *MockproductDBAccessorMockRecorder) getProductByID(ctx, productID any) *MockproductDBAccessorgetProductByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getProductByID", reflect.TypeOf((*MockproductDBAccessor)(nil).getProductByID), ctx, productID)
	return &MockproductDBAccessorgetProductByIDCall{Call: call}
}

// MockproductDBAccessorgetProductByIDCall wrap *gomock.Call
type MockproductDBAccessorgetProductByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductDBAccessorgetProductByIDCall) Return(arg0 *Product, arg1 error) *MockproductDBAccessorgetProductByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductDBAccessorgetProductByIDCall) Do(f func(context.Context, string) (*Product, error)) *MockproductDBAccessorgetProductByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductDBAccessorgetProductByIDCall) DoAndReturn(f func(context.Context, string) (*Product, error)) *MockproductDBAccessorgetProductByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// getProductCategoryByID mocks base method.
func (m *MockproductDBAccessor) getProductCategoryByID(ctx context.Context, pvID string) (*ProductCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getProductCategoryByID", ctx, pvID)
	ret0, _ := ret[0].(*ProductCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getProductCategoryByID indicates an expected call of getProductCategoryByID.
func (mr *MockproductDBAccessorMockRecorder) getProductCategoryByID(ctx, pvID any) *MockproductDBAccessorgetProductCategoryByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getProductCategoryByID", reflect.TypeOf((*MockproductDBAccessor)(nil).getProductCategoryByID), ctx, pvID)
	return &MockproductDBAccessorgetProductCategoryByIDCall{Call: call}
}

// MockproductDBAccessorgetProductCategoryByIDCall wrap *gomock.Call
type MockproductDBAccessorgetProductCategoryByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockproductDBAccessorgetProductCategoryByIDCall) Return(arg0 *ProductCategory, arg1 error) *MockproductDBAccessorgetProductCategoryByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockproductDBAccessorgetProductCategoryByIDCall) Do(f func(context.Context, string) (*ProductCategory, error)) *MockproductDBAccessorgetProductCategoryByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockproductDBAccessorgetProductCategoryByIDCall) DoAndReturn(f func(context.Context, string) (*ProductCategory, error)) *MockproductDBAccessorgetProductCategoryByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
