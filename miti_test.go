package miti_test

import (
	"math"
	"slices"
	"testing"

	"github.com/PritamDasTharu/miti"
)

func TestStartAndEndOfMonth(t *testing.T) {
	d, err := miti.New(2081, miti.Baisakh, 15)
	if err != nil {
		t.Fatalf("failed New: %v", err)
	}

	start := d.StartOfMonth()
	if start.Day() != 1 || start.Month() != miti.Baisakh || start.Year() != 2081 {
		t.Errorf("StartOfMonth failed, got %v", start)
	}

	end := d.EndOfMonth()
	if end.Day() != 31 || end.Month() != miti.Baisakh || end.Year() != 2081 {
		t.Errorf("EndOfMonth failed, got %v", end)
	}
}

func TestIsWeekend(t *testing.T) {
	// 2081-01-01 was Saturday (verified in golden anchors)
	newYear, err := miti.New(2081, miti.Baisakh, 1)
	if err != nil {
		t.Fatalf("failed New: %v", err)
	}
	if !newYear.IsWeekend() {
		t.Errorf("2081-01-01 should be weekend (Saturday)")
	}

	// 2081-01-02 was Sunday
	sunday, err := newYear.AddDays(1)
	if err != nil {
		t.Fatalf("failed AddDays(1): %v", err)
	}
	if sunday.IsWeekend() {
		t.Errorf("2081-01-02 should not be weekend (Sunday)")
	}
}

func TestCompareBeforeAfterEqual(t *testing.T) {
	d1, _ := miti.New(2080, miti.Chaitra, 30)
	d2, _ := miti.New(2081, miti.Baisakh, 1)
	d3, _ := miti.New(2081, miti.Baisakh, 1)

	if !d1.Before(d2) {
		t.Errorf("%v should be before %v", d1, d2)
	}
	if !d2.After(d1) {
		t.Errorf("%v should be after %v", d2, d1)
	}
	if !d2.Equal(d3) {
		t.Errorf("%v should equal %v", d2, d3)
	}
	if d2.Compare(d3) != 0 {
		t.Errorf("expected Compare=0, got %d", d2.Compare(d3))
	}

	// Test sorting with slices.SortFunc
	dates := []miti.Date{d2, miti.Date{}, d1}
	slices.SortFunc(dates, func(a, b miti.Date) int {
		return a.Compare(b)
	})

	if !dates[0].IsZero() || dates[1] != d1 || dates[2] != d2 {
		t.Errorf("SortFunc failed, got: %v", dates)
	}
}

func TestSub(t *testing.T) {
	d1, _ := miti.New(2081, miti.Baisakh, 1)
	d2, _ := miti.New(2081, miti.Baisakh, 15)

	if diff := d2.Sub(d1); diff != 14 {
		t.Errorf("expected 14 days difference, got %d", diff)
	}
	if diff := d1.Sub(d2); diff != -14 {
		t.Errorf("expected -14 days difference, got %d", diff)
	}

	// Zero Date{} handling
	if diff := d1.Sub(miti.Date{}); diff != 0 {
		t.Errorf("expected 0 for zero Date operand, got %d", diff)
	}
}

func TestAddDays(t *testing.T) {
	d, _ := miti.New(2081, miti.Baisakh, 1)

	// Add 30 days
	next, err := d.AddDays(30)
	if err != nil {
		t.Fatalf("AddDays(30) failed: %v", err)
	}
	if next.Day() != 31 || next.Month() != miti.Baisakh {
		t.Errorf("expected 2081-01-31, got %v", next)
	}

	// Cross month boundary (+1 day -> 2081-02-01)
	crossMonth, err := next.AddDays(1)
	if err != nil {
		t.Fatalf("AddDays(1) failed: %v", err)
	}
	if crossMonth.Day() != 1 || crossMonth.Month() != miti.Jestha || crossMonth.Year() != 2081 {
		t.Errorf("expected 2081-02-01, got %v", crossMonth)
	}

	// Extreme bounds and overflow protection
	if _, err := d.AddDays(math.MaxInt); err != miti.ErrOutOfRange {
		t.Errorf("expected ErrOutOfRange for math.MaxInt, got %v", err)
	}
	if _, err := d.AddDays(math.MinInt); err != miti.ErrOutOfRange {
		t.Errorf("expected ErrOutOfRange for math.MinInt, got %v", err)
	}
	if _, err := (miti.Date{}).AddDays(1); err != miti.ErrInvalidDate {
		t.Errorf("expected ErrInvalidDate for zero Date, got %v", err)
	}
}

func TestConsecutiveDayInvariant(t *testing.T) {
	// Sample every 50 days across the entire supported range
	// Verifies: d.AddDays(1).ToTime().Equal(d.ToTime().AddDate(0, 0, 1))
	start, _ := miti.New(miti.MinYear, miti.Baisakh, 1)
	curr := start

	for i := 0; i < 47848; i += 47 {
		next, err := curr.AddDays(1)
		if err != nil {
			t.Fatalf("AddDays(1) failed on %v: %v", curr, err)
		}

		t1 := curr.ToTime().AddDate(0, 0, 1)
		t2 := next.ToTime()

		if !t1.Equal(t2) {
			t.Fatalf("consecutive invariant failed at %v: %v != %v", curr, t1, t2)
		}

		// Advance curr by 47 days
		curr, err = curr.AddDays(47)
		if err != nil {
			break
		}
	}
}

func TestAddMonthsClamping(t *testing.T) {
	// Year 1970 BS: Ashadh (month 3) has 32 days, Shrawan (month 4) has 31 days.
	d32, err := miti.New(1970, miti.Ashadh, 32)
	if err != nil {
		t.Fatalf("failed New(1970, 3, 32): %v", err)
	}

	// Adding 1 month clamps to 31 Shrawan 1970
	clamped, err := d32.AddMonths(1)
	if err != nil {
		t.Fatalf("AddMonths(1) failed: %v", err)
	}

	if clamped.Month() != miti.Shrawan || clamped.Day() != 31 {
		t.Errorf("AddMonths clamping failed: expected 1970-04-31, got %v", clamped)
	}

	// Subtracting 1 month returns to Ashadh 31 (since clamped date was 31st)
	back, err := clamped.AddMonths(-1)
	if err != nil {
		t.Fatalf("AddMonths(-1) failed: %v", err)
	}
	if back.Month() != miti.Ashadh || back.Day() != 31 {
		t.Errorf("expected 1970-03-31, got %v", back)
	}

	// Overflow and zero value tests
	if _, err := d32.AddMonths(math.MaxInt); err != miti.ErrOutOfRange {
		t.Errorf("expected ErrOutOfRange for math.MaxInt months, got %v", err)
	}
	if _, err := (miti.Date{}).AddMonths(1); err != miti.ErrInvalidDate {
		t.Errorf("expected ErrInvalidDate for zero Date, got %v", err)
	}
}
