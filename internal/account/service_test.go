package account

import (
	"context"
	"errors"
	"testing"

	"bou.ke/monkey"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAccountService_RegisterAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		mockAccountDBAccessor *MockaccountDBAccessor
	}

	type args struct {
		ctx  context.Context
		spec RegisterAccountSpec
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				mockAccountDBAccessor: NewMockaccountDBAccessor(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: RegisterAccountSpec{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid email",
			args: args{
				ctx: context.Background(),
				spec: RegisterAccountSpec{
					Email:    "invalid-email",
					Password: "password123",
				},
			},
			wantErr: errors.New("invalid email"),
		},
		{
			name: "failed to hash password",
			fields: fields{
				mockAccountDBAccessor: NewMockaccountDBAccessor(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: RegisterAccountSpec{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			wantErr: errors.New("failed to hash password"),
		},
		{
			name: "failed to generate ID",
			fields: fields{
				mockAccountDBAccessor: NewMockaccountDBAccessor(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: RegisterAccountSpec{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			wantErr: errors.New("failed to generate ID"),
		},
		{
			name: "database error",
			fields: fields{
				mockAccountDBAccessor: NewMockaccountDBAccessor(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: RegisterAccountSpec{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			a := &AccountService{
				accountDBAccessor: tt.fields.mockAccountDBAccessor,
			}

			if tt.name == "success" {
				tt.fields.mockAccountDBAccessor.EXPECT().
					RegisterAccount(tt.args.ctx, gomock.Any()).
					Return(nil)
			} else if tt.name == "database error" {
				tt.fields.mockAccountDBAccessor.EXPECT().
					RegisterAccount(tt.args.ctx, gomock.Any()).
					Return(errors.New("database error"))
			} else if tt.name == "failed to hash password" {
				// Mock bcrypt.GenerateFromPassword to return an error
				monkey.Patch(bcrypt.GenerateFromPassword, func([]byte, int) ([]byte, error) {
					return nil, errors.New("failed to hash password")
				})
				defer monkey.Unpatch(bcrypt.GenerateFromPassword)
			} else if tt.name == "failed to generate ID" {
				// Mock GenerateRandomID to return an error
				monkey.Patch(GenerateRandomID, func() (string, error) {
					return "", errors.New("failed to generate ID")
				})
				defer monkey.Unpatch(GenerateRandomID)
			}

			err := a.RegisterAccount(tt.args.ctx, tt.args.spec)

			if tt.wantErr == nil {
				g.Expect(err).To(gomega.BeNil())
			} else {
				g.Expect(err).ToNot(gomega.BeNil())
				g.Expect(err.Error()).To(gomega.ContainSubstring(tt.wantErr.Error()))
			}
		})
	}
}
