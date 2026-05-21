package engine

import (
	"time"

	forms "google.golang.org/api/forms/v1"
)

// choiceTypeToAPI converts a ChoiceType to the Google Forms API string.
func choiceTypeToAPI(t ChoiceType) string {
	switch t {
	case ChoiceCheckbox:
		return "CHECKBOX"
	case ChoiceDropdown:
		return "DROP_DOWN"
	default:
		return "RADIO"
	}
}

// choiceTypeFromAPI converts a Google Forms API choice type string to ChoiceType.
func choiceTypeFromAPI(s string) ChoiceType {
	switch s {
	case "CHECKBOX":
		return ChoiceCheckbox
	case "DROP_DOWN":
		return ChoiceDropdown
	default:
		return ChoiceRadio
	}
}

// specToCreateRequests builds batchUpdate requests for a newly created form.
// Title is already set via Create(); this sets description and items.
func specToCreateRequests(spec *FormSpec) []*forms.Request {
	var requests []*forms.Request

	if spec.Description != "" {
		requests = append(requests, &forms.Request{
			UpdateFormInfo: &forms.UpdateFormInfoRequest{
				Info:       &forms.Info{Description: spec.Description},
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
// Strategy: update info, delete all items (reverse order), recreate all.
func specToUpdateRequests(spec *FormSpec, remote *forms.Form) []*forms.Request {
	var requests []*forms.Request

	requests = append(requests, &forms.Request{
		UpdateFormInfo: &forms.UpdateFormInfoRequest{
			Info: &forms.Info{
				Title:       spec.Title,
				Description: spec.Description,
			},
			UpdateMask: "title,description",
		},
	})

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
	case ItemShortAnswer:
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required:     item.Required,
				TextQuestion: &forms.TextQuestion{Paragraph: false},
			},
		}

	case ItemParagraph:
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required:     item.Required,
				TextQuestion: &forms.TextQuestion{Paragraph: true},
			},
		}

	case ItemChoice:
		if item.Choice == nil {
			return nil
		}
		options := make([]*forms.Option, len(item.Choice.Options))
		for i, v := range item.Choice.Options {
			options[i] = &forms.Option{Value: v}
		}
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required: item.Required,
				ChoiceQuestion: &forms.ChoiceQuestion{
					Type:    choiceTypeToAPI(item.Choice.Type),
					Options: options,
				},
			},
		}

	case ItemScale:
		if item.Scale == nil {
			return nil
		}
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required: item.Required,
				ScaleQuestion: &forms.ScaleQuestion{
					Low:       item.Scale.Low,
					High:      item.Scale.High,
					LowLabel:  item.Scale.LowLabel,
					HighLabel: item.Scale.HighLabel,
				},
			},
		}

	case ItemDate:
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required:     item.Required,
				DateQuestion: &forms.DateQuestion{},
			},
		}

	case ItemTime:
		apiItem.QuestionItem = &forms.QuestionItem{
			Question: &forms.Question{
				Required:     item.Required,
				TimeQuestion: &forms.TimeQuestion{},
			},
		}

	case ItemPageBreak:
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
		spec.Type = ItemPageBreak
		return spec
	}

	if item.QuestionItem == nil || item.QuestionItem.Question == nil {
		return spec
	}

	q := item.QuestionItem.Question
	spec.Required = q.Required

	switch {
	case q.TextQuestion != nil:
		if q.TextQuestion.Paragraph {
			spec.Type = ItemParagraph
		} else {
			spec.Type = ItemShortAnswer
		}

	case q.ChoiceQuestion != nil:
		spec.Type = ItemChoice
		options := make([]string, len(q.ChoiceQuestion.Options))
		for i, opt := range q.ChoiceQuestion.Options {
			options[i] = opt.Value
		}
		spec.Choice = &ChoiceSpec{
			Type:    choiceTypeFromAPI(q.ChoiceQuestion.Type),
			Options: options,
		}

	case q.ScaleQuestion != nil:
		spec.Type = ItemScale
		spec.Scale = &ScaleSpec{
			Low:       q.ScaleQuestion.Low,
			High:      q.ScaleQuestion.High,
			LowLabel:  q.ScaleQuestion.LowLabel,
			HighLabel: q.ScaleQuestion.HighLabel,
		}

	case q.DateQuestion != nil:
		spec.Type = ItemDate

	case q.TimeQuestion != nil:
		spec.Type = ItemTime
	}

	return spec
}

// buildState creates a State from a remote form response.
func buildState(form *forms.Form) *State {
	state := &State{
		FormID:       form.FormId,
		ResponderURL: form.ResponderUri,
		ItemIDs:      make([]string, len(form.Items)),
		QuestionIDs:  make([]string, len(form.Items)),
		LastApplied:  time.Now(),
	}

	for i, item := range form.Items {
		state.ItemIDs[i] = item.ItemId
		if item.QuestionItem != nil && item.QuestionItem.Question != nil {
			state.QuestionIDs[i] = item.QuestionItem.Question.QuestionId
		}
	}

	return state
}
