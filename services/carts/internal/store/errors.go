package store

import (
	"errors"
	"fmt"

	"github.com/amane15/ecom_microservice/services/carts/internal/data"
	"github.com/jackc/pgx/v5/pgconn"
)

func wrapError(message string, pgError *pgconn.PgError, err error) error {
	pgErrorFormat := fmt.Sprintf("(constraint=%s, column=%s, table=%s)",
		pgError.ConstraintName,
		pgError.ColumnName,
		pgError.TableName,
	)
	format := "%s %s: %w"
	return fmt.Errorf(format,
		message,
		pgErrorFormat,
		err)
}

func handlePostgresErrors(err error, pgError *pgconn.PgError) error {
	switch pgError.Code {
	case "23505":
		return wrapError("unique violation", pgError, data.ErrItemAlreadyExists)
	default:
		return fmt.Errorf("postgres error (code=%s): %w", pgError.Code, err)
	}
}

func checkAndHandlePostgresErrors(err error) error {
	var pgError *pgconn.PgError
	if ok := errors.As(err, &pgError); ok {
		return handlePostgresErrors(err, pgError)
	}

	return err
}
