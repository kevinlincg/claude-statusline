package themes

import (
	"strings"
	"testing"
)

func TestFormatAheadBehind(t *testing.T) {
	green := "\033[32m"
	orange := "\033[33m"

	tests := []struct {
		name        string
		ahead       int
		behind      int
		aheadColor  string
		behindColor string
		wantSubstrs []string // must all be present
		wantEmpty   bool
	}{
		{name: "both, colored", ahead: 2, behind: 1, aheadColor: green, behindColor: orange,
			wantSubstrs: []string{"↑2", "↓1", green, orange, Reset}},
		{name: "ahead only", ahead: 3, behind: 0, aheadColor: green, behindColor: orange,
			wantSubstrs: []string{"↑3"}},
		{name: "behind only", ahead: 0, behind: 5, aheadColor: green, behindColor: orange,
			wantSubstrs: []string{"↓5"}},
		{name: "level is empty", ahead: 0, behind: 0, aheadColor: green, behindColor: orange,
			wantEmpty: true},
		{name: "no color leaves no ANSI", ahead: 1, behind: 1, aheadColor: "", behindColor: "",
			wantSubstrs: []string{"↑1", "↓1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatAheadBehind(tt.ahead, tt.behind, tt.aheadColor, tt.behindColor)
			if tt.wantEmpty {
				if got != "" {
					t.Errorf("expected empty, got %q", got)
				}
				return
			}
			for _, s := range tt.wantSubstrs {
				if !strings.Contains(got, s) {
					t.Errorf("FormatAheadBehind(%d,%d,...) = %q; missing %q",
						tt.ahead, tt.behind, got, s)
				}
			}
			// behind-only must not emit an ahead arrow, and vice versa
			if tt.ahead == 0 && strings.Contains(got, "↑") {
				t.Errorf("unexpected ↑ when ahead=0: %q", got)
			}
			if tt.behind == 0 && strings.Contains(got, "↓") {
				t.Errorf("unexpected ↓ when behind=0: %q", got)
			}
			// no-color case must contain no ANSI reset
			if tt.aheadColor == "" && tt.behindColor == "" && strings.Contains(got, Reset) {
				t.Errorf("expected no ANSI reset in uncolored output: %q", got)
			}
		})
	}
}

func TestFormatGitExtras(t *testing.T) {
	green, orange, dim := "\033[32m", "\033[33m", "\033[2m"

	t.Run("all segments present", func(t *testing.T) {
		d := StatusData{GitAhead: 2, GitBehind: 1, GitStash: 3, GitSHA: "a1b2c3d"}
		got := FormatGitExtras(d, green, orange, dim)
		for _, s := range []string{"↑2", "↓1", "⚑3", "@a1b2c3d"} {
			if !strings.Contains(got, s) {
				t.Errorf("FormatGitExtras missing %q in %q", s, got)
			}
		}
	})

	t.Run("clean repo yields empty", func(t *testing.T) {
		d := StatusData{GitSHA: ""} // no ahead/behind, no stash, no sha
		if got := FormatGitExtras(d, green, orange, dim); got != "" {
			t.Errorf("expected empty for clean repo, got %q", got)
		}
	})

	t.Run("stash and sha only", func(t *testing.T) {
		d := StatusData{GitStash: 1, GitSHA: "deadbee"}
		got := FormatGitExtras(d, green, orange, dim)
		if strings.Contains(got, "↑") || strings.Contains(got, "↓") {
			t.Errorf("unexpected ahead/behind arrows: %q", got)
		}
		if !strings.Contains(got, "⚑1") || !strings.Contains(got, "@deadbee") {
			t.Errorf("missing stash/sha: %q", got)
		}
	})

	t.Run("empty dim color emits no reset for stash/sha", func(t *testing.T) {
		d := StatusData{GitStash: 2, GitSHA: "abc1234"}
		got := FormatGitExtras(d, "", "", "")
		if strings.Contains(got, Reset) {
			t.Errorf("expected no ANSI reset with empty colors: %q", got)
		}
		if !strings.Contains(got, "⚑2") || !strings.Contains(got, "@abc1234") {
			t.Errorf("missing stash/sha: %q", got)
		}
	})
}

