package moment

import "time"

func BeginOfYear(year int, location *time.Location) time.Time {
	if location == nil {
		location = time.Local
	}
	return time.Date(year, 1, 1, 0, 0, 0, 0, location)
}

func EndOfYear(year int, location *time.Location) time.Time {
	if location == nil {
		location = time.Local
	}
	return time.Date(year, 12, 31, 23, 59, 59, 999999999, location)
}

func BeginOfMonth(year, month int, location *time.Location) time.Time {
	if location == nil {
		location = time.Local
	}
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, location)
}

func EndOfMonth(year, month int, location *time.Location) time.Time {
	if location == nil {
		location = time.Local
	}
	return time.Date(year, time.Month(month), DayOfMonth(year, month), 23, 59, 59, 999999999, location)
}

func BeginOfDay(year, month, day int, location *time.Location) time.Time {
	if location == nil {
		location = time.Local
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)
}

func EndOfDay(year, month, day int, location *time.Location) time.Time {
	if location == nil {
		location = time.Local
	}
	return time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, location)
}

func DayOfYear(year int) int {
	if IsLeapYear(year) {
		return 366
	}
	return 365
}

func DayOfMonth(year int, month int) int {
	if month < 1 || month > 12 {
		return 0
	}
	switch month {
	case 2:
		if IsLeapYear(year) {
			return 29
		}
		return 28
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	default:
		return 30
	}
}

func IsLeapYear(year int) bool {
	if year%4 == 0 && year%100 != 0 || year%400 == 0 {
		return true
	}
	return false
}
