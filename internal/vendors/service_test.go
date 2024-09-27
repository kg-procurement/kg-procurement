package vendors

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestNewVendorService(t *testing.T) {
	type args struct {
		vendorDBAccessor vendorDBAccessor
	}
	tests := []struct {
		name string
		args args
		want *VendorService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVendorService(tt.args.vendorDBAccessor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVendorService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVendorService_GetSomeStuff(t *testing.T) {
	type fields struct {
		vendorDBAccessor vendorDBAccessor
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VendorService{
				vendorDBAccessor: tt.fields.vendorDBAccessor,
			}
			if _, err := v.GetSomeStuff(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("GetSomeStuff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVendorService_GetAll(t *testing.T) {
	constant := []Vendor{
		{
			Id:            1,
			Name:          "name",
			BpId:          1,
			BpName:        "bp_name",
			Rating:        1,
			AreaGroupId:   1,
			AreaGroupName: "group_name",
			SapCode:       "sap_code",
			ModifiedDate:  time.Now(),
			ModifiedBy:    1,
			Date:          time.Now(),
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		mockVendorDBAccessor *MockvendorDBAccessor
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Vendor
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test_GetAllVendor Success",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args:    args{ctx: context.Background()},
			want:    constant,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVendorService(tt.fields.mockVendorDBAccessor)

			tt.fields.mockVendorDBAccessor.EXPECT().
				GetAll(tt.args.ctx).
				Return(tt.want, nil)

			vendors, err := v.GetAll(tt.args.ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(vendors, tt.want) {
				t.Errorf("GetAll() = %v, wantErr %v", vendors, tt.wantErr)
			}
		})
	}
}
