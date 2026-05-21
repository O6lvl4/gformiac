package engine

import (
	"encoding/json"
	"os"
	"time"
)

type State struct {
	FormID       string    `json:"form_id"`
	ResponderURL string    `json:"responder_url,omitempty"`
	ItemIDs      []string  `json:"item_ids"`
	QuestionIDs  []string  `json:"question_ids"`
	LastApplied  time.Time `json:"last_applied"`
}

func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
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

func SaveState(path string, state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
