package miti

//go:generate go run ./internal/gen/gen_table.go

import (
	"errors"
	"fmt"
)

// Month represents a month in the Bikram Sambat calendar (1 = Baisakh, 12 = Chaitra).
type Month uint8

const (
	Baisakh Month = 1 + iota
	Jestha
	Ashadh
	Shrawan
	Bhadra
	Ashoj
	Kartik
	Mangsir
	Poush
	Magh
	Falgun
	Chaitra
)

var (
	// ErrOutOfRange is returned when a date falls outside the supported range (1970-2100 BS).
	ErrOutOfRange = errors.New("miti: date out of supported range (1970-2100 BS)")

	// ErrInvalidDate is returned when a date's components are invalid (e.g. day exceeds month length).
	ErrInvalidDate = errors.New("miti: invalid date")

	// ErrInvalidFormat is returned when a layout string or value cannot be parsed.
	ErrInvalidFormat = errors.New("miti: invalid date format")
)

// Date represents an immutable calendar date in Bikram Sambat (BS).
// It occupies exactly 4 bytes of memory and is comparable with ==.
// The zero value Date{} represents an unset, invalid date.
type Date struct {
	year  int16
	month uint8
	day   uint8
}

// New creates and validates a new Date. Returns ErrOutOfRange or ErrInvalidDate if invalid.
func New(year int, month Month, day int) (Date, error) {
	if year < MinYear || year > MaxYear {
		return Date{}, ErrOutOfRange
	}
	if month < Baisakh || month > Chaitra {
		return Date{}, ErrInvalidDate
	}
	maxDays := int(monthDays[year-MinYear][month-1])
	if day < 1 || day > maxDays {
		return Date{}, ErrInvalidDate
	}
	return Date{
		year:  int16(year),
		month: uint8(month),
		day:   uint8(day),
	}, nil
}

// IsZero reports whether d is the zero value Date{}.
// Compatible with Go 1.24+ `omitzero` struct tag.
func (d Date) IsZero() bool {
	return d.year == 0 && d.month == 0 && d.day == 0
}

// IsValid reports whether d is a valid, initialized date within the supported range.
func (d Date) IsValid() bool {
	if d.year < MinYear || d.year > MaxYear {
		return false
	}
	if d.month < uint8(Baisakh) || d.month > uint8(Chaitra) {
		return false
	}
	maxDays := monthDays[int(d.year)-MinYear][d.month-1]
	return d.day >= 1 && d.day <= maxDays
}

// Year returns the Bikram Sambat year (e.g. 2083), or 0 if d is the zero value.
func (d Date) Year() int {
	return int(d.year)
}

// Month returns the Bikram Sambat month (1..12), or 0 if d is the zero value.
func (d Date) Month() Month {
	return Month(d.month)
}

// Day returns the day of the month (1..32), or 0 if d is the zero value.
func (d Date) Day() int {
	return int(d.day)
}

// DaysInMonth returns the total number of days in d's month (29..32), or 0 if d is the zero value.
func (d Date) DaysInMonth() int {
	if !d.IsValid() {
		return 0
	}
	return int(monthDays[int(d.year)-MinYear][d.month-1])
}

// StartOfMonth returns the first day of the month for d.
// Returns Date{} if d is the zero value Date{}.
func (d Date) StartOfMonth() Date {
	if !d.IsValid() {
		return Date{}
	}
	return Date{
		year:  d.year,
		month: d.month,
		day:   1,
	}
}

// EndOfMonth returns the last day of the month for d.
// Returns Date{} if d is the zero value Date{}.
func (d Date) EndOfMonth() Date {
	if !d.IsValid() {
		return Date{}
	}
	maxDays := monthDays[int(d.year)-MinYear][d.month-1]
	return Date{
		year:  d.year,
		month: d.month,
		day:   maxDays,
	}
}

// IsWeekend reports whether d falls on a standard Nepali weekend (Saturday).
// Returns false for the zero value Date{}.
func (d Date) IsWeekend() bool {
	if !d.IsValid() {
		return false
	}
	return d.Weekday() == 6 // time.Saturday
}

