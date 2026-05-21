package engine

import (
	"testing"
)

func TestValidate_Valid(t *testing.T) {
	spec := &FormSpec{
		Title: "OK",
		Items: []ItemSpec{
			{Title: "Q1", Type: "short_answer"},
			{Title: "Q2", Type: "paragraph"},
			{Title: "Q3", Type: "choice", Choice: &ChoiceSpec{Type: "radio", Options: []string{"A"}}},
			{Title: "Q4", Type: "scale", Scale: &ScaleSpec{Low: 1, High: 5}},
			{Title: "Q5", Type: "date"},
			{Title: "Q6", Type: "time"},
			{Title: "", Type: "page_break"},
		},
	}
	if err := Validate(spec); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_EmptyTitle(t *testing.T) {
	spec := &FormSpec{
		Items: []ItemSpec{{Title: "Q", Type: "short_answer"}},
	}
	assertValidationErrors(t, spec, 1, "title は必須")
}

func TestValidate_EmptyItems(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{}}
	assertValidationErrors(t, spec, 1, "items は1つ以上")
}

func TestValidate_UnknownType(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "video"}},
	}
	assertValidationErrors(t, spec, 1, "不明な type")
}

func TestValidate_MissingType(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q"}},
	}
	assertValidationErrors(t, spec, 1, "type は必須")
}

func TestValidate_MissingItemTitle(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Type: "short_answer"}},
	}
	assertValidationErrors(t, spec, 1, "title は必須")
}

func TestValidate_PageBreakAllowsEmptyTitle(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Type: "page_break"}},
	}
	if err := Validate(spec); err != nil {
		t.Errorf("page_break should allow empty title: %v", err)
	}
}

func TestValidate_ChoiceMissingSpec(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "choice"}},
	}
	assertValidationErrors(t, spec, 1, "choice フィールドが必須")
}

func TestValidate_ChoiceInvalidType(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "choice", Choice: &ChoiceSpec{
			Type: "multiselect", Options: []string{"A"},
		}}},
	}
	assertValidationErrors(t, spec, 1, "不明な choice.type")
}

func TestValidate_ChoiceEmptyOptions(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "choice", Choice: &ChoiceSpec{
			Type: "radio", Options: []string{},
		}}},
	}
	assertValidationErrors(t, spec, 1, "options は1つ以上")
}

func TestValidate_ChoiceEmptyOptionValue(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "choice", Choice: &ChoiceSpec{
			Type: "radio", Options: []string{"A", ""},
		}}},
	}
	assertValidationErrors(t, spec, 1, "空にできません")
}

func TestValidate_ScaleMissingSpec(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "scale"}},
	}
	assertValidationErrors(t, spec, 1, "scale フィールドが必須")
}

func TestValidate_ScaleLowInvalid(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "scale", Scale: &ScaleSpec{Low: 2, High: 5}}},
	}
	assertValidationErrors(t, spec, 1, "scale.low は 0 または 1")
}

func TestValidate_ScaleHighTooLow(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "scale", Scale: &ScaleSpec{Low: 0, High: 1}}},
	}
	assertValidationErrors(t, spec, 1, "scale.high は 2〜10")
}

func TestValidate_ScaleHighTooHigh(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "scale", Scale: &ScaleSpec{Low: 1, High: 11}}},
	}
	assertValidationErrors(t, spec, 1, "scale.high は 2〜10")
}

func TestValidate_ScaleLowGteHigh(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "scale", Scale: &ScaleSpec{Low: 1, High: 1}}},
	}
	err := Validate(spec)
	if err == nil {
		t.Fatal("expected error")
	}
	ve := err.(*ValidationError)
	// scale.high は 2〜10 + low >= high
	if len(ve.Errors) < 1 {
		t.Errorf("expected errors, got %d", len(ve.Errors))
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	spec := &FormSpec{
		Items: []ItemSpec{
			{Type: "short_answer"},           // missing title
			{Title: "Q", Type: "choice"},     // missing choice spec
			{Title: "Q", Type: "scale"},      // missing scale spec
			{Title: "Q", Type: "alien_type"}, // unknown type
		},
	}
	err := Validate(spec)
	if err == nil {
		t.Fatal("expected errors")
	}
	ve := err.(*ValidationError)
	// title必須 + items[0]title + choice必須 + scale必須 + unknown type = 5
	if len(ve.Errors) < 5 {
		t.Errorf("expected at least 5 errors, got %d: %v", len(ve.Errors), ve.Errors)
	}
}

func assertValidationErrors(t *testing.T, spec *FormSpec, minErrors int, substr string) {
	t.Helper()
	err := Validate(spec)
	if err == nil {
		t.Fatal("expected validation error")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errors) < minErrors {
		t.Errorf("expected at least %d errors, got %d: %v", minErrors, len(ve.Errors), ve.Errors)
	}
	found := false
	for _, e := range ve.Errors {
		if contains(e, substr) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("no error containing %q in: %v", substr, ve.Errors)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
