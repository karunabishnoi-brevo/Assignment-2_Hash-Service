package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hash-service/backend/internal/model"
	"hash-service/backend/internal/service"
)

func newTestHandler() *HashHandler {
	svc := service.NewHashService()
	return NewHashHandler(svc)
}

func TestHandleHash(t *testing.T) {
	h := newTestHandler()

	tests := []struct {
		name           string
		method         string
		body           string
		wantStatus     int
		wantHash       string
		wantErrContain string
	}{
		{
			name:       "valid alphanumeric input",
			method:     http.MethodPost,
			body:       `{"input": "abc123"}`,
			wantStatus: http.StatusOK,
			wantHash:   "6ca13d52ca",
		},
		{
			name:       "valid single character",
			method:     http.MethodPost,
			body:       `{"input": "a"}`,
			wantStatus: http.StatusOK,
			wantHash:   "ca978112ca",
		},
		{
			name:           "empty input rejected",
			method:         http.MethodPost,
			body:           `{"input": ""}`,
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "input is required",
		},
		{
			name:           "non-alphanumeric input rejected",
			method:         http.MethodPost,
			body:           `{"input": "abc-123"}`,
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "input must be alphanumeric",
		},
		{
			name:           "special characters rejected",
			method:         http.MethodPost,
			body:           `{"input": "hello world"}`,
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "input must be alphanumeric",
		},
		{
			name:           "invalid JSON body",
			method:         http.MethodPost,
			body:           `{bad json}`,
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "invalid JSON body",
		},
		{
			name:           "GET method not allowed",
			method:         http.MethodGet,
			body:           "",
			wantStatus:     http.StatusMethodNotAllowed,
			wantErrContain: "method not allowed",
		},
		{
			name:       "OPTIONS returns no content",
			method:     http.MethodOptions,
			body:       "",
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, "/api/hash", bytes.NewBufferString(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, "/api/hash", nil)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			h.HandleHash(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusNoContent {
				return
			}

			if tt.wantHash != "" {
				var resp model.HashResponse
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Hash != tt.wantHash {
					t.Errorf("hash = %q, want %q", resp.Hash, tt.wantHash)
				}
				if len(resp.Hash) != 10 {
					t.Errorf("hash length = %d, want 10", len(resp.Hash))
				}
			}

			if tt.wantErrContain != "" {
				var errResp model.ErrorResponse
				if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if errResp.Error != tt.wantErrContain {
					t.Errorf("error = %q, want %q", errResp.Error, tt.wantErrContain)
				}
			}
		})
	}
}

func TestHandleHealth(t *testing.T) {
	h := newTestHandler()

	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{
			name:       "GET returns ok",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST not allowed",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "OPTIONS returns no content",
			method:     http.MethodOptions,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/health", nil)
			rr := httptest.NewRecorder()

			h.HandleHealth(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var resp map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp["status"] != "ok" {
					t.Errorf("status = %q, want %q", resp["status"], "ok")
				}
			}
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := CORSMiddleware(inner)

	tests := []struct {
		name       string
		method     string
		wantStatus int
		checkCORS  bool
	}{
		{
			name:       "GET request includes CORS headers",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			checkCORS:  true,
		},
		{
			name:       "OPTIONS preflight returns no content",
			method:     http.MethodOptions,
			wantStatus: http.StatusNoContent,
			checkCORS:  true,
		},
		{
			name:       "POST request includes CORS headers",
			method:     http.MethodPost,
			wantStatus: http.StatusOK,
			checkCORS:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			rr := httptest.NewRecorder()

			wrapped.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, tt.wantStatus)
			}

			if tt.checkCORS {
				origin := rr.Header().Get("Access-Control-Allow-Origin")
				if origin != "*" {
					t.Errorf("Access-Control-Allow-Origin = %q, want %q", origin, "*")
				}
				methods := rr.Header().Get("Access-Control-Allow-Methods")
				if methods != "GET, POST, OPTIONS" {
					t.Errorf("Access-Control-Allow-Methods = %q, want %q", methods, "GET, POST, OPTIONS")
				}
				headers := rr.Header().Get("Access-Control-Allow-Headers")
				if headers != "Content-Type" {
					t.Errorf("Access-Control-Allow-Headers = %q, want %q", headers, "Content-Type")
				}
			}
		})
	}
}
