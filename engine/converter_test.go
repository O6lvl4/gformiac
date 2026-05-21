package engine

import (
	"testing"

	forms "google.golang.org/api/forms/v1"
)

func TestItemSpecToRequest_ShortAnswer(t *testing.T) {
	item := ItemSpec{Title: "Name", Type: ItemShortAnswer, Required: true}
	req := itemSpecToRequest(item, 0)

	if req == nil || req.CreateItem == nil {
		t.Fatal("expected CreateItem request")
	}
	ci := req.CreateItem
	if ci.Item.Title != "Name" {
		t.Errorf("title = %q", ci.Item.Title)
	}
	q := ci.Item.QuestionItem.Question
	if !q.Required {
		t.Error("expected required=true")
	}
	if q.TextQuestion == nil || q.TextQuestion.Paragraph {
		t.Error("expected short answer (paragraph=false)")
	}
	if ci.Location.Index != 0 {
		t.Errorf("location index = %d", ci.Location.Index)
	}
}

func TestItemSpecToRequest_Paragraph(t *testing.T) {
	item := ItemSpec{Title: "Comment", Type: ItemParagraph}
	req := itemSpecToRequest(item, 2)

	q := req.CreateItem.Item.QuestionItem.Question
	if q.TextQuestion == nil || !q.TextQuestion.Paragraph {
		t.Error("expected paragraph=true")
	}
	if req.CreateItem.Location.Index != 2 {
		t.Errorf("location index = %d", req.CreateItem.Location.Index)
	}
}

func TestItemSpecToRequest_ChoiceRadio(t *testing.T) {
	item := ItemSpec{
		Title: "Pick", Type: ItemChoice,
		Choice: &ChoiceSpec{Type: ChoiceRadio, Options: []string{"A", "B"}},
	}
	req := itemSpecToRequest(item, 0)

	q := req.CreateItem.Item.QuestionItem.Question
	if q.ChoiceQuestion == nil {
		t.Fatal("expected ChoiceQuestion")
	}
	if q.ChoiceQuestion.Type != "RADIO" {
		t.Errorf("choice type = %q", q.ChoiceQuestion.Type)
	}
	if len(q.ChoiceQuestion.Options) != 2 {
		t.Errorf("options count = %d", len(q.ChoiceQuestion.Options))
	}
}

func TestItemSpecToRequest_ChoiceCheckbox(t *testing.T) {
	item := ItemSpec{
		Title: "Multi", Type: ItemChoice,
		Choice: &ChoiceSpec{Type: ChoiceCheckbox, Options: []string{"X"}},
	}
	req := itemSpecToRequest(item, 0)
	if req.CreateItem.Item.QuestionItem.Question.ChoiceQuestion.Type != "CHECKBOX" {
		t.Error("expected CHECKBOX")
	}
}

func TestItemSpecToRequest_ChoiceDropdown(t *testing.T) {
	item := ItemSpec{
		Title: "Select", Type: ItemChoice,
		Choice: &ChoiceSpec{Type: ChoiceDropdown, Options: []string{"X"}},
	}
	req := itemSpecToRequest(item, 0)
	if req.CreateItem.Item.QuestionItem.Question.ChoiceQuestion.Type != "DROP_DOWN" {
		t.Error("expected DROP_DOWN")
	}
}

func TestItemSpecToRequest_Scale(t *testing.T) {
	item := ItemSpec{
		Title: "Rate", Type: ItemScale,
		Scale: &ScaleSpec{Low: 1, High: 10, LowLabel: "Bad", HighLabel: "Good"},
	}
	req := itemSpecToRequest(item, 0)

	q := req.CreateItem.Item.QuestionItem.Question
	if q.ScaleQuestion == nil {
		t.Fatal("expected ScaleQuestion")
	}
	if q.ScaleQuestion.Low != 1 || q.ScaleQuestion.High != 10 {
		t.Errorf("scale = %d-%d", q.ScaleQuestion.Low, q.ScaleQuestion.High)
	}
	if q.ScaleQuestion.LowLabel != "Bad" || q.ScaleQuestion.HighLabel != "Good" {
		t.Error("scale labels mismatch")
	}
}

func TestItemSpecToRequest_Date(t *testing.T) {
	req := itemSpecToRequest(ItemSpec{Title: "When", Type: ItemDate}, 0)
	if req.CreateItem.Item.QuestionItem.Question.DateQuestion == nil {
		t.Error("expected DateQuestion")
	}
}

