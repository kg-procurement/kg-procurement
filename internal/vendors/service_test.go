package vendors

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/onsi/gomega"
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
	sampleData := []Vendor{
		{
			ID:            "1",
			Name:          "name",
			Description:   "description",
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
		{
			name: "success",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args:    args{ctx: context.Background()},
			want:    sampleData,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args:    args{ctx: context.Background()},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			v := NewVendorService(tt.fields.mockVendorDBAccessor)

			tt.fields.mockVendorDBAccessor.EXPECT().
				GetAll(tt.args.ctx).
				Return(tt.want, nil)

			res, err := v.GetAll(tt.args.ctx)

			g.Expect(err).To(gomega.BeNil())
			g.Expect(res).To(gomega.Equal(tt.want))
		})
	}
}

func TestVendorService_GetByLocation(t *testing.T) {
	location := "Indonesia"
	sampleData := []Vendor{
		{
			ID:            "1",
			Name:          "name",
			Description:   "description",
			BpId:          1,
			BpName:        "bp_name",
			Rating:        1,
			AreaGroupId:   1,
			AreaGroupName: location,
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
		{
			name: "success",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args:    args{ctx: context.Background()},
			want:    sampleData,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args:    args{ctx: context.Background()},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			v := NewVendorService(tt.fields.mockVendorDBAccessor)

			tt.fields.mockVendorDBAccessor.EXPECT().
				GetByLocation(tt.args.ctx, location).
				Return(tt.want, nil)

			res, err := v.GetByLocation(tt.args.ctx, location)

			g.Expect(err).To(gomega.BeNil())
			g.Expect(res).To(gomega.Equal(tt.want))
		})
	}
}
