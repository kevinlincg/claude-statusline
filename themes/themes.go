package themes

import (
	"fmt"
	"strings"
)

// ANSI color definitions
const (
	Reset = "\033[0m"
	Bold  = "\033[1m"
	Dim   = "\033[2m"

	// Basic colors
	ColorGold   = "\033[38;2;195;158;83m"
	ColorCyan   = "\033[38;2;118;170;185m"
	ColorPink   = "\033[38;2;255;182;193m"
	ColorGreen  = "\033[38;2;152;195;121m"
	ColorGray   = "\033[38;2;64;64;64m"
	ColorSilver = "\033[38;2;192;192;192m"
	ColorOrange = "\033[38;2;255;165;0m"
	ColorPurple = "\033[38;2;186;133;217m"
	ColorBlue   = "\033[38;2;100;149;237m"
	ColorRed    = "\033[38;2;220;88;88m"
	ColorDim    = "\033[38;2;170;170;170m"
	ColorYellow = "\033[38;2;255;215;0m"

	// Bright colors
	ColorBrightGreen  = "\033[38;2;80;255;100m"
	ColorBrightCyan   = "\033[38;2;0;255;255m"
	ColorBrightYellow = "\033[38;2;255;220;60m"
	ColorNeonGreen    = "\033[38;2;0;255;136m"
	ColorNeonPink     = "\033[38;2;255;0;255m"
	ColorNeonOrange   = "\033[38;2;255;150;50m"

	// Context colors
	ColorCtxGreen = "\033[38;2;108;167;108m"
	ColorCtxGold  = "\033[38;2;188;155;83m"
	ColorCtxRed   = "\033[38;2;185;102;82m"

	// Frame colors
	ColorFrame    = "\033[38;2;60;60;60m"
	ColorFrameDim = "\033[38;2;50;50;50m"
	ColorLabel    = "\033[38;2;140;140;140m"
	ColorLabelDim = "\033[38;2;100;100;100m"
	ColorTreeDim  = "\033[38;2;100;100;100m"

	// Glow bar background colors
	BgGreenGlow  = "\033[48;2;20;55;25m"
	BgYellowGlow = "\033[48;2;55;50;15m"
	BgCyanGlow   = "\033[48;2;0;60;60m"
	BgRedGlow    = "\033[48;2;60;20;20m"
)

// StatusData contains all status data to display
type StatusData struct {
	// Model info
	ModelName  string
	ModelType  string // Opus, Sonnet, Haiku
	ModelIcon  string
	ModelColor string

	// Version info
	Version         string
	UpdateAvailable bool

	// Workspace info
	ProjectPath string
	GitBranch   string
	GitStaged   int
	GitDirty    int
	GitAhead    int    // commits ahead of upstream (to push)
	GitBehind   int    // commits behind upstream (to pull)
	GitStash    int    // number of stash entries
	GitSHA      string // short commit SHA of HEAD (e.g. "a1b2c3d")

	// Session stats
	TokenCount   int64
	MessageCount int
	SessionTime  string
	CacheHitRate int
	TokensPerSec float64 // session token throughput (tokens/second)

	// Cost
	SessionCost float64
	DayCost     float64
	MonthCost   float64
	WeekCost    float64
	BurnRate    float64

	// Context
	ContextUsed    int
	ContextPercent int

	// API limits
	API5hrPercent   int
	API5hrTimeLeft  string
	API7dayPercent  int
	API7dayTimeLeft string

	// Code changes this session (from Claude Code's cost.* fields)
	LinesAdded   int
	LinesRemoved int
}

// Theme interface definition
type Theme interface {
	Name() string
	Description() string
	Render(data StatusData) string
}

// ThemeRegistry stores all registered themes
var ThemeRegistry = make(map[string]Theme)

// RegisterTheme registers a theme
func RegisterTheme(theme Theme) {
	ThemeRegistry[theme.Name()] = theme
}

// GetTheme retrieves a theme by name
func GetTheme(name string) (Theme, bool) {
	theme, ok := ThemeRegistry[name]
	return theme, ok
}

// ListThemes returns all registered themes
func ListThemes() []Theme {
	themes := make([]Theme, 0, len(ThemeRegistry))
	for _, theme := range ThemeRegistry {
		themes = append(themes, theme)
	}
	return themes
}

// Helper functions

