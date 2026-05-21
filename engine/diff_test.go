package engine

import (
	"strings"
	"testing"

	"github.com/O6lvl4/gformiac/locale"
)

func TestDiff_NoChanges(t *testing.T) {
	spec := &FormSpec{
		Title: "Same",
		Items: []ItemSpec{{Title: "Q1", Type: ItemShortAnswer}},
	}
	remote := &FormSpec{
		Title: "Same",
		Items: []ItemSpec{{Title: "Q1", Type: ItemShortAnswer}},
	}

	diff := Diff(spec, remote, nil)
	if diff.HasChanges() {
		t.Error("expected no changes")
	}
	if diff.String() != locale.M.NoChanges {
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
			{Title: "Q1", Type: ItemShortAnswer},
			{Title: "Q2", Type: ItemParagraph},
		},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q1", Type: ItemShortAnswer}},
	}

	diff := Diff(local, remote, nil)
	if len(diff.Changes) != 1 {
		t.Fatalf("changes = %d, want 1", len(diff.Changes))
	}
	if diff.Changes[0].Type != ChangeCreate || diff.Changes[0].New.Title != "Q2" {
		t.Errorf("got %+v", diff.Changes[0])
	}
}

func TestDiff_ItemDeleted(t *testing.T) {
	local := &FormSpec{Title: "T"}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q1", Type: ItemShortAnswer}},
	}

	diff := Diff(local, remote, nil)
	if len(diff.Changes) != 1 || diff.Changes[0].Type != ChangeDelete {
		t.Errorf("expected 1 delete, got %v", diff.Changes)
	}
}

func TestDiff_ItemUpdated(t *testing.T) {
	local := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q1 Updated", Type: ItemParagraph, Required: true}},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q1", Type: ItemShortAnswer}},
	}

	diff := Diff(local, remote, nil)
	if len(diff.Changes) != 1 || diff.Changes[0].Type != ChangeUpdate {
		t.Fatalf("expected 1 update, got %v", diff.Changes)
	}
	c := diff.Changes[0]
	if c.Old.Title != "Q1" || c.New.Title != "Q1 Updated" {
		t.Errorf("old=%q new=%q", c.Old.Title, c.New.Title)
	}
}

func TestDiff_ChoiceOptionsChanged(t *testing.T) {
	local := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemChoice, Choice: &ChoiceSpec{Type: ChoiceRadio, Options: []string{"A", "B", "C"}}}},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemChoice, Choice: &ChoiceSpec{Type: ChoiceRadio, Options: []string{"A", "B"}}}},
	}

	diff := Diff(local, remote, nil)
	if len(diff.Changes) != 1 || diff.Changes[0].Type != ChangeUpdate {
		t.Errorf("expected 1 update, got %v", diff.Changes)
	}
}

func TestDiff_ScaleChanged(t *testing.T) {
	local := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemScale, Scale: &ScaleSpec{Low: 1, High: 10}}},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemScale, Scale: &ScaleSpec{Low: 1, High: 5}}},
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
			{Title: "New", Type: ItemShortAnswer},
			{Title: "Changed", Type: ItemParagraph},
		},
	}
	remote := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Old", Type: ItemShortAnswer}},
	}

	s := Diff(local, remote, nil).String()
	if !strings.Contains(s, "+1") || !strings.Contains(s, "~1") {
		t.Errorf("summary: %s", s)
	}
}

func TestNewFormSummary(t *testing.T) {
	spec := &FormSpec{
		Title:       "My Form",
		Description: "Desc",
		Items: []ItemSpec{
			{Title: "Q1", Type: ItemShortAnswer, Required: true},
			{Title: "Q2", Type: ItemChoice, Choice: &ChoiceSpec{Type: ChoiceCheckbox, Options: []string{"A", "B"}}},
		},
	}

	s := NewFormSummary(spec)
	if !strings.Contains(s, "My Form") || !strings.Contains(s, "checkbox") {
		t.Errorf("summary: %s", s)
	}
}

func TestChangeType_String(t *testing.T) {
	cases := []struct {
		t    ChangeType
		want string
	}{
		{ChangeCreate, "create"},
		{ChangeUpdate, "update"},
		{ChangeDelete, "delete"},
		{ChangeType(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.t.String(); got != c.want {
			t.Errorf("ChangeType(%d).String() = %q, want %q", c.t, got, c.want)
		}
	}
}
