package loan

import (
	"amartha/config"
	modelAdmin "amartha/model/admin"
	modelApi "amartha/model/api"
	modelBorrower "amartha/model/borrower"
	modelInvestor "amartha/model/investor"
	modelLoan "amartha/model/loan"
	"context"
	"database/sql"
	"errors"
	"os"
	reflect "reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func Test_GetAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoLoan := NewMockRepoLoanInterface(ctrl)
	mockRepoUpload := NewMockRepoUploadInterface(ctrl)
	mockRepoNotification := NewMockRepoNotificationInterface(ctrl)

	ctx := context.Background()

	cfg := &config.Config{
		Settings: config.Settings{
			SecretKey: "amartha",
		},
	}

	type fields struct {
		repoLoan         RepoLoanInterface
		repoUpload       RepoUploadInterface
		repoNotification RepoNotificationInterface
	}
	type args struct {
		ctx context.Context
		req modelApi.GetAccessToken
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mockFn  func()
		wantErr bool
	}{
		{
			"Success get access token",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctx,
				req: modelApi.GetAccessToken{
					Email: "samsul@gmail",
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetAdmin(ctx, "samsul@gmail").Return(&modelAdmin.Admin{
					ID:   1,
					Name: "samsul",
				}, nil)
			},
			false,
		},
		{
			"Error get access token - get admin",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctx,
				req: modelApi.GetAccessToken{
					Email: "samsul@gmail",
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetAdmin(ctx, "samsul@gmail").Return(nil, errors.New("error"))
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			uc := &Usecase{
				ctx:              ctx,
				cfg:              cfg,
				repoLoan:         tt.fields.repoLoan,
				repoUpload:       tt.fields.repoUpload,
				repoNotification: tt.fields.repoNotification,
			}
			_, err := uc.GetAccessToken(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_LoanSubmit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoLoan := NewMockRepoLoanInterface(ctrl)
	mockRepoUpload := NewMockRepoUploadInterface(ctrl)
	mockRepoNotification := NewMockRepoNotificationInterface(ctrl)

	ctx := context.Background()

	cfg := &config.Config{
		Settings: config.Settings{
			SecretKey: "amartha",
		},
	}

	type fields struct {
		repoLoan         RepoLoanInterface
		repoUpload       RepoUploadInterface
		repoNotification RepoNotificationInterface
	}
	type args struct {
		ctx context.Context
		req modelApi.LoanSubmit
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mockFn   func()
		wantData modelApi.LoanResponse
		wantErr  bool
	}{
		{
			"Success submit loan",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctx,
				req: modelApi.LoanSubmit{
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetBorrower(ctx, int64(1)).Return(&modelBorrower.Borrower{
					ID: int64(1),
				}, nil)
				mockRepoLoan.EXPECT().InsertLoan(ctx, &modelLoan.Loan{
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Status:          modelLoan.STATUS_PROPOSED,
					Rate:            float64(2),
					Roi:             float64(3),
				}).Return(int64(1), nil)
			},
			modelApi.LoanResponse{
				LoanID:          int64(1),
				PrincipalAmount: int64(1000000),
				RateInPercent:   float64(2),
				RoiInPercent:    float64(3),
				Status:          modelLoan.STATUS_PROPOSED,
			},
			false,
		},
		{
			"Error get borrower",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctx,
				req: modelApi.LoanSubmit{
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetBorrower(ctx, int64(1)).Return(nil, errors.New("borrower not found"))
			},
			modelApi.LoanResponse{},
			true,
		},
		{
			"Error insert loan",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctx,
				req: modelApi.LoanSubmit{
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetBorrower(ctx, gomock.Any()).Return(&modelBorrower.Borrower{
					ID: int64(1),
				}, nil)
				mockRepoLoan.EXPECT().InsertLoan(ctx, gomock.Any()).Return(int64(0), errors.New("error insert loan"))
			},
			modelApi.LoanResponse{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			uc := &Usecase{
				ctx:              ctx,
				cfg:              cfg,
				repoLoan:         tt.fields.repoLoan,
				repoUpload:       tt.fields.repoUpload,
				repoNotification: tt.fields.repoNotification,
			}
			resp, err := uc.LoanSubmit(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.LoanSubmit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(resp, tt.wantData) {
				t.Errorf("Usecase.LoanSubmit() = %+v, want %+v", resp, tt.wantData)
			}
		})
	}
}

func Test_LoanApprovalt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoLoan := NewMockRepoLoanInterface(ctrl)
	mockRepoUpload := NewMockRepoUploadInterface(ctrl)
	mockRepoNotification := NewMockRepoNotificationInterface(ctrl)

	ctx := context.Background()
	ctxWithVal := context.WithValue(ctx, "adminID", "1")
	ctxWithVal = context.WithValue(ctxWithVal, "adminName", "Ana")

	cfg := &config.Config{
		Settings: config.Settings{
			SecretKey:           "amartha",
			AgreementLetterHtml: "agreementletter.html",
		},
	}

	templateContent := `
	<html>
		<head><title>{{.Title}}</title></head>
		<body>
			<p>Date: {{.Date}}</p>
			<p>Admin: {{.Admin}}</p>
			<p>Borrower Name: {{.BorrowerName}}</p>
			<p>Borrower Address: {{.BorrowerAddress}}</p>
			<p>Principal Amount: {{.PrincipalAmount}}</p>
			<p>Rate: {{.Rate}}</p>
		</body>
	</html>
	`

	templateFile := "agreementletter.html"
	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}
	defer func() {
		os.Remove(templateFile)
	}()

	type fields struct {
		repoLoan         RepoLoanInterface
		repoUpload       RepoUploadInterface
		repoNotification RepoNotificationInterface
	}
	type args struct {
		ctx      context.Context
		req      modelApi.LoanApproval
		fileName string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mockFn   func()
		wantData modelApi.LoanResponse
		wantErr  bool
	}{
		{
			"Success aprrove loan",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctxWithVal,
				req: modelApi.LoanApproval{
					LoanID:         int64(1),
					ValidatorPhoto: []byte("validator photo"),
				},
				fileName: "validator_photo",
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctxWithVal, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_PROPOSED,
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("https://example.com", nil)
				mockRepoLoan.EXPECT().GetBorrower(ctxWithVal, int64(1)).Return(&modelBorrower.Borrower{
					ID:      int64(1),
					Name:    "test",
					Address: "test",
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("https://example.com", nil)
				mockRepoLoan.EXPECT().UpdateApproval(ctxWithVal, gomock.Any()).Return(nil)
			},
			modelApi.LoanResponse{
				LoanID:          int64(1),
				PrincipalAmount: int64(1000000),
				RateInPercent:   float64(2),
				RoiInPercent:    float64(3),
				Status:          modelLoan.STATUS_APPROVED,
				AgreementLetter: "https://example.com",
			},
			false,
		},
		{
			"Error aprrove loan - getloan",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctxWithVal,
				req: modelApi.LoanApproval{
					LoanID:         int64(1),
					ValidatorPhoto: []byte("validator photo"),
				},
				fileName: "validator_photo",
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctxWithVal, gomock.Any()).Return(nil, errors.New("error"))
			},
			modelApi.LoanResponse{},
			true,
		},
		{
			"Error aprrove loan - upload file photo",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctxWithVal,
				req: modelApi.LoanApproval{
					LoanID:         int64(1),
					ValidatorPhoto: []byte("validator photo"),
				},
				fileName: "validator_photo",
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctxWithVal, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_PROPOSED,
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("", errors.New("error"))
			},
			modelApi.LoanResponse{},
			true,
		},
		{
			"Error aprrove loan - get borrower",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctxWithVal,
				req: modelApi.LoanApproval{
					LoanID:         int64(1),
					ValidatorPhoto: []byte("validator photo"),
				},
				fileName: "validator_photo",
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctxWithVal, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_PROPOSED,
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("https://example.com", nil)
				mockRepoLoan.EXPECT().GetBorrower(ctxWithVal, int64(1)).Return(nil, errors.New("error"))
			},
			modelApi.LoanResponse{},
			true,
		},
		{
			"Error aprrove loan - upload file agreement letter",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctxWithVal,
				req: modelApi.LoanApproval{
					LoanID:         int64(1),
					ValidatorPhoto: []byte("validator photo"),
				},
				fileName: "validator_photo",
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctxWithVal, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_PROPOSED,
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("https://example.com", nil)
				mockRepoLoan.EXPECT().GetBorrower(ctxWithVal, int64(1)).Return(&modelBorrower.Borrower{
					ID:      int64(1),
					Name:    "test",
					Address: "test",
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("", errors.New("error"))
			},
			modelApi.LoanResponse{},
			true,
		},
		{
			"Error aprrove loan - update approved",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctxWithVal,
				req: modelApi.LoanApproval{
					LoanID:         int64(1),
					ValidatorPhoto: []byte("validator photo"),
				},
				fileName: "validator_photo",
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctxWithVal, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_PROPOSED,
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("https://example.com", nil)
				mockRepoLoan.EXPECT().GetBorrower(ctxWithVal, int64(1)).Return(&modelBorrower.Borrower{
					ID:      int64(1),
					Name:    "test",
					Address: "test",
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("https://example.com", nil)
				mockRepoLoan.EXPECT().UpdateApproval(ctxWithVal, gomock.Any()).Return(errors.New("error"))
			},
			modelApi.LoanResponse{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			uc := &Usecase{
				ctx:              tt.args.ctx,
				cfg:              cfg,
				repoLoan:         tt.fields.repoLoan,
				repoUpload:       tt.fields.repoUpload,
				repoNotification: tt.fields.repoNotification,
			}
			resp, err := uc.LoanApproval(tt.args.ctx, tt.args.req, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.LoanApproval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(resp, tt.wantData) {
				t.Errorf("Usecase.LoanApproval() = %+v, want %+v", resp, tt.wantData)
			}
		})
	}
}

func Test_LoanInvestment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoLoan := NewMockRepoLoanInterface(ctrl)
	mockRepoUpload := NewMockRepoUploadInterface(ctrl)
	mockRepoNotification := NewMockRepoNotificationInterface(ctrl)

	ctx := context.Background()
	cfg := &config.Config{
		Settings: config.Settings{
			SecretKey: "amartha",
		},
	}

	type fields struct {
		repoLoan         RepoLoanInterface
		repoUpload       RepoUploadInterface
		repoNotification RepoNotificationInterface
	}
	type args struct {
		ctx      context.Context
		req      modelApi.LoanInvestment
		fileName string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mockFn   func()
		wantData modelApi.LoanResponse
		wantErr  bool
	}{
		{
			"Success submit investment",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctx,
				req: modelApi.LoanInvestment{
					LoanID:         int64(1),
					InvestorID:     int64(1),
					InvestedAmount: int64(1000000),
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctx, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_APPROVED,
					AgreementLetter: sql.NullString{
						String: "https://example.com",
						Valid:  true,
					},
				}, nil)
				mockRepoLoan.EXPECT().GetInvestor(ctx, gomock.Any()).Return(&modelInvestor.Investor{
					ID:    int64(1),
					Name:  "test",
					Email: "test@gmail.com",
				}, nil)
				mockRepoLoan.EXPECT().GetLoanDetail(ctx, gomock.Any()).Return([]*modelLoan.LoanDetail{}, nil)
				mockRepoLoan.EXPECT().InsertLoanDetail(ctx, gomock.Any()).Return(nil)
				mockRepoLoan.EXPECT().GetAllInvestorEmail(ctx, gomock.Any()).Return([]*modelInvestor.Investor{
					{
						Email: "test@gmail.com",
					},
				}, nil)
				mockRepoNotification.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepoLoan.EXPECT().UpdateInvested(ctx, gomock.Any()).Return(nil)
			},
			modelApi.LoanResponse{
				LoanID:          int64(1),
				PrincipalAmount: int64(1000000),
				RateInPercent:   float64(2),
				RoiInPercent:    float64(3),
				Status:          modelLoan.STATUS_INVESTED,
				AgreementLetter: "https://example.com",
			},
			false,
		},
		{
			"Error submit investment - total invested amount is greater than principal amount",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctx,
				req: modelApi.LoanInvestment{
					LoanID:         int64(1),
					InvestorID:     int64(1),
					InvestedAmount: int64(1000000),
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctx, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_APPROVED,
					AgreementLetter: sql.NullString{
						String: "https://example.com",
						Valid:  true,
					},
				}, nil)
				mockRepoLoan.EXPECT().GetInvestor(ctx, gomock.Any()).Return(&modelInvestor.Investor{
					ID:    int64(1),
					Name:  "test",
					Email: "test@gmail.com",
				}, nil)
				mockRepoLoan.EXPECT().GetLoanDetail(ctx, gomock.Any()).Return([]*modelLoan.LoanDetail{
					{
						LoanID:         int64(1),
						InvestorID:     int64(1),
						InvestedAmount: int64(500000),
					},
					{
						LoanID:         int64(1),
						InvestorID:     int64(2),
						InvestedAmount: int64(300000),
					},
				}, nil)
			},
			modelApi.LoanResponse{},
			true,
		},
		{
			"Error submit investment with empty loan detail - total invested amount is greater than principal amount",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctx,
				req: modelApi.LoanInvestment{
					LoanID:         int64(1),
					InvestorID:     int64(1),
					InvestedAmount: int64(2000000),
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctx, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_APPROVED,
					AgreementLetter: sql.NullString{
						String: "https://example.com",
						Valid:  true,
					},
				}, nil)
				mockRepoLoan.EXPECT().GetInvestor(ctx, gomock.Any()).Return(&modelInvestor.Investor{
					ID:    int64(1),
					Name:  "test",
					Email: "test@gmail.com",
				}, nil)
				mockRepoLoan.EXPECT().GetLoanDetail(ctx, gomock.Any()).Return([]*modelLoan.LoanDetail{}, nil)
			},
			modelApi.LoanResponse{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			uc := &Usecase{
				ctx:              tt.args.ctx,
				cfg:              cfg,
				repoLoan:         tt.fields.repoLoan,
				repoUpload:       tt.fields.repoUpload,
				repoNotification: tt.fields.repoNotification,
			}
			resp, err := uc.LoanInvestment(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.LoanInvestment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(resp, tt.wantData) {
				t.Errorf("Usecase.LoanInvestment() = %+v, want %+v", resp, tt.wantData)
			}
		})
	}
}

