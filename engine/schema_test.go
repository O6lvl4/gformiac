package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSpec(t *testing.T) {
	yaml := `
title: "Test Form"
description: "A test"
items:
  - title: "Name"
    type: short_answer
    required: true
  - title: "Rating"
    type: choice
    choice:
      type: radio
      options:
        - "Good"
        - "Bad"
  - title: "Score"
    type: scale
    scale:
      low: 1
      high: 5
      low_label: "Low"
      high_label: "High"
`
	path := filepath.Join(t.TempDir(), "form.yaml")
	os.WriteFile(path, []byte(yaml), 0644)

	spec, err := LoadSpec(path)
	if err != nil {
		t.Fatalf("LoadSpec failed: %v", err)
	}

	if spec.Title != "Test Form" {
		t.Errorf("title = %q, want %q", spec.Title, "Test Form")
	}
	if spec.Description != "A test" {
		t.Errorf("description = %q, want %q", spec.Description, "A test")
	}
	if len(spec.Items) != 3 {
		t.Fatalf("items count = %d, want 3", len(spec.Items))
	}

	// short_answer
	if spec.Items[0].Type != "short_answer" || !spec.Items[0].Required {
		t.Errorf("item[0] = %+v", spec.Items[0])
	}

	// choice
	item1 := spec.Items[1]
	if item1.Type != "choice" || item1.Choice == nil {
		t.Fatalf("item[1] type/choice unexpected: %+v", item1)
	}
	if item1.Choice.Type != "radio" || len(item1.Choice.Options) != 2 {
		t.Errorf("item[1].choice = %+v", item1.Choice)
	}

	// scale
	item2 := spec.Items[2]
	if item2.Type != "scale" || item2.Scale == nil {
		t.Fatalf("item[2] type/scale unexpected: %+v", item2)
	}
	if item2.Scale.Low != 1 || item2.Scale.High != 5 {
		t.Errorf("item[2].scale = %+v", item2.Scale)
	}
}

func TestLoadSpec_MissingTitle(t *testing.T) {
	yaml := `description: "no title"`
	path := filepath.Join(t.TempDir(), "form.yaml")
	os.WriteFile(path, []byte(yaml), 0644)

	_, err := LoadSpec(path)
	if err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestLoadSpec_FileNotFound(t *testing.T) {
	_, err := LoadSpec("/nonexistent/form.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	spec := &FormSpec{
		Title:       "Round Trip",
		Description: "Test",
		Items: []ItemSpec{
			{Title: "Q1", Type: "short_answer", Required: true},
			{Title: "Q2", Type: "choice", Choice: &ChoiceSpec{
				Type:    "checkbox",
				Options: []string{"A", "B", "C"},
			}},
		},
	}

	path := filepath.Join(t.TempDir(), "out.yaml")
	if err := SaveSpec(path, spec); err != nil {
		t.Fatalf("SaveSpec failed: %v", err)
	}

	loaded, err := LoadSpec(path)
	if err != nil {
		t.Fatalf("LoadSpec failed: %v", err)
	}

	if loaded.Title != spec.Title || len(loaded.Items) != 2 {
		t.Errorf("round trip mismatch: got %+v", loaded)
	}
	if loaded.Items[1].Choice.Type != "checkbox" || len(loaded.Items[1].Choice.Options) != 3 {
		t.Errorf("choice round trip mismatch: got %+v", loaded.Items[1].Choice)
	}
}
