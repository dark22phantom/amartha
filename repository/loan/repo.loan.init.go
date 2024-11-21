package loan

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository struct {
	db   *sqlx.DB
	stmt *statementQuery
}

func New(
	db *sqlx.DB,
) (*Repository, error) {
	repo := &Repository{
		db: db,
	}

	if err := repo.prepareStatements(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) prepareStatements() error {
	statements := &statementQuery{}
	var err error

	statements.qryGetAdmin, err = r.db.Preparex(sqlGetAdmin)
	if err != nil {
		return errors.Wrap(err, "error prepare statement query sqlGetAdmin")
	}

	statements.qryGetBorrower, err = r.db.Preparex(sqlGetBorrower)
	if err != nil {
		return errors.Wrap(err, "error prepare statement query sqlGetBorrower")
	}

	statements.qryGetLoan, err = r.db.Preparex(sqlGetLoan)
	if err != nil {
		return errors.Wrap(err, "error prepare statement query sqlGetLoan")
	}

	statements.qryGetLoanDetail, err = r.db.Preparex(sqlGetLoanDetail)
	if err != nil {
		return errors.Wrap(err, "error prepare statement query sqlGetLoanDetail")
	}

	statements.qryGetInvestor, err = r.db.Preparex(sqlGetInvestor)
	if err != nil {
		return errors.Wrap(err, "error prepare statement query sqlGetInvestor")
	}

	statements.qryGetAllInvestorEmail, err = r.db.Preparex(sqlGetAllInvestorEmail)
	if err != nil {
		return errors.Wrap(err, "error prepare statement query sqlGetInvestorEmail")
	}

	r.stmt = statements
	return err
}
