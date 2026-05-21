package engine

import (
	"fmt"
	"strings"
)

// Google Forms API の制約に基づくバリデーション
// ref: https://developers.google.com/forms/api/reference/rest/v1/forms

var validItemTypes = map[string]bool{
	"short_answer": true,
	"paragraph":    true,
	"choice":       true,
	"scale":        true,
	"date":         true,
	"time":         true,
	"page_break":   true,
}

var validChoiceTypes = map[string]bool{
	"radio":    true,
	"checkbox": true,
	"dropdown": true,
}

type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("バリデーションエラー (%d件):\n  %s", len(e.Errors), strings.Join(e.Errors, "\n  "))
}

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

		if item.Title == "" && item.Type != "page_break" {
			errs = append(errs, fmt.Sprintf("%s: title は必須です", prefix))
		}

		if item.Type == "" {
			errs = append(errs, fmt.Sprintf("%s: type は必須です", prefix))
			continue
		}

		if !validItemTypes[item.Type] {
			errs = append(errs, fmt.Sprintf("%s: 不明な type %q (有効値: %s)",
				prefix, item.Type, joinKeys(validItemTypes)))
			continue
		}

		switch item.Type {
		case "choice":
			errs = append(errs, validateChoice(prefix, item)...)
		case "scale":
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

	if !validChoiceTypes[item.Choice.Type] {
		errs = append(errs, fmt.Sprintf("%s: 不明な choice.type %q (有効値: %s)",
			prefix, item.Choice.Type, joinKeys(validChoiceTypes)))
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

	// Google Forms API: low は 0 または 1、high は 2〜10
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

func joinKeys(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}
