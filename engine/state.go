package engine

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"time"
)

// State tracks the mapping between a local form definition and its remote Google Form.
type State struct {
	FormID       string    `json:"form_id"`
	ResponderURL string    `json:"responder_url,omitempty"`
	ItemIDs      []string  `json:"item_ids"`
	QuestionIDs  []string  `json:"question_ids"`
	LastApplied  time.Time `json:"last_applied"`
}

// LoadState reads the state file. Returns (nil, nil) if the file does not exist,
// signaling that no form has been created yet.
func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// SaveState writes the state to disk as formatted JSON.
func SaveState(path string, state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
