package engine

import (
	"path/filepath"
	"testing"
	"time"
)

func TestLoadState_NotFound(t *testing.T) {
	state, err := LoadState("/nonexistent/state.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != nil {
		t.Error("expected nil for missing file")
	}
}

func TestSaveAndLoadState(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	state := &State{
		FormID:       "form_abc",
		ResponderURL: "https://example.com/form",
		ItemIDs:      []string{"i1", "i2", "i3"},
		QuestionIDs:  []string{"q1", "q2", ""},
		LastApplied:  now,
	}

	path := filepath.Join(t.TempDir(), "state.json")
	if err := SaveState(path, state); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	loaded, err := LoadState(path)
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}

	if loaded.FormID != "form_abc" {
		t.Errorf("FormID = %q", loaded.FormID)
	}
	if loaded.ResponderURL != "https://example.com/form" {
		t.Errorf("ResponderURL = %q", loaded.ResponderURL)
	}
	if len(loaded.ItemIDs) != 3 || loaded.ItemIDs[2] != "i3" {
		t.Errorf("ItemIDs = %v", loaded.ItemIDs)
	}
	if len(loaded.QuestionIDs) != 3 || loaded.QuestionIDs[2] != "" {
		t.Errorf("QuestionIDs = %v", loaded.QuestionIDs)
	}
	if !loaded.LastApplied.Equal(now) {
		t.Errorf("LastApplied = %v, want %v", loaded.LastApplied, now)
	}
}