// FormatAheadBehind returns a compact " ↑N ↓N" fragment describing how many
// commits the current branch is ahead of / behind its upstream. Each half is
// emitted only when non-zero and wrapped in the supplied ANSI color (followed
// by Reset). When a color is empty no color/Reset is added — useful for themes
// (e.g. powerline) that manage their own coloring. Returns "" when the branch
// is level with upstream or has no upstream.
func FormatAheadBehind(ahead, behind int, aheadColor, behindColor string) string {
	var b strings.Builder
	if ahead > 0 {
		if aheadColor != "" {
			b.WriteString(fmt.Sprintf(" %s↑%d%s", aheadColor, ahead, Reset))
		} else {
			b.WriteString(fmt.Sprintf(" ↑%d", ahead))
		}
	}
	if behind > 0 {
		if behindColor != "" {
			b.WriteString(fmt.Sprintf(" %s↓%d%s", behindColor, behind, Reset))
		} else {
			b.WriteString(fmt.Sprintf(" ↓%d", behind))
		}
	}
	return b.String()
}

// FormatGitExtras renders the compact git suffix shown after the branch and
// staged/dirty counts: commits ahead/behind upstream (↑N ↓N, using aheadColor/
// behindColor), stash entry count (⚑N) and the short HEAD SHA (@abc1234), the
// latter two in dimColor. Any color may be "" to emit that segment without ANSI
// (e.g. powerline themes that manage their own coloring). Each segment appears
// only when meaningful, so a clean repo with no stashes yields "".
func FormatGitExtras(data StatusData, aheadColor, behindColor, dimColor string) string {
	var b strings.Builder
	b.WriteString(FormatAheadBehind(data.GitAhead, data.GitBehind, aheadColor, behindColor))
	if data.GitStash > 0 {
		if dimColor != "" {
			b.WriteString(fmt.Sprintf(" %s⚑%d%s", dimColor, data.GitStash, Reset))
		} else {
			b.WriteString(fmt.Sprintf(" ⚑%d", data.GitStash))
		}
	}
	if data.GitSHA != "" {
		if dimColor != "" {
			b.WriteString(fmt.Sprintf(" %s@%s%s", dimColor, data.GitSHA, Reset))
		} else {
			b.WriteString(fmt.Sprintf(" @%s", data.GitSHA))
		}
	}
	return b.String()
}

// FormatLinesChanged renders session code churn (from Claude Code's cost.*
// line counters) as "+A -R", using addColor/removeColor, or "" when nothing
// changed this session.
func FormatLinesChanged(added, removed int, addColor, removeColor string) string {
	if added <= 0 && removed <= 0 {
		return ""
	}
	return fmt.Sprintf("%s+%d%s %s-%d%s", addColor, added, Reset, removeColor, removed, Reset)
}

// FormatTokensPerSec renders session token throughput compactly, e.g. "1.2k/s",
// or "" when the rate is zero (single message / no measurable span).
func FormatTokensPerSec(rate float64) string {
	if rate <= 0 {
		return ""
	}
	return FormatTokens(int64(rate)) + "/s"
}

// FormatTokens formats token count with K/M suffix
func FormatTokens(tokens int64) string {
	if tokens >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(tokens)/1000000)
	} else if tokens >= 1000 {
		return fmt.Sprintf("%.1fk", float64(tokens)/1000)
	}
	return fmt.Sprintf("%d", tokens)
}

// FormatTokensFixed formats token count with fixed width
func FormatTokensFixed(tokens int64, width int) string {
	s := FormatTokens(tokens)
	return PadLeft(s, width)
}

// FormatCost formats cost value
func FormatCost(cost float64) string {
	if cost >= 100 {
		return fmt.Sprintf("$%.0f", cost)
	} else if cost >= 10 {
		return fmt.Sprintf("$%.1f", cost)
	} else if cost >= 1 {
		return fmt.Sprintf("$%.2f", cost)
	}
	return fmt.Sprintf("$%.2f", cost)
}

// FormatCostShort formats cost value (short form)
func FormatCostShort(cost float64) string {
	if cost >= 100 {
		return fmt.Sprintf("$%.0f", cost)
	} else if cost >= 10 {
		return fmt.Sprintf("$%.0f", cost)
	}
	return fmt.Sprintf("$%.2f", cost)
}

// FormatPercent formats percentage
func FormatPercent(pct int) string {
	return fmt.Sprintf("%d%%", pct)
}

