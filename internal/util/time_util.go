package util

import "time"

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

func FormatTime(t time.Time, layout string) string {
	if t.IsZero() {
		return ""
	}

	return t.Format(layout)
}
