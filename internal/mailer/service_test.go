package mailer

import (
	"context"
	"errors"
	"kg/procurement/internal/common/database"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewEmailStatusService(t *testing.T) {
	_ = NewEmailStatusService(nil, nil)
}

func TestEmailStatusService_WriteEmailStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		now       = time.Now()
		writeSpec = EmailStatus{
			ID:           "some_id",
			EmailTo:      "email@email.com",
			Status:       "sent",
			ModifiedDate: now,
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			g                       = gomega.NewWithT(t)
			ctx                     = context.Background()
			mockCtrl                = gomock.NewController(t)
			mockEmailStatusAccessor = NewMockemailStatusDBAccessor(mockCtrl)
		)

		svc := &EmailStatusService{
			emailStatusDBAccessor: mockEmailStatusAccessor,
		}

		mockEmailStatusAccessor.EXPECT().WriteEmailStatus(ctx, writeSpec).Return(nil)

		err := svc.WriteEmailStatus(ctx, writeSpec)
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("returns error on accessor failure", func(t *testing.T) {
		var (
			g                       = gomega.NewWithT(t)
			ctx                     = context.Background()
			mockCtrl                = gomock.NewController(t)
			mockEmailStatusAccessor = NewMockemailStatusDBAccessor(mockCtrl)
		)

		svc := &EmailStatusService{
			emailStatusDBAccessor: mockEmailStatusAccessor,
		}

		mockEmailStatusAccessor.EXPECT().WriteEmailStatus(ctx, writeSpec).Return(errors.New("create error"))

		err := svc.WriteEmailStatus(ctx, writeSpec)
		g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func TestEmailService_GetAll(t *testing.T) {
	// Sample data to be returned by the mock accessor
	sampleData := []EmailStatus{
		{
			ID:           "1",
			EmailTo:      "test@example.com",
			Status:       "sent",
			ModifiedDate: time.Now(),
		},
	}

	data := &AccessorGetAllPaginationData{
		EmailStatus: sampleData,
		Metadata: database.PaginationMetadata{
			TotalPage:   1,
			CurrentPage: 1,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		mockEmailStatusDBAccessor *MockemailStatusDBAccessor
	}

	type args struct {
		ctx  context.Context
		spec GetAllEmailStatusSpec
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
				mockEmailStatusDBAccessor: NewMockemailStatusDBAccessor(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: GetAllEmailStatusSpec{
					PaginationSpec: database.PaginationSpec{
						Limit: 10,
						Order: "DESC",
						Page:  1,
					},
				},
			},
			want: data,
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			// Initialize the EmailService with the mock accessor
			e := &EmailStatusService{
				emailStatusDBAccessor: tt.fields.mockEmailStatusDBAccessor,
			}

			// Prepare the accessor spec based on the test args
			accessorSpec := GetAllEmailStatusSpec{
				PaginationSpec: database.PaginationSpec{
					Limit:   tt.args.spec.Limit,
					Page:    tt.args.spec.Page,
					Order:   tt.args.spec.Order,
					OrderBy: tt.args.spec.OrderBy,
				},
			}

			// Set up the mock expectation
			tt.fields.mockEmailStatusDBAccessor.EXPECT().
				GetAll(tt.args.ctx, accessorSpec).
				Return(tt.want, tt.err)

			// Call the method under test
			res, err := e.GetAll(tt.args.ctx, tt.args.spec)

			// Assert the results
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
