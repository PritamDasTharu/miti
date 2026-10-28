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
