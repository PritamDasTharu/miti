package miti_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/PritamDasTharu/miti"
)

type GoldenBaisakh1Entry struct {
	BSYear   int    `json:"bs_year"`
	ADDate   string `json:"ad_date"`
	Weekday  string `json:"weekday"`
	Source   string `json:"source"`
	Verified bool   `json:"verified"`
}

type GoldenAnchorEntry struct {
	BSDate      string `json:"bs_date"`
	ADDate      string `json:"ad_date"`
	Weekday     string `json:"weekday"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Verified    bool   `json:"verified"`
}

func TestGoldenBaisakh1(t *testing.T) {
	path := filepath.Join("data", "golden_baisakh_1.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read golden_baisakh_1.json: %v", err)
	}

	var entries []GoldenBaisakh1Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("failed to parse golden_baisakh_1.json: %v", err)
	}

	if len(entries) != 131 {
		t.Fatalf("expected 131 entries, got %d", len(entries))
	}

	for _, e := range entries {
		// Provenance enforcement: must have non-empty source
		if e.Source == "" {
			t.Errorf("entry for BS year %d is missing source citation", e.BSYear)
		}

		bsDate, err := miti.New(e.BSYear, miti.Baisakh, 1)
		if err != nil {
			t.Fatalf("failed to create BS date for %d-01-01: %v", e.BSYear, err)
		}

		adTime := bsDate.ToTime()
		gotAD := adTime.Format("2006-01-02")
		if gotAD != e.ADDate {
			t.Errorf("BS %d-01-01: expected AD %s, got %s", e.BSYear, e.ADDate, gotAD)
		}

		// Also verify reverse conversion
		backBS, err := miti.FromTime(adTime)
		if err != nil {
			t.Errorf("FromTime(%s) returned error: %v", gotAD, err)
		} else if backBS != bsDate {
			t.Errorf("FromTime(%s) = %v, expected %v", gotAD, backBS, bsDate)
		}
	}
}

func TestMidYearDecadeAnchors(t *testing.T) {
	path := filepath.Join("data", "golden_anchors.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read golden_anchors.json: %v", err)
	}

	var anchors []GoldenAnchorEntry
	if err := json.Unmarshal(data, &anchors); err != nil {
		t.Fatalf("failed to parse golden_anchors.json: %v", err)
	}

	for _, a := range anchors {
		if a.Source == "" {
			t.Errorf("anchor %s (%s) is missing source citation", a.BSDate, a.Description)
		}

		// Parse BS date
		adParsed, err := time.Parse("2006-01-02", a.ADDate)
		if err != nil {
			t.Fatalf("invalid AD date in anchor %s: %v", a.ADDate, err)
		}

		// Convert AD to BS
		bsDate, err := miti.FromTime(adParsed)
		if err != nil {
			t.Fatalf("FromTime(%s) for %s failed: %v", a.ADDate, a.Description, err)
		}

		if bsDate.String() != a.BSDate {
			t.Errorf("%s: FromTime(%s) = %s, expected BS %s", a.Description, a.ADDate, bsDate, a.BSDate)
		}

		// Convert back BS to AD
		adResult := bsDate.ToTime()
		if adResult.Format("2006-01-02") != a.ADDate {
			t.Errorf("%s: ToTime() = %s, expected AD %s", a.Description, adResult.Format("2006-01-02"), a.ADDate)
		}

		// Check weekday
		if bsDate.Weekday().String() != a.Weekday {
			t.Errorf("%s: Weekday() = %s, expected %s", a.Description, bsDate.Weekday(), a.Weekday)
		}
	}
}
