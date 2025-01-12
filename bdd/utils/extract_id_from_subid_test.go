package utils

import "testing"

func TestExtractIdFromSubid(t *testing.T) {
	subid := "1234567890-1234567890"
	id := ExtractIdFromSubid(subid)
	if id != "1234567890" {
		t.Fatalf("expected 1234567890, got %s", id)
	}
}

func TestExtractIdFromSubidWithNoDash(t *testing.T) {
	subid := "1234567890"
	id := ExtractIdFromSubid(subid)
	if id != "1234567890" {
		t.Fatalf("expected 1234567890, got %s", id)
	}
}

func TestExtractIdFromSubidWithDashAtTheEnd(t *testing.T) {
	subid := "1234567890-"
	id := ExtractIdFromSubid(subid)
	if id != "1234567890" {
		t.Fatalf("expected 1234567890, got %s", id)
	}
}

func TestExtractIdFromSubidWithDashAtTheBeginning(t *testing.T) {
	subid := "-1234567890"
	id := ExtractIdFromSubid(subid)
	if id != "" {
		t.Fatalf("expected empty string, got %s", id)
	}
}
