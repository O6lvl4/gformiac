package engine

import (
	"testing"

	forms "google.golang.org/api/forms/v1"
)

func TestItemSpecToRequest_ShortAnswer(t *testing.T) {
	req := itemSpecToRequest(ItemSpec{Title: "Name", Type: ItemShortAnswer, Required: true}, 0)
	assertCreateItem(t, req)

	q := req.CreateItem.Item.QuestionItem.Question
	if !q.Required {
		t.Error("expected required=true")
	}
	if q.TextQuestion == nil || q.TextQuestion.Paragraph {
		t.Error("expected paragraph=false")
	}
}

func TestItemSpecToRequest_Paragraph(t *testing.T) {
	req := itemSpecToRequest(ItemSpec{Title: "Comment", Type: ItemParagraph}, 2)
	assertCreateItem(t, req)

	if !req.CreateItem.Item.QuestionItem.Question.TextQuestion.Paragraph {
		t.Error("expected paragraph=true")
	}
	if req.CreateItem.Location.Index != 2 {
		t.Errorf("index = %d", req.CreateItem.Location.Index)
	}
}

func TestItemSpecToRequest_ChoiceTypes(t *testing.T) {
	cases := []struct {
		ct   ChoiceType
		want string
	}{
		{ChoiceRadio, "RADIO"},
		{ChoiceCheckbox, "CHECKBOX"},
		{ChoiceDropdown, "DROP_DOWN"},
	}
	for _, c := range cases {
		req := itemSpecToRequest(ItemSpec{
			Title: "Q", Type: ItemChoice,
			Choice: &ChoiceSpec{Type: c.ct, Options: []string{"A"}},
		}, 0)
		assertCreateItem(t, req)
		got := req.CreateItem.Item.QuestionItem.Question.ChoiceQuestion.Type
		if got != c.want {
			t.Errorf("ChoiceType %q: got %q, want %q", c.ct, got, c.want)
		}
	}
}

func TestItemSpecToRequest_Scale(t *testing.T) {
	req := itemSpecToRequest(ItemSpec{
		Title: "Rate", Type: ItemScale,
		Scale: &ScaleSpec{Low: 1, High: 10, LowLabel: "Bad", HighLabel: "Good"},
	}, 0)
	assertCreateItem(t, req)

	s := req.CreateItem.Item.QuestionItem.Question.ScaleQuestion
	if s.Low != 1 || s.High != 10 || s.LowLabel != "Bad" || s.HighLabel != "Good" {
		t.Errorf("scale = %+v", s)
	}
}

func TestItemSpecToRequest_DateTimePageBreak(t *testing.T) {
	cases := []struct {
		typ  ItemType
		check func(*forms.Request) bool
	}{
		{ItemDate, func(r *forms.Request) bool { return r.CreateItem.Item.QuestionItem.Question.DateQuestion != nil }},
		{ItemTime, func(r *forms.Request) bool { return r.CreateItem.Item.QuestionItem.Question.TimeQuestion != nil }},
		{ItemPageBreak, func(r *forms.Request) bool { return r.CreateItem.Item.PageBreakItem != nil }},
	}
	for _, c := range cases {
		req := itemSpecToRequest(ItemSpec{Title: "Q", Type: c.typ}, 0)
		assertCreateItem(t, req)
		if !c.check(req) {
			t.Errorf("%s: unexpected structure", c.typ)
		}
	}
}

func TestItemSpecToRequest_ReturnsNil(t *testing.T) {
	cases := []struct {
		name string
		item ItemSpec
	}{
		{"unknown type", ItemSpec{Title: "X", Type: "unknown_type"}},
		{"choice without spec", ItemSpec{Title: "X", Type: ItemChoice}},
		{"scale without spec", ItemSpec{Title: "X", Type: ItemScale}},
	}
	for _, c := range cases {
		if req := itemSpecToRequest(c.item, 0); req != nil {
			t.Errorf("%s: expected nil, got %+v", c.name, req)
		}
	}
}

