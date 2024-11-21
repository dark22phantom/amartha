package http

import (
	modelApi "amartha/model/api"
	"context"
)

type ucLoanInterface interface {
	GetAccessToken(ctx context.Context, reqBody modelApi.GetAccessToken) (modelApi.GetAccessTokenResponse, error)
	LoanSubmit(ctx context.Context, req modelApi.LoanSubmit) (modelApi.LoanResponse, error)
	LoanApproval(ctx context.Context, req modelApi.LoanApproval, fileName string) (modelApi.LoanResponse, error)
	LoanInvestment(ctx context.Context, req modelApi.LoanInvestment) (modelApi.LoanResponse, error)
	LoanDisbursement(ctx context.Context, req modelApi.LoanDisbursement, fileName string) (modelApi.LoanResponse, error)
}
