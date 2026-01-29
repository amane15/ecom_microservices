package dbutils

import (
	"database/sql"
	"time"
)

func MapNullString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func MapNullInt64(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

func MapNullTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

func MapNullBool(nb sql.NullBool) *bool {
	if nb.Valid {
		return &nb.Valid
	}
	return nil
}

func MapStringPtr(s *string) sql.NullString {
	ns := sql.NullString{Valid: false}
	if s != nil {
		ns.String = *s
		ns.Valid = true
	}
	return ns
}

func MapInt64Ptr(i *int64) sql.NullInt64 {
	ni := sql.NullInt64{Valid: false}
	if i != nil {
		ni.Int64 = *i
		ni.Valid = true
	}
	return ni
}

func MapBoolPtr(b *bool) sql.NullBool {
	nb := sql.NullBool{Valid: false}
	if b != nil {
		nb.Bool = *b
		nb.Valid = true
	}
	return nb
}
