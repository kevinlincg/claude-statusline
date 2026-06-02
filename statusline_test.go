package main

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"
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

// makeInputWithRateLimits is a small helper to build an Input carrying
// rate_limits, used by the apiUsageFromInput tests.
func makeInputWithRateLimits(fivePct float64, fiveReset int64, sevenPct float64, sevenReset int64) Input {
	var in Input
	in.RateLimits.FiveHour.UsedPercentage = fivePct
	in.RateLimits.FiveHour.ResetsAt = fiveReset
	in.RateLimits.SevenDay.UsedPercentage = sevenPct
	in.RateLimits.SevenDay.ResetsAt = sevenReset
	return in
}

func TestApiUsageFromInput(t *testing.T) {
	t.Run("absent returns nil", func(t *testing.T) {
		if got := apiUsageFromInput(Input{}); got != nil {
			t.Errorf("apiUsageFromInput(empty) = %+v, want nil", got)
		}
	})

	t.Run("both windows present", func(t *testing.T) {
		in := makeInputWithRateLimits(23.5, 1738425600, 41.2, 1738857600)
		got := apiUsageFromInput(in)
		if got == nil {
			t.Fatal("apiUsageFromInput returned nil, want non-nil")
		}
		if got.FiveHour.Utilization != 23.5 {
			t.Errorf("FiveHour.Utilization = %v, want 23.5", got.FiveHour.Utilization)
		}
		if got.FiveHour.ResetsAt != "1738425600" {
			t.Errorf("FiveHour.ResetsAt = %q, want %q", got.FiveHour.ResetsAt, "1738425600")
		}
		if got.SevenDay.Utilization != 41.2 {
			t.Errorf("SevenDay.Utilization = %v, want 41.2", got.SevenDay.Utilization)
		}
		if got.SevenDay.ResetsAt != "1738857600" {
			t.Errorf("SevenDay.ResetsAt = %q, want %q", got.SevenDay.ResetsAt, "1738857600")
		}
	})

	t.Run("only five hour present", func(t *testing.T) {
		in := makeInputWithRateLimits(10, 1738425600, 0, 0)
		got := apiUsageFromInput(in)
		if got == nil {
			t.Fatal("apiUsageFromInput returned nil, want non-nil")
		}
		if got.FiveHour.ResetsAt != "1738425600" {
			t.Errorf("FiveHour.ResetsAt = %q, want %q", got.FiveHour.ResetsAt, "1738425600")
		}
		// Seven-day absent: resets_at stays empty so the display shows a placeholder.
		if got.SevenDay.ResetsAt != "" {
			t.Errorf("SevenDay.ResetsAt = %q, want empty", got.SevenDay.ResetsAt)
		}
	})

	t.Run("only seven day present", func(t *testing.T) {
		in := makeInputWithRateLimits(0, 0, 55, 1738857600)
		got := apiUsageFromInput(in)
		if got == nil {
			t.Fatal("apiUsageFromInput returned nil, want non-nil")
		}
		if got.FiveHour.ResetsAt != "" {
			t.Errorf("FiveHour.ResetsAt = %q, want empty", got.FiveHour.ResetsAt)
		}
		if got.SevenDay.ResetsAt != "1738857600" {
			t.Errorf("SevenDay.ResetsAt = %q, want %q", got.SevenDay.ResetsAt, "1738857600")
		}
	})
}

// TestInputJSONRateLimits guards the json struct tags: a status line payload
// from Claude Code must deserialize rate_limits into the Input struct so the
// network round-trip can be skipped.
func TestInputJSONRateLimits(t *testing.T) {
	payload := `{
		"model": {"display_name": "Claude Opus 4.5"},
		"session_id": "abc123",
		"workspace": {"current_dir": "/tmp/project"},
		"context_window": {"context_window_size": 200000, "total_input_tokens": 15000, "total_output_tokens": 1200, "used_percentage": 8},
		"rate_limits": {
			"five_hour": {"used_percentage": 23.5, "resets_at": 1738425600},
			"seven_day": {"used_percentage": 41.2, "resets_at": 1738857600}
		}
	}`

	var in Input
	if err := json.Unmarshal([]byte(payload), &in); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if in.RateLimits.FiveHour.UsedPercentage != 23.5 {
		t.Errorf("five_hour.used_percentage = %v, want 23.5", in.RateLimits.FiveHour.UsedPercentage)
	}
	if in.RateLimits.FiveHour.ResetsAt != 1738425600 {
		t.Errorf("five_hour.resets_at = %v, want 1738425600", in.RateLimits.FiveHour.ResetsAt)
	}
	if in.RateLimits.SevenDay.ResetsAt != 1738857600 {
		t.Errorf("seven_day.resets_at = %v, want 1738857600", in.RateLimits.SevenDay.ResetsAt)
	}

	// The decoded rate_limits should flow straight into an APIUsage.
	if got := apiUsageFromInput(in); got == nil {
		t.Error("apiUsageFromInput(decoded) = nil, want non-nil")
	}
}

