package engine

import (
	"fmt"
	"slices"
	"strings"
)

// ChangeType classifies a diff entry.
type ChangeType int

const (
	ChangeCreate ChangeType = iota
	ChangeUpdate
	ChangeDelete
)

// String returns a label for the change type.
func (t ChangeType) String() string {
	switch t {
	case ChangeCreate:
		return "create"
	case ChangeUpdate:
		return "update"
	case ChangeDelete:
		return "delete"
	default:
		return "unknown"
	}
}

// Change represents a single item-level difference.
type Change struct {
	Type  ChangeType
	Index int
	Old   *ItemSpec
	New   *ItemSpec
}

// DiffResult holds the full set of differences between local and remote forms.
type DiffResult struct {
	InfoChanged bool
	InfoDetails []string
	Changes     []Change
}

// Diff computes the differences between a local spec and a remote spec.
func Diff(local, remote *FormSpec, _ *State) *DiffResult {
	result := &DiffResult{}
	diffInfo(result, local, remote)
	diffItems(result, local.Items, remote.Items)
	return result
}

func diffInfo(result *DiffResult, local, remote *FormSpec) {
	if local.Title != remote.Title {
		result.InfoChanged = true
		result.InfoDetails = append(result.InfoDetails,
			fmt.Sprintf("  ~ title: %q -> %q", remote.Title, local.Title))
	}
	if local.Description != remote.Description {
		result.InfoChanged = true
		result.InfoDetails = append(result.InfoDetails,
			fmt.Sprintf("  ~ description: %q -> %q", remote.Description, local.Description))
	}
}

func diffItems(result *DiffResult, local, remote []ItemSpec) {
	for i := range max(len(local), len(remote)) {
		switch {
		case i >= len(local):
			old := remote[i]
			result.Changes = append(result.Changes, Change{
				Type: ChangeDelete, Index: i, Old: &old,
			})
		case i >= len(remote):
			item := local[i]
			result.Changes = append(result.Changes, Change{
				Type: ChangeCreate, Index: i, New: &item,
			})
		default:
			old := remote[i]
			item := local[i]
			if !itemsEqual(old, item) {
				result.Changes = append(result.Changes, Change{
					Type: ChangeUpdate, Index: i, Old: &old, New: &item,
				})
			}
		}
	}
}

// HasChanges reports whether there are any differences.
func (d *DiffResult) HasChanges() bool {
	return d.InfoChanged || len(d.Changes) > 0
}

// String formats the diff as a human-readable summary.
func (d *DiffResult) String() string {
	if !d.HasChanges() {
		return "変更なし"
	}

	var lines []string

	if d.InfoChanged {
		lines = append(lines, "フォーム情報:")
		lines = append(lines, d.InfoDetails...)
	}

	var creates, updates, deletes int
	for _, c := range d.Changes {
		switch c.Type {
		case ChangeCreate:
			creates++
			lines = append(lines, fmt.Sprintf("  + [%d] %s (%s)", c.Index, c.New.Title, c.New.Type))
		case ChangeDelete:
			deletes++
			lines = append(lines, fmt.Sprintf("  - [%d] %s (%s)", c.Index, c.Old.Title, c.Old.Type))
		case ChangeUpdate:
			updates++
			lines = append(lines, formatUpdate(c)...)
		}
	}

	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("合計: +%d ~%d -%d", creates, updates, deletes))

	return strings.Join(lines, "\n")
}

// NewFormSummary returns a human-readable plan for creating a new form.
func NewFormSummary(spec *FormSpec) string {
	var lines []string
	lines = append(lines, "新規フォーム作成:")
	lines = append(lines, fmt.Sprintf("  title: %s", spec.Title))
	if spec.Description != "" {
		lines = append(lines, fmt.Sprintf("  description: %s", spec.Description))
	}
	lines = append(lines, "")
	for i, item := range spec.Items {
		lines = append(lines, formatNewItem(i, item))
	}
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("合計: %d項目を作成", len(spec.Items)))
	return strings.Join(lines, "\n")
}

func formatNewItem(i int, item ItemSpec) string {
	detail := string(item.Type)
	if item.Choice != nil {
		detail = fmt.Sprintf("%s/%s [%d options]", item.Type, item.Choice.Type, len(item.Choice.Options))
	}
	req := ""
	if item.Required {
		req = " *"
	}
	return fmt.Sprintf("  + [%d] %s (%s)%s", i, item.Title, detail, req)
}

func formatUpdate(c Change) []string {
	var lines []string
	lines = append(lines, fmt.Sprintf("  ~ [%d] %s", c.Index, c.New.Title))
	if c.Old.Title != c.New.Title {
		lines = append(lines, fmt.Sprintf("      title: %q -> %q", c.Old.Title, c.New.Title))
	}
	if c.Old.Type != c.New.Type {
		lines = append(lines, fmt.Sprintf("      type: %s -> %s", c.Old.Type, c.New.Type))
	}
	if c.Old.Required != c.New.Required {
		lines = append(lines, fmt.Sprintf("      required: %v -> %v", c.Old.Required, c.New.Required))
	}
	if c.Old.Description != c.New.Description {
		lines = append(lines, fmt.Sprintf("      description: %q -> %q", c.Old.Description, c.New.Description))
	}
	return lines
}

func itemsEqual(a, b ItemSpec) bool {
	if a.Title != b.Title || a.Type != b.Type || a.Required != b.Required || a.Description != b.Description {
		return false
	}
	return choiceEqual(a.Choice, b.Choice) && scaleEqual(a.Scale, b.Scale)
}

func choiceEqual(a, b *ChoiceSpec) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return a.Type == b.Type && slices.Equal(a.Options, b.Options)
}

func scaleEqual(a, b *ScaleSpec) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return *a == *b
}
