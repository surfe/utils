package pgtypex

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TextToStringPointer(str pgtype.Text) *string {
	if !str.Valid {
		return nil
	}

	return &str.String
}

func Int4ToPointer(num pgtype.Int4) *int {
	if !num.Valid {
		return nil
	}

	n := int(num.Int32)

	return &n
}

func BoolToBoolPointer(b pgtype.Bool) *bool {
	if !b.Valid {
		return nil
	}

	return &b.Bool
}

func UUIDsToPgUUIDs(ids []uuid.UUID) []pgtype.UUID {
	uuids := make([]pgtype.UUID, 0, len(ids))
	for _, id := range ids {
		uuids = append(uuids, pgtype.UUID{Bytes: id, Valid: true})
	}

	return uuids
}
