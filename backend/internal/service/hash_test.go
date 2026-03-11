package service

import "testing"

func TestGenerateHash(t *testing.T) {
	svc := NewHashService()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple string abc123",
			input: "abc123",
			want:  "6ca13d52ca",
		},
		{
			name:  "deterministic same input produces same output",
			input: "hello",
			want:  "2cf24dba5f",
		},
		{
			name:  "single character",
			input: "a",
			want:  "ca978112ca",
		},
		{
			name:  "numeric only",
			input: "12345",
			want:  "5994471abb",
		},
		{
			name:  "mixed case alphanumeric",
			input: "AbCdEf123",
			want:  "b4f7a3c119",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.GenerateHash(tt.input)
			if len(got) != 10 {
				t.Fatalf("expected hash length 10, got %d", len(got))
			}
			if got != tt.want {
				t.Errorf("GenerateHash(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGenerateHash_Deterministic(t *testing.T) {
	svc := NewHashService()
	input := "testinput42"

	first := svc.GenerateHash(input)
	second := svc.GenerateHash(input)

	if first != second {
		t.Errorf("expected deterministic output: first=%q, second=%q", first, second)
	}
}

func TestGenerateHash_DifferentInputs(t *testing.T) {
	svc := NewHashService()

	hash1 := svc.GenerateHash("input1")
	hash2 := svc.GenerateHash("input2")

	if hash1 == hash2 {
		t.Errorf("different inputs produced same hash: %q", hash1)
	}
}
