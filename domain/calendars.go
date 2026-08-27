package domain

import "time"

type Calendar struct {
	Year  int
	Month time.Month
	Days  int
}

func NewCalendar(year int, month time.Month) Calendar {
	days := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	return Calendar{Year: year, Month: month, Days: days}
}
func (c Calendar) Contains(t time.Time) bool { return t.Year() == c.Year && t.Month() == c.Month }
func (c Calendar) Workdays() int {
	n := 0
	for d := 1; d <= c.Days; d++ {
		w := time.Date(c.Year, c.Month, d, 0, 0, 0, 0, time.UTC).Weekday()
		if w != time.Saturday && w != time.Sunday {
			n++
		}
	}
	return n
}
func (c Calendar) Label() string {
	return time.Date(c.Year, c.Month, 1, 0, 0, 0, 0, time.UTC).Format("2006-01")
}
