package loan

import "github.com/jmoiron/sqlx"

const (
	sqlGetAdmin = `
		SELECT 
			id,
			name
		FROM
			admin
		WHERE
			email = $1
			AND deleted_at IS NULL
	`
	sqlGetBorrower = `
		SELECT 
			id,
			name,
			email,
			phone,
			address
		FROM
			borrower
		WHERE
			id = $1
			AND deleted_at IS NULL
	`
	sqlInsertLoan = `
		INSERT INTO
			loan (
				borrower_id,
				principal_amount,
				status,
				rate,
				roi
			)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id
	`
	sqlGetLoan = `
		SELECT 
			id,
			borrower_id,
			principal_amount,
			status,
			rate,
			roi,
			agreement_letter
		FROM
			loan
		WHERE
			id = $1
	`
	sqlUpdateApproval = `
		UPDATE
			loan
		SET
			status = $1,
			approved_by = $2,
			approved_at = $3,
			agreement_letter = $4,
			field_validator = $5,
			updated_at = now()
		WHERE
			id = $6
	`
	sqlGetLoanDetail = `
		SELECT
			id,
			loan_id,
			investor_id,
			invested_amount
		FROM
			loan_detail
		WHERE
			loan_id = $1
	`
	sqlGetInvestor = `
		SELECT
			id,
			name,
			email,
			phone
		FROM
			investor
		WHERE
			id = $1
			AND deleted_at IS NULL
	`
	sqlInsertLoanDetail = `
		INSERT INTO
			loan_detail (
				loan_id,
				investor_id,
				invested_amount
			)
		VALUES
			($1, $2, $3)
	`
	sqlGetAllInvestorEmail = `
		SELECT
			i.email 
		FROM
			loan_detail ld 
			LEFT JOIN investor i ON ld.investor_id = i.id
		WHERE
			ld.loan_id = $1
			AND i.deleted_at IS NULL
		GROUP BY i.id
	`
	sqlUpdateInvested = `
		UPDATE
			loan
		SET
			status = $1,
			updated_at = now()
		WHERE
			id = $2
	`
	sqlUpdateDisbursed = `
		UPDATE
			loan
		SET
			status = $1,
			agreement_letter = $2,
			disbursed_by = $3,
			disbursed_at = $4,
			updated_at = now()
		WHERE
			id = $5
	`
)

type statementQuery struct {
	qryGetAdmin            *sqlx.Stmt
	qryGetBorrower         *sqlx.Stmt
	qryGetLoan             *sqlx.Stmt
	qryGetLoanDetail       *sqlx.Stmt
	qryGetInvestor         *sqlx.Stmt
	qryGetAllInvestorEmail *sqlx.Stmt
}
