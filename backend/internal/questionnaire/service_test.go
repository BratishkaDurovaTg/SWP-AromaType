package questionnaire

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"testing"
)

type fakeFragranceRepository struct {
	questions       []Question
	optionWeights   map[string]map[string]int
	tagNames        map[string]string
	fragranceTags   []FragranceTagRow
	fragrances      map[string]Fragrance
	created         Fragrance
	createdTagIDs   []string
	createShouldErr error
}

func (r *fakeFragranceRepository) GetQuestions(_ context.Context) ([]Question, error) {
	return r.questions, nil
}

func (r *fakeFragranceRepository) GetOptionTagWeights(_ context.Context, optionID string) (map[string]int, error) {
	return r.optionWeights[optionID], nil
}

func (r *fakeFragranceRepository) GetTagNames(_ context.Context, tagIDs []string) (map[string]string, error) {
	result := make(map[string]string, len(tagIDs))
	for _, id := range tagIDs {
		result[id] = r.tagNames[id]
	}
	return result, nil
}

func (r *fakeFragranceRepository) GetActiveFragranceTags(_ context.Context) ([]FragranceTagRow, error) {
	return r.fragranceTags, nil
}

func (r *fakeFragranceRepository) GetFragranceByID(_ context.Context, id string) (Fragrance, error) {
	fragrance, ok := r.fragrances[id]
	if !ok {
		return Fragrance{}, ErrFragranceNotFound
	}
	return fragrance, nil
}

func (r *fakeFragranceRepository) CreateFragrance(_ context.Context, fragrance Fragrance, tagIDs []string) (Fragrance, error) {
	if r.createShouldErr != nil {
		return Fragrance{}, r.createShouldErr
	}
	r.created = fragrance
	r.createdTagIDs = tagIDs
	return fragrance, nil
}

func TestRecommendRanksFragrancesAndBuildsProfile(t *testing.T) {
	repo := &fakeFragranceRepository{
		optionWeights: map[string]map[string]int{
			"answer-focus": {"psych_focus": 3},
		},
		tagNames: map[string]string{
			"psych_focus": "Интеллект / Фокус",
			"psych_drive": "Драйв / Экстраверсия",
		},
		fragranceTags: []FragranceTagRow{
			{
				ID:          "fragrance-1",
				Name:        "Clean Tea",
				Brand:       "Aroma Lab",
				ImageURL:    "/uploads/clean-tea.png",
				Price:       "1500",
				TopNotes:    []byte(`["Зеленый чай","Бергамот"]`),
				MiddleNotes: []byte(`["Ирис"]`),
				BaseNotes:   []byte(`["Белый мускус"]`),
				MainAccords: []byte(`["Свежий","Зеленый"]`),
				Psychotype:  PsychotypeFocus,
				PsychotypeScores: []byte(`{
					"drive": 20,
					"focus": 100,
					"aesthetic": 45,
					"power": 20
				}`),
				TagID:   "psych_focus",
				TagName: "Интеллект / Фокус",
				Weight:  3,
			},
			{
				ID:          "fragrance-1",
				Name:        "Clean Tea",
				Brand:       "Aroma Lab",
				ImageURL:    "/uploads/clean-tea.png",
				Price:       "1500",
				TopNotes:    []byte(`["Зеленый чай","Бергамот"]`),
				MiddleNotes: []byte(`["Ирис"]`),
				BaseNotes:   []byte(`["Белый мускус"]`),
				MainAccords: []byte(`["Свежий","Зеленый"]`),
				Psychotype:  PsychotypeFocus,
				PsychotypeScores: []byte(`{
					"drive": 20,
					"focus": 100,
					"aesthetic": 45,
					"power": 20
				}`),
				TagID:   "psych_drive",
				TagName: "Драйв / Экстраверсия",
				Weight:  1,
			},
		},
	}

	result, err := NewService(repo).Recommend(context.Background(), []string{"answer-focus", "answer-focus", " "})
	if err != nil {
		t.Fatalf("Recommend returned error: %v", err)
	}

	if result.Profile.Name != "Интеллект и фокус" {
		t.Fatalf("expected focus profile, got %q", result.Profile.Name)
	}
	if result.TotalItems != 1 {
		t.Fatalf("expected 1 recommendation, got %d", result.TotalItems)
	}
	if result.Items[0].ID != "fragrance-1" {
		t.Fatalf("expected top recommendation fragrance-1, got %q", result.Items[0].ID)
	}
	if result.Items[0].MatchPercent != 99 {
		t.Fatalf("expected top match percent 99, got %d", result.Items[0].MatchPercent)
	}
	if !reflect.DeepEqual(result.Items[0].KeyNotes, []string{"Зеленый чай", "Бергамот", "Ирис"}) {
		t.Fatalf("unexpected key notes: %#v", result.Items[0].KeyNotes)
	}
	if result.Items[0].Reason == "" {
		t.Fatal("recommendation reason must be present")
	}
	if result.Items[0].Psychotype != PsychotypeFocus {
		t.Fatalf("expected focus psychotype, got %q", result.Items[0].Psychotype)
	}
	if result.Items[0].PsychotypeScores.Focus != 100 {
		t.Fatalf("expected focus score 100, got %#v", result.Items[0].PsychotypeScores)
	}
}

