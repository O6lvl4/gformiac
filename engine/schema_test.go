package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSpec_BasicFields(t *testing.T) {
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
`
	spec := loadSpecFromString(t, yaml)

	if spec.Title != "Test Form" {
		t.Errorf("title = %q, want %q", spec.Title, "Test Form")
	}
	if spec.Description != "A test" {
		t.Errorf("description = %q, want %q", spec.Description, "A test")
	}
	if len(spec.Items) != 3 {
		t.Fatalf("items count = %d, want 3", len(spec.Items))
	}
}

func TestLoadSpec_ShortAnswer(t *testing.T) {
	spec := loadSpecFromString(t, `
title: "T"
items:
  - title: "Q"
    type: short_answer
    required: true
`)
	if spec.Items[0].Type != ItemShortAnswer || !spec.Items[0].Required {
		t.Errorf("item = %+v", spec.Items[0])
	}
}

func TestLoadSpec_Choice(t *testing.T) {
	spec := loadSpecFromString(t, `
title: "T"
items:
  - title: "Q"
    type: choice
    choice:
      type: radio
      options: ["A", "B"]
`)
	item := spec.Items[0]
	if item.Type != ItemChoice || item.Choice == nil {
		t.Fatalf("unexpected: %+v", item)
	}
	if item.Choice.Type != ChoiceRadio || len(item.Choice.Options) != 2 {
		t.Errorf("choice = %+v", item.Choice)
	}
}

func TestLoadSpec_Scale(t *testing.T) {
	spec := loadSpecFromString(t, `
title: "T"
items:
  - title: "Q"
    type: scale
    scale:
      low: 1
      high: 5
`)
	item := spec.Items[0]
	if item.Type != ItemScale || item.Scale == nil {
		t.Fatalf("unexpected: %+v", item)
	}
	if item.Scale.Low != 1 || item.Scale.High != 5 {
		t.Errorf("scale = %+v", item.Scale)
	}
}

func TestLoadSpec_MissingTitle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "form.yaml")
	os.WriteFile(path, []byte(`description: "no title"`), 0644)
	if _, err := LoadSpec(path); err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestLoadSpec_FileNotFound(t *testing.T) {
	if _, err := LoadSpec("/nonexistent/form.yaml"); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	spec := &FormSpec{
		Title:       "Round Trip",
		Description: "Test",
		Items: []ItemSpec{
			{Title: "Q1", Type: ItemShortAnswer, Required: true},
			{Title: "Q2", Type: ItemChoice, Choice: &ChoiceSpec{
				Type: ChoiceCheckbox, Options: []string{"A", "B", "C"},
			}},
		},
	}

	path := filepath.Join(t.TempDir(), "out.yaml")
	if err := SaveSpec(path, spec); err != nil {
		t.Fatalf("SaveSpec: %v", err)
	}
	loaded, err := LoadSpec(path)
	if err != nil {
		t.Fatalf("LoadSpec: %v", err)
	}

	if loaded.Title != spec.Title || len(loaded.Items) != 2 {
		t.Errorf("mismatch: got %+v", loaded)
	}
	if loaded.Items[1].Choice.Type != ChoiceCheckbox {
		t.Errorf("choice type mismatch: %+v", loaded.Items[1].Choice)
	}
}

func TestItemType_IsValid(t *testing.T) {
	for _, v := range []ItemType{ItemShortAnswer, ItemParagraph, ItemChoice, ItemScale, ItemDate, ItemTime, ItemPageBreak} {
		if !v.IsValid() {
			t.Errorf("%q should be valid", v)
		}
	}
	if ItemType("bogus").IsValid() {
		t.Error("bogus should be invalid")
	}
}

func TestChoiceType_IsValid(t *testing.T) {
	for _, v := range []ChoiceType{ChoiceRadio, ChoiceCheckbox, ChoiceDropdown} {
		if !v.IsValid() {
			t.Errorf("%q should be valid", v)
		}
	}
	if ChoiceType("multi").IsValid() {
		t.Error("multi should be invalid")
	}
}

func loadSpecFromString(t *testing.T, yaml string) *FormSpec {
	t.Helper()
	path := filepath.Join(t.TempDir(), "form.yaml")
	os.WriteFile(path, []byte(yaml), 0644)
	spec, err := LoadSpec(path)
	if err != nil {
		t.Fatalf("LoadSpec: %v", err)
	}
	return spec
}
