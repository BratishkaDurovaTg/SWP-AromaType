package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/config"
	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/questionnaire"
)

type questionnaireTestRepository struct {
	questions     []questionnaire.Question
	optionWeights map[string]map[string]int
	tagNames      map[string]string
	fragranceTags []questionnaire.FragranceTagRow
	created       questionnaire.Fragrance
}

func (r *questionnaireTestRepository) GetQuestions(_ context.Context) ([]questionnaire.Question, error) {
	return r.questions, nil
}

func (r *questionnaireTestRepository) GetOptionTagWeights(_ context.Context, optionID string) (map[string]int, error) {
	return r.optionWeights[optionID], nil
}

func (r *questionnaireTestRepository) GetTagNames(_ context.Context, tagIDs []string) (map[string]string, error) {
	result := make(map[string]string, len(tagIDs))
	for _, id := range tagIDs {
		result[id] = r.tagNames[id]
	}
	return result, nil
}

func (r *questionnaireTestRepository) GetActiveFragranceTags(_ context.Context) ([]questionnaire.FragranceTagRow, error) {
	return r.fragranceTags, nil
}

func (r *questionnaireTestRepository) GetFragranceByID(_ context.Context, id string) (questionnaire.Fragrance, error) {
	if id == "fragrance-1" {
		return questionnaire.Fragrance{
			ID:          "fragrance-1",
			Name:        "Miami Shake",
			Brand:       "Juliette Has A Gun",
			ImageURL:    "/uploads/miami.png",
			Price:       "8393",
			Description: "Summer fragrance",
			IsActive:    true,
		}, nil
	}
	return questionnaire.Fragrance{}, questionnaire.ErrFragranceNotFound
}

func (r *questionnaireTestRepository) CreateFragrance(_ context.Context, fragrance questionnaire.Fragrance, _ []string) (questionnaire.Fragrance, error) {
	r.created = fragrance
	return fragrance, nil
}

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()

	questionnaireRepo := &questionnaireTestRepository{
		questions: []questionnaire.Question{
			{
				ID:        "q1",
				Text:      "Вы входите в незнакомую компанию или на встречу. Ваше действие?",
				Type:      "single_choice",
				SortOrder: 1,
				Options: []questionnaire.AnswerOption{
					{ID: "a1", QuestionID: "q1", Text: "Легко привлеку внимание, я в центре событий.", Value: "drive_attention", SortOrder: 1},
				},
			},
		},
		optionWeights: map[string]map[string]int{"a1": {"psych_drive": 3}},
		tagNames:      map[string]string{"psych_drive": "Драйв / Экстраверсия"},
		fragranceTags: []questionnaire.FragranceTagRow{
			{
				ID:          "fragrance-1",
				Name:        "Miami Shake",
				Brand:       "Juliette Has A Gun",
				ImageURL:    "/uploads/miami.png",
				Price:       "8393",
				TopNotes:    []byte(`["Клубника"]`),
				MiddleNotes: []byte(`["Мороженое"]`),
				BaseNotes:   []byte(`["Абсолют ванили"]`),
				MainAccords: []byte(`["Сладкий"]`),
				Psychotype:  questionnaire.PsychotypeDrive,
				PsychotypeScores: []byte(`{
					"drive": 100,
					"focus": 20,
					"aesthetic": 40,
					"power": 35
				}`),
				TagID:   "psych_drive",
				TagName: "Драйв / Экстраверсия",
				Weight:  3,
			},
		},
	}

	cfg := config.Config{
		AppEnv:      "test",
		OpenAPIPath: "../../docs/api/openapi.yaml",
		UploadDir:   t.TempDir(),
		CORSOrigins: []string{"*"},
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	return NewRouter(cfg, logger, questionnaire.NewService(questionnaireRepo))
}

func TestHealthEndpointReturnsOK(t *testing.T) {
	router := newTestRouter(t)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	var payload map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("health response is not JSON: %v", err)
	}
	if payload["status"] != "ok" {
		t.Fatalf("unexpected health payload: %#v", payload)
	}
}

func TestPublicQuestionnaireFlow(t *testing.T) {
	router := newTestRouter(t)

	questionsResponse := httptest.NewRecorder()
	router.ServeHTTP(questionsResponse, httptest.NewRequest(http.MethodGet, "/api/questions", nil))
	if questionsResponse.Code != http.StatusOK {
		t.Fatalf("expected questions status 200, got %d", questionsResponse.Code)
	}

	recommendBody := bytes.NewBufferString(`{"answerOptionIds":["a1"]}`)
	recommendResponse := httptest.NewRecorder()
	router.ServeHTTP(recommendResponse, httptest.NewRequest(http.MethodPost, "/api/recommendations", recommendBody))
	if recommendResponse.Code != http.StatusOK {
		t.Fatalf("expected recommendations status 200, got %d: %s", recommendResponse.Code, recommendResponse.Body.String())
	}

	var recommendation questionnaire.RecommendationResponse
	if err := json.Unmarshal(recommendResponse.Body.Bytes(), &recommendation); err != nil {
		t.Fatalf("recommendation response is not JSON: %v", err)
	}
	if recommendation.TotalItems != 1 || recommendation.Items[0].ID != "fragrance-1" {
		t.Fatalf("unexpected recommendation response: %#v", recommendation)
	}
}

func TestAuthEndpointIsNotRegistered(t *testing.T) {
	router := newTestRouter(t)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(
		http.MethodPost,
		"/api/auth/register",
		bytes.NewBufferString(`{"email":"user@example.com","password":"secret123"}`),
	))

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected removed auth endpoint to return 404, got %d", response.Code)
	}
}
