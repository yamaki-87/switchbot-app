package utils_test

import (
	"testing"
	"time"

	"github.com/yamaki-87/switchbot-app/src/internal/utils"
)

func unwarpTimeParse(layout string, value string) time.Time {
	result, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return result
}

func TestResolveProcessOneDate(t *testing.T) {
	from := utils.CorrectTimeToDate(unwarpTimeParse(utils.YYYYMMDD, "20260920"))
	to := utils.CorrectTimeToDate(unwarpTimeParse(utils.YYYYMMDD, "20260921"))
	dates := utils.ResolveProcessDate(from, to)
	if len(dates) != 1 {
		t.Fatalf("Expected 1 dates, got %d", len(dates))
	}

	if !dates[0].Equal(from) {
		t.Errorf("Expected %s, got %s", to, dates[0])
	}
}

func TestResolveProcessTwoDates(t *testing.T) {
	from := utils.CorrectTimeToDate(unwarpTimeParse(utils.YYYYMMDD, "20260920"))
	to := utils.CorrectTimeToDate(unwarpTimeParse(utils.YYYYMMDD, "20260922"))
	dates := utils.ResolveProcessDate(from, to)
	if len(dates) != 2 {
		t.Fatalf("Expected 2 dates, got %d", len(dates))
	}

	if !dates[0].Equal(from) {
		t.Errorf("Expected %s, got %s", from, dates[0])
	}

	dayafter := utils.AddDays(from, 1)
	if !dates[1].Equal(dayafter) {
		t.Errorf("Expected %s, got %s", dayafter, dates[1])
	}
}

func TestResolveProccessDateNextMonth(t *testing.T) {
	from := utils.CorrectTimeToDate(unwarpTimeParse(utils.YYYYMMDD, "20260920"))
	to := utils.CorrectTimeToDate(unwarpTimeParse(utils.YYYYMMDD, "20261020"))
	dates := utils.ResolveProcessDate(from, to)
	if len(dates) != 30 {
		t.Fatalf("Expected 31 dates, got %d", len(dates))
	}

	for _, date := range dates {
		if !date.Equal(from) {
			t.Errorf("Expected %s, got %s", from, date)
		}
		from = utils.AddDays(from, 1)
	}
}

func TestResoleveProcessDateSameDay(t *testing.T) {
	from := utils.CorrectTimeToDate(unwarpTimeParse(utils.YYYYMMDD, "20260920"))
	to := utils.CorrectTimeToDate(unwarpTimeParse(utils.YYYYMMDD, "20260920"))
	dates := utils.ResolveProcessDate(from, to)
	if len(dates) != 0 {
		t.Fatalf("Expected 1 dates, got %d", len(dates))
	}
}

func TestResolveProcessDateFromAfterTo(t *testing.T) {
	from := utils.CorrectTimeToDate(unwarpTimeParse(utils.YYYYMMDD, "20260920"))
	to := utils.CorrectTimeToDate(unwarpTimeParse(utils.YYYYMMDD, "20260919"))
	dates := utils.ResolveProcessDate(from, to)
	if len(dates) != 0 {
		t.Fatalf("Expected 0 dates, got %d", len(dates))
	}
}
