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
	_ = NewVendorService(config.Application{}, nil, nil, nil, nil)
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
		mockEmailStatusSvc *MockemailStatusSvc
		subject            *VendorService
	)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		ctrl := gomock.NewController(t)
		mockVendorAccessor = NewMockvendorDBAccessor(ctrl)
		mockEmailProvider = mailer.NewMockEmailProvider(ctrl)
		mockEmailStatusSvc = NewMockemailStatusSvc(ctrl)

		subject = &VendorService{
			cfg:              config.Application{},
			vendorDBAccessor: mockVendorAccessor,
			smtpProvider:     mockEmailProvider,
			emailStatusSvc:   mockEmailStatusSvc,
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

		mockEmailStatusSvc.EXPECT().
			WriteEmailStatus(ctx, gomock.Any()).
			Return(nil).
			Times(2)

		errList, err := subject.BlastEmail(ctx, vendorIDs, mailer.Email{
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

		errList, err := subject.BlastEmail(ctx, vendorIDs, mailer.Email{
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

		mockEmailStatusSvc.EXPECT().
			WriteEmailStatus(ctx, gomock.Any()).
			Return(nil).
			Times(2)

		errList, err := subject.BlastEmail(ctx, vendorIDs, mailer.Email{
			Subject: "test",
			Body:    "email body here uwaa",
		})
		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(errList).To(gomega.HaveLen(1))
	})

	t.Run("error writing email status", func(t *testing.T) {
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

		// simulate error when inserting to email_status
		mockEmailStatusSvc.EXPECT().
			WriteEmailStatus(ctx, gomock.Any()).
			Return(errors.New("write error")).
			Times(2)

		errList, err := subject.BlastEmail(ctx, vendorIDs, mailer.Email{
			Subject: "Test Subject",
			Body:    "Test Body",
		})

		g.Expect(err).To(gomega.BeNil())
		g.Expect(errList).To(gomega.BeNil())
	})
}

func TestVendorService_AutomatedBlastEmail(t *testing.T) {
	t.Parallel()

	var (
		mockVendorAccessor *MockvendorDBAccessor
		mockEmailProvider  *mailer.MockEmailProvider
		mockEmailStatusSvc *MockemailStatusSvc
		service            *VendorService
	)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		ctrl := gomock.NewController(t)
		mockVendorAccessor = NewMockvendorDBAccessor(ctrl)
		mockEmailProvider = mailer.NewMockEmailProvider(ctrl)
		mockEmailStatusSvc = NewMockemailStatusSvc(ctrl)

		service = &VendorService{
			cfg:              config.Application{},
			vendorDBAccessor: mockVendorAccessor,
			smtpProvider:     mockEmailProvider,
			emailStatusSvc:   mockEmailStatusSvc,
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

		mockEmailStatusSvc.EXPECT().
			WriteEmailStatus(ctx, gomock.Any()).
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

		mockEmailStatusSvc.EXPECT().
			WriteEmailStatus(ctx, gomock.Any()).
			Return(nil).
			Times(2)

		errList, err := service.AutomatedEmailBlast(ctx, product_name)
		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(errList).To(gomega.HaveLen(1))
	})
}

func TestVendorService_CreateEvaluation(t *testing.T) {
	t.Parallel()

	var (
		mockVendorAccessor *MockvendorDBAccessor
		service            *VendorService
	)

	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		ctrl := gomock.NewController(t)
		mockVendorAccessor = NewMockvendorDBAccessor(ctrl)

		service = &VendorService{
			cfg:              config.Application{},
			vendorDBAccessor: mockVendorAccessor,
		}

		return gomega.NewWithT(t)
	}

	var (
		vendorEvaluation = VendorEvaluation{
			VendorID:                         "1",
			KesesuaianProduk:                 1,
			KualitasProduk:                   1,
			KetepatanWaktuPengiriman:         1,
			KompetitifitasHarga:              1,
			ResponsivitasKemampuanKomunikasi: 1,
			KemampuanDalamMenanganiMasalah:   1,
			KelengkapanBarang:                1,
			Harga:                            1,
			TermOfPayment:                    1,
			Reputasi:                         1,
			KetersediaanBarang:               1,
			KualitasLayananAfterServices:     1,
		}

		expectation = VendorEvaluation{
			ID:                               "H58S2LBQblHMjce",
			VendorID:                         "1",
			KesesuaianProduk:                 1,
			KualitasProduk:                   1,
			KetepatanWaktuPengiriman:         1,
			KompetitifitasHarga:              1,
			ResponsivitasKemampuanKomunikasi: 1,
			KemampuanDalamMenanganiMasalah:   1,
			KelengkapanBarang:                1,
			Harga:                            1,
			TermOfPayment:                    1,
			Reputasi:                         1,
			KetersediaanBarang:               1,
			KualitasLayananAfterServices:     1,
			ModifiedDate:                     fixedTime,
		}
	)

	t.Run("success", func(t *testing.T) {
		g := setup(t)
		ctx := context.Background()

		mockVendorAccessor.EXPECT().
			CreateEvaluation(ctx, &vendorEvaluation).
			Return(&expectation, nil)

		result, err := service.CreateEvaluation(ctx, &vendorEvaluation)
		g.Expect(err).To(gomega.BeNil())
		g.Expect(result).To(gomega.Equal(&expectation))
	})
	t.Run("error", func(t *testing.T) {
		g := setup(t)
		ctx := context.Background()

		mockVendorAccessor.EXPECT().
			CreateEvaluation(ctx, &vendorEvaluation).Return(nil, errors.New("error"))

		result, err := service.CreateEvaluation(ctx, &vendorEvaluation)
		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(result).To(gomega.BeNil())
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
		temp := mailer.Email{
			Body: "this is body",
		}
		service.applyDefaultEmailTemplate(&temp)
		g.Expect(temp.Subject).ToNot(gomega.BeEmpty())
		g.Expect(temp.Body).To(gomega.Equal("this is body"))
	})

	t.Run("body is empty", func(t *testing.T) {
		g := setup(t)
		temp := mailer.Email{
			Subject: "this is subject",
		}
		service.applyDefaultEmailTemplate(&temp)
		g.Expect(temp.Body).ToNot(gomega.BeEmpty())
		g.Expect(temp.Subject).To(gomega.Equal("this is subject"))
	})

	t.Run("both are empty", func(t *testing.T) {
		g := setup(t)
		temp := mailer.Email{}
		service.applyDefaultEmailTemplate(&temp)
		g.Expect(temp.Body).ToNot(gomega.BeEmpty())
		g.Expect(temp.Subject).ToNot(gomega.BeEmpty())
	})

	t.Run("both are filled", func(t *testing.T) {
		g := setup(t)
		temp := mailer.Email{
			Subject: "this is subject",
			Body:    "this is body",
		}
		service.applyDefaultEmailTemplate(&temp)
		g.Expect(temp.Body).To(gomega.Equal("this is body"))
		g.Expect(temp.Subject).To(gomega.Equal("this is subject"))
	})
}

