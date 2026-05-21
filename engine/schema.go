package engine

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ItemType represents the type of a form item.
type ItemType string

const (
	ItemShortAnswer ItemType = "short_answer"
	ItemParagraph   ItemType = "paragraph"
	ItemChoice      ItemType = "choice"
	ItemScale       ItemType = "scale"
	ItemDate        ItemType = "date"
	ItemTime        ItemType = "time"
	ItemPageBreak   ItemType = "page_break"
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

const (
	ChoiceRadio    ChoiceType = "radio"
	ChoiceCheckbox ChoiceType = "checkbox"
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
		return nil, fmt.Errorf("ファイル読み込み失敗 %s: %w", path, err)
	}
	var spec FormSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("YAML解析失敗: %w", err)
	}
	if spec.Title == "" {
		return nil, fmt.Errorf("titleは必須です")
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
