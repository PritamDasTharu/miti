package miti_test

import (
	"testing"

	"github.com/PritamDasTharu/miti"
)

func TestFiscalYear(t *testing.T) {
	// 1. Shrawan 1 starts new FY: 2080-04-01 -> "2080/81"
	d1, _ := miti.New(2080, miti.Shrawan, 1)
	if fy := d1.FiscalYear(); fy != "2080/81" {
		t.Errorf("2080-04-01 FiscalYear = %q, expected '2080/81'", fy)
	}

	// 2. Chaitra in same calendar year belongs to "2080/81"
	d2, _ := miti.New(2080, miti.Chaitra, 30)
	if fy := d2.FiscalYear(); fy != "2080/81" {
		t.Errorf("2080-12-30 FiscalYear = %q, expected '2080/81'", fy)
	}

	// 3. Ashadh 31 ends FY: 2081-03-31 belongs to "2080/81"
	d3, _ := miti.New(2081, miti.Ashadh, 31)
	if fy := d3.FiscalYear(); fy != "2080/81" {
		t.Errorf("2081-03-31 FiscalYear = %q, expected '2080/81'", fy)
	}

	// 4. Baisakh rule check: Baisakh 2083 belongs to "2082/83", NOT "2083/84"
	baisakh2083, _ := miti.New(2083, miti.Baisakh, 1)
	if fy := baisakh2083.FiscalYear(); fy != "2082/83" {
		t.Errorf("2083-01-01 FiscalYear = %q, expected '2082/83'", fy)
	}

	// 5. Century rollover: year 2099 -> "2099/00"
	d99, _ := miti.New(2099, miti.Shrawan, 1)
	if fy := d99.FiscalYear(); fy != "2099/00" {
		t.Errorf("2099-04-01 FiscalYear = %q, expected '2099/00'", fy)
	}

	d100, _ := miti.New(2100, miti.Baisakh, 1)
	if fy := d100.FiscalYear(); fy != "2099/00" {
		t.Errorf("2100-01-01 FiscalYear = %q, expected '2099/00'", fy)
	}

	// 6. Zero Date returns empty string
	var zero miti.Date
	if fy := zero.FiscalYear(); fy != "" {
		t.Errorf("zero.FiscalYear() = %q, expected empty", fy)
	}
}

func TestFiscalYearRange(t *testing.T) {
	// Normal mid-range date
	d, _ := miti.New(2081, miti.Baisakh, 15)
	start, end, err := d.FiscalYearRange()
	if err != nil {
		t.Fatalf("FiscalYearRange failed: %v", err)
	}

	// FY 2080/81 starts 2080-04-01 and ends 2081-03-31
	if start.Year() != 2080 || start.Month() != miti.Shrawan || start.Day() != 1 {
		t.Errorf("FY start got %v, expected 2080-04-01", start)
	}
	if end.Year() != 2081 || end.Month() != miti.Ashadh || end.Day() != 31 {
		t.Errorf("FY end got %v, expected 2081-03-31", end)
	}

	// Table boundary checks:
	// 1. Baisakh 1970 (FY began in 1969 BS, out of table) -> ErrOutOfRange
	edgeLow, _ := miti.New(miti.MinYear, miti.Baisakh, 1)
	if _, _, err := edgeLow.FiscalYearRange(); err != miti.ErrOutOfRange {
		t.Errorf("expected ErrOutOfRange for Baisakh 1970, got %v", err)
	}

	// Valid low boundary: Shrawan 1970 starts valid FY 1970/71
	validLow, _ := miti.New(miti.MinYear, miti.Shrawan, 1)
	if s, e, err := validLow.FiscalYearRange(); err != nil || s.Year() != 1970 || e.Year() != 1971 {
		t.Errorf("validLow FiscalYearRange failed: s=%v, e=%v, err=%v", s, e, err)
	}

	// 2. Shrawan 2100 (FY ends in 2101 BS, out of table) -> ErrOutOfRange
	edgeHigh, _ := miti.New(miti.MaxYear, miti.Shrawan, 1)
	if _, _, err := edgeHigh.FiscalYearRange(); err != miti.ErrOutOfRange {
		t.Errorf("expected ErrOutOfRange for Shrawan 2100, got %v", err)
	}

	// Valid high boundary: Baisakh 2100 belongs to FY 2099/00
	validHigh, _ := miti.New(miti.MaxYear, miti.Baisakh, 1)
	if s, e, err := validHigh.FiscalYearRange(); err != nil || s.Year() != 2099 || e.Year() != 2100 {
		t.Errorf("validHigh FiscalYearRange failed: s=%v, e=%v, err=%v", s, e, err)
	}

	// Zero Date returns ErrInvalidDate
	var zero miti.Date
	if _, _, err := zero.FiscalYearRange(); err != miti.ErrInvalidDate {
		t.Errorf("expected ErrInvalidDate for zero Date, got %v", err)
	}
}

func TestQuarter(t *testing.T) {
	// Q1: Shrawan-Ashoj (Months 4-6)
	q1, _ := miti.New(2081, miti.Shrawan, 1)
	if q := q1.Quarter(); q != 1 {
		t.Errorf("Shrawan expected Q1, got %d", q)
	}

	// Q2: Kartik-Poush (Months 7-9)
	q2, _ := miti.New(2081, miti.Kartik, 1)
	if q := q2.Quarter(); q != 2 {
		t.Errorf("Kartik expected Q2, got %d", q)
	}

	// Q3: Magh-Chaitra (Months 10-12)
	q3, _ := miti.New(2081, miti.Magh, 1)
	if q := q3.Quarter(); q != 3 {
		t.Errorf("Magh expected Q3, got %d", q)
	}

	// Q4: Baisakh-Ashadh (Months 1-3)
	q4, _ := miti.New(2081, miti.Baisakh, 1)
	if q := q4.Quarter(); q != 4 {
		t.Errorf("Baisakh expected Q4, got %d", q)
	}

	// Zero Date returns 0
	var zero miti.Date
	if q := zero.Quarter(); q != 0 {
		t.Errorf("zero Date expected Q0, got %d", q)
	}
}
