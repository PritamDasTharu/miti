package miti_test

import (
	"testing"
	"time"

	"github.com/PritamDasTharu/miti"
)

func TestAllDaysRoundTrip(t *testing.T) {
	// Loop through every valid day in 1970-2100 BS
	count := 0
	for y := miti.MinYear; y <= miti.MaxYear; y++ {
		for m := miti.Baisakh; m <= miti.Chaitra; m++ {
			// Find month days
			sampleDate, err := miti.New(y, m, 1)
			if err != nil {
				t.Fatalf("failed to create date for %d-%d-01: %v", y, m, err)
			}
			daysInMonth := sampleDate.DaysInMonth()

			for d := 1; d <= daysInMonth; d++ {
				bsDate, err := miti.New(y, m, d)
				if err != nil {
					t.Fatalf("failed New(%d, %d, %d): %v", y, m, d, err)
				}

				adTime := bsDate.ToTime()
				backBS, err := miti.FromTime(adTime)
				if err != nil {
					t.Fatalf("FromTime(%v) failed for %v: %v", adTime, bsDate, err)
				}

				if backBS != bsDate {
					t.Fatalf("roundtrip mismatch: original %v != converted %v (AD: %v)", bsDate, backBS, adTime)
				}

				// Weekday invariant check: independent Weekday() == ToTime().Weekday()
				if bsDate.Weekday() != adTime.Weekday() {
					t.Fatalf("weekday mismatch on %v: bs.Weekday()=%v, ad.Weekday()=%v",
						bsDate, bsDate.Weekday(), adTime.Weekday())
				}

				count++
			}
		}
	}

	if count != 47849 {
		t.Errorf("expected 47,849 round-trip days, got %d", count)
	}
}

func TestFromTimeNepalServerSafety(t *testing.T) {
	// 2026-11-03 at 20:00:00 UTC
	// In Nepal (+05:45), this is 2026-11-04 01:45:00 (+05:45)
	utcTime := time.Date(2026, time.November, 3, 20, 0, 0, 0, time.UTC)

	// Calling FromTime on UTC instant gives 2026-11-03 in UTC
	bsFromUTC, err := miti.FromTime(utcTime)
	if err != nil {
		t.Fatalf("FromTime failed: %v", err)
	}

	// Calling FromTimeNepal correctly advances into Nepal's calendar day (2026-11-04)
	bsFromNepal, err := miti.FromTimeNepal(utcTime)
	if err != nil {
		t.Fatalf("FromTimeNepal failed: %v", err)
	}

	if bsFromUTC == bsFromNepal {
		t.Errorf("expected FromTime and FromTimeNepal to differ during 20:00 UTC, but got identical: %v", bsFromUTC)
	}

	// Check that Today() does not return zero value
	today := miti.Today()
	if today.IsZero() || !today.IsValid() {
		t.Errorf("Today() returned invalid date: %v", today)
	}
}

func TestZeroValueNoPanic(t *testing.T) {
	var zero miti.Date

	if !zero.IsZero() {
		t.Errorf("zero.IsZero() expected true, got false")
	}
	if zero.IsValid() {
		t.Errorf("zero.IsValid() expected false, got true")
	}
	if zero.Year() != 0 {
		t.Errorf("zero.Year() expected 0, got %d", zero.Year())
	}
	if zero.Month() != 0 {
		t.Errorf("zero.Month() expected 0, got %v", zero.Month())
	}
	if zero.Day() != 0 {
		t.Errorf("zero.Day() expected 0, got %d", zero.Day())
	}
	if zero.DaysInMonth() != 0 {
		t.Errorf("zero.DaysInMonth() expected 0, got %d", zero.DaysInMonth())
	}
	if zero.String() != "" {
		t.Errorf("zero.String() expected empty, got %q", zero.String())
	}
	if zero.Weekday() != time.Sunday {
		t.Errorf("zero.Weekday() expected Sunday, got %v", zero.Weekday())
	}
	if !zero.ToTime().IsZero() {
		t.Errorf("zero.ToTime() expected zero time, got %v", zero.ToTime())
	}
	if zero.StartOfMonth() != (miti.Date{}) {
		t.Errorf("zero.StartOfMonth() expected zero Date, got %v", zero.StartOfMonth())
	}
	if zero.EndOfMonth() != (miti.Date{}) {
		t.Errorf("zero.EndOfMonth() expected zero Date, got %v", zero.EndOfMonth())
	}
	if zero.IsWeekend() {
		t.Errorf("zero.IsWeekend() expected false, got true")
	}
	if zero.Sub(zero) != 0 {
		t.Errorf("zero.Sub(zero) expected 0, got %d", zero.Sub(zero))
	}
	if zero.Compare(zero) != 0 {
		t.Errorf("zero.Compare(zero) expected 0, got %d", zero.Compare(zero))
	}
	if _, err := zero.AddDays(1); err != miti.ErrInvalidDate {
		t.Errorf("zero.AddDays() expected ErrInvalidDate, got %v", err)
	}
	if zero.FiscalYear() != "" {
		t.Errorf("zero.FiscalYear() expected empty string, got %q", zero.FiscalYear())
	}
	if _, _, err := zero.FiscalYearRange(); err != miti.ErrInvalidDate {
		t.Errorf("zero.FiscalYearRange() expected ErrInvalidDate, got %v", err)
	}
	if zero.Quarter() != 0 {
		t.Errorf("zero.Quarter() expected 0, got %d", zero.Quarter())
	}
}

func FuzzConversion(f *testing.F) {
	// Seed known timestamps
	f.Add(int64(0))           // 1970-01-01 UTC
	f.Add(int64(-1790035200)) // 1913-04-13 UTC
	f.Add(int64(1442707200))  // 2015-09-20 UTC
	f.Add(int64(1793750400))  // 2026-11-03 UTC

	f.Fuzz(func(t *testing.T, timestamp int64) {
		tm := time.Unix(timestamp, 0).UTC()
		bs, err := miti.FromTime(tm)
		if err != nil {
			return // Expected for dates outside 1970-2100 BS
		}
		// If conversion succeeded, must roundtrip to exact same calendar day in AD
		ad := bs.ToTime()
		if ad.Year() != tm.Year() || ad.Month() != tm.Month() || ad.Day() != tm.Day() {
			t.Fatalf("fuzz roundtrip mismatch: tm=%v, bs=%v, ad=%v", tm, bs, ad)
		}
	})
}
