package loan

import (
	modelAdmin "amartha/model/admin"
	modelBorrower "amartha/model/borrower"
	modelInvestor "amartha/model/investor"
	modelLoan "amartha/model/loan"
	"context"
)

type RepoLoanInterface interface {
	GetAdmin(ctx context.Context, email string) (*modelAdmin.Admin, error)
	GetBorrower(ctx context.Context, borrowerID int64) (*modelBorrower.Borrower, error)
	InsertLoan(ctx context.Context, loan *modelLoan.Loan) (int64, error)
	GetLoan(ctx context.Context, loanID int64) (*modelLoan.Loan, error)
	UpdateApproval(ctx context.Context, loan *modelLoan.Loan) error
	GetLoanDetail(ctx context.Context, loanID int64) ([]*modelLoan.LoanDetail, error)
	GetInvestor(ctx context.Context, investorID int64) (*modelInvestor.Investor, error)
	InsertLoanDetail(ctx context.Context, loanDetail *modelLoan.LoanDetail) error
	GetAllInvestorEmail(ctx context.Context, loanID int64) ([]*modelInvestor.Investor, error)
	UpdateInvested(ctx context.Context, loan *modelLoan.Loan) error
	UpdateDisbursed(ctx context.Context, loan *modelLoan.Loan) error
}

type RepoUploadInterface interface {
	UploadFile(ctx context.Context, file []byte, path string) (string, error)
}

type RepoNotificationInterface interface {
	SendEmail(email string, subject string, message string) error
}
