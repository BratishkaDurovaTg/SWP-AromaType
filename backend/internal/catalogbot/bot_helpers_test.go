package catalogbot

import (
	"strings"
	"testing"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/questionnaire"
)

func TestKeyboardBuilders(t *testing.T) {
	main := mainMenuKeyboard()
	if !main.ResizeKeyboard || len(main.Keyboard) != 2 {
		t.Fatalf("unexpected main menu keyboard: %#v", main)
	}

	items := make([]questionnaire.Fragrance, 45)
	for index := range items {
		items[index] = questionnaire.Fragrance{
			ID:    "item",
			Brand: "Brand",
			Name:  "Name",
		}
	}
	list := listKeyboard(items)
	if got := len(list.InlineKeyboard); got != listButtonLimit+1 {
		t.Fatalf("expected limited list keyboard rows, got %d", got)
	}

	item := itemKeyboard("miami-shake")
	if len(item.InlineKeyboard) != 2 || item.InlineKeyboard[0][0].CallbackData != "edit:miami-shake" {
		t.Fatalf("unexpected item keyboard: %#v", item)
	}

	edit := editFieldKeyboard("miami-shake")
	if len(edit.InlineKeyboard) < 2 {
		t.Fatalf("expected edit field keyboard rows, got %#v", edit)
	}
}

func TestPromptsCoverKnownFieldsAndAddSteps(t *testing.T) {
	for step := addStepID; step <= addStepPhoto; step++ {
		if addPrompt(step) == "" {
			t.Fatalf("empty add prompt for step %d", step)
		}
	}
	if addPrompt(999) == "" {
		t.Fatal("expected fallback add prompt")
	}

	for _, field := range []string{
		"name", "brand", "gender", "price", "volumes", "price_3ml", "price_5ml", "price_10ml", "description", "accords",
		"top", "middle", "base", "psychotype", "scores", "active", "image_url", "unknown",
	} {
		if fieldPrompt(field) == "" {
			t.Fatalf("empty field prompt for %q", field)
		}
	}
}

func TestOptionalList(t *testing.T) {
	if got := optionalList("-"); len(got) != 0 {
		t.Fatalf("expected empty list for dash, got %#v", got)
	}
	got := optionalList("чай, чай, мускус")
	want := []string{"чай", "мускус"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("unexpected optional list: got %#v want %#v", got, want)
	}
}

func TestApplyFieldUpdatesAllEditableFields(t *testing.T) {
	item := questionnaire.Fragrance{
		ID:         "miami-shake",
		Name:       "Old",
		Brand:      "Old Brand",
		Gender:     questionnaire.GenderUnisex,
		Price:      "100",
		Psychotype: "balanced",
		PsychotypeScores: questionnaire.PsychotypeScores{
			Drive: 25, Focus: 25, Aesthetic: 25, Power: 25,
		},
		IsActive: true,
	}

	updates := []struct {
		field string
		value string
	}{
		{field: "name", value: "Miami Shake"},
		{field: "brand", value: "Juliette Has A Gun"},
		{field: "gender", value: "female"},
		{field: "price", value: "8393"},
		{field: "volumes", value: "3:8393, 5:12990, 10:18990"},
		{field: "price_5ml", value: "13990"},
		{field: "description", value: "Summer fragrance"},
		{field: "top", value: "Клубника"},
		{field: "middle", value: "Мороженое"},
		{field: "base", value: "Абсолют ванили"},
		{field: "accords", value: "Сладкий, Фруктовый"},
		{field: "psychotype", value: "aesthetic"},
		{field: "scores", value: "20,20,40,20"},
		{field: "active", value: "no"},
		{field: "image_url", value: "/uploads/miami.jpg"},
	}

	for _, update := range updates {
		if err := applyField(&item, update.field, update.value); err != nil {
			t.Fatalf("applyField(%q) returned error: %v", update.field, err)
		}
	}

	if item.Name != "Miami Shake" || item.Brand != "Juliette Has A Gun" || item.Price != "8393" || item.Gender != "female" {
		t.Fatalf("required fields were not updated: %#v", item)
	}
	if len(item.VolumeOptions) != 3 || item.VolumeOptions[0].VolumeML != 3 || item.VolumeOptions[1].Price != 13990 {
		t.Fatalf("volume options were not updated: %#v", item.VolumeOptions)
	}
	if item.Psychotype != "aesthetic" || item.PsychotypeScores.Aesthetic != 40 || item.IsActive {
		t.Fatalf("psychotype fields were not updated: %#v", item)
	}
	if item.ImageURL != "/uploads/miami.jpg" || len(item.MainAccords) != 2 {
		t.Fatalf("catalog fields were not updated: %#v", item)
	}
}

func TestApplyFieldRejectsInvalidValues(t *testing.T) {
	item := questionnaire.Fragrance{
		ID:         "miami-shake",
		Name:       "Miami Shake",
		Brand:      "Juliette Has A Gun",
		Gender:     questionnaire.GenderFemale,
		Price:      "8393",
		Psychotype: "aesthetic",
		PsychotypeScores: questionnaire.PsychotypeScores{
			Drive: 20, Focus: 20, Aesthetic: 40, Power: 20,
		},
	}

	for field, value := range map[string]string{
		"price":      "-1",
		"volumes":    "3",
		"gender":     "unknown",
		"psychotype": "unknown",
		"scores":     "30,30,30,30",
		"active":     "maybe",
		"unknown":    "value",
	} {
		if err := applyField(&item, field, value); err == nil {
			t.Fatalf("expected error for field %q value %q", field, value)
		}
	}
}

func TestValidateAndFormatFragrance(t *testing.T) {
	item := questionnaire.Fragrance{
		ID:            "miami-shake",
		Name:          "Miami Shake",
		Brand:         "Juliette Has A Gun",
		Gender:        questionnaire.GenderFemale,
		Price:         "8393",
		VolumeOptions: []questionnaire.VolumeOption{{VolumeML: 3, Price: 8393}},
		Description:   "Summer fragrance",
		TopNotes:      []string{"Клубника"},
		MiddleNotes:   []string{"Мороженое"},
		BaseNotes:     []string{"Абсолют ванили"},
		MainAccords:   []string{"Сладкий"},
		Psychotype:    "aesthetic",
		PsychotypeScores: questionnaire.PsychotypeScores{
			Drive: 20, Focus: 20, Aesthetic: 40, Power: 20,
		},
		IsActive: true,
	}
	if err := validateFragrance(item); err != nil {
		t.Fatalf("validateFragrance returned error: %v", err)
	}

	text := formatFragrance(item)
	for _, expected := range []string{"miami-shake", "Miami Shake", "3:8393", "Клубника"} {
		if !strings.Contains(text, expected) {
			t.Fatalf("formatted fragrance does not contain %q: %s", expected, text)
		}
	}
	if emptyDash("") != "-" || emptyDash("photo") != "photo" {
		t.Fatal("emptyDash returned unexpected value")
	}
}
