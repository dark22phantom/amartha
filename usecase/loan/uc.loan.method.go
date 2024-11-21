package loan

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"amartha/internal/pkg/helper"
	modelApi "amartha/model/api"
	modelBorrower "amartha/model/borrower"
	modelLoan "amartha/model/loan"
	modelUpload "amartha/model/upload"
)

func (uc *Usecase) GetAccessToken(ctx context.Context, req modelApi.GetAccessToken) (modelApi.GetAccessTokenResponse, error) {
	resp := modelApi.GetAccessTokenResponse{}

	// get admin by email
	admin, err := uc.repoLoan.GetAdmin(ctx, req.Email)
	if err != nil {
		return resp, err
	}

	// generate access token by adminID
	secretKey := uc.cfg.Settings.SecretKey
	adminID := strconv.Itoa(int(admin.ID))
	accessToken, err := helper.GenerateAccessToken(secretKey, adminID, admin.Name)
	if err != nil {
		return resp, err
	}

	resp.AccessToken = accessToken
	return resp, nil
}

func (uc *Usecase) getBorrower(ctx context.Context, borrowerID int64) (*modelBorrower.Borrower, error) {
	borrower, err := uc.repoLoan.GetBorrower(ctx, borrowerID)
	if err != nil {
		return nil, err
	}
	return borrower, nil
}

func (uc *Usecase) insertLoan(ctx context.Context, req modelApi.LoanSubmit, rate, roi float64) (int64, error) {
	loan := &modelLoan.Loan{
		BorrowerID:      req.BorrowerID,
		PrincipalAmount: req.PrincipalAmount,
		Status:          modelLoan.STATUS_PROPOSED,
		Rate:            rate,
		Roi:             roi,
	}
	loanID, err := uc.repoLoan.InsertLoan(ctx, loan)
	if err != nil {
		return loanID, err
	}
	return loanID, nil
}

func (uc *Usecase) LoanSubmit(ctx context.Context, req modelApi.LoanSubmit) (modelApi.LoanResponse, error) {
	resp := modelApi.LoanResponse{}

	// check borrowerID
	_, err := uc.getBorrower(ctx, req.BorrowerID)
	if err != nil {
		return resp, err
	}

	// insert loan
	rate := float64(2)
	roi := float64(3)
	loanID, err := uc.insertLoan(ctx, req, rate, roi)
	if err != nil {
		return resp, err
	}

	resp = modelApi.LoanResponse{
		LoanID:          loanID,
		PrincipalAmount: req.PrincipalAmount,
		RateInPercent:   rate,
		RoiInPercent:    roi,
		Status:          modelLoan.STATUS_PROPOSED,
	}
	return resp, nil
}

func (uc *Usecase) validateLoan(ctx context.Context, loanID int64, nextStatus string) (bool, *modelLoan.Loan, error) {
	loan, err := uc.repoLoan.GetLoan(ctx, loanID)
	if err != nil {
		return false, loan, err
	}

	// mapping status
	status := map[string]string{
		modelLoan.STATUS_APPROVED:  modelLoan.STATUS_PROPOSED,
		modelLoan.STATUS_INVESTED:  modelLoan.STATUS_APPROVED,
		modelLoan.STATUS_DISBURSED: modelLoan.STATUS_INVESTED,
	}

	if value, ok := status[nextStatus]; ok {
		if value != loan.Status {
			return false, loan, nil
		}
	}
	return true, loan, nil
}

func (uc *Usecase) uploadFile(ctx context.Context, file []byte, path string) (string, error) {
	// upload file
	link, err := uc.repoUpload.UploadFile(ctx, file, path)
	if err != nil {
		return "", err
	}
	return link, nil
}