func Test_LoanDisbursement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoLoan := NewMockRepoLoanInterface(ctrl)
	mockRepoUpload := NewMockRepoUploadInterface(ctrl)
	mockRepoNotification := NewMockRepoNotificationInterface(ctrl)

	ctx := context.Background()
	ctxWithVal := context.WithValue(ctx, "adminID", "1")
	ctxWithVal = context.WithValue(ctxWithVal, "adminName", "Ana")

	cfg := &config.Config{
		Settings: config.Settings{
			SecretKey:           "amartha",
			AgreementLetterHtml: "agreementletter.html",
		},
	}

	type fields struct {
		repoLoan         RepoLoanInterface
		repoUpload       RepoUploadInterface
		repoNotification RepoNotificationInterface
	}
	type args struct {
		ctx      context.Context
		req      modelApi.LoanDisbursement
		fileName string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mockFn   func()
		wantData modelApi.LoanResponse
		wantErr  bool
	}{
		{
			"Success loan disbursement",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctxWithVal,
				req: modelApi.LoanDisbursement{
					LoanID:          int64(1),
					AgreementLetter: []byte("test"),
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctxWithVal, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_INVESTED,
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("https://example.com", nil)
				mockRepoLoan.EXPECT().UpdateDisbursed(ctxWithVal, gomock.Any()).Return(nil)
			},
			modelApi.LoanResponse{
				LoanID:          int64(1),
				PrincipalAmount: int64(1000000),
				RateInPercent:   float64(2),
				RoiInPercent:    float64(3),
				Status:          modelLoan.STATUS_DISBURSED,
				AgreementLetter: "https://example.com",
			},
			false,
		},
		{
			"Error loan disbursement - status not valid",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctxWithVal,
				req: modelApi.LoanDisbursement{
					LoanID:          int64(1),
					AgreementLetter: []byte("test"),
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctxWithVal, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_APPROVED,
				}, nil)
			},
			modelApi.LoanResponse{},
			true,
		},
		{
			"Error loan disbursement - error upload file",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctxWithVal,
				req: modelApi.LoanDisbursement{
					LoanID:          int64(1),
					AgreementLetter: []byte("test"),
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctxWithVal, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_INVESTED,
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("", errors.New("error"))
			},
			modelApi.LoanResponse{},
			true,
		},
		{
			"Error loan disbursement - error update disbursed",
			fields{
				repoLoan:         mockRepoLoan,
				repoUpload:       mockRepoUpload,
				repoNotification: mockRepoNotification,
			},
			args{
				ctx: ctxWithVal,
				req: modelApi.LoanDisbursement{
					LoanID:          int64(1),
					AgreementLetter: []byte("test"),
				},
			},
			func() {
				mockRepoLoan.EXPECT().GetLoan(ctxWithVal, gomock.Any()).Return(&modelLoan.Loan{
					ID:              int64(1),
					BorrowerID:      int64(1),
					PrincipalAmount: int64(1000000),
					Rate:            float64(2),
					Roi:             float64(3),
					Status:          modelLoan.STATUS_INVESTED,
				}, nil)
				mockRepoUpload.EXPECT().UploadFile(ctxWithVal, gomock.Any(), gomock.Any()).Return("https://example.com", nil)
				mockRepoLoan.EXPECT().UpdateDisbursed(ctxWithVal, gomock.Any()).Return(errors.New("error"))
			},
			modelApi.LoanResponse{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			uc := &Usecase{
				ctx:              tt.args.ctx,
				cfg:              cfg,
				repoLoan:         tt.fields.repoLoan,
				repoUpload:       tt.fields.repoUpload,
				repoNotification: tt.fields.repoNotification,
			}
			resp, err := uc.LoanDisbursement(tt.args.ctx, tt.args.req, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.LoanDisbursement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(resp, tt.wantData) {
				t.Errorf("Usecase.LoanDisbursement() = %+v, want %+v", resp, tt.wantData)
			}
		})
	}
}
