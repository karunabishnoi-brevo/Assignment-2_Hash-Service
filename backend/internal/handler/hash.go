package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"hash-service/backend/internal/model"
	"hash-service/backend/internal/service"
)

var alphanumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// HashHandler handles HTTP requests for the hash service.
type HashHandler struct {
	service *service.HashService
}

// NewHashHandler creates a new HashHandler with the given service.
func NewHashHandler(svc *service.HashService) *HashHandler {
	return &HashHandler{service: svc}
}

// HandleHash handles POST /api/hash requests.
// It reads a JSON body with an "input" field, validates that the input is
// alphanumeric, generates a 10-character hash, and returns the result.
func (h *HashHandler) HandleHash(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, model.ErrorResponse{
			Error: "method not allowed",
		})
		return
	}

	var req model.HashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "invalid JSON body",
		})
		return
	}

	if req.Input == "" {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "input is required",
		})
		return
	}

	if !alphanumericRegex.MatchString(req.Input) {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "input must be alphanumeric",
		})
		return
	}

	hash := h.service.GenerateHash(req.Input)

	writeJSON(w, http.StatusOK, model.HashResponse{
		Input: req.Input,
		Hash:  hash,
	})
}

// HandleHealth handles GET /api/health requests.
func (h *HashHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, model.ErrorResponse{
			Error: "method not allowed",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// corsMiddleware sets CORS headers on every response and handles preflight requests.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// writeJSON encodes a value as JSON and writes it to the response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
