package models

import "time"

type Verification struct {
	ID               int       `db:"id"`
	ContactID        int       `db:"id_contacts"`
	VerificationCode string    `db:"sms_code_verification"`
	CodeExpiration   time.Time `db:"sms_code_expiration"`
	BlockExpiration  time.Time `db:"block_expiration"`
	WrongAttempts    int       `db:"wrong_attempts"`
}