func TestFormatTimeLeftShortDurations(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		when     time.Time
		expected string
	}{
		// 30s buffer absorbs test execution time before truncation.
		{"minutes only", now.Add(7*time.Minute + 30*time.Second), "7m"},
		{"hours and minutes", now.Add(2*time.Hour + 5*time.Minute + 30*time.Second), "2h5m"},
		{"days and hours", now.Add(27*time.Hour + 30*time.Second), "1d3h"},
		{"in the past", now.Add(-1 * time.Hour), "now"},
	}

	for _, tt := range tests {
		t.Run("epoch/"+tt.name, func(t *testing.T) {
			epoch := strconv.FormatInt(tt.when.Unix(), 10)
			if got := formatTimeLeftShort(epoch); got != tt.expected {
				t.Errorf("formatTimeLeftShort(%q) = %q, want %q", epoch, got, tt.expected)
			}
		})
		t.Run("rfc3339/"+tt.name, func(t *testing.T) {
			iso := tt.when.Format(time.RFC3339)
			if got := formatTimeLeftShort(iso); got != tt.expected {
				t.Errorf("formatTimeLeftShort(%q) = %q, want %q", iso, got, tt.expected)
			}
		})
	}
}

func TestCalculateCostWithCache(t *testing.T) {
	// 1M cache-read tokens at Sonnet cache-read rate ($0.30/1M).
	usage := SessionUsageResult{CacheReadTokens: 1000000}
	got := calculateCost(usage, "Sonnet")
	if got < 0.29 || got > 0.31 {
		t.Errorf("calculateCost(cache read 1M, Sonnet) = %v, want ~0.30", got)
	}

	// Unknown model type falls back to Sonnet pricing (output $15/1M).
	out := SessionUsageResult{OutputTokens: 1000000}
	fallback := calculateCost(out, "Mystery")
	sonnet := calculateCost(out, "Sonnet")
	if fallback != sonnet {
		t.Errorf("unknown model cost = %v, want Sonnet fallback %v", fallback, sonnet)
	}

	// Combined token kinds sum together.
	mixed := SessionUsageResult{InputTokens: 1000000, OutputTokens: 1000000}
	if c := calculateCost(mixed, "Haiku"); c < 5.9 || c > 6.1 {
		// Haiku: input $1 + output $5 per 1M = $6.
		t.Errorf("calculateCost(mixed, Haiku) = %v, want ~6.0", c)
	}
}

func TestFormatProjectPathHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		t.Skip("no home dir available")
	}
	got := formatProjectPath(home + "/work/repo")
	if got != "~/work/repo" {
		t.Errorf("formatProjectPath(home/work/repo) = %q, want %q", got, "~/work/repo")
	}
}

func TestGetModelTypeFromInput(t *testing.T) {
	mk := func(id, display string) Input {
		var in Input
		in.Model.ID = id
		in.Model.DisplayName = display
		return in
	}

	tests := []struct {
		name string
		in   Input
		want string
	}{
		{"id opus 4.8", mk("claude-opus-4-8", "Whatever"), "Opus"},
		{"id opus fast", mk("claude-opus-4-8-fast", "Opus 4.8"), "Opus"},
		{"id sonnet", mk("claude-sonnet-4-6", "x"), "Sonnet"},
		{"id haiku", mk("claude-haiku-4-5", "x"), "Haiku"},
		{"no id falls back to display name", mk("", "Claude Opus 4.8"), "Opus"},
		{"unknown id falls back to display name", mk("gpt-9", "Claude Haiku 4.5"), "Haiku"},
		{"both unknown defaults to Sonnet", mk("gpt-9", "Mystery"), "Sonnet"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getModelTypeFromInput(tt.in); got != tt.want {
				t.Errorf("getModelTypeFromInput(%+v) = %q, want %q", tt.in.Model, got, tt.want)
			}
		})
	}
}

func TestIsFastMode(t *testing.T) {
	mk := func(id, display string) Input {
		var in Input
		in.Model.ID = id
		in.Model.DisplayName = display
		return in
	}

	tests := []struct {
		name string
		in   Input
		want bool
	}{
		{"id suffix fast", mk("claude-opus-4-8-fast", "Opus 4.8"), true},
		{"display contains Fast", mk("claude-opus-4-8", "Opus 4.8 (Fast)"), true},
		{"standard", mk("claude-opus-4-8", "Opus 4.8"), false},
		{"empty", mk("", ""), false},
		{"substring false positive in display", mk("claude-sonnet-4-6", "Steadfast Sonnet"), false},
		{"substring false positive in id", mk("claude-fastpath-1", "x"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFastMode(tt.in); got != tt.want {
				t.Errorf("isFastMode(%+v) = %v, want %v", tt.in.Model, got, tt.want)
			}
		})
	}
}

