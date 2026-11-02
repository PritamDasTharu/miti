//go:build !race

package miti_test

import (
	"testing"
	"time"

	"github.com/PritamDasTharu/miti"
)

func TestZeroAllocations(t *testing.T) {
	d, _ := miti.New(2081, miti.Baisakh, 15)
	now := time.Date(2024, 4, 27, 10, 0, 0, 0, time.UTC)
	preSizedBuf := make([]byte, 0, 32)

	// 1. FromTime zero allocs
	allocs := testing.AllocsPerRun(1000, func() {
		_, _ = miti.FromTime(now)
	})
	if allocs > 0 {
		t.Errorf("FromTime allocated %v allocs/run, expected 0", allocs)
	}

	// 2. FromTimeNepal zero allocs
	allocs = testing.AllocsPerRun(1000, func() {
		_, _ = miti.FromTimeNepal(now)
	})
	if allocs > 0 {
		t.Errorf("FromTimeNepal allocated %v allocs/run, expected 0", allocs)
	}

	// 3. ToTime zero allocs
	allocs = testing.AllocsPerRun(1000, func() {
		_ = d.ToTime()
	})
	if allocs > 0 {
		t.Errorf("ToTime allocated %v allocs/run, expected 0", allocs)
	}

	// 4. AddDays zero allocs
	allocs = testing.AllocsPerRun(1000, func() {
		_, _ = d.AddDays(10)
	})
	if allocs > 0 {
		t.Errorf("AddDays allocated %v allocs/run, expected 0", allocs)
	}

	// 5. DaysInMonth zero allocs
	allocs = testing.AllocsPerRun(1000, func() {
		_ = d.DaysInMonth()
	})
	if allocs > 0 {
		t.Errorf("DaysInMonth allocated %v allocs/run, expected 0", allocs)
	}

	// 6. Pre-sized AppendFormat zero allocs
	allocs = testing.AllocsPerRun(1000, func() {
		b := preSizedBuf[:0]
		_ = d.AppendFormat(b, "YYYY-MM-DD")
	})
	if allocs > 0 {
		t.Errorf("AppendFormat allocated %v allocs/run, expected 0", allocs)
	}
}
