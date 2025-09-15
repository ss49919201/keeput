package date

import "time"

func BeginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func LocationJST() time.Location {
	return *jst
}
