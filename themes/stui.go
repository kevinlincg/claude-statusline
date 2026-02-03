package themes

import (
	"fmt"
	"strings"
	"unicode"
)

// StuiTheme s-tui 壓力測試監視器風格
type StuiTheme struct{}

func init() {
	RegisterTheme(&StuiTheme{})
}

func (t *StuiTheme) Name() string {
	return "stui"
}

func (t *StuiTheme) Description() string {
	return "s-tui：CPU 壓力測試監視器，頻率溫度圖風格"
}

const (
	StuiGreen     = "\033[38;2;0;255;0m"
	StuiDarkGreen = "\033[38;2;0;180;0m"
	StuiDimGreen  = "\033[38;2;0;100;0m"
	StuiYellow    = "\033[38;2;255;255;0m"
	StuiOrange    = "\033[38;2;255;165;0m"
	StuiRed       = "\033[38;2;255;0;0m"
	StuiCyan      = "\033[38;2;0;255;255m"
	StuiWhite     = "\033[38;2;255;255;255m"
	StuiGray      = "\033[38;2;128;128;128m"
	StuiDark      = "\033[38;2;64;64;64m"
)

func (t *StuiTheme) Render(data StatusData) string {
	var sb strings.Builder
	width := 80

	// Title bar (s-tui style)
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	title := fmt.Sprintf("%s┌─ %s%s%s %s%s%s ─",
		StuiDimGreen,
		modelColor, modelIcon, data.ModelName, Reset,
		StuiGray, data.Version)
	if data.UpdateAvailable {
		title += StuiYellow + " [UP]" + Reset
	}
	titleVisLen := stuiDisplayWidth(title)
	titlePad := width - titleVisLen - 1 // -1 for the closing ┐
	if titlePad < 0 {
		titlePad = 0
	}
	sb.WriteString(title + StuiDimGreen + strings.Repeat("─", titlePad) + "┐" + Reset + "\n")

	// Frequency graph style header
	sb.WriteString(stuiPadLine(StuiDimGreen+"│"+Reset+" "+StuiGreen+"Utilization"+Reset, width, StuiDimGreen+"│"+Reset))

	// Context usage as "frequency" graph
	graphWidth := 60
	ctxGraph := t.generateStuiGraph(data.ContextPercent, graphWidth)
	ctxColor := StuiGreen
	if data.ContextPercent > 75 {
		ctxColor = StuiRed
	} else if data.ContextPercent > 50 {
		ctxColor = StuiYellow
	}
	line1 := fmt.Sprintf("%s│%s %sCTX%s %s%s%s %s%3d%%%s",
		StuiDimGreen, Reset,
		StuiCyan, Reset,
		ctxColor, ctxGraph, Reset,
		StuiWhite, data.ContextPercent, Reset)
	sb.WriteString(stuiPadLine(line1, width, StuiDimGreen+"│"+Reset))

	// API 5hr as "temperature" graph
	api5Graph := t.generateStuiGraph(data.API5hrPercent, graphWidth)
	line2 := fmt.Sprintf("%s│%s %s5HR%s %s%s%s %s%3d%%%s",
		StuiDimGreen, Reset,
		StuiCyan, Reset,
		StuiYellow, api5Graph, Reset,
		StuiWhite, data.API5hrPercent, Reset)
	sb.WriteString(stuiPadLine(line2, width, StuiDimGreen+"│"+Reset))

	// API 7day as "power" graph
	api7Graph := t.generateStuiGraph(data.API7dayPercent, graphWidth)
	line3 := fmt.Sprintf("%s│%s %s7DY%s %s%s%s %s%3d%%%s",
		StuiDimGreen, Reset,
		StuiCyan, Reset,
		StuiOrange, api7Graph, Reset,
		StuiWhite, data.API7dayPercent, Reset)
	sb.WriteString(stuiPadLine(line3, width, StuiDimGreen+"│"+Reset))

	// Separator
	sb.WriteString(StuiDimGreen + "├" + strings.Repeat("─", width-2) + "┤" + Reset + "\n")

	// Summary section header
	sb.WriteString(stuiPadLine(StuiDimGreen+"│"+Reset+" "+StuiGreen+"Summary"+Reset, width, StuiDimGreen+"│"+Reset))

	// Path and git info
	line4 := fmt.Sprintf("%s│%s %sPath:%s %s%s%s",
		StuiDimGreen, Reset,
		StuiDarkGreen, Reset,
		StuiWhite, ShortenPath(data.ProjectPath, 30), Reset)
	if data.GitBranch != "" {
		line4 += fmt.Sprintf("  %s[%s]%s", StuiCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line4 += fmt.Sprintf(" %s+%d%s", StuiGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line4 += fmt.Sprintf(" %s*%d%s", StuiYellow, data.GitDirty, Reset)
		}
	}
	sb.WriteString(stuiPadLine(line4, width, StuiDimGreen+"│"+Reset))

	// Stats row 1
	line5 := fmt.Sprintf("%s│%s %sTokens:%s %s%s%s  %sMsgs:%s %s%d%s  %sTime:%s %s%s%s  %sHit:%s %s%d%%%s",
		StuiDimGreen, Reset,
		StuiDarkGreen, Reset, StuiWhite, FormatTokens(data.TokenCount), Reset,
		StuiDarkGreen, Reset, StuiWhite, data.MessageCount, Reset,
		StuiDarkGreen, Reset, StuiWhite, data.SessionTime, Reset,
		StuiDarkGreen, Reset, StuiCyan, data.CacheHitRate, Reset)
	sb.WriteString(stuiPadLine(line5, width, StuiDimGreen+"│"+Reset))

	// Stats row 2: costs
	line6 := fmt.Sprintf("%s│%s %sSession:%s %s%s%s  %sRate:%s %s%s/h%s  %sDay:%s %s%s%s  %sLeft:%s %s%s%s",
		StuiDimGreen, Reset,
		StuiDarkGreen, Reset, StuiGreen, FormatCostShort(data.SessionCost), Reset,
		StuiDarkGreen, Reset, StuiRed, FormatCostShort(data.BurnRate), Reset,
		StuiDarkGreen, Reset, StuiYellow, FormatCostShort(data.DayCost), Reset,
		StuiDarkGreen, Reset, StuiGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(stuiPadLine(line6, width, StuiDimGreen+"│"+Reset))

	// Bottom border
	sb.WriteString(StuiDimGreen + "└" + strings.Repeat("─", width-2) + "┘" + Reset + "\n")

	return sb.String()
}

func stuiPadLine(line string, targetWidth int, suffix string) string {
	visible := stuiDisplayWidth(line)
	suffixLen := stuiDisplayWidth(suffix)
	padding := targetWidth - visible - suffixLen
	if padding < 0 {
		padding = 0
	}
	return line + strings.Repeat(" ", padding) + suffix + "\n"
}

// stuiDisplayWidth calculates display width accounting for:
// - ANSI escape codes (0 width)
// - Emojis and wide characters (2 width)
// - Regular ASCII and box drawing (1 width)
func stuiDisplayWidth(s string) int {
	inEscape := false
	width := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
		} else if inEscape {
			if r == 'm' {
				inEscape = false
			}
		} else {
			if stuiIsWideChar(r) {
				width += 2
			} else {
				width += 1
			}
		}
	}
	return width
}

