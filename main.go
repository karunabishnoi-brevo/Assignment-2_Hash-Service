package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
)

const (
	charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hashLen = 10
)

var alphanumericRe = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

var (
	hashCache   = map[string]string{}
	hashCacheMu sync.Mutex
)

// these are 3 structures
type hashRequest struct {
	Input string `json:"input"`
}

type hashResponse struct {
	Input string `json:"input"`
	Hash  string `json:"hash"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// - two calls with the same input return same hashes (unique per input)
// - the input is bound into the hash (not just ignored)
func generateHash(input string) (string, error) {
	hashCacheMu.Lock()

	if h, ok := hashCache[input]; ok {
		hashCacheMu.Unlock()
		return h, nil
	}
	hashCacheMu.Unlock()

	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("rand read: %w", err)
	}

	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(input))
	digest := h.Sum(nil) // 32 bytes

	//converting first 8 bytes into uint64
	//Why?
	//To create a large random number
	//To encode into base62
	// To generate 10-character hash
	//Because 64 bits is mathematically enough
	var num uint64
	for i := 0; i < 8; i++ {
		num = num<<8 | uint64(digest[i])
	}

	out := make([]byte, hashLen)
	base := uint64(len(charset))
	for i := hashLen - 1; i >= 0; i-- {
		out[i] = charset[num%base]
		num /= base
	}
	result := string(out)

	hashCacheMu.Lock()
	hashCache[input] = result
	hashCacheMu.Unlock()

	return result, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// line 101 — if not POST → return 405 error
// line 107 — read and parse JSON body
// line 108 — if bad JSON → return 400 error
// line 112 — if input empty → return 400 error
// line 116 — if not alphanumeric → return 400 error
// line 121 — all good → call generateHash()

func hashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed, use POST"})
		return
	}

	var req hashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	if req.Input == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "input is required"})
		return
	}
	if !alphanumericRe.MatchString(req.Input) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "input must be alphanumeric (a-z, A-Z, 0-9)"})
		return
	}

	hash, err := generateHash(req.Input)
	if err != nil {
		log.Printf("hash generation failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "hash generation failed"})
		return
	}

	writeJSON(w, http.StatusOK, hashResponse{Input: req.Input, Hash: hash})
}

// func healthHandler(w http.ResponseWriter, _ *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	_, _ = fmt.Fprint(w, "ok")
// }

// browser opens — serves static/index.html from disk
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	mux := http.NewServeMux() // It creates a router
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/hash", hashHandler)
	// mux.HandleFunc("/health", healthHandler)

	addr := ":9000"
	log.Printf("hash service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