func TestFormatLinesChanged(t *testing.T) {
	green, red := "\033[32m", "\033[31m"
	if got := FormatLinesChanged(0, 0, green, red); got != "" {
		t.Errorf("expected empty when no lines changed, got %q", got)
	}
	got := FormatLinesChanged(120, 45, green, red)
	if !strings.Contains(got, "+120") || !strings.Contains(got, "-45") {
		t.Errorf("FormatLinesChanged = %q; want +120 and -45", got)
	}
	// only additions still renders (removed shows 0)
	if got := FormatLinesChanged(5, 0, green, red); !strings.Contains(got, "+5") {
		t.Errorf("FormatLinesChanged(5,0) = %q; want +5", got)
	}
}

func TestFormatTokensPerSec(t *testing.T) {
	tests := []struct {
		rate float64
		want string
	}{
		{0, ""},
		{-1, ""},
		{1250, "1.2k/s"},
		{500, "500/s"},
		{2000000, "2.0M/s"},
	}
	for _, tt := range tests {
		if got := FormatTokensPerSec(tt.rate); got != tt.want {
			t.Errorf("FormatTokensPerSec(%v) = %q; want %q", tt.rate, got, tt.want)
		}
	}
}

func TestFormatTokens(t *testing.T) {
	tests := []struct {
		tokens   int64
		expected string
	}{
		{0, "0"},
		{500, "500"},
		{999, "999"},
		{1000, "1.0k"},
		{1500, "1.5k"},
		{10000, "10.0k"},
		{999999, "1000.0k"},
		{1000000, "1.0M"},
		{1500000, "1.5M"},
		{10000000, "10.0M"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatTokens(tt.tokens)
			if result != tt.expected {
				t.Errorf("FormatTokens(%d) = %q, want %q", tt.tokens, result, tt.expected)
			}
		})
	}
}

func TestFormatCost(t *testing.T) {
	tests := []struct {
		cost     float64
		expected string
	}{
		{0, "$0.00"},
		{0.01, "$0.01"},
		{0.99, "$0.99"},
		{1.00, "$1.00"},
		{5.50, "$5.50"},
		{10.00, "$10.0"},
		{99.99, "$100.0"},
		{100.00, "$100"},
		{999.99, "$1000"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatCost(tt.cost)
			if result != tt.expected {
				t.Errorf("FormatCost(%v) = %q, want %q", tt.cost, result, tt.expected)
			}
		})
	}
}

func TestVisibleWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"empty", "", 0},
		{"ascii", "hello", 5},
		{"with ansi", "\033[31mred\033[0m", 3},
		{"cjk", "中文", 4},
		{"mixed", "a中b", 4},
		{"complex ansi", "\033[38;2;100;100;100mtext\033[0m", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VisibleWidth(tt.input)
			if result != tt.expected {
				t.Errorf("VisibleWidth(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRuneWidth(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		expected int
	}{
		{"ascii", 'a', 1},
		{"digit", '0', 1},
		{"cjk", '中', 2},
		{"emoji", '😀', 2},
		{"variation selector", '\uFE0F', 0},
		{"zero width", '\u200B', 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RuneWidth(tt.r)
			if result != tt.expected {
				t.Errorf("RuneWidth(%q) = %d, want %d", tt.r, result, tt.expected)
			}
		})
	}
}

func TestGetTheme(t *testing.T) {
	// Test that default theme exists
	theme, ok := GetTheme("classic_framed")
	if !ok {
		t.Error("classic_framed theme should exist")
	}
	if theme == nil {
		t.Error("theme should not be nil")
	}

	// Test non-existent theme
	_, ok = GetTheme("non_existent_theme_12345")
	if ok {
		t.Error("non-existent theme should return false")
	}
}