func TestItemSpecToRequest_Time(t *testing.T) {
	req := itemSpecToRequest(ItemSpec{Title: "When", Type: ItemTime}, 0)
	if req.CreateItem.Item.QuestionItem.Question.TimeQuestion == nil {
		t.Error("expected TimeQuestion")
	}
}

func TestItemSpecToRequest_PageBreak(t *testing.T) {
	req := itemSpecToRequest(ItemSpec{Title: "Section", Type: ItemPageBreak}, 0)
	if req.CreateItem.Item.PageBreakItem == nil {
		t.Error("expected PageBreakItem")
	}
}

func TestItemSpecToRequest_Unknown(t *testing.T) {
	req := itemSpecToRequest(ItemSpec{Title: "X", Type: "unknown_type"}, 0)
	if req != nil {
		t.Error("expected nil for unknown type")
	}
}

func TestItemSpecToRequest_ChoiceNilSpec(t *testing.T) {
	req := itemSpecToRequest(ItemSpec{Title: "X", Type: ItemChoice}, 0)
	if req != nil {
		t.Error("expected nil for choice without spec")
	}
}

func TestFormToSpec_RoundTrip(t *testing.T) {
	form := &forms.Form{
		FormId: "abc123",
		Info:   &forms.Info{Title: "Test", Description: "Desc"},
		Items: []*forms.Item{
			{
				ItemId: "i1", Title: "Name",
				QuestionItem: &forms.QuestionItem{
					Question: &forms.Question{
						QuestionId: "q1", Required: true,
						TextQuestion: &forms.TextQuestion{Paragraph: false},
					},
				},
			},
			{
				ItemId: "i2", Title: "Pick",
				QuestionItem: &forms.QuestionItem{
					Question: &forms.Question{
						QuestionId: "q2",
						ChoiceQuestion: &forms.ChoiceQuestion{
							Type:    "CHECKBOX",
							Options: []*forms.Option{{Value: "A"}, {Value: "B"}},
						},
					},
				},
			},
			{
				ItemId: "i3", Title: "Rate",
				QuestionItem: &forms.QuestionItem{
					Question: &forms.Question{
						QuestionId: "q3",
						ScaleQuestion: &forms.ScaleQuestion{
							Low: 1, High: 5, LowLabel: "L", HighLabel: "H",
						},
					},
				},
			},
			{ItemId: "i4", Title: "Section", PageBreakItem: &forms.PageBreakItem{}},
			{
				ItemId: "i5", Title: "Date",
				QuestionItem: &forms.QuestionItem{
					Question: &forms.Question{QuestionId: "q5", DateQuestion: &forms.DateQuestion{}},
				},
			},
			{
				ItemId: "i6", Title: "Time",
				QuestionItem: &forms.QuestionItem{
					Question: &forms.Question{QuestionId: "q6", TimeQuestion: &forms.TimeQuestion{}},
				},
			},
		},
	}

	spec := formToSpec(form)

	if spec.Title != "Test" || spec.Description != "Desc" {
		t.Errorf("info mismatch: %q %q", spec.Title, spec.Description)
	}
	if len(spec.Items) != 6 {
		t.Fatalf("items = %d, want 6", len(spec.Items))
	}

	checks := []struct {
		index    int
		title    string
		itemType ItemType
	}{
		{0, "Name", ItemShortAnswer},
		{1, "Pick", ItemChoice},
		{2, "Rate", ItemScale},
		{3, "Section", ItemPageBreak},
		{4, "Date", ItemDate},
		{5, "Time", ItemTime},
	}
	for _, c := range checks {
		got := spec.Items[c.index]
		if got.Title != c.title || got.Type != c.itemType {
			t.Errorf("[%d] got title=%q type=%q, want %q %q",
				c.index, got.Title, got.Type, c.title, c.itemType)
		}
	}

	if spec.Items[1].Choice == nil || spec.Items[1].Choice.Type != ChoiceCheckbox {
		t.Errorf("choice type mismatch: %+v", spec.Items[1].Choice)
	}
	if len(spec.Items[1].Choice.Options) != 2 {
		t.Errorf("choice options = %d", len(spec.Items[1].Choice.Options))
	}
	if spec.Items[2].Scale == nil || spec.Items[2].Scale.High != 5 {
		t.Errorf("scale mismatch: %+v", spec.Items[2].Scale)
	}
}

