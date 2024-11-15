package vendors

import (
	"context"
	"errors"
	"kg/procurement/cmd/config"
	"kg/procurement/internal/common/database"
	"kg/procurement/internal/mailer"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewVendorService(t *testing.T) {
	_ = NewVendorService(config.Application{}, nil, nil, nil)
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
			name: "success with order by, location and product",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: GetAllVendorSpec{
					PaginationSpec: database.PaginationSpec{
						Limit:   10,
						Order:   "DESC",
						OrderBy: "rating",
						Page:    1,
					},
					Location: "Indonesia",
					Product:  "test product",
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
					Limit:   tt.args.spec.Limit,
					Page:    tt.args.spec.Page,
					Order:   tt.args.spec.Order,
					OrderBy: tt.args.spec.OrderBy,
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

func TestVendorService_GetById(t *testing.T) {
	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)

	data := &Vendor{
		ID:            "1",
		Name:          "name",
		Description:   "description",
		BpID:          "1",
		BpName:        "bp_name",
		Rating:        1,
		AreaGroupID:   "1",
		AreaGroupName: "group_name",
		SapCode:       "sap_code",
		ModifiedDate:  fixedTime,
		ModifiedBy:    "1",
		Date:          fixedTime,
	}

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	type fields struct {
		mockVendorDBAccessor *MockvendorDBAccessor
	}

	type args struct {
		ctx context.Context
		id  string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Vendor
		err    error
	}{
		{
			name: "success",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args: args{ctx: context.Background(), id: "ID"},
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

			tt.fields.mockVendorDBAccessor.EXPECT().
				GetById(tt.args.ctx, tt.args.id).
				Return(tt.want, tt.err)

			res, err := v.GetById(tt.args.ctx, tt.args.id)

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

func TestVendorService_UpdateDetail(t *testing.T) {
	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)
	updatedFixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 1, time.UTC)
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
		ModifiedBy:    "ID",
		Date:          time.Time{},
	}

	UpdatedVendorData := &Vendor{
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
		name    string
		fields  fields
		args    args
		want    *Vendor
		wantErr error
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
			want: UpdatedVendorData,
		},
	}
	for _, tt := range tests {
		g := gomega.NewWithT(t)
		v := VendorService{
			vendorDBAccessor: tt.fields.mockvendorDBAccessor,
		}

		tt.fields.mockvendorDBAccessor.
			EXPECT().
			UpdateDetail(tt.args.ctx, tt.args.spec).
			Return(tt.want, tt.wantErr)

		res, err := v.UpdateDetail(tt.args.ctx, tt.args.spec)

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(tt.want))

	}
}

func TestVendorService_GetLocations(t *testing.T) {
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
		want    []string
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []string{"Location1", "Location2"},
			wantErr: nil,
		},
		{
			name: "database error",
			fields: fields{
				mockVendorDBAccessor: NewMockvendorDBAccessor(ctrl),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			v := &VendorService{
				vendorDBAccessor: tt.fields.mockVendorDBAccessor,
			}

			if tt.name == "success" {
				tt.fields.mockVendorDBAccessor.EXPECT().
					GetAllLocations(tt.args.ctx).
					Return(tt.want, nil)
			} else if tt.name == "database error" {
				tt.fields.mockVendorDBAccessor.EXPECT().
					GetAllLocations(tt.args.ctx).
					Return(nil, errors.New("database error"))
			}

			got, err := v.GetLocations(tt.args.ctx)

			if tt.wantErr == nil {
				g.Expect(err).To(gomega.BeNil())
				g.Expect(got).To(gomega.Equal(tt.want))
			} else {
				g.Expect(err).ToNot(gomega.BeNil())
				g.Expect(err.Error()).To(gomega.ContainSubstring(tt.wantErr.Error()))
			}
		})
	}
}

func TestVendorService_BlastEmail(t *testing.T) {
	t.Parallel()

	var (
		mockVendorAccessor *MockvendorDBAccessor
		mockEmailProvider  *mailer.MockEmailProvider
		subject            *VendorService
	)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		ctrl := gomock.NewController(t)
		mockVendorAccessor = NewMockvendorDBAccessor(ctrl)
		mockEmailProvider = mailer.NewMockEmailProvider(ctrl)
		subject = &VendorService{
			cfg:              config.Application{},
			vendorDBAccessor: mockVendorAccessor,
			smtpProvider:     mockEmailProvider,
		}

		return gomega.NewWithT(t)
	}

	var (
		vendors = []Vendor{
			{
				ID:    "1111",
				Email: "valenganteng@gmail.com",
			},
			{
				ID:    "2222",
				Email: "ferryganteng@gmail.com",
			},
		}
	)

	t.Run("success", func(t *testing.T) {
		g := setup(t)
		ctx := context.Background()

		vendorIDs := []string{"1111", "2222"}
		mockVendorAccessor.EXPECT().
			BulkGetByIDs(ctx, vendorIDs).
			Return(vendors, nil)

		mockEmailProvider.EXPECT().
			SendEmail(gomock.Any()).
			Return(nil).
			Times(2)

		errList, err := subject.BlastEmail(ctx, vendorIDs, emailTemplate{
			Subject: "test",
			Body:    "email body here uwaa",
		})
		g.Expect(err).To(gomega.BeNil())
		g.Expect(errList).To(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		g := setup(t)
		ctx := context.Background()

		vendorIDs := []string{"1111", "2222"}
		mockVendorAccessor.EXPECT().
			BulkGetByIDs(ctx, vendorIDs).
			Return(nil, errors.New("oh noo"))

		errList, err := subject.BlastEmail(ctx, vendorIDs, emailTemplate{
			Subject: "test",
			Body:    "email body here uwaa",
		})
		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(errList).To(gomega.BeNil())
	})

	t.Run("does not return error even if sending email fails", func(t *testing.T) {
		g := setup(t)
		ctx := context.Background()

		vendorIDs := []string{"1111", "2222"}
		mockVendorAccessor.EXPECT().
			BulkGetByIDs(ctx, vendorIDs).
			Return(vendors, nil)

		mockEmailProvider.EXPECT().
			SendEmail(gomock.Any()).
			Return(errors.New("oh nooo"))

		mockEmailProvider.EXPECT().
			SendEmail(gomock.Any()).
			Return(nil)

		errList, err := subject.BlastEmail(ctx, vendorIDs, emailTemplate{
			Subject: "test",
			Body:    "email body here uwaa",
		})
		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(errList).To(gomega.HaveLen(1))
	})
}

