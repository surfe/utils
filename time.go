package utils

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/jinzhu/now"
	"github.com/surfe/logger/v2"
)

func TimeToString(t time.Time) string {
	return t.Format(time.DateTime)
}

func TimeToStringOrEmpty(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format(time.DateTime)
}

func UTCTimestampToDateOrDefault(s string) string {
	if s == "" {
		return ""
	}

	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t.Format("02/01/2006")
	}

	return s
}

func Iso8601ToDate(s string) string {
	if s == "" {
		return ""
	}

	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return ""
	}

	return t.Format("01/02/2006")
}

func DMYtoDate(s string) string {
	if s == "" {
		return ""
	}

	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		return ""
	}

	return t.Format("01/02/2006")
}

func UnixMsecTimeToDate(s string) string {
	if s == "" {
		return ""
	}

	timeInMsec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return ""
	}

	t := TimestampToTime(timeInMsec, true)

	return t.Format("01/02/2006")
}

func DMYtoUnixMsec(s string) string {
	if s == "" {
		return ""
	}

	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		return ""
	}

	return strconv.FormatInt(t.UnixMilli(), 10)
}

func TimeToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}

	s := t.Format(time.DateTime)

	return &s
}

func TimeWithTimezoneToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}

	s := t.Format(time.RFC3339)

	return &s
}

func TimeToTimestamp(t time.Time, ms bool) int64 {
	ts := t.Unix()
	if ms {
		ts *= 1000
	}

	return ts
}

func TimeToTimestampPtr(t *time.Time, ms bool) *int64 {
	if t == nil {
		return nil
	}

	ts := t.Unix()
	if ms {
		ts *= 1000
	}

	return &ts
}

func StringToTime(s string) time.Time {
	t, _ := time.Parse(time.DateTime, s)

	return t
}

func StringToTime2(s string) time.Time {
	t, _ := time.Parse(time.DateOnly, s)

	return t
}

func StringToTime2Safe(s string) *time.Time {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil || t.IsZero() {
		return nil
	}

	return &t
}

func Iso8601ToTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05-0700", s)

	return t
}

func Iso8601ZToTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05Z0700", s)

	return t
}

func TimestampToTime(ts int64, ms bool) time.Time {
	if ms {
		ts /= 1000
	}

	return time.Unix(ts, 0)
}

func TimestampToTimePtr(ts *int64, ms bool) *time.Time {
	if ts == nil {
		return nil
	}

	t := TimestampToTime(*ts, ms)

	return &t
}

func TimestampStrToTime(ts string, ms bool) time.Time {
	i, _ := strconv.ParseInt(ts, 10, 64)

	return TimestampToTime(i, ms)
}

func ExpiryDate(s int) time.Time {
	return time.Now().Add(time.Duration(s) * time.Second)
}

func SameDay(date1, date2 time.Time, tz string) bool {
	// Load user's location for formatting
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}

	y1, m1, d1 := date1.In(loc).Date()
	y2, m2, d2 := date2.In(loc).Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

func GetLastWeekRange() (time.Time, time.Time) {
	var beginningOfWeek, endOfWeek time.Time

	lastWeek := time.Now().AddDate(0, 0, -7)

	// Set the start day of the week as Monday
	now.WeekStartDay = time.Monday
	beginningOfWeek = now.With(lastWeek).BeginningOfWeek()
	endOfWeek = now.With(lastWeek).EndOfWeek()

	return beginningOfWeek, endOfWeek
}

func DaysPassedSince(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

// IsSameTimeNoTimezone checks if two dates are roughly equal if no timezone considered.
func IsSameTimeNoTimezone(t1, t2 time.Time) bool {
	absoluteDiff := math.Abs(t1.Sub(t2).Minutes())
	minsDiff := math.Mod(absoluteDiff, 60)
	minsDiff = math.Min(minsDiff, 60-minsDiff)

	return minsDiff <= 2 && absoluteDiff <= 26*60
}

// DifferenceBetweenTimesToString calculates the time difference in the context of the time https://stackoverflow.com/a/36531443
func DifferenceBetweenTimesToString(a, b time.Time) string {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}

	if a.After(b) {
		a, b = b, a
	}

	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	years := y2 - y1
	months := int(M2 - M1)
	day := d2 - d1

	// Normalize negative values
	if day < 0 {
		months--
	}

	if months < 0 {
		months += 12
		years--
	}

	// Avoid zero value
	if months == 0 && years == 0 {
		months = 1
	}

	var s string
	if years > 0 {
		s = fmt.Sprintf("%d years", years)
		if months == 0 {
			return s
		}

		s += ", "
	}

	if months == 1 {
		s += fmt.Sprintf("%d month", months)
	} else {
		s += fmt.Sprintf("%d months", months)
	}

	return s
}

func DaysUntil(t time.Time) int {
	return int(math.Ceil(time.Until(t).Hours() / 24))
}

func GetTimezone(ctx context.Context, timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logger.Log(ctx).Err(err).Errorf("LoadLocation for %v", timezone)

		return time.UTC.String()
	}

	if loc != nil {
		return loc.String()
	}

	return timezone
}

