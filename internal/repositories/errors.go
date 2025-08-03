package repositories

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrInternal = errors.New("internal database error")
)

func extractPQError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}

func isInternalError(err error) bool {
	pgErr := extractPQError(err)
	if err != nil {
		return false
	}
	return pgErr.Code == pgerrcode.InternalError
}
