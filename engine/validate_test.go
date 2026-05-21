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
	spec := &FormSpec{
		Items: []ItemSpec{{Title: "Q", Type: ItemShortAnswer}},
	}
	assertValidationError(t, spec, "title は必須")
}

func TestValidate_EmptyItems(t *testing.T) {
	spec := &FormSpec{Title: "T", Items: []ItemSpec{}}
	assertValidationError(t, spec, "items は1つ以上")
}

func TestValidate_UnknownType(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: "video"}},
	}
	assertValidationError(t, spec, "不明な type")
}

func TestValidate_MissingType(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q"}},
	}
	assertValidationError(t, spec, "type は必須")
}

func TestValidate_MissingItemTitle(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Type: ItemShortAnswer}},
	}
	assertValidationError(t, spec, "title は必須")
}

func TestValidate_PageBreakAllowsEmptyTitle(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Type: ItemPageBreak}},
	}
	if err := Validate(spec); err != nil {
		t.Errorf("page_break should allow empty title: %v", err)
	}
}

func TestValidate_ChoiceMissingSpec(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemChoice}},
	}
	assertValidationError(t, spec, "choice フィールドが必須")
}

func TestValidate_ChoiceInvalidType(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemChoice, Choice: &ChoiceSpec{
			Type: "multiselect", Options: []string{"A"},
		}}},
	}
	assertValidationError(t, spec, "不明な choice.type")
}

func TestValidate_ChoiceEmptyOptions(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemChoice, Choice: &ChoiceSpec{
			Type: ChoiceRadio,
		}}},
	}
	assertValidationError(t, spec, "options は1つ以上")
}

func TestValidate_ChoiceEmptyOptionValue(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemChoice, Choice: &ChoiceSpec{
			Type: ChoiceRadio, Options: []string{"A", ""},
		}}},
	}
	assertValidationError(t, spec, "空にできません")
}

func TestValidate_ScaleMissingSpec(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemScale}},
	}
	assertValidationError(t, spec, "scale フィールドが必須")
}

func TestValidate_ScaleLowInvalid(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemScale, Scale: &ScaleSpec{Low: 2, High: 5}}},
	}
	assertValidationError(t, spec, "scale.low は 0 または 1")
}

func TestValidate_ScaleHighTooLow(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemScale, Scale: &ScaleSpec{Low: 0, High: 1}}},
	}
	assertValidationError(t, spec, "scale.high は 2〜10")
}

func TestValidate_ScaleHighTooHigh(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemScale, Scale: &ScaleSpec{Low: 1, High: 11}}},
	}
	assertValidationError(t, spec, "scale.high は 2〜10")
}

func TestValidate_ScaleLowGteHigh(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q", Type: ItemScale, Scale: &ScaleSpec{Low: 1, High: 1}}},
	}
	err := Validate(spec)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	spec := &FormSpec{
		Items: []ItemSpec{
			{Type: ItemShortAnswer},        // missing title
			{Title: "Q", Type: ItemChoice}, // missing choice spec
			{Title: "Q", Type: ItemScale},  // missing scale spec
			{Title: "Q", Type: "alien"},    // unknown type
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