func TestVendorService_GetPopulatedEmailStatus(t *testing.T) {
	t.Parallel()

	var (
		mockVendorAccessor *MockvendorDBAccessor
		mockEmailStatusSvc *MockemailStatusSvc
		service            *VendorService
	)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		ctrl := gomock.NewController(t)
		mockVendorAccessor = NewMockvendorDBAccessor(ctrl)
		mockEmailStatusSvc = NewMockemailStatusSvc(ctrl)

		service = &VendorService{
			vendorDBAccessor: mockVendorAccessor,
			emailStatusSvc:   mockEmailStatusSvc,
		}

		return gomega.NewWithT(t)
	}

	t.Run("successfully populates vendor names", func(t *testing.T) {
		g := setup(t)

		sampleEmailStatus := []mailer.EmailStatus{
			{
				ID:           "1",
				EmailTo:      "vendor1@example.com",
				Status:       "sent",
				VendorID:     "vendor1",
				DateSent:     time.Now(),
				ModifiedDate: time.Now(),
			},
			{
				ID:           "2",
				EmailTo:      "vendor2@example.com",
				Status:       "sent",
				VendorID:     "vendor2",
				DateSent:     time.Now(),
				ModifiedDate: time.Now(),
			},
		}

		sampleVendors := []Vendor{
			{
				ID:   "vendor1",
				Name: "Vendor 1",
			},
			{
				ID:   "vendor2",
				Name: "Vendor 2",
			},
		}

		mockEmailStatusSvc.EXPECT().GetAllEmailStatus(gomock.Any(), gomock.Any()).Return(
			&mailer.AccessorGetEmailStatusPaginationData{
				EmailStatus: sampleEmailStatus,
				Metadata: database.PaginationMetadata{
					TotalPage:   1,
					CurrentPage: 1,
				},
			}, nil)

		mockVendorAccessor.EXPECT().BulkGetByIDs(gomock.Any(), []string{"vendor1", "vendor2"}).Return(sampleVendors, nil)

		result, err := service.GetPopulatedEmailStatus(context.Background(), mailer.GetAllEmailStatusSpec{})

		g.Expect(err).To(gomega.BeNil())
		g.Expect(result).ToNot(gomega.BeNil())
		g.Expect(result.EmailStatus).To(gomega.HaveLen(2))

		g.Expect(result.EmailStatus[0].VendorName).To(gomega.Equal("Vendor 1"))
		g.Expect(result.EmailStatus[1].VendorName).To(gomega.Equal("Vendor 2"))
	})

	t.Run("vendor not found", func(t *testing.T) {
		g := setup(t)

		sampleEmailStatus := []mailer.EmailStatus{
			{
				ID:           "1",
				EmailTo:      "vendor1@example.com",
				Status:       "sent",
				VendorID:     "vendor1",
				DateSent:     time.Now(),
				ModifiedDate: time.Now(),
			},
		}

		sampleVendors := []Vendor{}

		mockEmailStatusSvc.EXPECT().GetAllEmailStatus(gomock.Any(), gomock.Any()).Return(
			&mailer.AccessorGetEmailStatusPaginationData{
				EmailStatus: sampleEmailStatus,
				Metadata: database.PaginationMetadata{
					TotalPage:   1,
					CurrentPage: 1,
				},
			}, nil)

		mockVendorAccessor.EXPECT().BulkGetByIDs(gomock.Any(), []string{"vendor1"}).Return(sampleVendors, nil)

		result, err := service.GetPopulatedEmailStatus(context.Background(), mailer.GetAllEmailStatusSpec{})

		g.Expect(err).To(gomega.BeNil())
		g.Expect(result).ToNot(gomega.BeNil())
		g.Expect(result.EmailStatus).To(gomega.HaveLen(1))

		g.Expect(result.EmailStatus[0].VendorName).To(gomega.Equal("Unknown Vendor"))
	})

	t.Run("fail when getting vendors by ID", func(t *testing.T) {
		g := setup(t)

		sampleEmailStatus := []mailer.EmailStatus{
			{
				ID:           "1",
				EmailTo:      "vendor1@example.com",
				Status:       "sent",
				VendorID:     "vendor1",
				DateSent:     time.Now(),
				ModifiedDate: time.Now(),
			},
		}

		mockEmailStatusSvc.EXPECT().GetAllEmailStatus(gomock.Any(), gomock.Any()).Return(
			&mailer.AccessorGetEmailStatusPaginationData{
				EmailStatus: sampleEmailStatus,
				Metadata: database.PaginationMetadata{
					TotalPage:   1,
					CurrentPage: 1,
				},
			}, nil)

		mockVendorAccessor.EXPECT().BulkGetByIDs(gomock.Any(), []string{"vendor1"}).Return(nil, errors.New("failed to fetch vendors"))

		result, err := service.GetPopulatedEmailStatus(context.Background(), mailer.GetAllEmailStatusSpec{})

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(err.Error()).To(gomega.Equal("failed to fetch vendors"))
		g.Expect(result).To(gomega.BeNil())
	})

	t.Run("fail when getting email status", func(t *testing.T) {
		g := setup(t)

		mockEmailStatusSvc.EXPECT().GetAllEmailStatus(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to fetch email status"))

		result, err := service.GetPopulatedEmailStatus(context.Background(), mailer.GetAllEmailStatusSpec{})

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(err.Error()).To(gomega.Equal("failed to fetch email status"))
		g.Expect(result).To(gomega.BeNil())
	})
}
