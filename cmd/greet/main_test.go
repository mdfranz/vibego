package main

import (
	"errors"
	"testing"
)

func TestGreet(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid name",
			input:   "Gemini",
			want:    "Hello, Gemini!",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Greet(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Greet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, ErrEmptyName) {
				t.Errorf("Greet() error = %v, want ErrEmptyName", err)
			}
			if got != tt.want {
				t.Errorf("Greet() got = %v, want %v", got, tt.want)
			}
		})
	}
}
