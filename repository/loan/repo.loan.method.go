package loan

import (
	modelAdmin "amartha/model/admin"
	modelBorrower "amartha/model/borrower"
	modelInvestor "amartha/model/investor"
	modelLoan "amartha/model/loan"
	"context"
	"database/sql"
	"errors"
)

func (r *Repository) GetAdmin(ctx context.Context, email string) (*modelAdmin.Admin, error) {
	admin := &modelAdmin.Admin{}
	if err := r.stmt.qryGetAdmin.GetContext(ctx, admin, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}
	return admin, nil
}

func (r *Repository) GetBorrower(ctx context.Context, borrowerID int64) (*modelBorrower.Borrower, error) {
	borrower := &modelBorrower.Borrower{}
	if err := r.stmt.qryGetBorrower.GetContext(ctx, borrower, borrowerID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("borrower not found")
		}
		return nil, err
	}
	return borrower, nil
}

func (r *Repository) InsertLoan(ctx context.Context, loan *modelLoan.Loan) (int64, error) {
	var insertedID int64

	tx, err := r.db.Beginx()
	if err != nil {
		return insertedID, err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(sqlInsertLoan)
	if err != nil {
		tx.Rollback()
		return insertedID, err
	}

	err = stmt.QueryRowContext(ctx,
		loan.BorrowerID,
		loan.PrincipalAmount,
		loan.Status,
		loan.Rate,
		loan.Roi,
	).Scan(&insertedID)
	if err != nil {
		tx.Rollback()
		return insertedID, err
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return insertedID, err
		}
		return insertedID, err
	}

	return insertedID, nil
}

func (r *Repository) GetLoan(ctx context.Context, loanID int64) (*modelLoan.Loan, error) {
	loan := &modelLoan.Loan{}
	if err := r.stmt.qryGetLoan.GetContext(ctx, loan, loanID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("loan not found")
		}
		return nil, err
	}
	return loan, nil
}

func (r *Repository) UpdateApproval(ctx context.Context, loan *modelLoan.Loan) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(sqlUpdateApproval)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.ExecContext(ctx,
		loan.Status,
		loan.ApprovedBy,
		loan.ApprovedAt,
		loan.AgreementLetter,
		loan.FieldValidator,
		loan.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (r *Repository) GetLoanDetail(ctx context.Context, loanID int64) ([]*modelLoan.LoanDetail, error) {
	loanDetail := []*modelLoan.LoanDetail{}
	if err := r.stmt.qryGetLoanDetail.SelectContext(ctx, &loanDetail, loanID); err != nil {
		return nil, err
	}
	return loanDetail, nil
}

func (r *Repository) InsertLoanDetail(ctx context.Context, loanDetail *modelLoan.LoanDetail) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(sqlInsertLoanDetail)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.ExecContext(ctx,
		loanDetail.LoanID,
		loanDetail.InvestorID,
		loanDetail.InvestedAmount,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (r *Repository) GetInvestor(ctx context.Context, investorID int64) (*modelInvestor.Investor, error) {
	investor := &modelInvestor.Investor{}
	if err := r.stmt.qryGetInvestor.GetContext(ctx, investor, investorID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("investor not found")
		}
		return nil, err
	}
	return investor, nil
}

func (r *Repository) GetAllInvestorEmail(ctx context.Context, loanID int64) ([]*modelInvestor.Investor, error) {
	email := []*modelInvestor.Investor{}
	if err := r.stmt.qryGetAllInvestorEmail.SelectContext(ctx, &email, loanID); err != nil {
		return nil, err
	}
	return email, nil
}

func (r *Repository) UpdateInvested(ctx context.Context, loan *modelLoan.Loan) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(sqlUpdateInvested)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.ExecContext(ctx,
		loan.Status,
		loan.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (r *Repository) UpdateDisbursed(ctx context.Context, loan *modelLoan.Loan) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(sqlUpdateDisbursed)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.ExecContext(ctx,
		loan.Status,
		loan.AgreementLetter,
		loan.DisbursedBy,
		loan.DisbursedAt,
		loan.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return nil
}
