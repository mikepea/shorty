package cuid

import (
	"testing"
)

func TestNew(t *testing.T) {
	id := New()
	if len(id) != 24 {
		t.Errorf("expected CUID length 24, got %d", len(id))
	}
	if !IsValid(id) {
		t.Errorf("generated CUID %q is not valid", id)
	}
}

func TestNew_Unique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := New()
		if seen[id] {
			t.Errorf("duplicate CUID generated: %s", id)
		}
		seen[id] = true
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name  string
		id    string
		valid bool
	}{
		{"valid CUID", "clh3fm6b90000qf6p8h6w7y1k", true},
		{"generated CUID", New(), true},
		{"empty string", "", false},
		{"short valid format", "abc123", true}, // library accepts shorter CUIDs
		{"starts with number", "1lh3fm6b90000qf6p8h6w7y1", false},
		{"uppercase", "CLH3FM6B90000QF6P8H6W7Y1K", false},
		{"integer string", "123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.id); got != tt.valid {
				t.Errorf("IsValid(%q) = %v, want %v", tt.id, got, tt.valid)
			}
		})
	}
}
