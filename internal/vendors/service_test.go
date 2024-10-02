package vendors

import (
	"context"
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

	data := &AccessorGetAllPaginationData{
		Vendors: sampleData,
		Metadata: database.PaginationMetadata{
			TotalPage:   1,
			CurrentPage: 1,
		},
	}

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	type fields struct {
		mockVendorDBAccessor *MockvendorDBAccessor
		mockDBConnector      *database.MockDBConnector
	}

	type args struct {
		ctx  context.Context
		spec GetAllVendorSpec
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *AccessorGetAllPaginationData
		err    error
	}{
		{
			name: "success",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: GetAllVendorSpec{
					PaginationSpec: database.PaginationSpec{Limit: 10, Order: "DESC", Page: 1},
				},
			},
			want: data,
			err:  nil,
		},
		{
			name: "success with location and product",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: GetAllVendorSpec{
					PaginationSpec: database.PaginationSpec{Limit: 10, Order: "DESC", Page: 1},
					Location:       "Indonesia",
					Product:        "test product",
				},
			},
			want: data,
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			v := &VendorService{
				vendorDBAccessor: tt.fields.mockVendorDBAccessor,
			}

			accessorSpec := GetAllVendorSpec{
				PaginationSpec: database.PaginationSpec{
					Limit: tt.args.spec.Limit,
					Page:  tt.args.spec.Page,
					Order: tt.args.spec.Order,
				},
				Location: tt.args.spec.Location,
				Product:  tt.args.spec.Product,
			}

			tt.fields.mockVendorDBAccessor.EXPECT().
				GetAll(tt.args.ctx, accessorSpec).
				Return(tt.want, tt.err)

			res, err := v.GetAll(tt.args.ctx, tt.args.spec)

			if tt.err == nil {
				g.Expect(err).To(gomega.BeNil())
				g.Expect(res).To(gomega.Equal(tt.want))
			} else {
				g.Expect(err).ToNot(gomega.BeNil())
				g.Expect(res).To(gomega.BeNil())
			}
		})
	}
}
