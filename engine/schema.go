// Package engine implements the core logic for gformiac: schema types, state
// management, Google Forms API interactions, diff computation, and validation.
package engine

import (
	"fmt"
	"os"

	"github.com/O6lvl4/gformiac/locale"
	"gopkg.in/yaml.v3"
)

// ItemType represents the type of a form item.
type ItemType string

// Recognized item types that map directly to Google Forms question types.
const (
	// ItemShortAnswer is a single-line text question.
	ItemShortAnswer ItemType = "short_answer"
	// ItemParagraph is a multi-line text question.
	ItemParagraph ItemType = "paragraph"
	// ItemChoice is a question with selectable options (radio, checkbox, or dropdown).
	ItemChoice ItemType = "choice"
	// ItemScale is a linear-scale question with a numeric range.
	ItemScale ItemType = "scale"
	// ItemDate is a date-picker question.
	ItemDate ItemType = "date"
	// ItemTime is a time-picker question.
	ItemTime ItemType = "time"
	// ItemPageBreak inserts a visual page break between sections.
	ItemPageBreak ItemType = "page_break"
)

// IsValid reports whether t is a recognized item type.
func (t ItemType) IsValid() bool {
	switch t {
	case ItemShortAnswer, ItemParagraph, ItemChoice, ItemScale,
		ItemDate, ItemTime, ItemPageBreak:
		return true
	}
	return false
}

// ValidItemTypes returns the recognized item types as a display string.
func ValidItemTypes() string {
	return "short_answer, paragraph, choice, scale, date, time, page_break"
}

// ChoiceType represents the sub-type of a choice question.
type ChoiceType string

// Recognized choice sub-types for ItemChoice questions.
const (
	// ChoiceRadio renders options as radio buttons (single-select).
	ChoiceRadio ChoiceType = "radio"
	// ChoiceCheckbox renders options as checkboxes (multi-select).
	ChoiceCheckbox ChoiceType = "checkbox"
	// ChoiceDropdown renders options as a dropdown list (single-select).
	ChoiceDropdown ChoiceType = "dropdown"
)

// IsValid reports whether t is a recognized choice type.
func (t ChoiceType) IsValid() bool {
	switch t {
	case ChoiceRadio, ChoiceCheckbox, ChoiceDropdown:
		return true
	}
	return false
}

// ValidChoiceTypes returns the recognized choice types as a display string.
func ValidChoiceTypes() string {
	return "radio, checkbox, dropdown"
}

// FormSpec is the top-level YAML definition of a Google Form.
type FormSpec struct {
	Title       string     `yaml:"title"`
	Description string     `yaml:"description,omitempty"`
	Items       []ItemSpec `yaml:"items"`
}

// ItemSpec defines a single item (question, page break, etc.) in a form.
type ItemSpec struct {
	Title       string      `yaml:"title"`
	Description string      `yaml:"description,omitempty"`
	Type        ItemType    `yaml:"type"`
	Required    bool        `yaml:"required,omitempty"`
	Choice      *ChoiceSpec `yaml:"choice,omitempty"`
	Scale       *ScaleSpec  `yaml:"scale,omitempty"`
}

// ChoiceSpec defines the options for a choice question.
type ChoiceSpec struct {
	Type    ChoiceType `yaml:"type"`
	Options []string   `yaml:"options"`
}

// ScaleSpec defines the range and labels for a linear scale question.
type ScaleSpec struct {
	Low       int64  `yaml:"low"`
	High      int64  `yaml:"high"`
	LowLabel  string `yaml:"low_label,omitempty"`
	HighLabel string `yaml:"high_label,omitempty"`
}

// LoadSpec reads and parses a YAML form definition file.
func LoadSpec(path string) (*FormSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", locale.M.ErrReadFile, path, err)
	}
	var spec FormSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("%s: %w", locale.M.ErrParseYAML, err)
	}
	if spec.Title == "" {
		return nil, fmt.Errorf("%s", locale.M.ErrTitleReq)
	}
	return &spec, nil
}

// SaveSpec writes a FormSpec as YAML to the given path.
func SaveSpec(path string, spec *FormSpec) error {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
