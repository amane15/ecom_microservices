package dbutils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func StringToPtr(ns pgtype.Text) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func Int8ToPtr(ni pgtype.Int8) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

func TimeToPtr(nt pgtype.Timestamptz) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

func BoolToPtr(nb pgtype.Bool) *bool {
	if nb.Valid {
		return &nb.Valid
	}
	return nil
}

func PtrToString(s *string) pgtype.Text {
	ns := pgtype.Text{Valid: false}
	if s != nil {
		ns.String = *s
		ns.Valid = true
	}
	return ns
}

func PtrToInt8(i *int64) pgtype.Int8 {
	ni := pgtype.Int8{Valid: false}
	if i != nil {
		ni.Int64 = *i
		ni.Valid = true
	}
	return ni
}

func PtrToBool(b *bool) pgtype.Bool {
	nb := pgtype.Bool{Valid: false}
	if b != nil {
		nb.Bool = *b
		nb.Valid = true
	}
	return nb
}

func DecimalPtrToNumeric(d *decimal.Decimal) pgtype.Numeric {
	n := pgtype.Numeric{Valid: false}
	if d == nil {
		return n
	}

	n.Int = d.Coefficient()
	n.Exp = d.Exponent()
	n.Valid = true
	return n
}

func NumericToDecimalPtr(n pgtype.Numeric) *decimal.Decimal {
	if !n.Valid || n.Int == nil {
		return nil
	}

	d := decimal.NewFromBigInt(n.Int, n.Exp)
	return &d
}