func GetTimeFromLinkedInSentTimeLabelAndUserTimeZone(ctx context.Context, sentTimeLabel, timezone string) *time.Time {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logger.Log(ctx).Err(err).Errorf("LoadLocation for %v", timezone)
	}

	timeNow := time.Now().UTC().In(loc)
	number := Atoi(ExtractNumbersFromString(sentTimeLabel))

	sentTime := timeNow

	switch {
	case strings.Contains(sentTimeLabel, "today"):
		sentTime = timeNow
	case strings.Contains(sentTimeLabel, "yesterday"):
		sentTime = timeNow.AddDate(0, 0, -1)
	case strings.Contains(sentTimeLabel, "month"):
		sentTime = timeNow.AddDate(0, -number, 0)
	case strings.Contains(sentTimeLabel, "week"):
		sentTime = timeNow.AddDate(0, 0, -7*number)
	case strings.Contains(sentTimeLabel, "day"):
		sentTime = timeNow.AddDate(0, 0, -number)
	default:
		logger.Log(ctx).Warnf("Unknown sentTimeLabel: %v", sentTimeLabel)
	}

	return &sentTime
}

// ParseISODuration parses an ISO 8601 duration string.
// https://en.wikipedia.org/wiki/ISO_8601#Durations
// It handles both date (PnYnMnWnD) and time (PTnHnMnS) components,
// as well as combinations (PnYnMnWnDTnHnMnS). It does NOT handle
// fractions (like "PT0.5S").
func ParseISODuration(s string) (time.Duration, error) {
	if !strings.HasPrefix(s, "P") {
		return 0, errors.New("invalid duration: must start with P")
	}

	s = s[1:]
	if len(s) == 0 {
		return 0, errors.New("invalid date duration format: empty duration")
	}

	timePartIndex := strings.Index(s, "T")
	datePart := s
	timePart := ""

	if timePartIndex != -1 {
		datePart = s[:timePartIndex]
		timePart = s[timePartIndex+1:]

		if len(timePart) == 0 {
			return 0, errors.New("invalid time duration format: empty time part")
		}
	}

	years, months, weeks, days, err := parseIsoDurationDatePart(datePart)
	if err != nil {
		return 0, err
	}

	hours, minutes, seconds, err := parseIsoDurationTimePart(timePart)
	if err != nil {
		return 0, err
	}

	if years == 0 && months == 0 && weeks == 0 && days == 0 &&
		hours == 0 && minutes == 0 && seconds == 0 {
		if len(timePart) > 0 {
			return 0, fmt.Errorf("invalid time duration format: %s", timePart)
		}

		return 0, fmt.Errorf("invalid date duration format: %s", datePart)
	}

	return calculateIsoDuration(years, months, weeks, days, hours, minutes, seconds), nil
}

func parseIsoDurationDatePart(datePart string) (float64, float64, float64, float64, error) {
	var years, months, weeks, days float64

	if len(datePart) == 0 {
		return 0, 0, 0, 0, nil
	}

	regexDate := regexp.MustCompile(`^(?P<years>\d+Y)?(?P<months>\d+M)?(?P<weeks>\d+W)?(?P<days>\d+D)?$`)
	match := regexDate.FindStringSubmatch(datePart)

	if len(match) == 0 || (match[0] == "" && len(datePart) > 0) {
		return 0, 0, 0, 0, fmt.Errorf("invalid date duration format: %s", datePart)
	}

	for i, name := range regexDate.SubexpNames() {
		if i != 0 && match[i] != "" {
			value, err := strconv.ParseFloat(match[i][:len(match[i])-1], 64)
			if err != nil {
				return 0, 0, 0, 0, fmt.Errorf("invalid numeric value in date: %w", err)
			}

			switch name {
			case "years":
				years = value
			case "months":
				months = value
			case "weeks":
				weeks = value
			case "days":
				days = value
			}
		}
	}

	return years, months, weeks, days, nil
}

func parseIsoDurationTimePart(timePart string) (float64, float64, float64, error) {
	var hours, minutes, seconds float64

	if len(timePart) == 0 {
		return 0, 0, 0, nil
	}

	regexTime := regexp.MustCompile(`^(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?$`)
	match := regexTime.FindStringSubmatch(timePart)

	if len(match) == 0 || (match[0] == "" && len(timePart) > 0) {
		return 0, 0, 0, fmt.Errorf("invalid time duration format: %s", timePart)
	}

	for i, name := range regexTime.SubexpNames() {
		if i != 0 && match[i] != "" {
			value, err := strconv.ParseFloat(match[i][:len(match[i])-1], 64)
			if err != nil {
				return 0, 0, 0, fmt.Errorf("invalid numeric value in time: %w", err)
			}

			switch name {
			case "hours":
				hours = value
			case "minutes":
				minutes = value
			case "seconds":
				seconds = value
			}
		}
	}

	return hours, minutes, seconds, nil
}

func calculateIsoDuration(years, months, weeks, days, hours, minutes, seconds float64) time.Duration {
	var duration time.Duration

	duration += time.Duration(years * 365 * 24 * float64(time.Hour))
	duration += time.Duration(months * 30 * 24 * float64(time.Hour))
	duration += time.Duration(weeks * 7 * 24 * float64(time.Hour))
	duration += time.Duration(days * 24 * float64(time.Hour))
	duration += time.Duration(hours * float64(time.Hour))
	duration += time.Duration(minutes * float64(time.Minute))
	duration += time.Duration(seconds * float64(time.Second))

	return duration
}
