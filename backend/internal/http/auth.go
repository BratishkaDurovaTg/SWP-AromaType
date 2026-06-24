package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/auth"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken string       `json:"accessToken"`
	User        userResponse `json:"user"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (r *Router) handleRegister(w http.ResponseWriter, req *http.Request) {
	var payload registerRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON.")
		return
	}

	result, err := r.authService.Register(req.Context(), payload.Email, payload.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, "invalid_input", "Email and password with at least 6 characters are required.")
		case errors.Is(err, auth.ErrEmailAlreadyExists):
			writeError(w, http.StatusConflict, "email_already_exists", "User with this email already exists.")
		default:
			r.logger.Error("failed to register user", "error", err)
			writeError(w, http.StatusInternalServerError, "internal_error", "Failed to register user.")
		}
		return
	}

	writeJSON(w, http.StatusCreated, newAuthResponse(result))
}

func (r *Router) handleLogin(w http.ResponseWriter, req *http.Request) {
	var payload loginRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON.")
		return
	}

	result, err := r.authService.Login(req.Context(), payload.Email, payload.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password.")
			return
		}

		r.logger.Error("failed to login user", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to login user.")
		return
	}

	writeJSON(w, http.StatusOK, newAuthResponse(result))
}

func newAuthResponse(result auth.AuthResult) authResponse {
	return authResponse{
		AccessToken: result.AccessToken,
		User: userResponse{
			ID:    result.User.ID,
			Email: result.User.Email,
			Role:  result.User.Role,
		},
	}
}