// FormatPercentFixed formats percentage with fixed width
func FormatPercentFixed(pct int, width int) string {
	s := fmt.Sprintf("%d%%", pct)
	return PadLeft(s, width)
}

// FormatNumber formats number with K/M suffix
func FormatNumber(num int) string {
	if num >= 1000000 {
		return fmt.Sprintf("%dM", num/1000000)
	} else if num >= 1000 {
		return fmt.Sprintf("%dk", num/1000)
	}
	return fmt.Sprintf("%d", num)
}

// ShortenPath shortens a path to fit within maxLen
func ShortenPath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return "~/" + parts[len(parts)-1]
	}
	return path
}

// GenerateBar generates a progress bar
func GenerateBar(percent, width int, filledChar, emptyChar string, filledColor, emptyColor string) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(filledColor)
		bar.WriteString(strings.Repeat(filledChar, filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(emptyColor)
		bar.WriteString(strings.Repeat(emptyChar, empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}

// GenerateGlowBar generates a glowing progress bar
func GenerateGlowBar(percent, width int, color, bgColor string) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(bgColor)
		bar.WriteString(Bold)
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString("\033[38;2;35;35;35m")
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}

// GetBarColor returns color based on percentage
func GetBarColor(percent int) (string, string) {
	if percent < 50 {
		return ColorBrightGreen, BgGreenGlow
	} else if percent < 75 {
		return ColorBrightYellow, BgYellowGlow
	}
	return ColorRed, BgRedGlow
}

// GetContextColor returns color based on context percentage
func GetContextColor(percent int) string {
	if percent < 60 {
		return ColorCtxGreen
	} else if percent < 80 {
		return ColorCtxGold
	}
	return ColorCtxRed
}

// PadLeft pads string on the left
func PadLeft(s string, width int) string {
	visible := VisibleWidth(s)
	if visible >= width {
		return s
	}
	return strings.Repeat(" ", width-visible) + s
}

// PadRight pads string on the right
func PadRight(s string, width int) string {
	visible := VisibleWidth(s)
	if visible >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visible)
}

// PadCenter centers string within given width
func PadCenter(s string, width int) string {
	visible := VisibleWidth(s)
	if visible >= width {
		return s
	}
	left := (width - visible) / 2
	right := width - visible - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// VisibleWidth calculates visible width (excluding ANSI codes)
func VisibleWidth(s string) int {
	// Remove ANSI escape codes
	clean := s
	for {
		start := strings.Index(clean, "\033[")
		if start == -1 {
			break
		}
		end := strings.Index(clean[start:], "m")
		if end == -1 {
			break
		}
		clean = clean[:start] + clean[start+end+1:]
	}

	width := 0
	for _, r := range clean {
		w := RuneWidth(r)
		width += w
	}
	return width
}

// RuneWidth calculates display width of a single rune
func RuneWidth(r rune) int {
	// Variation selectors - zero width
	if r >= 0xFE00 && r <= 0xFE0F {
		return 0
	}
	// Zero-width characters
	if r == 0x200B || r == 0x200C || r == 0x200D || r == 0xFEFF {
		return 0
	}
	// Combining characters - zero width
	if r >= 0x0300 && r <= 0x036F {
		return 0
	}

	// Wide characters (2 cells)
	if r >= 0x1F300 && r <= 0x1FAFF {
		return 2
	}
	if r >= 0x2300 && r <= 0x23FF {
		return 2
	}
	if r >= 0x2600 && r <= 0x26FF {
		return 2
	}
	if r >= 0x2700 && r <= 0x27BF {
		return 2
	}
	if r >= 0x2B50 && r <= 0x2B55 {
		return 2
	}
	if r >= 0x4E00 && r <= 0x9FFF {
		return 2
	}
	if r >= 0x3000 && r <= 0x303F {
		return 2
	}
	if r >= 0xFF00 && r <= 0xFFEF {
		return 2
	}

	return 1
}

// GetModelConfig returns color and icon for model type
func GetModelConfig(modelType string) (color string, icon string) {
	switch modelType {
	case "Opus":
		return ColorGold, "💛"
	case "Sonnet":
		return ColorCyan, "💠"
	case "Haiku":
		return ColorPink, "🌸"
	default:
		return ColorCyan, "◆"
	}
}
