package engine

import (
	"fmt"
	"strings"

	"github.com/O6lvl4/gformiac/locale"
)

// ValidationError collects multiple validation failures.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(locale.M.ValErrors, len(e.Errors), strings.Join(e.Errors, "\n  "))
}

// Validate checks a FormSpec against Google Forms API constraints.
func Validate(spec *FormSpec) error {
	var errs []string

	if spec.Title == "" {
		errs = append(errs, locale.M.ValTitle)
	}
	if len(spec.Items) == 0 {
		errs = append(errs, locale.M.ValItems)
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
		errs = append(errs, fmt.Sprintf(locale.M.ValItemTitle, prefix))
	}

	if item.Type == "" {
		return append(errs, fmt.Sprintf(locale.M.ValItemType, prefix))
	}

	if !item.Type.IsValid() {
		return append(errs, fmt.Sprintf(locale.M.ValTypeUnknown,
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
		return []string{fmt.Sprintf(locale.M.ValChoiceReq, prefix)}
	}

	var errs []string
	if !item.Choice.Type.IsValid() {
		errs = append(errs, fmt.Sprintf(locale.M.ValChoiceType,
			prefix, item.Choice.Type, ValidChoiceTypes()))
	}
	if len(item.Choice.Options) == 0 {
		errs = append(errs, fmt.Sprintf(locale.M.ValChoiceOpts, prefix))
	}
	errs = append(errs, validateOptionValues(prefix, item.Choice.Options)...)
	return errs
}

func validateOptionValues(prefix string, options []string) []string {
	var errs []string
	for j, opt := range options {
		if opt == "" {
			errs = append(errs, fmt.Sprintf(locale.M.ValChoiceEmpty, prefix, j))
		}
	}
	return errs
}

func validateScale(prefix string, item ItemSpec) []string {
	if item.Scale == nil {
		return []string{fmt.Sprintf(locale.M.ValScaleReq, prefix)}
	}
	return validateScaleRange(prefix, item.Scale)
}

func validateScaleRange(prefix string, s *ScaleSpec) []string {
	var errs []string
	if s.Low != 0 && s.Low != 1 {
		errs = append(errs, fmt.Sprintf(locale.M.ValScaleLow, prefix, s.Low))
	}
	if s.High < 2 || s.High > 10 {
		errs = append(errs, fmt.Sprintf(locale.M.ValScaleHigh, prefix, s.High))
	}
	if s.Low >= s.High {
		errs = append(errs, fmt.Sprintf(locale.M.ValScaleRange, prefix, s.Low, s.High))
	}
	return errs
}
