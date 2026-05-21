package engine

import (
	"strings"
	"testing"
)

func TestValidate_Valid(t *testing.T) {
	spec := &FormSpec{
		Title: "OK",
		Items: []ItemSpec{
			{Title: "Q1", Type: ItemShortAnswer},
			{Title: "Q2", Type: ItemParagraph},
			{Title: "Q3", Type: ItemChoice, Choice: &ChoiceSpec{Type: ChoiceRadio, Options: []string{"A"}}},
			{Title: "Q4", Type: ItemScale, Scale: &ScaleSpec{Low: 1, High: 5}},
			{Title: "Q5", Type: ItemDate},
			{Title: "Q6", Type: ItemTime},
			{Type: ItemPageBreak},
		},
	}
	if err := Validate(spec); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_EmptyTitle(t *testing.T) {
	spec := &FormSpec{Items: []ItemSpec{{Title: "Q", Type: ItemShortAnswer}}}
	assertValidationError(t, spec, "title")
}

func TestValidate_EmptyItems(t *testing.T) {
	assertValidationError(t, &FormSpec{Title: "T"}, "items")
}

func TestValidate_UnknownType(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Title: "Q", Type: "video"}}}
	assertValidationError(t, spec, "unknown type")
}

func TestValidate_MissingType(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Title: "Q"}}}
	assertValidationError(t, spec, "type")
}

func TestValidate_MissingItemTitle(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Type: ItemShortAnswer}}}
	assertValidationError(t, spec, "title")
}

func TestValidate_PageBreakAllowsEmptyTitle(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Type: ItemPageBreak}}}
	if err := Validate(spec); err != nil {
		t.Errorf("page_break should allow empty title: %v", err)
	}
}

func TestValidate_ChoiceMissingSpec(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Title: "Q", Type: ItemChoice}}}
	assertValidationError(t, spec, "choice")
}

func TestValidate_ChoiceInvalidType(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Title: "Q", Type: ItemChoice, Choice: &ChoiceSpec{
		Type: "multiselect", Options: []string{"A"},
	}}}}
	assertValidationError(t, spec, "choice.type")
}

func TestValidate_ChoiceEmptyOptions(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Title: "Q", Type: ItemChoice, Choice: &ChoiceSpec{
		Type: ChoiceRadio,
	}}}}
	assertValidationError(t, spec, "options")
}

func TestValidate_ChoiceEmptyOptionValue(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Title: "Q", Type: ItemChoice, Choice: &ChoiceSpec{
		Type: ChoiceRadio, Options: []string{"A", ""},
	}}}}
	assertValidationError(t, spec, "options[1]")
}

func TestValidate_ScaleMissingSpec(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Title: "Q", Type: ItemScale}}}
	assertValidationError(t, spec, "scale")
}

func TestValidate_ScaleLowInvalid(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Title: "Q", Type: ItemScale, Scale: &ScaleSpec{Low: 2, High: 5}}}}
	assertValidationError(t, spec, "scale.low")
}

func TestValidate_ScaleHighOutOfRange(t *testing.T) {
	cases := []ScaleSpec{
		{Low: 0, High: 1},
		{Low: 1, High: 11},
	}
	for _, s := range cases {
		spec := &FormSpec{Title: "T", Items: []ItemSpec{{Title: "Q", Type: ItemScale, Scale: &s}}}
		assertValidationError(t, spec, "scale.high")
	}
}

func TestValidate_ScaleLowGteHigh(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{{Title: "Q", Type: ItemScale, Scale: &ScaleSpec{Low: 1, High: 1}}}}
	if Validate(spec) == nil {
		t.Fatal("expected error")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	spec := &FormSpec{
		Items: []ItemSpec{
			{Type: ItemShortAnswer},
			{Title: "Q", Type: ItemChoice},
			{Title: "Q", Type: ItemScale},
			{Title: "Q", Type: "alien"},
		},
	}
	err := Validate(spec)
	if err == nil {
		t.Fatal("expected errors")
	}
	ve := err.(*ValidationError)
	if len(ve.Errors) < 5 {
		t.Errorf("expected >= 5 errors, got %d: %v", len(ve.Errors), ve.Errors)
	}
}

func assertValidationError(t *testing.T, spec *FormSpec, substr string) {
	t.Helper()
	err := Validate(spec)
	if err == nil {
		t.Fatal("expected validation error")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	for _, e := range ve.Errors {
		if strings.Contains(e, substr) {
			return
		}
	}
	t.Errorf("no error containing %q in: %v", substr, ve.Errors)
}
