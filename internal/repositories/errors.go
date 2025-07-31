package repositories

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound           = errors.New("record not found")
	ErrAuthorIdNotFound   = errors.New("provided author doesn't exist")
	ErrMemberIdNotFound   = errors.New("provided member doesn't exist")
	ErrBookIdNotFound     = errors.New("provided book doesn't exist")
	ErrInvalidPhoneNumber = errors.New("phone number is invalid")
	ErrInternal           = errors.New("internal database error")
)

func extractPQError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}
