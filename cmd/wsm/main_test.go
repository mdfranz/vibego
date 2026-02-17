package main

import (
	"regexp"
	"testing"
	"time"
)

func TestGenerateFolderName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *regexp.Regexp
	}{
		{
			name:  "simple title",
			input: "My Project",
			want:  regexp.MustCompile(`^\d{4}-\d{2}-\d{2}_my-project$`),
		},
		{
			name:  "special characters",
			input: "Project #1! (test)",
			want:  regexp.MustCompile(`^\d{4}-\d{2}-\d{2}_project-1-test$`),
		},
		{
			name:  "multiple spaces",
			input: "  some   title  ",
			want:  regexp.MustCompile(`^\d{4}-\d{2}-\d{2}_some-title$`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateFolderName(tt.input)
			if !tt.want.MatchString(got) {
				t.Errorf("generateFolderName() = %v, want to match %v", got, tt.want)
			}
		})
	}
}

func TestGenerateFolderNameDate(t *testing.T) {
	input := "test"
	got := generateFolderName(input)
	datePrefix := time.Now().Format("2006-01-02")
	expected := datePrefix + "_test"
	if got != expected {
		t.Errorf("generateFolderName() = %v, want %v", got, expected)
	}
}
