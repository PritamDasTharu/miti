package miti

import (
	"encoding"
	"encoding/json"
	"fmt"
)

var (
	_ json.Marshaler           = Date{}
	_ json.Unmarshaler         = (*Date)(nil)
	_ encoding.TextMarshaler   = Date{}
	_ encoding.TextUnmarshaler = (*Date)(nil)
)

// MarshalJSON serializes d as a JSON string "YYYY-MM-DD".
// The zero value Date{} serializes as null.
func (d Date) MarshalJSON() ([]byte, error) {
	if !d.IsValid() {
		return []byte("null"), nil
	}
	var buf [12]byte
	buf[0] = '"'
	copy(buf[1:11], d.String())
	buf[11] = '"'
	return buf[:], nil
}

// UnmarshalJSON deserializes d from JSON string "YYYY-MM-DD".
// In accordance with standard library time.Time convention, unmarshaling "null" is a no-op.
func (d *Date) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" {
		return nil
	}
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return fmt.Errorf("miti: cannot unmarshal %s into Date: expected string", s)
	}

	content := s[1 : len(s)-1]
	parsed, err := Parse(StandardDate, content)
	if err != nil {
		return fmt.Errorf("miti: cannot unmarshal JSON date %q: %w", content, err)
	}
	*d = parsed
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// Returns an empty slice for the zero Date{}.
func (d Date) MarshalText() ([]byte, error) {
	if !d.IsValid() {
		return []byte(""), nil
	}
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// An empty text sets d to Date{}.
func (d *Date) UnmarshalText(text []byte) error {
	s := string(text)
	if len(s) == 0 {
		*d = Date{}
		return nil
	}
	parsed, err := Parse(StandardDate, s)
	if err != nil {
		return fmt.Errorf("miti: cannot unmarshal text date %q: %w", s, err)
	}
	*d = parsed
	return nil
}
