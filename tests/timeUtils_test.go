package tests

import (
	"testing"
	"time"

	"github.com/ALiwoto/ssg/ssg/timeUtils"
)

func TestParseDuration(t *testing.T) {
	day, err := timeUtils.ParseDuration("1d")
	if err != nil {
		t.Error(err)
		return
	}
	correctDay := 24 * time.Hour
	if day != correctDay {
		t.Error("Expected ", correctDay, " but got ", day)
		return
	}

	week, err := timeUtils.ParseDuration("1w")
	if err != nil {
		t.Error(err)
		return
	}
	correctWeek := 7 * day
	if week != correctWeek {
		t.Error("Expected ", correctWeek, " but got ", week)
		return
	}

	month, err := timeUtils.ParseDuration("1mo")
	if err != nil {
		t.Error(err)
		return
	}
	correctMonth := 30 * day
	if month != correctMonth {
		t.Error("Expected ", correctMonth, " but got ", month)
		return
	}

	year, err := timeUtils.ParseDuration("1y")
	if err != nil {
		t.Error(err)
		return
	}
	correctYear := 365 * day
	if year != correctYear {
		t.Error("Expected ", correctYear, " but got ", year)
		return
	}

	var testCases = map[string]time.Duration{
		"1y2w":         year + (week * 2),
		"1y2w3d":       year + (week * 2) + (day * 3),
		"1y2w3d4h":     year + (week * 2) + (day * 3) + (4 * time.Hour),
		"1y2w3d4h5m":   year + (week * 2) + (day * 3) + (4 * time.Hour) + (5 * time.Minute),
		"3d4h5m6s":     (day * 3) + (4 * time.Hour) + (5 * time.Minute) + (6 * time.Second),
		"1s2w3y":       (3 * year) + (week * 2) + time.Second,
		"10h2mo1y605s": year + (2 * month) + (10 * time.Hour) + (605 * time.Second),
	}

	for input, correctValue := range testCases {
		value, err := timeUtils.ParseDuration(input)
		if err != nil {
			t.Error(err)
			return
		}
		if value != correctValue {
			t.Error("For input", input, "expected", correctValue, "but got", value)
			return
		}
	}
}

func TestParseDurationWithDefault(t *testing.T) {
	value := timeUtils.ParseDurationWithDefault("123", 24*time.Hour)
	if value != 123*24*time.Hour {
		t.Error("Expected 24h but got", value)
		return
	}
}

func TestPrettyTimeDuration(t *testing.T) {
	duration := (2 * time.Hour) + (3 * time.Second)
	pretty := timeUtils.GetPrettyTimeDuration(duration, true)

	parsedDuration, err := timeUtils.ParseDuration(pretty)
	if err != nil {
		t.Error(err)
		return
	}
	if parsedDuration != duration {
		t.Error("Expected", duration, "but got", parsedDuration)
		return
	}
}
