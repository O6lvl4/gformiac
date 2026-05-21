package engine

import (
	"fmt"
	"strings"
)

type ChangeType int

const (
	ChangeNone ChangeType = iota
	ChangeCreate
	ChangeUpdate
	ChangeDelete
)

type Change struct {
	Type  ChangeType
	Index int
	Old   *ItemSpec
	New   *ItemSpec
}

type DiffResult struct {
	InfoChanged bool
	InfoDetails []string
	Changes     []Change
}

func Diff(local *FormSpec, remote *FormSpec, state *State) *DiffResult {
	result := &DiffResult{}

	// Form info diff
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

	// Item diff
	maxLen := len(local.Items)
	if len(remote.Items) > maxLen {
		maxLen = len(remote.Items)
	}

	for i := 0; i < maxLen; i++ {
		switch {
		case i >= len(local.Items):
			old := remote.Items[i]
			result.Changes = append(result.Changes, Change{
				Type: ChangeDelete, Index: i, Old: &old,
			})
		case i >= len(remote.Items):
			new := local.Items[i]
			result.Changes = append(result.Changes, Change{
				Type: ChangeCreate, Index: i, New: &new,
			})
		default:
			old := remote.Items[i]
			new := local.Items[i]
			if !itemsEqual(old, new) {
				result.Changes = append(result.Changes, Change{
					Type: ChangeUpdate, Index: i, Old: &old, New: &new,
				})
			}
		}
	}

	return result
}

func (d *DiffResult) HasChanges() bool {
	return d.InfoChanged || len(d.Changes) > 0
}

func (d *DiffResult) String() string {
	if !d.HasChanges() {
		return "変更なし"
	}

	var lines []string

	if d.InfoChanged {
		lines = append(lines, "フォーム情報:")
		lines = append(lines, d.InfoDetails...)
	}

	creates, updates, deletes := 0, 0, 0
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

// NewFormSummary returns a human-readable summary for creating a new form.
func NewFormSummary(spec *FormSpec) string {
	var lines []string
	lines = append(lines, "新規フォーム作成:")
	lines = append(lines, fmt.Sprintf("  title: %s", spec.Title))
	if spec.Description != "" {
		lines = append(lines, fmt.Sprintf("  description: %s", spec.Description))
	}
	lines = append(lines, "")
	for i, item := range spec.Items {
		detail := item.Type
		if item.Choice != nil {
			detail = fmt.Sprintf("%s/%s [%d options]", item.Type, item.Choice.Type, len(item.Choice.Options))
		}
		req := ""
		if item.Required {
			req = " *"
		}
		lines = append(lines, fmt.Sprintf("  + [%d] %s (%s)%s", i, item.Title, detail, req))
	}
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("合計: %d項目を作成", len(spec.Items)))
	return strings.Join(lines, "\n")
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
	if (a.Choice == nil) != (b.Choice == nil) {
		return false
	}
	if a.Choice != nil {
		if a.Choice.Type != b.Choice.Type {
			return false
		}
		if len(a.Choice.Options) != len(b.Choice.Options) {
			return false
		}
		for i := range a.Choice.Options {
			if a.Choice.Options[i] != b.Choice.Options[i] {
				return false
			}
		}
	}
	if (a.Scale == nil) != (b.Scale == nil) {
		return false
	}
	if a.Scale != nil && *a.Scale != *b.Scale {
		return false
	}
	return true
}
