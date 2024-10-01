package vendors

import (
	"context"
	"fmt"
	"kg/procurement/internal/common/database"
	"reflect"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewVendorService(t *testing.T) {
	_ = NewVendorService(nil)
}

func TestNewVendorService(t *testing.T) {
	type args struct {
		mockDBConnector *database.MockDBConnector
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
			if got := NewVendorService(tt.args.mockDBConnector); !reflect.DeepEqual(got, tt.want) {
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
			BpID:          "1",
			BpName:        "bp_name",
			Rating:        1,
			AreaGroupID:   "1",
			AreaGroupName: "group_name",
			SapCode:       "sap_code",
			ModifiedDate:  time.Now(),
			ModifiedBy:    1,
			Date:          time.Now(),
		},
	}
	accessorData := &AccessorGetAllPaginationData{
		Vendors:      sampleData,
		TotalEntries: 1,
	}
	payloadData := &ServiceGetAllPaginationData{
		Vendors:      sampleData,
		TotalEntries: 1,
		CurrentPage:  1,
		PreviousPage: nil,
		NextPage:     2,
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		mockVendorDBAccessor *MockvendorDBAccessor
		mockDBConnector      *database.MockDBConnector
	}
	type args struct {
		ctx  context.Context
		spec ServiceGetAllPaginationSpec
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantService  *ServiceGetAllPaginationData
		wantAccessor *AccessorGetAllPaginationData
	}{
		{
			name: "success",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args:         args{ctx: context.Background(), spec: ServiceGetAllPaginationSpec{Limit: 10, Order: "DESC", Page: 1}},
			wantService:  payloadData,
			wantAccessor: accessorData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			v := &VendorService{
				vendorDBAccessor: tt.fields.mockVendorDBAccessor,
			}

			offset := tt.args.spec.Limit * (tt.args.spec.Page - 1)

			accessorSpec := AccessorGetAllPaginationSpec{
				Limit:  tt.args.spec.Limit,
				Offset: offset,
				Order:  tt.args.spec.Order,
			}

			tt.fields.mockVendorDBAccessor.EXPECT().
				GetAll(tt.args.ctx, accessorSpec).
				Return(tt.wantAccessor, nil)

			res, err := v.GetAll(tt.args.ctx, tt.args.spec)

			g.Expect(err).To(gomega.BeNil())
			g.Expect(res).ToNot(gomega.BeNil())
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
			BpID:          "1",
			BpName:        "bp_name",
			Rating:        1,
			AreaGroupID:   "1",
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
			v := &VendorService{
				vendorDBAccessor: tt.fields.mockVendorDBAccessor,
			}

			if tt.wantErr {
				tt.fields.mockVendorDBAccessor.EXPECT().
					GetByLocation(tt.args.ctx, location).
					Return(nil, fmt.Errorf("some error"))
			} else {
				tt.fields.mockVendorDBAccessor.EXPECT().
					GetByLocation(tt.args.ctx, location).
					Return(tt.want, nil)
			}

			res, err := v.GetByLocation(tt.args.ctx, location)

			if tt.wantErr {
				g.Expect(err).ToNot(gomega.BeNil())
				g.Expect(res).To(gomega.BeNil())
			} else {
				g.Expect(err).To(gomega.BeNil())
				g.Expect(res).To(gomega.Equal(tt.want))
			}
		})
	}
}
