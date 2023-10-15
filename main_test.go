package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCountLines(t *testing.T) {
	tests := []struct {
		input  string
		output int
	}{
		{"", 0},
		{"hello", 1},
		{"hello\nworld", 2},
		{"hello\nworld\n", 2},
		{"hello\nworld\n\n", 3},
		{"\n", 1},
	}

	for _, test := range tests {
		r := []byte(test.input)
		result := countLines(bytes.NewReader(r))
		if result != test.output {
			t.Errorf("countLines(%q) = %d; want %d", test.input, result, test.output)
		}
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		input  string
		output int
	}{
		{"", 0},
		{"hello", 1},
		{"hello\nworld", 2},
		{"hello\nworld\n", 2},
		{"hello\nworld\n\n", 2},
		{"\n", 0},
	}

	for _, test := range tests {
		r := []byte(test.input)
		result := countWords(bytes.NewReader(r))
		if result != test.output {
			t.Errorf("countWords(%q) = %d, want %d", test.input, result, test.output)
		}
	}
}

func TestCountBytes(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"hello", 5},
		{"壞過凱婷", 12}, // Each character in "壞過凱婷" takes 3 bytes in UTF-8 encoding
		{"", 0},
		{"\n", 1},
		{"\t", 1},
		{"hello\nworld", 11},
		{"hello\nworld\n", 12},
	}

	for _, tt := range tests {
		result := countBytes(strings.NewReader(tt.input))
		if result != tt.expected {
			t.Errorf("countBytes(%q) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestCountCharacters(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{"壞過凱婷", 4}, // "壞過凱婷 is a Cantonese slang, and it has 4 characters.
		{"", 0},
		{"\n", 1},
		{"\t", 1},
		{"hello\nworld", 11},
		{"hello\nworld\n", 12},
	}

	for _, tt := range tests {
		result := countCharacters(strings.NewReader(tt.input))
		if result != tt.expected {
			t.Errorf("countCharacters(%q) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}
