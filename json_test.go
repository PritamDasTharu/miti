package miti_test

import (
	"encoding/json"
	"testing"

	"github.com/PritamDasTharu/miti"
)

type Event struct {
	Name string        `json:"name"`
	Date miti.Date     `json:"date"`
	BS   miti.BSSQLDate `json:"bs,omitempty"`
}

func TestJSONAndTextMarshaling(t *testing.T) {
	d, _ := miti.New(2081, miti.Baisakh, 1)
	event := Event{
		Name: "New Year",
		Date: d,
		BS:   miti.BSSQLDate(d),
	}

	// Marshal
	bytes, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := `{"name":"New Year","date":"2081-01-01","bs":"2081-01-01"}`
	if string(bytes) != expected {
		t.Errorf("Marshal got %s, expected %s", string(bytes), expected)
	}

	// Unmarshal
	var unmarshaled Event
	if err := json.Unmarshal(bytes, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if unmarshaled.Date != d || unmarshaled.BS.Date() != d {
		t.Errorf("Unmarshal mismatch: %+v", unmarshaled)
	}

	// Zero value marshals to null
	zeroEvent := Event{Name: "Zero"}
	zeroBytes, err := json.Marshal(zeroEvent)
	if err != nil {
		t.Fatalf("Marshal zero failed: %v", err)
	}
	expectedZero := `{"name":"Zero","date":null,"bs":null}`
	if string(zeroBytes) != expectedZero {
		t.Errorf("zero got %s, expected %s", string(zeroBytes), expectedZero)
	}

	// Unmarshaling "null" is a no-op (preserves existing struct fields)
	existingDate, _ := miti.New(2081, miti.Baisakh, 1)
	if err := json.Unmarshal([]byte("null"), &existingDate); err != nil {
		t.Fatalf("Unmarshal null failed: %v", err)
	}
	if existingDate.Day() != 1 {
		t.Errorf("Unmarshal null should be no-op, got %v", existingDate)
	}

	// Text marshaling
	txt, err := d.MarshalText()
	if err != nil || string(txt) != "2081-01-01" {
		t.Errorf("MarshalText got %s, err %v", string(txt), err)
	}

	var parsedText miti.Date
	if err := parsedText.UnmarshalText(txt); err != nil || parsedText != d {
		t.Errorf("UnmarshalText got %v, err %v", parsedText, err)
	}

	// Empty text unmarshals to zero Date{}
	var zeroText miti.Date
	if err := zeroText.UnmarshalText([]byte("")); err != nil || !zeroText.IsZero() {
		t.Errorf("UnmarshalText('') expected zero date, got %v", zeroText)
	}
}
