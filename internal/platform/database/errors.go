package database

import (
	"errors"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrDuplicateKey        = errors.New("unique key violation")
	ErrInternal            = errors.New("internal database error")
)