// stuiIsWideChar checks if a rune is a wide character (emoji or CJK)
func stuiIsWideChar(r rune) bool {
	// Emojis
	if r >= 0x1F300 && r <= 0x1F9FF {
		return true
	}
	if r >= 0x2600 && r <= 0x26FF {
		return true
	}
	if r >= 0x2700 && r <= 0x27BF {
		return true
	}
	// Box Drawing characters are 1 width
	if r >= 0x2500 && r <= 0x257F {
		return false
	}
	// Block elements are 1 width
	if r >= 0x2580 && r <= 0x259F {
		return false
	}
	// CJK characters
	if unicode.Is(unicode.Han, r) {
		return true
	}
	// Model icons (emojis)
	switch r {
	case '💛', '💙', '💚':
		return true
	}
	return false
}

func (t *StuiTheme) generateStuiGraph(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := percent * width / 100
	empty := width - filled

	var graph strings.Builder

	// Create filled portion with slight "noise" for graph effect
	for i := 0; i < filled; i++ {
		if i%3 == 0 {
			graph.WriteString("▓")
		} else if i%3 == 1 {
			graph.WriteString("█")
		} else {
			graph.WriteString("▒")
		}
	}
	if empty > 0 {
		graph.WriteString(StuiDark)
		graph.WriteString(strings.Repeat("░", empty))
		graph.WriteString(Reset)
	}

	return graph.String()
}
