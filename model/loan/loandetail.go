package loan

import (
	"database/sql"
	"time"
)

type LoanDetail struct {
	ID             int64        `db:"id"`
	LoanID         int64        `db:"loan_id"`
	InvestorID     int64        `db:"investor_id"`
	InvestedAmount int64        `db:"invested_amount"`
	Roi            int64        `db:"roi"`
	CreatedAt      time.Time    `db:"created_at"`
	UpdatedAt      sql.NullTime `db:"updated_at"`
	DeletedAt      sql.NullTime `db:"deleted_at"`
}
