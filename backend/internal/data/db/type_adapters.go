package db

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/jackc/pgx/v5/pgtype"
)

func IntoPgTimePrt(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{
			Time:             time.Time{},
			InfinityModifier: 0,
			Valid:            false,
		}
	}

	return pgtype.Timestamp{
		Time:             *t,
		InfinityModifier: 0,
		Valid:            true,
	}
}

func IntoCivilDate(d pgtype.Date) civil.Date {
	if !d.Valid {
		return civil.Date{}
	}

	return civil.Date{
		Year:  d.Time.Year(),
		Month: d.Time.Month(),
		Day:   d.Time.Day(),
	}
}

func IntoPgDate(d *civil.Date) pgtype.Date {
	if d == nil {
		return pgtype.Date{
			Time:             time.Time{},
			InfinityModifier: 0,
			Valid:            false,
		}
	}
	return pgtype.Date{
		Time:             time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC),
		InfinityModifier: 0,
		Valid:            true,
	}
}
