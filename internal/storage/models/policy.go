package models

import (
	"database/sql"
	"math/big"
	"time"
)

type UnissuedPolicy struct {
	IOType           string
	OrderID          int64
	ID               int64
	Name             string
	Email            string
	PolicyNumber     sql.NullString
	OrderSum         big.Float
	PolicySum        big.Float
	ExternalID       sql.NullString
	PaymentID        string
	InsuranceCompany string
	DateTime         time.Time
}
