package miti

import (
	"time"
)

// anchorUnixDays is the number of days from Unix Epoch (1970-01-01 UTC) to
// the anchor date 1 Baisakh 1970 BS (1913-04-13 UTC).
// 1913-04-13 is -20717 days from 1970-01-01.
const anchorUnixDays int64 = -20717

// daysFromAnchor converts year, month, day to days elapsed since 1970-01-01 BS.
func daysFromDate(d Date) int32 {
	yIdx := int(d.year) - MinYear
	days := yearCumulativeDays[yIdx]
	mIdx := int(d.month) - 1
	for m := 0; m < mIdx; m++ {
		days += int32(monthDays[yIdx][m])
	}
	days += int32(d.day - 1)
	return days
}

// dateFromDays converts days elapsed since 1970-01-01 BS into a valid Date.
// Expects 0 <= days < totalSupportedDays.
func dateFromDays(days int32) Date {
	// Approximate year index (average 365.25 days/year)
	yIdx := int(days / 365)
	if yIdx >= len(monthDays) {
		yIdx = len(monthDays) - 1
	}

	// Adjust estimate in O(1)
	for yIdx+1 < len(yearCumulativeDays) && yearCumulativeDays[yIdx+1] <= days {
		yIdx++
	}
	for yIdx > 0 && yearCumulativeDays[yIdx] > days {
		yIdx--
	}

	remDays := days - yearCumulativeDays[yIdx]
	var m uint8 = 0
	for m = 0; m < 12; m++ {
		mDays := int32(monthDays[yIdx][m])
		if remDays < mDays {
			break
		}
		remDays -= mDays
	}

	return Date{
		year:  int16(MinYear + yIdx),
		month: m + 1,
		day:   uint8(remDays + 1),
	}
}

// FromTime converts a Gregorian time.Time into its corresponding Bikram Sambat Date.
// It evaluates t.Date() in t's own location.
//
// WARNING: On servers running in UTC (such as most cloud hosts and Docker containers),
// between 18:15 UTC and 24:00 UTC, Nepal is already on the next calendar day (+05:45).
// To convert wall-clock time on a UTC host safely, use FromTimeNepal(t) or Today().
func FromTime(t time.Time) (Date, error) {
	y, m, d := t.Date()
	midnight := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	dayNumber := midnight.Unix() / 86400
	days := int32(dayNumber - anchorUnixDays)

	if days < 0 || days >= totalSupportedDays {
		return Date{}, ErrOutOfRange
	}
	return dateFromDays(days), nil
}

// FromTimeNepal converts an instant t to a Bikram Sambat Date using fixed Nepal Standard Time
// (UTC+05:45) without allocating any time.Location objects.
func FromTimeNepal(t time.Time) (Date, error) {
	// Offset t by +05:45 to get the calendar date in Nepal
	y, m, d := t.Add(5*time.Hour + 45*time.Minute).UTC().Date()
	midnight := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	dayNumber := midnight.Unix() / 86400
	days := int32(dayNumber - anchorUnixDays)

	if days < 0 || days >= totalSupportedDays {
		return Date{}, ErrOutOfRange
	}
	return dateFromDays(days), nil
}

// Today returns the current Bikram Sambat date in Nepal Standard Time (UTC+05:45).
// Returns Date{} if system clock is outside the supported range (1970-2100 BS).
func Today() Date {
	d, err := FromTimeNepal(time.Now())
	if err != nil {
		return Date{}
	}
	return d
}

// ToTime converts the Bikram Sambat Date to its Gregorian equivalent at midnight UTC.
// Returns time.Time{} for the zero Date{}.
func (d Date) ToTime() time.Time {
	if !d.IsValid() {
		return time.Time{}
	}
	days := daysFromDate(d)
	unixSec := (int64(days) + anchorUnixDays) * 86400
	return time.Unix(unixSec, 0).UTC()
}

// ToTimeAt converts the Bikram Sambat Date to a Gregorian time.Time with the specified
// hour, minute, second, nanosecond, and location.
// If loc is nil, time.UTC is used. Returns time.Time{} for the zero Date{}.
func (d Date) ToTimeAt(hour, min, sec, nsec int, loc *time.Location) time.Time {
	if !d.IsValid() {
		return time.Time{}
	}
	if loc == nil {
		loc = time.UTC
	}
	ad := d.ToTime()
	return time.Date(ad.Year(), ad.Month(), ad.Day(), hour, min, sec, nsec, loc)
}

// Weekday returns the day of the week for d.
// Computed independently from the cumulative day offset from the anchor weekday (Sunday).
// Returns time.Sunday for the zero Date{} (callers should verify d.IsValid() first).
func (d Date) Weekday() time.Weekday {
	if !d.IsValid() {
		return time.Sunday
	}
	days := daysFromDate(d)
	return time.Weekday(days % 7)
}
