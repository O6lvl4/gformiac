package engine

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type FormSpec struct {
	Title       string     `yaml:"title"`
	Description string     `yaml:"description,omitempty"`
	Items       []ItemSpec `yaml:"items"`
}

type ItemSpec struct {
	Title       string      `yaml:"title"`
	Description string      `yaml:"description,omitempty"`
	Type        string      `yaml:"type"`
	Required    bool        `yaml:"required,omitempty"`
	Choice      *ChoiceSpec `yaml:"choice,omitempty"`
	Scale       *ScaleSpec  `yaml:"scale,omitempty"`
}

type ChoiceSpec struct {
	Type    string   `yaml:"type"` // radio, checkbox, dropdown
	Options []string `yaml:"options"`
}

type ScaleSpec struct {
	Low       int64  `yaml:"low"`
	High      int64  `yaml:"high"`
	LowLabel  string `yaml:"low_label,omitempty"`
	HighLabel string `yaml:"high_label,omitempty"`
}

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

func SaveSpec(path string, spec *FormSpec) error {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
