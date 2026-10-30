package miti_test

import (
	"testing"

	"github.com/PritamDasTharu/miti"
)

func TestFormatAndAppendFormat(t *testing.T) {
	d, err := miti.New(2081, miti.Baisakh, 5)
	if err != nil {
		t.Fatalf("failed New: %v", err)
	}

	// Default StandardDate
	if s := d.Format("YYYY-MM-DD"); s != "2081-01-05" {
		t.Errorf("Format(YYYY-MM-DD) = %q, expected 2081-01-05", s)
	}
	if s := d.Format("YYYY/M/D"); s != "2081/1/5" {
		t.Errorf("Format(YYYY/M/D) = %q, expected 2081/1/5", s)
	}

	// AppendFormat pre-sized
	buf := make([]byte, 0, 16)
	buf = d.AppendFormat(buf, "YYYY-MM-DD")
	if string(buf) != "2081-01-05" {
		t.Errorf("AppendFormat = %q, expected 2081-01-05", string(buf))
	}

	// Zero value behavior
	var zero miti.Date
	if s := zero.Format("YYYY-MM-DD"); s != "" {
		t.Errorf("zero.Format() = %q, expected empty", s)
	}
	if b := zero.AppendFormat(nil, "YYYY-MM-DD"); len(b) != 0 {
		t.Errorf("zero.AppendFormat() = %v, expected nil/empty", b)
	}
}

func TestFormatInLocale(t *testing.T) {
	d, _ := miti.New(2081, miti.Baisakh, 1)

	// English format
	en := d.FormatIn("dddd, D MMMM YYYY", miti.LocaleEnglish)
	if en != "Saturday, 1 Baisakh 2081" {
		t.Errorf("English format = %q, expected Saturday, 1 Baisakh 2081", en)
	}

	// Nepali Devanagari format
	ne := d.FormatIn("dddd, D MMMM YYYY", miti.LocaleNepali)
	if ne != "शनिबार, १ वैशाख २०८१" {
		t.Errorf("Nepali format = %q, expected शनिबार, १ वैशाख २०८१", ne)
	}

	// Devanagari standard numeric date
	neNum := d.FormatIn("YYYY-MM-DD", miti.LocaleNepali)
	if neNum != "२०८१-०१-०१" {
		t.Errorf("Nepali numeric format = %q, expected २०८१-०१-०१", neNum)
	}
}

func TestLiteralEscaping(t *testing.T) {
	d, _ := miti.New(2081, miti.Baisakh, 15)

	// Escaped literals: [Month: ]MMMM, D[ of ]MMMM
	s1 := d.Format("[Month: ]MMMM")
	if s1 != "Month: Baisakh" {
		t.Errorf("Format([Month: ]MMMM) = %q, expected 'Month: Baisakh'", s1)
	}

	s2 := d.Format("D[ of ]MMMM, YYYY")
	if s2 != "15 of Baisakh, 2081" {
		t.Errorf("Format(D[ of ]MMMM, YYYY) = %q, expected '15 of Baisakh, 2081'", s2)
	}

	// Parse with literal escaping
	parsed, err := miti.Parse("D[ of ]MMMM, YYYY", "15 of Baisakh, 2081")
	if err != nil {
		t.Fatalf("Parse with escape failed: %v", err)
	}
	if parsed != d {
		t.Errorf("parsed date %v != original %v", parsed, d)
	}
}