func TestListThemes(t *testing.T) {
	themes := ListThemes()
	if len(themes) == 0 {
		t.Error("ListThemes should return at least one theme")
	}

	// Verify each theme has required methods
	for _, theme := range themes {
		if theme.Name() == "" {
			t.Error("theme name should not be empty")
		}
		if theme.Description() == "" {
			t.Error("theme description should not be empty")
		}
	}
}

func TestGenerateBar(t *testing.T) {
	tests := []struct {
		name        string
		percent     int
		width       int
		filledChar  string
		emptyChar   string
		filledColor string
		emptyColor  string
	}{
		{"empty bar", 0, 10, "█", "░", "", ""},
		{"full bar", 100, 10, "█", "░", "", ""},
		{"half bar", 50, 10, "█", "░", "", ""},
		{"over 100", 150, 10, "█", "░", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateBar(tt.percent, tt.width, tt.filledChar, tt.emptyChar, tt.filledColor, tt.emptyColor)
			if result == "" {
				t.Error("GenerateBar should not return empty string")
			}
		})
	}
}

func TestPadFunctions(t *testing.T) {
	t.Run("PadLeft", func(t *testing.T) {
		result := PadLeft("test", 10)
		if len(result) != 10 {
			t.Errorf("PadLeft result length = %d, want 10", len(result))
		}
		if !strings.HasSuffix(result, "test") {
			t.Errorf("PadLeft should end with 'test': %q", result)
		}
	})

	t.Run("PadRight", func(t *testing.T) {
		result := PadRight("test", 10)
		if len(result) != 10 {
			t.Errorf("PadRight result length = %d, want 10", len(result))
		}
		if !strings.HasPrefix(result, "test") {
			t.Errorf("PadRight should start with 'test': %q", result)
		}
	})

	t.Run("PadCenter", func(t *testing.T) {
		result := PadCenter("test", 10)
		if len(result) != 10 {
			t.Errorf("PadCenter result length = %d, want 10", len(result))
		}
		if !strings.Contains(result, "test") {
			t.Errorf("PadCenter should contain 'test': %q", result)
		}
	})
}

func TestGetModelConfig(t *testing.T) {
	tests := []struct {
		modelType     string
		expectedColor string
	}{
		{"Opus", ColorGold},
		{"Sonnet", ColorCyan},
		{"Haiku", ColorPink},
		{"Unknown", ColorCyan},
	}

	for _, tt := range tests {
		t.Run(tt.modelType, func(t *testing.T) {
			color, icon := GetModelConfig(tt.modelType)
			if color != tt.expectedColor {
				t.Errorf("GetModelConfig(%q) color = %q, want %q", tt.modelType, color, tt.expectedColor)
			}
			if icon == "" {
				t.Errorf("GetModelConfig(%q) icon should not be empty", tt.modelType)
			}
		})
	}
}

func TestGetBarColor(t *testing.T) {
	tests := []struct {
		percent       int
		expectedColor string
	}{
		{0, ColorBrightGreen},
		{25, ColorBrightGreen},
		{49, ColorBrightGreen},
		{50, ColorBrightYellow},
		{74, ColorBrightYellow},
		{75, ColorRed},
		{100, ColorRed},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.percent)), func(t *testing.T) {
			color, _ := GetBarColor(tt.percent)
			if color != tt.expectedColor {
				t.Errorf("GetBarColor(%d) = %q, want %q", tt.percent, color, tt.expectedColor)
			}
		})
	}
}

func TestGetContextColor(t *testing.T) {
	tests := []struct {
		percent       int
		expectedColor string
	}{
		{0, ColorCtxGreen},
		{59, ColorCtxGreen},
		{60, ColorCtxGold},
		{79, ColorCtxGold},
		{80, ColorCtxRed},
		{100, ColorCtxRed},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.percent)), func(t *testing.T) {
			color := GetContextColor(tt.percent)
			if color != tt.expectedColor {
				t.Errorf("GetContextColor(%d) = %q, want %q", tt.percent, color, tt.expectedColor)
			}
		})
	}
}