func TestFormToSpec_AllTypes(t *testing.T) {
	form := &forms.Form{
		Info: &forms.Info{Title: "Test", Description: "Desc"},
		Items: []*forms.Item{
			{Title: "Name", QuestionItem: &forms.QuestionItem{Question: &forms.Question{Required: true, TextQuestion: &forms.TextQuestion{}}}},
			{Title: "Pick", QuestionItem: &forms.QuestionItem{Question: &forms.Question{ChoiceQuestion: &forms.ChoiceQuestion{Type: "CHECKBOX", Options: []*forms.Option{{Value: "A"}, {Value: "B"}}}}}},
			{Title: "Rate", QuestionItem: &forms.QuestionItem{Question: &forms.Question{ScaleQuestion: &forms.ScaleQuestion{Low: 1, High: 5}}}},
			{Title: "Section", PageBreakItem: &forms.PageBreakItem{}},
			{Title: "Date", QuestionItem: &forms.QuestionItem{Question: &forms.Question{DateQuestion: &forms.DateQuestion{}}}},
			{Title: "Time", QuestionItem: &forms.QuestionItem{Question: &forms.Question{TimeQuestion: &forms.TimeQuestion{}}}},
		},
	}

	spec := formToSpec(form)

	if spec.Title != "Test" || spec.Description != "Desc" {
		t.Errorf("info: %q %q", spec.Title, spec.Description)
	}

	expected := []ItemType{ItemShortAnswer, ItemChoice, ItemScale, ItemPageBreak, ItemDate, ItemTime}
	for i, want := range expected {
		if spec.Items[i].Type != want {
			t.Errorf("[%d] type = %q, want %q", i, spec.Items[i].Type, want)
		}
	}
}

func TestFormToSpec_ChoiceDetails(t *testing.T) {
	form := &forms.Form{
		Info: &forms.Info{Title: "T"},
		Items: []*forms.Item{{
			Title: "Q",
			QuestionItem: &forms.QuestionItem{Question: &forms.Question{
				ChoiceQuestion: &forms.ChoiceQuestion{Type: "DROP_DOWN", Options: []*forms.Option{{Value: "X"}, {Value: "Y"}}},
			}},
		}},
	}
	spec := formToSpec(form)
	c := spec.Items[0].Choice
	if c == nil || c.Type != ChoiceDropdown || len(c.Options) != 2 {
		t.Errorf("choice = %+v", c)
	}
}

func TestSpecToCreateRequests(t *testing.T) {
	spec := &FormSpec{
		Title: "T", Description: "D",
		Items: []ItemSpec{{Title: "Q1", Type: ItemShortAnswer}, {Title: "Q2", Type: ItemParagraph}},
	}
	requests := specToCreateRequests(spec)

	if len(requests) != 3 {
		t.Fatalf("requests = %d, want 3", len(requests))
	}
	if requests[0].UpdateFormInfo == nil {
		t.Error("[0] expected UpdateFormInfo")
	}
}

func TestSpecToCreateRequests_NoDescription(t *testing.T) {
	requests := specToCreateRequests(&FormSpec{
		Title: "T",
		Items: []ItemSpec{{Title: "Q1", Type: ItemShortAnswer}},
	})
	if len(requests) != 1 || requests[0].CreateItem == nil {
		t.Errorf("expected 1 CreateItem, got %d", len(requests))
	}
}

func TestSpecToUpdateRequests(t *testing.T) {
	spec := &FormSpec{
		Title: "New", Description: "D",
		Items: []ItemSpec{{Title: "Q1", Type: ItemShortAnswer}},
	}
	remote := &forms.Form{Items: []*forms.Item{{ItemId: "a"}, {ItemId: "b"}}}

	requests := specToUpdateRequests(spec, remote)

	// 1 UpdateFormInfo + 2 Delete + 1 Create = 4
	if len(requests) != 4 {
		t.Fatalf("requests = %d, want 4", len(requests))
	}
	if requests[1].DeleteItem.Location.Index != 1 {
		t.Error("expected delete index 1 first")
	}
	if requests[2].DeleteItem.Location.Index != 0 {
		t.Error("expected delete index 0 second")
	}
}

func TestBuildState(t *testing.T) {
	form := &forms.Form{
		FormId:       "f1",
		ResponderUri: "https://example.com",
		Items: []*forms.Item{
			{ItemId: "i1", QuestionItem: &forms.QuestionItem{Question: &forms.Question{QuestionId: "q1"}}},
			{ItemId: "i2", PageBreakItem: &forms.PageBreakItem{}},
		},
	}
	state := buildState(form)

	if state.FormID != "f1" || state.ResponderURL != "https://example.com" {
		t.Errorf("basic fields: %+v", state)
	}
	if len(state.ItemIDs) != 2 || state.ItemIDs[0] != "i1" {
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
			t.Errorf("toAPI(%q) = %q, want %q", c.ct, got, c.api)
		}
		if got := choiceTypeFromAPI(c.api); got != c.ct {
			t.Errorf("fromAPI(%q) = %q, want %q", c.api, got, c.ct)
		}
	}
}

func assertCreateItem(t *testing.T, req *forms.Request) {
	t.Helper()
	if req == nil || req.CreateItem == nil || req.CreateItem.Item == nil {
		t.Fatal("expected non-nil CreateItem request")
	}
}
