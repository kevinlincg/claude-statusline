package main

import (
	"testing"
)

func TestGetModelType(t *testing.T) {
	tests := []struct {
		displayName string
		expected    string
	}{
		{"Claude Opus 4.5", "Opus"},
		{"Claude Sonnet 4", "Sonnet"},
		{"Claude Haiku 3.5", "Haiku"},
		{"Opus 4.5", "Opus"},
		{"Sonnet", "Sonnet"},
		{"Unknown Model", "Sonnet"}, // Default fallback
	}

	for _, tt := range tests {
		t.Run(tt.displayName, func(t *testing.T) {
			result := getModelType(tt.displayName)
			if result != tt.expected {
				t.Errorf("getModelType(%q) = %q, want %q", tt.displayName, result, tt.expected)
			}
		})
	}
}

func TestFormatModelName(t *testing.T) {
	tests := []struct {
		displayName string
		expected    string
	}{
		{"Claude Opus 4.5", "Opus 4.5"},
		{"Claude Sonnet 4", "Sonnet 4"},
		{"Claude Haiku 3.5", "Haiku 3.5"},
		{"Opus", "Opus"},
		{"Sonnet", "Sonnet"},
		{"Haiku", "Haiku"},
		{"Unknown", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.displayName, func(t *testing.T) {
			result := formatModelName(tt.displayName)
			if result != tt.expected {
				t.Errorf("formatModelName(%q) = %q, want %q", tt.displayName, result, tt.expected)
			}
		})
	}
}

func TestFormatProjectPath(t *testing.T) {
	tests := []struct {
		name     string
		fullPath string
		expected string
	}{
		{"absolute path", "/tmp/project", "/tmp/project"},
		{"root path", "/", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatProjectPath(tt.fullPath)
			// Just verify it returns something non-empty
			if result == "" {
				t.Errorf("formatProjectPath(%q) returned empty string", tt.fullPath)
			}
		})
	}
}

func TestFormatTimeLeftShort(t *testing.T) {
	tests := []struct {
		name     string
		isoTime  string
		expected string
	}{
		{"invalid time", "invalid", "?"},
		{"empty time", "", "?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimeLeftShort(tt.isoTime)
			if result != tt.expected {
				t.Errorf("formatTimeLeftShort(%q) = %q, want %q", tt.isoTime, result, tt.expected)
			}
		})
	}
}

func TestCalculateCost(t *testing.T) {
	tests := []struct {
		name      string
		usage     SessionUsageResult
		modelType string
		minCost   float64
		maxCost   float64
	}{
		{
			name: "zero usage",
			usage: SessionUsageResult{
				InputTokens:  0,
				OutputTokens: 0,
			},
			modelType: "Sonnet",
			minCost:   0,
			maxCost:   0,
		},
		{
			name: "some input tokens",
			usage: SessionUsageResult{
				InputTokens:  1000000, // 1M tokens
				OutputTokens: 0,
			},
			modelType: "Sonnet",
			minCost:   2.9, // Slightly less than $3
			maxCost:   3.1, // Slightly more than $3
		},
		{
			name: "opus pricing",
			usage: SessionUsageResult{
				InputTokens:  1000000,
				OutputTokens: 0,
			},
			modelType: "Opus",
			minCost:   4.9,
			maxCost:   5.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateCost(tt.usage, tt.modelType)
			if result < tt.minCost || result > tt.maxCost {
				t.Errorf("calculateCost() = %v, want between %v and %v", result, tt.minCost, tt.maxCost)
			}
		})
	}
}

func TestVersionVariables(t *testing.T) {
	// Verify version variables exist and have default values
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if Commit == "" {
		t.Error("Commit should not be empty")
	}
	if Date == "" {
		t.Error("Date should not be empty")
	}
}