func (uc *Usecase) LoanApproval(ctx context.Context, req modelApi.LoanApproval, fileName string) (modelApi.LoanResponse, error) {
	resp := modelApi.LoanResponse{}

	// get context value
	adminID := ctx.Value("adminID").(string)
	adminIDint, _ := strconv.ParseInt(adminID, 10, 64)
	adminName := ctx.Value("adminName").(string)

	// validate loan
	isValid, loan, err := uc.validateLoan(ctx, req.LoanID, modelLoan.STATUS_APPROVED)
	if err != nil {
		return resp, err
	}
	if !isValid {
		return resp, errors.New("loan status is not valid")
	}

	// upload validator photo
	validatorPhotoPath := modelUpload.VALIDATOR_PHOTO_FOLDER + strconv.Itoa(int(req.LoanID)) + "-" + fileName
	validatorPhotoLink, err := uc.uploadFile(ctx, req.ValidatorPhoto, validatorPhotoPath)
	if err != nil {
		return resp, err
	}

	borrower, err := uc.getBorrower(ctx, loan.BorrowerID)
	if err != nil {
		return resp, err
	}

	// generate agreement letter
	dataAgreementLetter := helper.TemplateDataAgreementLetter{
		Title:           "Agreement Letter",
		Date:            time.Now().Format("02-01-2006"),
		Admin:           adminName,
		BorrowerName:    borrower.Name,
		BorrowerAddress: borrower.Address,
		PrincipalAmount: strconv.Itoa(int(loan.PrincipalAmount)),
		Rate:            strconv.Itoa(int(loan.Rate)),
	}
	templatePath := uc.cfg.Settings.AgreementLetterHtml
	agreementLetter, err := dataAgreementLetter.GenerateAgreementLetter(templatePath)
	if err != nil {
		return resp, err
	}
	agreementLetterPath := modelUpload.AGREEMENT_LETTER_FOLDER + strconv.Itoa(int(req.LoanID)) + "-agreement-letter.pdf"
	agreementLetterLink, err := uc.uploadFile(ctx, agreementLetter, agreementLetterPath)
	if err != nil {
		return resp, err
	}

	// update loan
	loan.Status = modelLoan.STATUS_APPROVED
	loan.ApprovedBy = sql.NullInt64{Int64: adminIDint, Valid: true}
	loan.ApprovedAt = sql.NullTime{Time: time.Now(), Valid: true}
	loan.AgreementLetter = sql.NullString{String: agreementLetterLink, Valid: true}
	loan.FieldValidator = sql.NullString{String: validatorPhotoLink, Valid: true}
	if err := uc.repoLoan.UpdateApproval(ctx, loan); err != nil {
		return resp, err
	}

	resp = modelApi.LoanResponse{
		LoanID:          req.LoanID,
		PrincipalAmount: loan.PrincipalAmount,
		RateInPercent:   loan.Rate,
		RoiInPercent:    loan.Roi,
		Status:          modelLoan.STATUS_APPROVED,
		AgreementLetter: agreementLetterLink,
	}
	return resp, nil
}

func (uc *Usecase) calculateTotalInvestment(ctx context.Context, loanID, investedAmount int64) (int64, error) {
	totalInvestment := investedAmount
	loanDetail, err := uc.repoLoan.GetLoanDetail(ctx, loanID)
	if err != nil {
		return 0, err
	}

	for _, detail := range loanDetail {
		totalInvestment += detail.InvestedAmount
	}
	return totalInvestment, nil
}

func (uc *Usecase) sendEmailNotification(ctx context.Context, loan *modelLoan.Loan) error {
	// get all investor email
	emails, err := uc.repoLoan.GetAllInvestorEmail(ctx, loan.ID)
	if err != nil {
		return err
	}

	// send email notif to all investor
	subject := "Loan Invested"
	message := "Your loan has been invested with agreement letter detail: " + loan.AgreementLetter.String
	for _, email := range emails {
		if err := uc.repoNotification.SendEmail(email.Email, subject, message); err != nil {
			return err
		}
	}
	return nil
}

