package api

import (
	"encoding/json"
	"testing"
)

func TestFlexBoolUnmarshalBoolean(t *testing.T) {
	var s struct {
		Flag FlexBool `json:"flag"`
	}

	if err := json.Unmarshal([]byte(`{"flag":true}`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Flag {
		t.Error("expected true")
	}

	if err := json.Unmarshal([]byte(`{"flag":false}`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Flag {
		t.Error("expected false")
	}
}

func TestFlexBoolUnmarshalInteger(t *testing.T) {
	var s struct {
		Flag FlexBool `json:"flag"`
	}

	// SQLite/Sequelize returns 0/1
	if err := json.Unmarshal([]byte(`{"flag":1}`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Flag {
		t.Error("expected true for integer 1")
	}

	if err := json.Unmarshal([]byte(`{"flag":0}`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Flag {
		t.Error("expected false for integer 0")
	}
}

func TestFlexBoolUnmarshalFloat(t *testing.T) {
	var s struct {
		Flag FlexBool `json:"flag"`
	}

	if err := json.Unmarshal([]byte(`{"flag":1.0}`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Flag {
		t.Error("expected true for float 1.0")
	}

	if err := json.Unmarshal([]byte(`{"flag":0.0}`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Flag {
		t.Error("expected false for float 0.0")
	}
}

func TestFlexBoolUnmarshalInvalid(t *testing.T) {
	var s struct {
		Flag FlexBool `json:"flag"`
	}

	if err := json.Unmarshal([]byte(`{"flag":"yes"}`), &s); err == nil {
		t.Error("expected error for string value")
	}

	if err := json.Unmarshal([]byte(`{"flag":2}`), &s); err != nil {
		t.Fatalf("unexpected error for integer 2: %v", err)
	}
	if !s.Flag {
		t.Error("expected truthy for non-zero integer")
	}
}

func TestFlexBoolInLinkStruct(t *testing.T) {
	// Simulate API response from SQLite with integer isHidden
	data := []byte(`{"url":"https://example.com","isHidden":1}`)
	var link Link
	if err := json.Unmarshal(data, &link); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !link.IsHidden {
		t.Error("expected IsHidden=true for integer 1")
	}

	// Simulate normal boolean response
	data = []byte(`{"url":"https://example.com","isHidden":false}`)
	if err := json.Unmarshal(data, &link); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link.IsHidden {
		t.Error("expected IsHidden=false for boolean false")
	}
}
