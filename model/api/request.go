package api

import "github.com/go-playground/validator/v10"

type GetAccessToken struct {
	Email string `json:"email" validate:"required"`
}

type LoanSubmit struct {
	BorrowerID      int64 `json:"borrower_id" validate:"required"`
	PrincipalAmount int64 `json:"principal_amount" validate:"required"`
}

type LoanApproval struct {
	LoanID         int64  `json:"loan_id"`
	ValidatorPhoto []byte `json:"validator_photo"`
}

type LoanInvestment struct {
	LoanID         int64 `json:"loan_id" validate:"required"`
	InvestorID     int64 `json:"investor_id" validate:"required"`
	InvestedAmount int64 `json:"invested_amount" validate:"required"`
}

type LoanDisbursement struct {
	LoanID          int64  `json:"loan_id"`
	AgreementLetter []byte `json:"agreement_letter"`
}

var Validate = validator.New()