func TestSpecToCreateRequests(t *testing.T) {
	spec := &FormSpec{
		Title:       "T",
		Description: "D",
		Items: []ItemSpec{
			{Title: "Q1", Type: ItemShortAnswer},
			{Title: "Q2", Type: ItemParagraph},
		},
	}

	requests := specToCreateRequests(spec)

	if requests[0].UpdateFormInfo == nil {
		t.Fatal("first request should be UpdateFormInfo")
	}
	if requests[0].UpdateFormInfo.Info.Description != "D" {
		t.Errorf("description = %q", requests[0].UpdateFormInfo.Info.Description)
	}
	if len(requests) != 3 {
		t.Fatalf("requests = %d, want 3", len(requests))
	}
	for i := 1; i < len(requests); i++ {
		if requests[i].CreateItem == nil {
			t.Errorf("request[%d] should be CreateItem", i)
		}
	}
}

func TestSpecToCreateRequests_NoDescription(t *testing.T) {
	spec := &FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q1", Type: ItemShortAnswer}},
	}

	requests := specToCreateRequests(spec)
	if len(requests) != 1 {
		t.Fatalf("requests = %d, want 1", len(requests))
	}
	if requests[0].CreateItem == nil {
		t.Error("expected CreateItem")
	}
}

func TestSpecToUpdateRequests(t *testing.T) {
	spec := &FormSpec{
		Title:       "New Title",
		Description: "New Desc",
		Items:       []ItemSpec{{Title: "Q1", Type: ItemShortAnswer}},
	}
	remoteForm := &forms.Form{
		Items: []*forms.Item{
			{ItemId: "old1", Title: "Old Q1"},
			{ItemId: "old2", Title: "Old Q2"},
		},
	}

	requests := specToUpdateRequests(spec, remoteForm)

	// 1 UpdateFormInfo + 2 DeleteItem + 1 CreateItem = 4
	if len(requests) != 4 {
		t.Fatalf("requests = %d, want 4", len(requests))
	}

	if requests[0].UpdateFormInfo == nil {
		t.Error("[0] expected UpdateFormInfo")
	}
	if requests[1].DeleteItem == nil || requests[1].DeleteItem.Location.Index != 1 {
		t.Errorf("[1] expected delete index 1, got %+v", requests[1])
	}
	if requests[2].DeleteItem == nil || requests[2].DeleteItem.Location.Index != 0 {
		t.Errorf("[2] expected delete index 0, got %+v", requests[2])
	}
	if requests[3].CreateItem == nil {
		t.Error("[3] expected CreateItem")
	}
}

func TestBuildState(t *testing.T) {
	form := &forms.Form{
		FormId:       "form123",
		ResponderUri: "https://docs.google.com/forms/d/e/form123/viewform",
		Items: []*forms.Item{
			{
				ItemId: "item1",
				QuestionItem: &forms.QuestionItem{
					Question: &forms.Question{QuestionId: "q1"},
				},
			},
			{ItemId: "item2", PageBreakItem: &forms.PageBreakItem{}},
		},
	}

	state := buildState(form)

	if state.FormID != "form123" {
		t.Errorf("FormID = %q", state.FormID)
	}
	if state.ResponderURL != form.ResponderUri {
		t.Errorf("ResponderURL = %q", state.ResponderURL)
	}
	if len(state.ItemIDs) != 2 || state.ItemIDs[0] != "item1" || state.ItemIDs[1] != "item2" {
		t.Errorf("ItemIDs = %v", state.ItemIDs)
	}
	if state.QuestionIDs[0] != "q1" || state.QuestionIDs[1] != "" {
		t.Errorf("QuestionIDs = %v", state.QuestionIDs)
	}
}

func TestChoiceTypeMapping(t *testing.T) {
	cases := []struct {
		ct  ChoiceType
		api string
	}{
		{ChoiceRadio, "RADIO"},
		{ChoiceCheckbox, "CHECKBOX"},
		{ChoiceDropdown, "DROP_DOWN"},
	}
	for _, c := range cases {
		if got := choiceTypeToAPI(c.ct); got != c.api {
			t.Errorf("choiceTypeToAPI(%q) = %q, want %q", c.ct, got, c.api)
		}
		if got := choiceTypeFromAPI(c.api); got != c.ct {
			t.Errorf("choiceTypeFromAPI(%q) = %q, want %q", c.api, got, c.ct)
		}
	}
}
