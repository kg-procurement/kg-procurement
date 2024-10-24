package account

import (
	"context"
	"errors"
	"testing"

	"bou.ke/monkey"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"kg/procurement/internal/token"
)

func Test_NewAccountService(t *testing.T) {
	_ = NewAccountService(nil, nil, nil)
}

func TestAccountService_RegisterAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		mockAccountDBAccessor *MockaccountDBAccessor
	}

	type args struct {
		ctx  context.Context
		spec RegisterContract
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
				spec: RegisterContract{
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
				spec: RegisterContract{
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
				spec: RegisterContract{
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
				spec: RegisterContract{
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
				spec: RegisterContract{
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
				monkey.Patch(generateRandomID, func() (string, error) {
					return "", errors.New("failed to generate ID")
				})
				defer monkey.Unpatch(generateRandomID)
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

func TestAccountService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		mockAccountDBAccessor *MockaccountDBAccessor
		mockTokenService      *MocktokenService
	}

	type args struct {
		ctx  context.Context
		spec LoginContract
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   error
		wantToken string
	}{
		{
			name: "success",
			fields: fields{
				mockAccountDBAccessor: NewMockaccountDBAccessor(ctrl),
				mockTokenService:      NewMocktokenService(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: LoginContract{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			wantErr:   nil,
			wantToken: "mockToken",
		},
		{
			name: "account not found",
			fields: fields{
				mockAccountDBAccessor: NewMockaccountDBAccessor(ctrl),
				mockTokenService:      NewMocktokenService(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: LoginContract{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			wantErr: errors.New("login failed"),
		},
		{
			name: "invalid password",
			fields: fields{
				mockAccountDBAccessor: NewMockaccountDBAccessor(ctrl),
				mockTokenService:      NewMocktokenService(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: LoginContract{
					Email:    "test@example.com",
					Password: "wrongpassword",
				},
			},
			wantErr: errors.New("login failed"),
		},
		{
			name: "failed to generate token",
			fields: fields{
				mockAccountDBAccessor: NewMockaccountDBAccessor(ctrl),
				mockTokenService:      NewMocktokenService(ctrl),
			},
			args: args{
				ctx: context.Background(),
				spec: LoginContract{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			wantErr: errors.New("login failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			a := &AccountService{
				accountDBAccessor: tt.fields.mockAccountDBAccessor,
				tokenService:      tt.fields.mockTokenService,
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
			if err != nil {
				t.Fatalf("failed to hash password: %v", err)
			}

			if tt.name == "success" {
				// Mock successful account retrieval and password verification
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				account := &Account{
					ID:       "1",
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}
				tt.fields.mockAccountDBAccessor.EXPECT().
					FindAccountByEmail(tt.args.ctx, tt.args.spec.Email).
					Return(account, nil)

				// Mock successful token generation
				tt.fields.mockTokenService.EXPECT().
					GenerateToken(token.ClaimSpec{UserID: account.ID}).
					Return(tt.wantToken, nil)
			} else if tt.name == "account not found" {
				// Mock account not found error
				tt.fields.mockAccountDBAccessor.EXPECT().
					FindAccountByEmail(tt.args.ctx, tt.args.spec.Email).
					Return(nil, errors.New("account not found"))
			} else if tt.name == "invalid password" {
				// Mock successful account retrieval
				account := &Account{
					ID:       "1",
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}
				tt.fields.mockAccountDBAccessor.EXPECT().
					FindAccountByEmail(tt.args.ctx, tt.args.spec.Email).
					Return(account, nil)
			} else if tt.name == "failed to generate token" {
				// Mock successful account retrieval
				account := &Account{
					ID:       "1",
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}
				tt.fields.mockAccountDBAccessor.EXPECT().
					FindAccountByEmail(tt.args.ctx, tt.args.spec.Email).
					Return(account, nil)

				// Mock token generation error
				tt.fields.mockTokenService.EXPECT().
					GenerateToken(token.ClaimSpec{UserID: account.ID}).
					Return("", errors.New("token generation error"))
			}

			token, err := a.Login(tt.args.ctx, tt.args.spec)

			if tt.wantErr == nil {
				g.Expect(err).To(gomega.BeNil())
				g.Expect(token).To(gomega.Equal(tt.wantToken))
			} else {
				g.Expect(err).ToNot(gomega.BeNil())
				g.Expect(err.Error()).To(gomega.ContainSubstring(tt.wantErr.Error()))
			}
		})
	}
}
