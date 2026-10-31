package miti_test

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/PritamDasTharu/miti"
)

func TestSQLValuerAndScanner(t *testing.T) {
	// 1. Zero value returns nil (SQL NULL)
	var zero miti.Date
	val, err := zero.Value()
	if err != nil {
		t.Fatalf("zero.Value() error: %v", err)
	}
	if val != nil {
		t.Errorf("zero.Value() expected nil, got %v", val)
	}

	// 2. Scan nil sets Date{}
	var scanned miti.Date
	if err := scanned.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) error: %v", err)
	}
	if !scanned.IsZero() {
		t.Errorf("Scan(nil) expected zero date, got %v", scanned)
	}

	// 3. Valid date Value() returns AD time.Time at midnight UTC
	d, _ := miti.New(2072, miti.Ashoj, 3) // Constitution Day = 2015-09-20 AD
	v, err := d.Value()
	if err != nil {
		t.Fatalf("d.Value() error: %v", err)
	}

	tm, ok := v.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time from Value(), got %T", v)
	}
	if tm.Format("2006-01-02") != "2015-09-20" {
		t.Errorf("expected 2015-09-20, got %v", tm)
	}

	// 4. Scan(time.Time) converts AD back to BS
	var roundtrip miti.Date
	if err := roundtrip.Scan(tm); err != nil {
		t.Fatalf("Scan(time.Time) error: %v", err)
	}
	if roundtrip != d {
		t.Errorf("roundtrip mismatch: got %v, expected %v", roundtrip, d)
	}
}

func TestSQLDriverSymmetry(t *testing.T) {
	// Drivers like MySQL (parseTime=false) and SQLite return SQL DATE as text bytes/strings
	// Value() returned AD 2015-09-20 -> driver returns []byte("2015-09-20")
	orig, _ := miti.New(2072, miti.Ashoj, 3)

	// Simulated text returned from database DATE column
	driverOutput := []byte("2015-09-20")

	var scanned miti.Date
	if err := scanned.Scan(driverOutput); err != nil {
		t.Fatalf("Scan(driverOutput) failed: %v", err)
	}

	if scanned != orig {
		t.Errorf("symmetry failed: scanned %v, expected %v", scanned, orig)
	}
}

func TestBSSQLDate(t *testing.T) {
	d, _ := miti.New(2081, miti.Baisakh, 1)
	bs := miti.BSSQLDate(d)

	// Value returns string
	val, err := bs.Value()
	if err != nil {
		t.Fatalf("bs.Value() error: %v", err)
	}
	if s, ok := val.(string); !ok || s != "2081-01-01" {
		t.Errorf("expected string '2081-01-01', got %v", val)
	}

	// Scan parses BS string
	var scanned miti.BSSQLDate
	if err := scanned.Scan("2081-01-01"); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	if scanned.Date() != d {
		t.Errorf("scanned date %v != original %v", scanned.Date(), d)
	}

	// Zero value
	var zero miti.BSSQLDate
	zv, err := zero.Value()
	if err != nil || zv != nil {
		t.Errorf("expected nil for zero BSSQLDate, got %v", zv)
	}

	var scannedZero miti.BSSQLDate
	if err := scannedZero.Scan(nil); err != nil || !scannedZero.Date().IsZero() {
		t.Errorf("expected zero BSSQLDate for Scan(nil), got %v", scannedZero)
	}

	// Check interface assertions
	var _ driver.Valuer = miti.BSSQLDate{}
	var _ = (*miti.BSSQLDate)(nil)
}
