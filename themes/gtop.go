package themes

import (
	"fmt"
	"strings"
)

// GtopTheme gtop 簡約系統監視器風格
type GtopTheme struct{}

func init() {
	RegisterTheme(&GtopTheme{})
}

func (t *GtopTheme) Name() string {
	return "gtop"
}

func (t *GtopTheme) Description() string {
	return "gtop：簡約系統監視器，火花圖與乾淨排版"
}

const (
	GtopGreen      = "\033[38;2;98;214;164m"
	GtopCyan       = "\033[38;2;137;221;255m"
	GtopMagenta    = "\033[38;2;255;121;198m"
	GtopYellow     = "\033[38;2;241;250;140m"
	GtopRed        = "\033[38;2;255;85;85m"
	GtopBlue       = "\033[38;2;139;233;253m"
	GtopWhite      = "\033[38;2;248;248;242m"
	GtopGray       = "\033[38;2;98;114;164m"
	GtopDark       = "\033[38;2;68;71;90m"
	GtopBrightGreen = "\033[38;2;80;250;123m"
)

func (t *GtopTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Simple top border
	sb.WriteString(GtopDark + "┌" + strings.Repeat("─", 78) + "┐" + Reset + "\n")

	// Header with model info
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	header := fmt.Sprintf("%s│%s %s%s%s%s%s  %s%s%s",
		GtopDark, Reset,
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		GtopGray, data.Version, Reset)
	if data.UpdateAvailable {
		header += GtopYellow + " ↑" + Reset
	}
	sb.WriteString(PadRight(header, 79))
	sb.WriteString(GtopDark + "│" + Reset + "\n")

	// Sparkline-style CPU graph
	cpuSparkline := t.generateSparkline(data.ContextPercent)
	cpuColor := GtopGreen
	if data.ContextPercent > 75 {
		cpuColor = GtopRed
	} else if data.ContextPercent > 50 {
		cpuColor = GtopYellow
	}

	line1 := fmt.Sprintf("%s│%s %sCPU%s %s%s%s %s%3d%%%s",
		GtopDark, Reset,
		GtopCyan, Reset,
		cpuColor, cpuSparkline, Reset,
		GtopWhite, data.ContextPercent, Reset)
	sb.WriteString(PadRight(line1, 79))
	sb.WriteString(GtopDark + "│" + Reset + "\n")

	// Memory sparkline
	memSparkline := t.generateSparkline(data.API5hrPercent)
	line2 := fmt.Sprintf("%s│%s %sMEM%s %s%s%s %s%3d%%%s  %s%s%s left",
		GtopDark, Reset,
		GtopMagenta, Reset,
		GtopMagenta, memSparkline, Reset,
		GtopWhite, data.API5hrPercent, Reset,
		GtopGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(PadRight(line2, 79))
	sb.WriteString(GtopDark + "│" + Reset + "\n")

	// Network sparkline
	netSparkline := t.generateSparkline(data.API7dayPercent)
	line3 := fmt.Sprintf("%s│%s %sNET%s %s%s%s %s%3d%%%s  %s%s%s left",
		GtopDark, Reset,
		GtopBlue, Reset,
		GtopBlue, netSparkline, Reset,
		GtopWhite, data.API7dayPercent, Reset,
		GtopGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(PadRight(line3, 79))
	sb.WriteString(GtopDark + "│" + Reset + "\n")

	// Separator
	sb.WriteString(GtopDark + "├" + strings.Repeat("─", 78) + "┤" + Reset + "\n")

	// Process info style
	line4 := fmt.Sprintf("%s│%s %sPROC%s %s%s%s",
		GtopDark, Reset,
		GtopGray, Reset,
		GtopWhite, ShortenPath(data.ProjectPath, 35), Reset)
	if data.GitBranch != "" {
		line4 += fmt.Sprintf("  %s⎇%s %s%s%s", GtopGray, Reset, GtopCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line4 += fmt.Sprintf(" %s+%d%s", GtopBrightGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line4 += fmt.Sprintf(" %s*%d%s", GtopYellow, data.GitDirty, Reset)
		}
	}
	sb.WriteString(PadRight(line4, 79))
	sb.WriteString(GtopDark + "│" + Reset + "\n")

	// Stats in clean columns
	line5 := fmt.Sprintf("%s│%s %sTOKENS%s %s%s%s  %sMSGS%s %s%d%s  %sTIME%s %s%s%s  %sHIT%s %s%d%%%s",
		GtopDark, Reset,
		GtopGray, Reset, GtopWhite, FormatTokens(data.TokenCount), Reset,
		GtopGray, Reset, GtopWhite, data.MessageCount, Reset,
		GtopGray, Reset, GtopWhite, data.SessionTime, Reset,
		GtopGray, Reset, GtopCyan, data.CacheHitRate, Reset)
	sb.WriteString(PadRight(line5, 79))
	sb.WriteString(GtopDark + "│" + Reset + "\n")

	// Cost row
	line6 := fmt.Sprintf("%s│%s %sSESSION%s %s%s%s  %sRATE%s %s%s/h%s  %sDAY%s %s%s%s",
		GtopDark, Reset,
		GtopGray, Reset, GtopGreen, FormatCostShort(data.SessionCost), Reset,
		GtopGray, Reset, GtopRed, FormatCostShort(data.BurnRate), Reset,
		GtopGray, Reset, GtopYellow, FormatCostShort(data.DayCost), Reset)
	sb.WriteString(PadRight(line6, 79))
	sb.WriteString(GtopDark + "│" + Reset + "\n")

	// Bottom border
	sb.WriteString(GtopDark + "└" + strings.Repeat("─", 78) + "┘" + Reset + "\n")

	return sb.String()
}

func (t *GtopTheme) generateSparkline(percent int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	// Sparkline characters: ▁▂▃▄▅▆▇█
	sparkChars := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	width := 20
	var spark strings.Builder

	// Generate a fake "history" based on current value with some variation
	for i := 0; i < width; i++ {
		// Create slight variation for visual interest
		variation := (i * 7) % 15 - 7
		val := percent + variation
		if val < 0 {
			val = 0
		}
		if val > 100 {
			val = 100
		}
		charIdx := val * 7 / 100
		if charIdx > 7 {
			charIdx = 7
		}
		spark.WriteRune(sparkChars[charIdx])
	}

	return spark.String()
}