func TestRecommendReturnsAtMostFiveItems(t *testing.T) {
	rows := make([]FragranceTagRow, 0, 8)
	for index := 0; index < 8; index++ {
		rows = append(rows, FragranceTagRow{
			ID:          "fragrance-" + strconv.Itoa(index),
			Name:        "Fragrance " + strconv.Itoa(index),
			Brand:       "Aroma Lab",
			Price:       "1500",
			TopNotes:    []byte(`["Бергамот"]`),
			MiddleNotes: []byte(`[]`),
			BaseNotes:   []byte(`[]`),
			MainAccords: []byte(`["Свежий"]`),
			Psychotype:  PsychotypeDrive,
			PsychotypeScores: []byte(`{
				"drive": 100,
				"focus": 10,
				"aesthetic": 10,
				"power": 10
			}`),
		})
	}

	repo := &fakeFragranceRepository{
		optionWeights: map[string]map[string]int{
			"answer-drive": {"psych_drive": 3},
		},
		tagNames:      map[string]string{"psych_drive": "Драйв / Экстраверсия"},
		fragranceTags: rows,
	}

	result, err := NewService(repo).Recommend(context.Background(), []string{"answer-drive"})
	if err != nil {
		t.Fatalf("Recommend returned error: %v", err)
	}
	if result.TotalItems != MaxRecommendedItems {
		t.Fatalf("expected %d recommendations, got %d", MaxRecommendedItems, result.TotalItems)
	}
}

func TestRecommendRequiresAtLeastOneAnswer(t *testing.T) {
	_, err := NewService(&fakeFragranceRepository{}).Recommend(context.Background(), []string{"", " "})
	if !errors.Is(err, ErrNoAnswers) {
		t.Fatalf("expected ErrNoAnswers, got %v", err)
	}
}

func TestCreateFragranceCleansPayloadAndUsesDefaultActive(t *testing.T) {
	repo := &fakeFragranceRepository{}
	payload := CreateFragranceRequest{
		Name:     "  Miami Shake  ",
		Brand:    "  Juliette Has A Gun ",
		ImageURL: " /uploads/miami.png ",
		Price:    8393,
		VolumeOptions: []VolumeOption{
			{VolumeML: 50, Price: 8393},
			{VolumeML: 50, Price: 8393},
			{VolumeML: 0, Price: 100},
			{VolumeML: 100, Price: -1},
		},
		TopNotes:    []string{" Клубника ", "Клубника", ""},
		MiddleNotes: []string{"Мороженое"},
		BaseNotes:   []string{"Абсолют ванили"},
		MainAccords: []string{"Сладкий", "Сладкий", "Фруктовый"},
		Psychotype:  PsychotypeAesthetic,
		PsychotypeScores: PsychotypeScores{
			Drive:     -10,
			Focus:     40,
			Aesthetic: 120,
			Power:     20,
		},
		TagIDs: []string{"psych_aesthetic", "psych_aesthetic", "sweet", ""},
	}

	fragrance, err := NewService(repo).CreateFragrance(context.Background(), payload)
	if err != nil {
		t.Fatalf("CreateFragrance returned error: %v", err)
	}

	if fragrance.ID == "" {
		t.Fatal("fragrance ID must be generated")
	}
	if fragrance.Name != "Miami Shake" || fragrance.Brand != "Juliette Has A Gun" {
		t.Fatalf("payload was not trimmed: %#v", fragrance)
	}
	if !fragrance.IsActive {
		t.Fatal("fragrance should be active by default")
	}
	if !reflect.DeepEqual(fragrance.VolumeOptions, []VolumeOption{{VolumeML: 50, Price: 8393}}) {
		t.Fatalf("unexpected volume options: %#v", fragrance.VolumeOptions)
	}
	if !reflect.DeepEqual(fragrance.TopNotes, []string{"Клубника"}) {
		t.Fatalf("unexpected top notes: %#v", fragrance.TopNotes)
	}
	if fragrance.Psychotype != PsychotypeAesthetic {
		t.Fatalf("expected aesthetic psychotype, got %q", fragrance.Psychotype)
	}
	expectedScores := PsychotypeScores{Drive: 0, Focus: 40, Aesthetic: 100, Power: 20}
	if fragrance.PsychotypeScores != expectedScores {
		t.Fatalf("unexpected psychotype scores: %#v", fragrance.PsychotypeScores)
	}
	if !reflect.DeepEqual(repo.createdTagIDs, []string{"psych_aesthetic", "sweet"}) {
		t.Fatalf("unexpected tag ids: %#v", repo.createdTagIDs)
	}
}

func TestCreateFragranceRejectsInvalidRequiredFields(t *testing.T) {
	_, err := NewService(&fakeFragranceRepository{}).CreateFragrance(context.Background(), CreateFragranceRequest{
		Name:  "",
		Brand: "Aroma Lab",
		Price: 1000,
	})
	if !errors.Is(err, ErrInvalidFragrance) {
		t.Fatalf("expected ErrInvalidFragrance, got %v", err)
	}
}
