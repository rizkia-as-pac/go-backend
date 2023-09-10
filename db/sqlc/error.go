package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Unlike lib/pq, pgx doesn't have a map from the error code number to its name, but it only returns the original code number from postgress
// so we need to make it ourself
const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = pgx.ErrNoRows // this way, if we ever have to change the db driver again, we only need to update its value in 1 single place.

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

// konversi error menjadi pgError, jika berhasil itu menandakan bahwa error yang diterima berasal dari database postgres
func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	// *pgconn.PgError struct has a code field of type string, and in fact it's the code name that Postgres server returnsto the DB driver to handle
	if errors.As(err, &pgErr) { // converting err to PgErr
		// println(">>", pgErr.ConstraintName) // untuk mengetahui error terjadi karena constraint yang mana
		return pgErr.Code // if converting success we return this
	}
	return ""
}