func TestParseCanonicalAndAliases(t *testing.T) {
	// Parse English Standard Date
	d1, err := miti.Parse("YYYY-MM-DD", "2081-01-15")
	if err != nil {
		t.Fatalf("Parse(YYYY-MM-DD) failed: %v", err)
	}
	if d1.Year() != 2081 || d1.Month() != miti.Baisakh || d1.Day() != 15 {
		t.Errorf("Parse result mismatch: %v", d1)
	}

	// Parse Devanagari numerals
	d2, err := miti.Parse("YYYY-MM-DD", "२०८१-०१-१५")
	if err != nil {
		t.Fatalf("Parse Devanagari numerals failed: %v", err)
	}
	if d2 != d1 {
		t.Errorf("Devanagari parsed %v != %v", d2, d1)
	}

	// Parse month aliases (e.g. Asar, Saun, Fagun)
	samples := []struct {
		input       string
		layout      string
		expectMonth miti.Month
	}{
		{"15 Baishakh 2081", "D MMMM YYYY", miti.Baisakh},
		{"10 Asar 2081", "D MMMM YYYY", miti.Ashadh},
		{"1 Saun 2081", "D MMMM YYYY", miti.Shrawan},
		{"25 Asoj 2081", "D MMMM YYYY", miti.Ashoj},
		{"5 Push 2081", "D MMMM YYYY", miti.Poush},
		{"12 Fagun 2081", "D MMMM YYYY", miti.Falgun},
		{"१५ असार २०८१", "D MMMM YYYY", miti.Ashadh},
		{"१ साउन २०८१", "D MMMM YYYY", miti.Shrawan},
	}

	for _, s := range samples {
		parsed, err := miti.Parse(s.layout, s.input)
		if err != nil {
			t.Errorf("Parse(%q) failed: %v", s.input, err)
			continue
		}
		if parsed.Month() != s.expectMonth {
			t.Errorf("Parse(%q) got month %v, expected %v", s.input, parsed.Month(), s.expectMonth)
		}
	}
}

func TestMonthNameRoundTrip(t *testing.T) {
	// Verify that every canonical month name in English and Nepali formats and parses cleanly
	for m := miti.Baisakh; m <= miti.Chaitra; m++ {
		d, err := miti.New(2081, m, 1)
		if err != nil {
			t.Fatalf("failed New for month %v: %v", m, err)
		}

		// English
		enStr := d.FormatIn("YYYY-MMMM-DD", miti.LocaleEnglish)
		enParsed, err := miti.Parse("YYYY-MMMM-DD", enStr)
		if err != nil {
			t.Errorf("roundtrip EN parse failed for %s: %v", enStr, err)
		} else if enParsed != d {
			t.Errorf("roundtrip EN mismatch: got %v, expected %v", enParsed, d)
		}

		// Nepali
		neStr := d.FormatIn("YYYY-MMMM-DD", miti.LocaleNepali)
		neParsed, err := miti.Parse("YYYY-MMMM-DD", neStr)
		if err != nil {
			t.Errorf("roundtrip NE parse failed for %s: %v", neStr, err)
		} else if neParsed != d {
			t.Errorf("roundtrip NE mismatch: got %v, expected %v", neParsed, d)
		}
	}
}

func TestParseRejectsInvalid(t *testing.T) {
	// 1. YY format-only rejection
	if _, err := miti.Parse("YY-MM-DD", "81-01-15"); err != miti.ErrInvalidFormat {
		t.Errorf("expected ErrInvalidFormat for YY, got %v", err)
	}

	// 2. Adjacent variable-width tokens without delimiter (MD)
	if _, err := miti.Parse("YYYYMD", "2081115"); err != miti.ErrInvalidFormat {
		t.Errorf("expected ErrInvalidFormat for adjacent MD, got %v", err)
	}

	// 3. Weekday mismatch
	// 2081-01-01 was Saturday, not Sunday
	if _, err := miti.Parse("dddd, YYYY-MM-DD", "Sunday, 2081-01-01"); err != miti.ErrInvalidFormat {
		t.Errorf("expected ErrInvalidFormat for mismatched weekday, got %v", err)
	}

	// 4. Unknown tokens in parse
	if _, err := miti.Parse("YYYY-MMMMM-DD", "2081-01-15"); err != miti.ErrInvalidFormat {
		t.Errorf("expected ErrInvalidFormat for MMMMM, got %v", err)
	}
}

