package service

import (
	"crypto/sha256"
	"fmt"
)

// HashService provides hash generation operations.
type HashService struct{}

// NewHashService creates a new HashService.
func NewHashService() *HashService {
	return &HashService{}
}

// GenerateHash computes the SHA-256 hash of the input and returns
// the first 10 characters of the hex-encoded digest.
func (s *HashService) GenerateHash(input string) string {
	sum := sha256.Sum256([]byte(input))
	hexStr := fmt.Sprintf("%x", sum)
	return hexStr[:10]
}
