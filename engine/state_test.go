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

func TestSaveAndLoadState_FormInfo(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	state := &State{
		FormID:       "form_abc",
		ResponderURL: "https://example.com/form",
		ItemIDs:      []string{"i1"},
		QuestionIDs:  []string{"q1"},
		LastApplied:  now,
	}

	loaded := saveAndReload(t, state)

	if loaded.FormID != "form_abc" {
		t.Errorf("FormID = %q", loaded.FormID)
	}
	if loaded.ResponderURL != "https://example.com/form" {
		t.Errorf("ResponderURL = %q", loaded.ResponderURL)
	}
	if !loaded.LastApplied.Equal(now) {
		t.Errorf("LastApplied = %v, want %v", loaded.LastApplied, now)
	}
}

func TestSaveAndLoadState_IDs(t *testing.T) {
	state := &State{
		FormID:      "f",
		ItemIDs:     []string{"i1", "i2", "i3"},
		QuestionIDs: []string{"q1", "q2", ""},
		LastApplied: time.Now(),
	}

	loaded := saveAndReload(t, state)

	if len(loaded.ItemIDs) != 3 || loaded.ItemIDs[2] != "i3" {
		t.Errorf("ItemIDs = %v", loaded.ItemIDs)
	}
	if len(loaded.QuestionIDs) != 3 || loaded.QuestionIDs[2] != "" {
		t.Errorf("QuestionIDs = %v", loaded.QuestionIDs)
	}
}

func saveAndReload(t *testing.T, state *State) *State {
	t.Helper()
	path := filepath.Join(t.TempDir(), "state.json")
	if err := SaveState(path, state); err != nil {
		t.Fatalf("SaveState: %v", err)
	}
	loaded, err := LoadState(path)
	if err != nil {
		t.Fatalf("LoadState: %v", err)
	}
	return loaded
}
