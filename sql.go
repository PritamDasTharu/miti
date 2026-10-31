package miti

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Value implements the database/sql/driver.Valuer interface for Date.
//
// By default, it returns the Gregorian (AD) time.Time at midnight UTC.
// This is safe for standard SQL DATE columns across PostgreSQL, MySQL, and SQLite,
// and avoids corrupting database columns with BS years or 32-day months.
//
// Returns nil (SQL NULL) if d is the zero value Date{}.
func (d Date) Value() (driver.Value, error) {
	if !d.IsValid() {
		return nil, nil
	}
	return d.ToTime(), nil
}

// Scan implements the database/sql.Scanner interface for Date.
//
// Semantics:
//   - nil: sets *d = Date{} (SQL NULL).
//   - time.Time: converts AD to BS using the timestamp's own location.
//   - string or []byte: parses as Gregorian (AD) ISO date "YYYY-MM-DD".
//     This maintains symmetry with Value() for database drivers (such as SQLite or MySQL
//     without parseTime=true) that return SQL DATE columns as raw text.
//
// If you explicitly wish to store and read Bikram Sambat strings in VARCHAR/TEXT columns,
// use the BSSQLDate type instead.
func (d *Date) Scan(src any) error {
	if src == nil {
		*d = Date{}
		return nil
	}

	switch v := src.(type) {
	case time.Time:
		bs, err := FromTime(v)
		if err != nil {
			return fmt.Errorf("miti: cannot scan time.Time %v: %w", v, err)
		}
		*d = bs
		return nil

	case []byte:
		return d.scanADString(string(v))

	case string:
		return d.scanADString(v)

	default:
		return fmt.Errorf("miti: unsupported Scan source type %T", src)
	}
}

func (d *Date) scanADString(s string) error {
	s = stringsTrimSpace(s)
	if len(s) == 0 {
		*d = Date{}
		return nil
	}
	// Parse AD date "2006-01-02"
	parsed, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("miti: cannot scan AD date string %q: %w", s, err)
	}
	bs, err := FromTime(parsed)
	if err != nil {
		return fmt.Errorf("miti: date %q out of BS range: %w", s, err)
	}
	*d = bs
	return nil
}

// BSSQLDate is an opt-in wrapper type for schemas that store Bikram Sambat dates
// as text strings ("YYYY-MM-DD") in VARCHAR or TEXT columns.
//
// Zero value BSSQLDate{} maps to SQL NULL.
type BSSQLDate Date

// Date returns the underlying Date.
func (b BSSQLDate) Date() Date {
	return Date(b)
}

// Value implements driver.Valuer, returning the Bikram Sambat date string "YYYY-MM-DD"
// or nil for the zero value.
func (b BSSQLDate) Value() (driver.Value, error) {
	d := Date(b)
	if !d.IsValid() {
		return nil, nil
	}
	return d.String(), nil
}

// Scan implements database/sql.Scanner, parsing Bikram Sambat date strings "YYYY-MM-DD".
func (b *BSSQLDate) Scan(src any) error {
	if src == nil {
		*b = BSSQLDate{}
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("miti: unsupported BSSQLDate Scan source type %T", src)
	}

	s = stringsTrimSpace(s)
	if len(s) == 0 {
		*b = BSSQLDate{}
		return nil
	}

	d, err := Parse(StandardDate, s)
	if err != nil {
		return fmt.Errorf("miti: cannot scan BSSQLDate %q: %w", s, err)
	}
	*b = BSSQLDate(d)
	return nil
}

// MarshalJSON serializes BSSQLDate as JSON string "YYYY-MM-DD" or null.
func (b BSSQLDate) MarshalJSON() ([]byte, error) {
	return Date(b).MarshalJSON()
}

// UnmarshalJSON deserializes BSSQLDate from JSON.
func (b *BSSQLDate) UnmarshalJSON(data []byte) error {
	d := Date(*b)
	if err := d.UnmarshalJSON(data); err != nil {
		return err
	}
	*b = BSSQLDate(d)
	return nil
}

// MarshalText implements encoding.TextMarshaler for BSSQLDate.
func (b BSSQLDate) MarshalText() ([]byte, error) {
	return Date(b).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler for BSSQLDate.
func (b *BSSQLDate) UnmarshalText(text []byte) error {
	d := Date(*b)
	if err := d.UnmarshalText(text); err != nil {
		return err
	}
	*b = BSSQLDate(d)
	return nil
}

func stringsTrimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
