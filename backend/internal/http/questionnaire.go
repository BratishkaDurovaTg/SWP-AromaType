package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/questionnaire"
)

func (r *Router) handleQuestions(w http.ResponseWriter, req *http.Request) {
	questions, err := r.questionnaireService.GetQuestions(req.Context())
	if err != nil {
		r.logger.Error("failed to load questions", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to load questions.")
		return
	}

	writeJSON(w, http.StatusOK, questions)
}

func (r *Router) handleRecommendations(w http.ResponseWriter, req *http.Request) {
	var payload questionnaire.RecommendationRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON.")
		return
	}

	result, err := r.questionnaireService.Recommend(req.Context(), payload.AnswerOptionIDs)
	if err != nil {
		if errors.Is(err, questionnaire.ErrNoAnswers) {
			writeError(w, http.StatusBadRequest, "no_answers", "At least one answer option is required.")
			return
		}

		r.logger.Error("failed to build recommendations", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to build recommendations.")
		return
	}

	writeJSON(w, http.StatusOK, result)
}
