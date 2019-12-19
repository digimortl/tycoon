package main

import (
	"strings"
	"testing"
	"time"
)

func TestUseCases(t *testing.T) {
	tests := []struct {
		dest string
		want time.Duration
	}{
		{"A", 5 * time.Hour},
		{"AB", 5 * time.Hour},
		{"BB", 5 * time.Hour},
		{"ABB", 7 * time.Hour},
		{"AABABBAB", 29 * time.Hour},
		{"ABBBABAAABBB", 41 * time.Hour},
	}
	for _, test := range tests {
		dest := strings.Split(test.dest, "")
		got := UseCase(dest...)
		if got != test.want {
			t.Errorf("Test %s got %d; want %d", test.dest, got, test.want)
		}
	}
}