// Compare compares d and other.
// Returns:
//   - -1 if d is before other
//   -  0 if d equals other
//   - +1 if d is after other
//
// The zero value Date{} sorts before all valid dates.
func (d Date) Compare(other Date) int {
	if d == other {
		return 0
	}
	if d.IsZero() {
		return -1
	}
	if other.IsZero() {
		return 1
	}

	if d.year != other.year {
		if d.year < other.year {
			return -1
		}
		return 1
	}
	if d.month != other.month {
		if d.month < other.month {
			return -1
		}
		return 1
	}
	if d.day < other.day {
		return -1
	}
	return 1
}

// Before reports whether d is before other.
func (d Date) Before(other Date) bool {
	return d.Compare(other) < 0
}

// After reports whether d is after other.
func (d Date) After(other Date) bool {
	return d.Compare(other) > 0
}

// Equal reports whether d equals other.
func (d Date) Equal(other Date) bool {
	return d.Compare(other) == 0
}

// Sub returns the number of calendar days between d and other (d - other).
// Returns 0 if either operand is the zero value Date{} (check !d.IsZero() first).
func (d Date) Sub(other Date) int {
	if !d.IsValid() || !other.IsValid() {
		return 0
	}
	return int(daysFromDate(d) - daysFromDate(other))
}

// AddDays returns the date d plus days.
// Returns (Date{}, ErrInvalidDate) if d is the zero Date{}.
// Returns (Date{}, ErrOutOfRange) if the result falls outside the supported range (1970-2100 BS)
// or if days causes integer overflow.
func (d Date) AddDays(days int) (Date, error) {
	if !d.IsValid() {
		return Date{}, ErrInvalidDate
	}

	// Overflow guard
	if days > int(totalSupportedDays) || days < -int(totalSupportedDays) {
		return Date{}, ErrOutOfRange
	}

	currDays := daysFromDate(d)
	newDays := currDays + int32(days)
	if newDays < 0 || newDays >= totalSupportedDays {
		return Date{}, ErrOutOfRange
	}
	return dateFromDays(newDays), nil
}

// AddMonths returns the date d plus months.
// If the target month has fewer days than d.Day(), the day is clamped to the last valid day
// of the target month (e.g., 32 Baisakh + 1 month clamps to the last day of Jestha).
// Note: This differs from time.AddDate, which overflows into the subsequent month.
// Returns (Date{}, ErrInvalidDate) if d is the zero Date{}.
// Returns (Date{}, ErrOutOfRange) if the result falls outside 1970-2100 BS or on integer overflow.
func (d Date) AddMonths(months int) (Date, error) {
	if !d.IsValid() {
		return Date{}, ErrInvalidDate
	}

	// Overflow check
	maxMonths := (MaxYear - MinYear + 1) * 12
	if months > maxMonths || months < -maxMonths {
		return Date{}, ErrOutOfRange
	}

	totalMonths := int(d.year)*12 + int(d.month) - 1 + months
	targetYear := totalMonths / 12
	targetMonthIdx := totalMonths % 12
	if targetMonthIdx < 0 {
		targetMonthIdx += 12
		targetYear--
	}
	targetMonth := Month(targetMonthIdx + 1)

	if targetYear < MinYear || targetYear > MaxYear {
		return Date{}, ErrOutOfRange
	}

	targetMaxDays := int(monthDays[targetYear-MinYear][targetMonth-1])
	targetDay := int(d.day)
	if targetDay > targetMaxDays {
		targetDay = targetMaxDays
	}

	return New(targetYear, targetMonth, targetDay)
}

// String returns the ISO representation "YYYY-MM-DD" of the Bikram Sambat date.
// Returns an empty string for the zero Date{}.
func (d Date) String() string {
	if !d.IsValid() {
		return ""
	}
	// Manual zero-allocation digit formatting
	var buf [10]byte
	y := int(d.year)
	buf[0] = byte('0' + y/1000)
	buf[1] = byte('0' + (y/100)%10)
	buf[2] = byte('0' + (y/10)%10)
	buf[3] = byte('0' + y%10)
	buf[4] = '-'
	buf[5] = byte('0' + d.month/10)
	buf[6] = byte('0' + d.month%10)
	buf[7] = '-'
	buf[8] = byte('0' + d.day/10)
	buf[9] = byte('0' + d.day%10)
	return string(buf[:])
}

// FormatStringer compatibility check
var _ fmt.Stringer = Date{}
