package nullvalue

import (
	"database/sql"
	"time"
)

func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func ToNullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: i, Valid: true}
}

func ToNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

func ToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func ToNullFloat64(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}

// StringToSqlString は文字列からdatabase/sql.NullStringに変換します
func StringToSqlString(s string) sql.NullString {
	return ToNullString(s)
}

// PointerToSqlString はポインタ型文字列からdatabase/sql.NullStringに変換します
func PointerToSqlString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// PointerToString はポインタ型文字列から文字列に変換します（nil の場合は空文字）
func PointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// SqlStringToPointer はdatabase/sql.NullStringからポインタ型文字列に変換します
func SqlStringToPointer(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// SqlStringToString はdatabase/sql.NullStringから文字列に変換します（無効な場合は空文字）
func SqlStringToString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return s.String
}