func TestVendorService_AutomatedBlastEmail(t *testing.T) {
	t.Parallel()

	var (
		mockVendorAccessor *MockvendorDBAccessor
		mockEmailProvider  *mailer.MockEmailProvider
		service            *VendorService
	)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		ctrl := gomock.NewController(t)
		mockVendorAccessor = NewMockvendorDBAccessor(ctrl)
		mockEmailProvider = mailer.NewMockEmailProvider(ctrl)
		service = &VendorService{
			cfg:              config.Application{},
			vendorDBAccessor: mockVendorAccessor,
			smtpProvider:     mockEmailProvider,
		}

		return gomega.NewWithT(t)
	}

	var (
		vendors = []Vendor{
			{
				ID:    "1",
				Email: "valerian@outlook.com",
			},
			{
				ID:    "2",
				Email: "salim@outlook.com",
			},
		}
		product_name = "Buku"
	)

	t.Run("success", func(t *testing.T) {
		g := setup(t)
		ctx := context.Background()

		mockVendorAccessor.EXPECT().
			BulkGetByProductName(ctx, product_name).
			Return(vendors, nil)

		mockEmailProvider.EXPECT().
			SendEmail(gomock.Any()).
			Return(nil).
			Times(2)

		errList, err := service.AutomatedEmailBlast(ctx, product_name)
		g.Expect(err).To(gomega.BeNil())
		g.Expect(errList).To(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		g := setup(t)
		ctx := context.Background()

		mockVendorAccessor.EXPECT().
			BulkGetByProductName(ctx, product_name).
			Return(nil, errors.New("error"))

		result, err := service.AutomatedEmailBlast(ctx, product_name)
		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(result).To(gomega.BeNil())
	})

	t.Run("not returning error even when email fails to send", func(t *testing.T) {
		g := setup(t)
		ctx := context.Background()

		mockVendorAccessor.EXPECT().
			BulkGetByProductName(ctx, product_name).
			Return(vendors, nil)

		mockEmailProvider.EXPECT().
			SendEmail(gomock.Any()).
			Return(errors.New("error"))

		mockEmailProvider.EXPECT().
			SendEmail(gomock.Any()).
			Return(nil)

		errList, err := service.AutomatedEmailBlast(ctx, product_name)
		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(errList).To(gomega.HaveLen(1))
	})
}

func TestVendorService_applyDefaultEmailTemplate(t *testing.T) {
	t.Parallel()

	var (
		mockVendorAccessor *MockvendorDBAccessor
		service            *VendorService
	)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		ctrl := gomock.NewController(t)
		mockVendorAccessor = NewMockvendorDBAccessor(ctrl)
		service = &VendorService{
			cfg:              config.Application{},
			vendorDBAccessor: mockVendorAccessor,
		}

		return gomega.NewWithT(t)
	}

	t.Run("subject is empty", func(t *testing.T) {
		g := setup(t)
		temp := emailTemplate{
			Body: "this is body",
		}
		service.applyDefaultEmailTemplate(&temp)
		g.Expect(temp.Subject).ToNot(gomega.BeEmpty())
		g.Expect(temp.Body).To(gomega.Equal("this is body"))
	})

	t.Run("body is empty", func(t *testing.T) {
		g := setup(t)
		temp := emailTemplate{
			Subject: "this is subject",
		}
		service.applyDefaultEmailTemplate(&temp)
		g.Expect(temp.Body).ToNot(gomega.BeEmpty())
		g.Expect(temp.Subject).To(gomega.Equal("this is subject"))
	})

	t.Run("both are empty", func(t *testing.T) {
		g := setup(t)
		temp := emailTemplate{}
		service.applyDefaultEmailTemplate(&temp)
		g.Expect(temp.Body).ToNot(gomega.BeEmpty())
		g.Expect(temp.Subject).ToNot(gomega.BeEmpty())
	})

	t.Run("both are filled", func(t *testing.T) {
		g := setup(t)
		temp := emailTemplate{
			Subject: "this is subject",
			Body:    "this is body",
		}
		service.applyDefaultEmailTemplate(&temp)
		g.Expect(temp.Body).To(gomega.Equal("this is body"))
		g.Expect(temp.Subject).To(gomega.Equal("this is subject"))
	})
}
