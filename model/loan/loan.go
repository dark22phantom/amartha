package loan

import (
	"database/sql"
	"time"
)

const (
	STATUS_PROPOSED  = "proposed"
	STATUS_APPROVED  = "approved"
	STATUS_INVESTED  = "invested"
	STATUS_DISBURSED = "disbursed"
)

type Loan struct {
	ID              int64          `db:"id"`
	BorrowerID      int64          `db:"borrower_id"`
	PrincipalAmount int64          `db:"principal_amount"`
	Status          string         `db:"status"`
	Rate            float64        `db:"rate"`
	Roi             float64        `db:"roi"`
	ApprovedBy      sql.NullInt64  `db:"approved_by"`
	ApprovedAt      sql.NullTime   `db:"approved_at"`
	AgreementLetter sql.NullString `db:"agreement_letter"`
	FieldValidator  sql.NullString `db:"field_validator"`
	DisbursedBy     sql.NullInt64  `db:"disbursed_by"`
	DisbursedAt     sql.NullTime   `db:"disbursed_at"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       sql.NullTime   `db:"updated_at"`
	DeletedAt       sql.NullTime   `db:"deleted_at"`
}
