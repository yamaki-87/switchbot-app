package utils

import "time"

const YYYYMMDD = "20060102"

func ResolveProcessDate(from time.Time, to time.Time) []time.Time {
	from = CorrectTimeToDate(from)
	to = CorrectTimeToDate(to)

	var dates []time.Time

	if from.After(to) {
		return dates
	}

	for !from.Equal(to) {
		dates = append(dates, from)
		from = AddDays(from, 1)
	}
	return dates
}

func CorrectTimeToDate(date time.Time) time.Time {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	return date
}

func AddDays(date time.Time, days int) time.Time {
	date = date.AddDate(0, 0, days)
	return date
}