func TestCalculateCostModeFast(t *testing.T) {
	usage := SessionUsageResult{InputTokens: 1000000, OutputTokens: 1000000}
	std := calculateCostMode(usage, "Opus", false)
	fast := calculateCostMode(usage, "Opus", true)

	// Opus standard: input $5 + output $25 per 1M = $30.
	if std < 29.9 || std > 30.1 {
		t.Errorf("standard Opus cost = %v, want ~30", std)
	}
	// Fast mode is ~2x standard.
	if fast < 2*std-0.01 || fast > 2*std+0.01 {
		t.Errorf("fast cost = %v, want ~2x standard (%v)", fast, 2*std)
	}
	// calculateCost (2-arg) must equal the non-fast path.
	if calculateCost(usage, "Opus") != std {
		t.Errorf("calculateCost != calculateCostMode(...false)")
	}
}

func TestApplyCostDelta(t *testing.T) {
	t.Run("monotonic rises accumulate once", func(t *testing.T) {
		var s UsageStats
		applyCostDelta(&s, "sess", 1.0)
		applyCostDelta(&s, "sess", 1.0) // same render again — no double count
		applyCostDelta(&s, "sess", 2.5) // rose by 1.5
		if s.TotalCost != 2.5 {
			t.Errorf("TotalCost = %v, want 2.5", s.TotalCost)
		}
		if s.SessionCosts["sess"] != 2.5 {
			t.Errorf("baseline = %v, want 2.5", s.SessionCosts["sess"])
		}
	})

	t.Run("two sessions sum independently", func(t *testing.T) {
		var s UsageStats
		applyCostDelta(&s, "a", 3.0)
		applyCostDelta(&s, "b", 4.0)
		applyCostDelta(&s, "a", 5.0) // a rose by 2
		if s.TotalCost != 9.0 {      // 3 + 4 + 2
			t.Errorf("TotalCost = %v, want 9.0", s.TotalCost)
		}
	})

	t.Run("lower incoming value is ignored, no corruption", func(t *testing.T) {
		var s UsageStats
		applyCostDelta(&s, "sess", 10.0) // transcript over-estimate
		applyCostDelta(&s, "sess", 6.0)  // authoritative lands lower — ignored
		if s.TotalCost != 10.0 {
			t.Errorf("TotalCost = %v, want 10.0 (no decrease)", s.TotalCost)
		}
		// Subsequent climb above the baseline resumes counting from it.
		applyCostDelta(&s, "sess", 12.0) // rose by 2 over the 10 baseline
		if s.TotalCost != 12.0 {
			t.Errorf("TotalCost = %v, want 12.0", s.TotalCost)
		}
	})

	t.Run("nil map is initialized", func(t *testing.T) {
		s := UsageStats{} // SessionCosts nil
		applyCostDelta(&s, "sess", 1.0)
		if s.SessionCosts == nil || s.SessionCosts["sess"] != 1.0 {
			t.Errorf("nil map not handled: %+v", s.SessionCosts)
		}
	})
}

func TestFormatCCVersion(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"2.1.90", "v2.1.90"},
		{"v2.1.90", "v2.1.90"},
		{"  2.1.158  ", "v2.1.158"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := formatCCVersion(tt.in); got != tt.want {
				t.Errorf("formatCCVersion(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestInputJSONCostAndModel guards the json tags for the cost/model.id/version
// fields that P1+P2 newly consume.
func TestInputJSONCostAndModel(t *testing.T) {
	payload := `{
		"model": {"display_name": "Opus 4.8", "id": "claude-opus-4-8-fast"},
		"version": "2.1.200",
		"cost": {"total_cost_usd": 4.37, "total_lines_added": 156, "total_lines_removed": 23},
		"session_id": "abc"
	}`

	var in Input
	if err := json.Unmarshal([]byte(payload), &in); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if in.Model.ID != "claude-opus-4-8-fast" {
		t.Errorf("model.id = %q", in.Model.ID)
	}
	if in.Version != "2.1.200" {
		t.Errorf("version = %q", in.Version)
	}
	if in.Cost.TotalCostUSD != 4.37 {
		t.Errorf("cost.total_cost_usd = %v", in.Cost.TotalCostUSD)
	}
	if in.Cost.TotalLinesAdded != 156 || in.Cost.TotalLinesRemoved != 23 {
		t.Errorf("lines = +%d/-%d", in.Cost.TotalLinesAdded, in.Cost.TotalLinesRemoved)
	}
	if !isFastMode(in) {
		t.Error("isFastMode(decoded) = false, want true")
	}
	if getModelTypeFromInput(in) != "Opus" {
		t.Errorf("getModelTypeFromInput(decoded) = %q, want Opus", getModelTypeFromInput(in))
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
