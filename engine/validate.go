package engine

import (
	"fmt"
	"strings"
)

// ValidationError collects multiple validation failures.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("バリデーションエラー (%d件):\n  %s", len(e.Errors), strings.Join(e.Errors, "\n  "))
}

// Validate checks a FormSpec against Google Forms API constraints.
// ref: https://developers.google.com/forms/api/reference/rest/v1/forms
func Validate(spec *FormSpec) error {
	var errs []string

	if spec.Title == "" {
		errs = append(errs, "title は必須です")
	}
	if len(spec.Items) == 0 {
		errs = append(errs, "items は1つ以上必要です")
	}

	for i, item := range spec.Items {
		errs = append(errs, validateItem(i, item)...)
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}

func validateItem(index int, item ItemSpec) []string {
	prefix := fmt.Sprintf("items[%d]", index)
	var errs []string

	if item.Title == "" && item.Type != ItemPageBreak {
		errs = append(errs, fmt.Sprintf("%s: title は必須です", prefix))
	}

	if item.Type == "" {
		return append(errs, fmt.Sprintf("%s: type は必須です", prefix))
	}

	if !item.Type.IsValid() {
		return append(errs, fmt.Sprintf("%s: 不明な type %q (有効値: %s)",
			prefix, item.Type, ValidItemTypes()))
	}

	switch item.Type {
	case ItemChoice:
		errs = append(errs, validateChoice(prefix, item)...)
	case ItemScale:
		errs = append(errs, validateScale(prefix, item)...)
	}
	return errs
}

func validateChoice(prefix string, item ItemSpec) []string {
	if item.Choice == nil {
		return []string{fmt.Sprintf("%s: type=choice には choice フィールドが必須です", prefix)}
	}

	var errs []string
	if !item.Choice.Type.IsValid() {
		errs = append(errs, fmt.Sprintf("%s: 不明な choice.type %q (有効値: %s)",
			prefix, item.Choice.Type, ValidChoiceTypes()))
	}
	if len(item.Choice.Options) == 0 {
		errs = append(errs, fmt.Sprintf("%s: choice.options は1つ以上必要です", prefix))
	}
	errs = append(errs, validateOptionValues(prefix, item.Choice.Options)...)
	return errs
}

func validateOptionValues(prefix string, options []string) []string {
	var errs []string
	for j, opt := range options {
		if opt == "" {
			errs = append(errs, fmt.Sprintf("%s: choice.options[%d] は空にできません", prefix, j))
		}
	}
	return errs
}

func validateScale(prefix string, item ItemSpec) []string {
	if item.Scale == nil {
		return []string{fmt.Sprintf("%s: type=scale には scale フィールドが必須です", prefix)}
	}
	return validateScaleRange(prefix, item.Scale)
}

func validateScaleRange(prefix string, s *ScaleSpec) []string {
	var errs []string
	if s.Low != 0 && s.Low != 1 {
		errs = append(errs, fmt.Sprintf("%s: scale.low は 0 または 1 のみ (got %d)", prefix, s.Low))
	}
	if s.High < 2 || s.High > 10 {
		errs = append(errs, fmt.Sprintf("%s: scale.high は 2〜10 の範囲 (got %d)", prefix, s.High))
	}
	if s.Low >= s.High {
		errs = append(errs, fmt.Sprintf("%s: scale.low (%d) は scale.high (%d) より小さくなければなりません",
			prefix, s.Low, s.High))
	}
	return errs
}
