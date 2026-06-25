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
	"time"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/auth"
	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/config"
	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/questionnaire"
)

type authTestRepository struct {
	usersByEmail map[string]auth.User
}

func newAuthTestRepository() *authTestRepository {
	return &authTestRepository{usersByEmail: map[string]auth.User{}}
}

func (r *authTestRepository) CreateUser(_ context.Context, user auth.User) (auth.User, error) {
	if _, exists := r.usersByEmail[user.Email]; exists {
		return auth.User{}, auth.ErrEmailAlreadyExists
	}
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	r.usersByEmail[user.Email] = user
	return user, nil
}

func (r *authTestRepository) UpsertAdmin(_ context.Context, user auth.User) (auth.User, error) {
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	r.usersByEmail[user.Email] = user
	return user, nil
}

func (r *authTestRepository) FindUserByEmail(_ context.Context, email string) (auth.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return auth.User{}, auth.ErrInvalidCredentials
	}
	return user, nil
}

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

func newTestRouter(t *testing.T) (http.Handler, *auth.Service) {
	t.Helper()

	authRepo := newAuthTestRepository()
	authService := auth.NewService(authRepo, "test-secret", time.Hour)
	if _, err := authService.EnsureAdmin(context.Background(), "admin@example.com", "very-strong-password"); err != nil {
		t.Fatalf("EnsureAdmin returned error: %v", err)
	}

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
				TagID:       "psych_drive",
				TagName:     "Драйв / Экстраверсия",
				Weight:      3,
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

	return NewRouter(cfg, logger, authService, questionnaire.NewService(questionnaireRepo)), authService
}

func TestHealthEndpointReturnsOK(t *testing.T) {
	router, _ := newTestRouter(t)

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
	router, _ := newTestRouter(t)

	registerBody := bytes.NewBufferString(`{"email":"user@example.com","password":"secret123"}`)
	registerResponse := httptest.NewRecorder()
	router.ServeHTTP(registerResponse, httptest.NewRequest(http.MethodPost, "/api/auth/register", registerBody))
	if registerResponse.Code != http.StatusCreated {
		t.Fatalf("expected register status 201, got %d: %s", registerResponse.Code, registerResponse.Body.String())
	}

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

func TestAdminCreateFragranceRequiresAdminToken(t *testing.T) {
	router, _ := newTestRouter(t)

	userRegister := httptest.NewRecorder()
	router.ServeHTTP(userRegister, httptest.NewRequest(
		http.MethodPost,
		"/api/auth/register",
		bytes.NewBufferString(`{"email":"user@example.com","password":"secret123"}`),
	))
	var userAuth authResponse
	if err := json.Unmarshal(userRegister.Body.Bytes(), &userAuth); err != nil {
		t.Fatalf("register response is not JSON: %v", err)
	}

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/admin/fragrances",
		bytes.NewBufferString(`{"name":"Miami Shake","brand":"Juliette Has A Gun","price":8393}`),
	)
	createRequest.Header.Set("Authorization", "Bearer "+userAuth.AccessToken)
	createResponse := httptest.NewRecorder()
	router.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusForbidden {
		t.Fatalf("expected user token to be forbidden, got %d", createResponse.Code)
	}
}
