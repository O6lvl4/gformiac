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
		prefix := fmt.Sprintf("items[%d]", i)

		if item.Title == "" && item.Type != ItemPageBreak {
			errs = append(errs, fmt.Sprintf("%s: title は必須です", prefix))
		}

		if item.Type == "" {
			errs = append(errs, fmt.Sprintf("%s: type は必須です", prefix))
			continue
		}

		if !item.Type.IsValid() {
			errs = append(errs, fmt.Sprintf("%s: 不明な type %q (有効値: %s)",
				prefix, item.Type, ValidItemTypes()))
			continue
		}

		switch item.Type {
		case ItemChoice:
			errs = append(errs, validateChoice(prefix, item)...)
		case ItemScale:
			errs = append(errs, validateScale(prefix, item)...)
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}

func validateChoice(prefix string, item ItemSpec) []string {
	var errs []string

	if item.Choice == nil {
		return []string{fmt.Sprintf("%s: type=choice には choice フィールドが必須です", prefix)}
	}

	if !item.Choice.Type.IsValid() {
		errs = append(errs, fmt.Sprintf("%s: 不明な choice.type %q (有効値: %s)",
			prefix, item.Choice.Type, ValidChoiceTypes()))
	}

	if len(item.Choice.Options) == 0 {
		errs = append(errs, fmt.Sprintf("%s: choice.options は1つ以上必要です", prefix))
	}

	for j, opt := range item.Choice.Options {
		if opt == "" {
			errs = append(errs, fmt.Sprintf("%s: choice.options[%d] は空にできません", prefix, j))
		}
	}

	return errs
}

func validateScale(prefix string, item ItemSpec) []string {
	var errs []string

	if item.Scale == nil {
		return []string{fmt.Sprintf("%s: type=scale には scale フィールドが必須です", prefix)}
	}

	if item.Scale.Low != 0 && item.Scale.Low != 1 {
		errs = append(errs, fmt.Sprintf("%s: scale.low は 0 または 1 のみ (got %d)", prefix, item.Scale.Low))
	}

	if item.Scale.High < 2 || item.Scale.High > 10 {
		errs = append(errs, fmt.Sprintf("%s: scale.high は 2〜10 の範囲 (got %d)", prefix, item.Scale.High))
	}

	if item.Scale.Low >= item.Scale.High {
		errs = append(errs, fmt.Sprintf("%s: scale.low (%d) は scale.high (%d) より小さくなければなりません",
			prefix, item.Scale.Low, item.Scale.High))
	}

	return errs
}
