package plformat

import (
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func FormatPgTypeNumeric(num pgtype.Numeric) (string, error) {
	f8, err := num.Float64Value()
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(f8.Float64, 'f', 1, 32), nil
}

func FormatTimestamp(time time.Time) string {
	return time.Format("02/01/2006 15:04")
}