func (uc *Usecase) validateInvestedAmount(ctx context.Context, loan *modelLoan.Loan, investedAmount, investorID int64) (*modelLoan.Loan, error) {
	totalInvestment := investedAmount

	// calculate total invested amount
	loanDetail, err := uc.repoLoan.GetLoanDetail(ctx, loan.ID)
	if err != nil {
		return loan, err
	}
	for _, detail := range loanDetail {
		totalInvestment += detail.InvestedAmount
	}

	// validate total invested amount
	if totalInvestment > loan.PrincipalAmount {
		return loan, errors.New("total invested amount is greater than principal amount")
	}

	// insert loan detail
	newDetail := &modelLoan.LoanDetail{
		LoanID:         loan.ID,
		InvestorID:     investorID,
		InvestedAmount: investedAmount,
	}
	if err := uc.repoLoan.InsertLoanDetail(ctx, newDetail); err != nil {
		return loan, err
	}

	// set status to invested & send email notif to all investor
	if totalInvestment == loan.PrincipalAmount {
		loan.Status = modelLoan.STATUS_INVESTED
		if err := uc.sendEmailNotification(ctx, loan); err != nil {
			return loan, err
		}
	}
	return loan, nil
}

func (uc *Usecase) LoanInvestment(ctx context.Context, req modelApi.LoanInvestment) (modelApi.LoanResponse, error) {
	resp := modelApi.LoanResponse{}

	// validate loan
	isValid, loan, err := uc.validateLoan(ctx, req.LoanID, modelLoan.STATUS_INVESTED)
	if err != nil {
		return resp, err
	}
	if !isValid {
		return resp, errors.New("loan status is not valid")
	}

	// validate investor
	if _, err := uc.repoLoan.GetInvestor(ctx, req.InvestorID); err != nil {
		return resp, err
	}

	// validate invested amount
	loan, err = uc.validateInvestedAmount(ctx, loan, req.InvestedAmount, req.InvestorID)
	if err != nil {
		return resp, err
	}

	// update loan
	if err := uc.repoLoan.UpdateInvested(ctx, loan); err != nil {
		return resp, err
	}

	resp = modelApi.LoanResponse{
		LoanID:          req.LoanID,
		PrincipalAmount: loan.PrincipalAmount,
		RateInPercent:   loan.Rate,
		RoiInPercent:    loan.Roi,
		Status:          loan.Status,
		AgreementLetter: loan.AgreementLetter.String,
	}

	return resp, nil
}

func (uc *Usecase) LoanDisbursement(ctx context.Context, req modelApi.LoanDisbursement, fileName string) (modelApi.LoanResponse, error) {
	resp := modelApi.LoanResponse{}

	// get context value
	adminID := ctx.Value("adminID").(string)
	adminIDint, _ := strconv.ParseInt(adminID, 10, 64)

	// validate loan
	isValid, loan, err := uc.validateLoan(ctx, req.LoanID, modelLoan.STATUS_DISBURSED)
	if err != nil {
		return resp, err
	}
	if !isValid {
		return resp, errors.New("loan status is not valid")
	}

	// upload signed agreement letter
	agreementLetterPath := modelUpload.AGREEMENT_LETTER_FOLDER + strconv.Itoa(int(req.LoanID)) + "-agreement-letter-signed.pdf"
	agreementLetterLink, err := uc.uploadFile(ctx, req.AgreementLetter, agreementLetterPath)
	if err != nil {
		return resp, err
	}

	// update loan
	loan.Status = modelLoan.STATUS_DISBURSED
	loan.AgreementLetter = sql.NullString{String: agreementLetterLink, Valid: true}
	loan.DisbursedBy = sql.NullInt64{Int64: adminIDint, Valid: true}
	loan.DisbursedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := uc.repoLoan.UpdateDisbursed(ctx, loan); err != nil {
		return resp, err
	}

	resp = modelApi.LoanResponse{
		LoanID:          req.LoanID,
		PrincipalAmount: loan.PrincipalAmount,
		RateInPercent:   loan.Rate,
		RoiInPercent:    loan.Roi,
		Status:          modelLoan.STATUS_DISBURSED,
		AgreementLetter: agreementLetterLink,
	}
	return resp, nil
}
