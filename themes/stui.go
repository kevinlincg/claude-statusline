package themes

import (
	"fmt"
	"strings"
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
	StuiGreen       = "\033[38;2;0;255;0m"
	StuiDarkGreen   = "\033[38;2;0;180;0m"
	StuiDimGreen    = "\033[38;2;0;100;0m"
	StuiYellow      = "\033[38;2;255;255;0m"
	StuiOrange      = "\033[38;2;255;165;0m"
	StuiRed         = "\033[38;2;255;0;0m"
	StuiCyan        = "\033[38;2;0;255;255m"
	StuiWhite       = "\033[38;2;255;255;255m"
	StuiGray        = "\033[38;2;128;128;128m"
	StuiDark        = "\033[38;2;64;64;64m"
	StuiBgGreen     = "\033[48;2;0;60;0m"
)

func (t *StuiTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Title bar (s-tui style)
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	title := fmt.Sprintf("%s┌─ %s%s%s%s %s%s%s ─",
		StuiDimGreen,
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		StuiGray, data.Version)
	if data.UpdateAvailable {
		title += StuiYellow + " [UPDATE]" + Reset
	}
	titlePad := 80 - 30 // approximate
	sb.WriteString(title + StuiDimGreen + strings.Repeat("─", titlePad) + "┐" + Reset + "\n")

	// Frequency graph style header
	sb.WriteString(StuiDimGreen + "│" + Reset + " " + StuiGreen + "Utilization" + Reset)
	sb.WriteString(strings.Repeat(" ", 67))
	sb.WriteString(StuiDimGreen + "│" + Reset + "\n")

	// Context usage as "frequency" graph
	ctxGraph := t.generateStuiGraph(data.ContextPercent, 70)
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
	sb.WriteString(PadRight(line1, 79))
	sb.WriteString(StuiDimGreen + "│" + Reset + "\n")

	// API 5hr as "temperature" graph
	api5Graph := t.generateStuiGraph(data.API5hrPercent, 70)
	line2 := fmt.Sprintf("%s│%s %s5HR%s %s%s%s %s%3d%%%s",
		StuiDimGreen, Reset,
		StuiCyan, Reset,
		StuiYellow, api5Graph, Reset,
		StuiWhite, data.API5hrPercent, Reset)
	sb.WriteString(PadRight(line2, 79))
	sb.WriteString(StuiDimGreen + "│" + Reset + "\n")

	// API 7day as "power" graph
	api7Graph := t.generateStuiGraph(data.API7dayPercent, 70)
	line3 := fmt.Sprintf("%s│%s %s7DY%s %s%s%s %s%3d%%%s",
		StuiDimGreen, Reset,
		StuiCyan, Reset,
		StuiOrange, api7Graph, Reset,
		StuiWhite, data.API7dayPercent, Reset)
	sb.WriteString(PadRight(line3, 79))
	sb.WriteString(StuiDimGreen + "│" + Reset + "\n")

	// Separator
	sb.WriteString(StuiDimGreen + "├" + strings.Repeat("─", 78) + "┤" + Reset + "\n")

	// Summary section header
	sb.WriteString(StuiDimGreen + "│" + Reset + " " + StuiGreen + "Summary" + Reset)
	sb.WriteString(strings.Repeat(" ", 71))
	sb.WriteString(StuiDimGreen + "│" + Reset + "\n")

	// Path and git info
	line4 := fmt.Sprintf("%s│%s %sPath:%s %s%s%s",
		StuiDimGreen, Reset,
		StuiDarkGreen, Reset,
		StuiWhite, ShortenPath(data.ProjectPath, 35), Reset)
	if data.GitBranch != "" {
		line4 += fmt.Sprintf("  %s[%s]%s", StuiCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line4 += fmt.Sprintf(" %s+%d%s", StuiGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line4 += fmt.Sprintf(" %s*%d%s", StuiYellow, data.GitDirty, Reset)
		}
	}
	sb.WriteString(PadRight(line4, 79))
	sb.WriteString(StuiDimGreen + "│" + Reset + "\n")

	// Stats row 1
	line5 := fmt.Sprintf("%s│%s %sTokens:%s %s%s%s  %sMsgs:%s %s%d%s  %sTime:%s %s%s%s  %sHit:%s %s%d%%%s",
		StuiDimGreen, Reset,
		StuiDarkGreen, Reset, StuiWhite, FormatTokens(data.TokenCount), Reset,
		StuiDarkGreen, Reset, StuiWhite, data.MessageCount, Reset,
		StuiDarkGreen, Reset, StuiWhite, data.SessionTime, Reset,
		StuiDarkGreen, Reset, StuiCyan, data.CacheHitRate, Reset)
	sb.WriteString(PadRight(line5, 79))
	sb.WriteString(StuiDimGreen + "│" + Reset + "\n")

	// Stats row 2: costs
	line6 := fmt.Sprintf("%s│%s %sSession:%s %s%s%s  %sRate:%s %s%s/h%s  %sDay:%s %s%s%s  %sLeft:%s %s%s%s",
		StuiDimGreen, Reset,
		StuiDarkGreen, Reset, StuiGreen, FormatCostShort(data.SessionCost), Reset,
		StuiDarkGreen, Reset, StuiRed, FormatCostShort(data.BurnRate), Reset,
		StuiDarkGreen, Reset, StuiYellow, FormatCostShort(data.DayCost), Reset,
		StuiDarkGreen, Reset, StuiGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(PadRight(line6, 79))
	sb.WriteString(StuiDimGreen + "│" + Reset + "\n")

	// Bottom border
	sb.WriteString(StuiDimGreen + "└" + strings.Repeat("─", 78) + "┘" + Reset + "\n")

	return sb.String()
}

func (t *StuiTheme) generateStuiGraph(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	// s-tui uses braille-like patterns for graphs
	// We'll simulate a simple bar with gradient
	filled := percent * width / 100
	empty := width - filled

	var graph strings.Builder

	// Create filled portion with slight "noise" for graph effect
	for i := 0; i < filled; i++ {
		// Alternate between full and partial blocks for texture
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
