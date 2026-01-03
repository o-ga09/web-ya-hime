package date

import "time"

const (
	// DefaultFormat はデフォルトの日時フォーマット
	DefaultFormat = "2006-01-02 15:04:05"
)

// Format はtime.Timeを指定されたフォーマットで文字列に変換します
func Format(t time.Time, format string) string {
	return t.Format(format)
}

// FormatDefault はtime.Timeをデフォルトフォーマット（YYYY-MM-DD HH:MM:SS）で文字列に変換します
func FormatDefault(t time.Time) string {
	return t.Format(DefaultFormat)
}
