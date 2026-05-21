package engine

import (
	"strings"
	"testing"
)

func TestDiff_NoChanges(t *testing.T) {
	spec := &FormSpec{
		Title: "Same",
		Items: []ItemSpec{
			{Title: "Q1", Type: "short_answer"},
		},
	}
	remote := &FormSpec{
		Title: "Same",
		Items: []ItemSpec{
			{Title: "Q1", Type: "short_answer"},
		},
	}

	diff := Diff(spec, remote, nil)
	if diff.HasChanges() {
		t.Error("expected no changes")
	}
	if diff.String() != "変更なし" {
		t.Errorf("String() = %q", diff.String())
	}
}

func TestDiff_InfoChanged(t *testing.T) {
	local := &FormSpec{Title: "New Title", Description: "New Desc"}
	remote := &FormSpec{Title: "Old Title", Description: "Old Desc"}

	diff := Diff(local, remote, nil)
	if !diff.InfoChanged {
		t.Fatal("expected info changed")
	}
	if len(diff.InfoDetails) != 2 {
		t.Errorf("InfoDetails count = %d, want 2", len(diff.InfoDetails))
	}
}

func TestDiff_ItemCreated(t *testing.T) {
	local := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "Q1", Type: "short_answer"},
			{Title: "Q2", Type: "paragraph"},
		},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "Q1", Type: "short_answer"},
		},
	}

	diff := Diff(local, remote, nil)
	if len(diff.Changes) != 1 {
		t.Fatalf("changes = %d, want 1", len(diff.Changes))
	}
	if diff.Changes[0].Type != ChangeCreate {
		t.Errorf("change type = %v, want Create", diff.Changes[0].Type)
	}
	if diff.Changes[0].New.Title != "Q2" {
		t.Errorf("new title = %q", diff.Changes[0].New.Title)
	}
}

func TestDiff_ItemDeleted(t *testing.T) {
	local := &FormSpec{
		Title: "T",
		Items: []ItemSpec{},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "Q1", Type: "short_answer"},
		},
	}

	diff := Diff(local, remote, nil)
	if len(diff.Changes) != 1 {
		t.Fatalf("changes = %d, want 1", len(diff.Changes))
	}
	if diff.Changes[0].Type != ChangeDelete {
		t.Errorf("change type = %v, want Delete", diff.Changes[0].Type)
	}
}

func TestDiff_ItemUpdated(t *testing.T) {
	local := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "Q1 Updated", Type: "paragraph", Required: true},
		},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "Q1", Type: "short_answer", Required: false},
		},
	}

	diff := Diff(local, remote, nil)
	if len(diff.Changes) != 1 {
		t.Fatalf("changes = %d, want 1", len(diff.Changes))
	}
	c := diff.Changes[0]
	if c.Type != ChangeUpdate {
		t.Errorf("change type = %v, want Update", c.Type)
	}
	if c.Old.Title != "Q1" || c.New.Title != "Q1 Updated" {
		t.Errorf("old=%q new=%q", c.Old.Title, c.New.Title)
	}
}

func TestDiff_ChoiceOptionsChanged(t *testing.T) {
	local := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "Q", Type: "choice", Choice: &ChoiceSpec{Type: "radio", Options: []string{"A", "B", "C"}}},
		},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "Q", Type: "choice", Choice: &ChoiceSpec{Type: "radio", Options: []string{"A", "B"}}},
		},
	}

	diff := Diff(local, remote, nil)
	if len(diff.Changes) != 1 || diff.Changes[0].Type != ChangeUpdate {
		t.Errorf("expected 1 update, got %v", diff.Changes)
	}
}

func TestDiff_ScaleChanged(t *testing.T) {
	local := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "Q", Type: "scale", Scale: &ScaleSpec{Low: 1, High: 10}},
		},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "Q", Type: "scale", Scale: &ScaleSpec{Low: 1, High: 5}},
		},
	}

	diff := Diff(local, remote, nil)
	if len(diff.Changes) != 1 || diff.Changes[0].Type != ChangeUpdate {
		t.Errorf("expected 1 update, got %v", diff.Changes)
	}
}

func TestDiff_String_Summary(t *testing.T) {
	local := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "New", Type: "short_answer"},
			{Title: "Changed", Type: "paragraph"},
		},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{
			{Title: "Old", Type: "short_answer"},
		},
	}

	diff := Diff(local, remote, nil)
	s := diff.String()

	if !strings.Contains(s, "+1") || !strings.Contains(s, "~1") {
		t.Errorf("summary missing counts: %s", s)
	}
}

func TestNewFormSummary(t *testing.T) {
	spec := &FormSpec{
		Title:       "My Form",
		Description: "Desc",
		Items: []ItemSpec{
			{Title: "Q1", Type: "short_answer", Required: true},
			{Title: "Q2", Type: "choice", Choice: &ChoiceSpec{Type: "checkbox", Options: []string{"A", "B"}}},
		},
	}

	s := NewFormSummary(spec)
	if !strings.Contains(s, "My Form") {
		t.Errorf("missing title in summary: %s", s)
	}
	if !strings.Contains(s, "2項目を作成") {
		t.Errorf("missing item count in summary: %s", s)
	}
	if !strings.Contains(s, "checkbox") {
		t.Errorf("missing choice detail in summary: %s", s)
	}
}
