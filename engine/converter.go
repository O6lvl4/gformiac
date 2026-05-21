package engine

import (
	"time"

	forms "google.golang.org/api/forms/v1"
)

// specToCreateRequests builds batchUpdate requests for a newly created form.
// The form title is already set via Create(); this sets description + items.
func specToCreateRequests(spec *FormSpec) []*forms.Request {
	var requests []*forms.Request

	if spec.Description != "" {
		requests = append(requests, &forms.Request{
			UpdateFormInfo: &forms.UpdateFormInfoRequest{
				Info: &forms.Info{
					Description: spec.Description,
				},
				UpdateMask: "description",
			},
		})
	}

	for i, item := range spec.Items {
		if req := itemSpecToRequest(item, i); req != nil {
			requests = append(requests, req)
		}
	}

	return requests
}

// specToUpdateRequests builds batchUpdate requests to reconcile a remote form.
// Strategy: update info, delete all items (reverse order), recreate all items.
func specToUpdateRequests(spec *FormSpec, remote *forms.Form) []*forms.Request {
	var requests []*forms.Request

	// Update form info
	requests = append(requests, &forms.Request{
		UpdateFormInfo: &forms.UpdateFormInfoRequest{
			Info: &forms.Info{
				Title:       spec.Title,
				Description: spec.Description,
			},
			UpdateMask: "title,description",
		},
	})

	// Delete existing items in reverse order
	for i := len(remote.Items) - 1; i >= 0; i-- {
		requests = append(requests, &forms.Request{
			DeleteItem: &forms.DeleteItemRequest{
				Location: &forms.Location{
					Index:           int64(i),
					ForceSendFields: []string{"Index"},
				},
			},
		})
	}

	// Recreate all items from spec
	for i, item := range spec.Items {
		if req := itemSpecToRequest(item, i); req != nil {
			requests = append(requests, req)
		}
	}

	return requests
}

func itemSpecToRequest(item ItemSpec, index int) *forms.Request {
	apiItem := &forms.Item{
		Title:       item.Title,
		Description: item.Description,
	}

	switch item.Type {
	case "short_answer":
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required:     item.Required,
				TextQuestion: &forms.TextQuestion{Paragraph: false},
			},
		}

	case "paragraph":
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required:     item.Required,
				TextQuestion: &forms.TextQuestion{Paragraph: true},
			},
		}

	case "choice":
		if item.Choice == nil {
			return nil
		}
		choiceType := "RADIO"
		switch item.Choice.Type {
		case "checkbox":
			choiceType = "CHECKBOX"
		case "dropdown":
			choiceType = "DROP_DOWN"
		}
		var options []*forms.Option
		for _, v := range item.Choice.Options {
			options = append(options, &forms.Option{Value: v})
		}
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required: item.Required,
				ChoiceQuestion: &forms.ChoiceQuestion{
					Type:    choiceType,
					Options: options,
				},
			},
		}

	case "scale":
		if item.Scale == nil {
			return nil
		}
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required: item.Required,
				ScaleQuestion: &forms.ScaleQuestion{
					Low:      item.Scale.Low,
					High:     item.Scale.High,
					LowLabel: item.Scale.LowLabel,
					HighLabel: item.Scale.HighLabel,
				},
			},
		}

	case "date":
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required:     item.Required,
				DateQuestion: &forms.DateQuestion{},
			},
		}

	case "time":
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required:     item.Required,
				TimeQuestion: &forms.TimeQuestion{},
			},
		}

	case "page_break":
		apiItem.PageBreakItem = &forms.PageBreakItem{}

	default:
		return nil
	}

	return &forms.Request{
		CreateItem: &forms.CreateItemRequest{
			Item: apiItem,
			Location: &forms.Location{
				Index:           int64(index),
				ForceSendFields: []string{"Index"},
			},
		},
	}
}

// formToSpec converts a Google Forms API response to a FormSpec.
func formToSpec(form *forms.Form) *FormSpec {
	spec := &FormSpec{}
	if form.Info != nil {
		spec.Title = form.Info.Title
		spec.Description = form.Info.Description
	}

	for _, item := range form.Items {
		spec.Items = append(spec.Items, apiItemToSpec(item))
	}

	return spec
}

func apiItemToSpec(item *forms.Item) ItemSpec {
	spec := ItemSpec{
		Title:       item.Title,
		Description: item.Description,
	}

	if item.PageBreakItem != nil {
		spec.Type = "page_break"
		return spec
	}

	if item.QuestionItem == nil || item.QuestionItem.Question == nil {
		spec.Type = "unknown"
		return spec
	}

	q := item.QuestionItem.Question
	spec.Required = q.Required

	switch {
	case q.TextQuestion != nil:
		if q.TextQuestion.Paragraph {
			spec.Type = "paragraph"
		} else {
			spec.Type = "short_answer"
		}

	case q.ChoiceQuestion != nil:
		spec.Type = "choice"
		choiceType := "radio"
		switch q.ChoiceQuestion.Type {
		case "CHECKBOX":
			choiceType = "checkbox"
		case "DROP_DOWN":
			choiceType = "dropdown"
		}
		var options []string
		for _, opt := range q.ChoiceQuestion.Options {
			options = append(options, opt.Value)
		}
		spec.Choice = &ChoiceSpec{Type: choiceType, Options: options}

	case q.ScaleQuestion != nil:
		spec.Type = "scale"
		spec.Scale = &ScaleSpec{
			Low:       q.ScaleQuestion.Low,
			High:      q.ScaleQuestion.High,
			LowLabel:  q.ScaleQuestion.LowLabel,
			HighLabel: q.ScaleQuestion.HighLabel,
		}

	case q.DateQuestion != nil:
		spec.Type = "date"

	case q.TimeQuestion != nil:
		spec.Type = "time"

	default:
		spec.Type = "unknown"
	}

	return spec
}

// buildState creates a State from a remote form response.
func buildState(form *forms.Form) *State {
	state := &State{
		FormID:      form.FormId,
		ResponderURL: form.ResponderUri,
		LastApplied: time.Now(),
	}

	for _, item := range form.Items {
		state.ItemIDs = append(state.ItemIDs, item.ItemId)
		qid := ""
		if item.QuestionItem != nil && item.QuestionItem.Question != nil {
			qid = item.QuestionItem.Question.QuestionId
		}
		state.QuestionIDs = append(state.QuestionIDs, qid)
	}

	return state
}
