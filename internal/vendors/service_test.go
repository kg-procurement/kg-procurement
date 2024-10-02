package vendors

import (
	"context"
	"fmt"
	"kg/procurement/internal/common/database"
	"reflect"
	"strings"
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
			ModifiedBy:    "1",
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
		spec database.PaginationSpec
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
			args: args{ctx: context.Background(), spec: database.PaginationSpec{Limit: 10, Order: "DESC", Page: 1}},
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

			accessorSpec := database.PaginationSpec{
				Limit: tt.args.spec.Limit,
				Page:  tt.args.spec.Page,
				Order: tt.args.spec.Order,
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
			ModifiedBy:    "1",
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

func TestVendorService_GetByProduct(t *testing.T) {
	product := "ProductA"
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
			ModifiedBy:    "1",
			Date:          time.Now(),
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		mockVendorDBAccessor *MockvendorDBAccessor
	}
	type args struct {
		ctx     context.Context
		product string
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
			args: args{
				ctx:     context.Background(),
				product: product,
			},
			want:    sampleData,
			wantErr: false,
		},
		{
			name: "failure",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args: args{
				ctx:     context.Background(),
				product: product,
			},
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

			productDescription := strings.Fields(tt.args.product)
			if tt.wantErr {
				tt.fields.mockVendorDBAccessor.EXPECT().
					GetByProductDescription(tt.args.ctx, productDescription).
					Return(nil, fmt.Errorf("some error"))
			} else {
				tt.fields.mockVendorDBAccessor.EXPECT().
					GetByProductDescription(tt.args.ctx, productDescription).
					Return(tt.want, nil)
			}

			res, err := v.GetByProduct(tt.args.ctx, tt.args.product)

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

func TestVendorService_Put(t *testing.T) {
	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)
	updatedFixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 1, time.UTC)

	existingVendorData := Vendor{
		ID:            "ID",
		Name:          "name",
		Description:   "description",
		BpID:          "BpID",
		BpName:        "BpName",
		Rating:        1,
		AreaGroupID:   "AreaGroupID",
		AreaGroupName: "AreaGroupName",
		SapCode:       "SapCode",
		ModifiedDate:  fixedTime,
		ModifiedBy:    "ID",
		Date:          fixedTime,
	}
	updateSpec := Vendor{
		ID:            "ID",
		Name:          "udpate",
		Description:   "udpate",
		BpID:          "udpate",
		BpName:        "udpate",
		Rating:        2,
		AreaGroupID:   "udpate",
		AreaGroupName: "udpate",
		SapCode:       "udpate",
		ModifiedDate:  time.Time{},
		ModifiedBy:    "",
		Date:          time.Time{},
	}

	UpdatedVendorData := Vendor{
		ID:            "ID",
		Name:          "udpate",
		Description:   "udpate",
		BpID:          "udpate",
		BpName:        "udpate",
		Rating:        2,
		AreaGroupID:   "udpate",
		AreaGroupName: "udpate",
		SapCode:       "udpate",
		ModifiedDate:  updatedFixedTime,
		ModifiedBy:    "UpdatedID",
		Date:          fixedTime,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		mockvendorDBAccessor *MockvendorDBAccessor
	}

	type args struct {
		ctx  context.Context
		spec Vendor
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantGetById *Vendor
		wantPut     *Vendor
		wantErr     error
	}{
		{
			name: "success",
			fields: fields{
				mockvendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args: args{
				ctx:  context.Background(),
				spec: updateSpec,
			},
			wantGetById: &existingVendorData,
			wantPut:     &UpdatedVendorData,
			wantErr:     nil,
		},
	}
	for _, tt := range tests {
		g := gomega.NewWithT(t)
		v := VendorService{
			vendorDBAccessor: tt.fields.mockvendorDBAccessor,
		}

		tt.fields.mockvendorDBAccessor.
			EXPECT().
			GetById(tt.args.ctx, tt.args.spec.ID).
			Return(tt.wantGetById, nil)

		newVendor := Vendor(*tt.wantGetById)
		newVendor.Name = tt.args.spec.Name
		newVendor.Description = tt.args.spec.Description
		newVendor.BpID = tt.args.spec.BpID
		newVendor.BpName = tt.args.spec.BpName
		newVendor.Rating = tt.args.spec.Rating
		newVendor.AreaGroupID = tt.args.spec.AreaGroupID
		newVendor.AreaGroupName = tt.args.spec.AreaGroupName
		newVendor.SapCode = tt.args.spec.SapCode

		tt.fields.mockvendorDBAccessor.
			EXPECT().
			Put(tt.args.ctx, newVendor).
			Return(tt.wantPut, nil)

		updatedVendor, err := v.Put(tt.args.ctx, newVendor)

		g.Expect(err).To(gomega.BeNil())
		g.Expect(updatedVendor).To(gomega.Equal(tt.wantPut))
	}
}
