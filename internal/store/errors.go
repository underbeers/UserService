package store

import "git.friends.com/PetLand/UserService/v2/internal/genErr"

var (
	ErrRecordNotFound            = genErr.New("record not found")
	ErrTransactionCreationFailed = genErr.New("transaction creation failed")
	ErrTransactionFailed         = genErr.New("transaction failed, rollback")
	ErrTransactionRollbackFailed = genErr.New("transaction rollback failed")
	ErrTransactionCommitFailed   = genErr.New("transaction commit failed")
	ErrScanStructFailed          = genErr.New("failed to Scan structure")
	ErrNotFound                  = genErr.New("not found ")
)
