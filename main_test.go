package main

import (
	"strings"
	"testing"
	"time"
)

func TestUseCase1(t *testing.T) {
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
		got := UseCase1(dest...)
		if got != test.want {
			t.Errorf("Test %s got %d; want %d", test.dest, got, test.want)
		}
	}
}

func TestUseCase2(t *testing.T) {
	tests := []struct {
		dest string
		want time.Duration
	}{
		{"ABBBABAAABBB", 39 * time.Hour},
	}
	for _, test := range tests {
		dest := strings.Split(test.dest, "")
		got := UseCase2(dest...)
		if got != test.want {
			t.Errorf("Test %s got %d; want %d", test.dest, got, test.want)
		}
	}
}
