package themes

import (
	"fmt"
	"strings"
	"unicode"
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
	GtopGreen       = "\033[38;2;98;214;164m"
	GtopCyan        = "\033[38;2;137;221;255m"
	GtopMagenta     = "\033[38;2;255;121;198m"
	GtopYellow      = "\033[38;2;241;250;140m"
	GtopRed         = "\033[38;2;255;85;85m"
	GtopBlue        = "\033[38;2;139;233;253m"
	GtopWhite       = "\033[38;2;248;248;242m"
	GtopGray        = "\033[38;2;98;114;164m"
	GtopDark        = "\033[38;2;68;71;90m"
	GtopBrightGreen = "\033[38;2;80;250;123m"
)

func (t *GtopTheme) Render(data StatusData) string {
	var sb strings.Builder
	width := 80

	// Simple top border
	sb.WriteString(GtopDark + "┌" + strings.Repeat("─", width-2) + "┐" + Reset + "\n")

	// Header with model info
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	header := fmt.Sprintf("%s│%s %s%s%s  %s%s%s",
		GtopDark, Reset,
		modelColor, modelIcon, data.ModelName, Reset,
		GtopGray, data.Version)
	if data.UpdateAvailable {
		header += GtopYellow + " ↑" + Reset
	}
	sb.WriteString(gtopPadLine(header, width, GtopDark+"│"+Reset))

	// Sparkline-style CPU graph
	cpuSparkline := t.generateSparkline(data.ContextPercent)
	cpuColor := GtopGreen
	if data.ContextPercent > 75 {
		cpuColor = GtopRed
	} else if data.ContextPercent > 50 {
		cpuColor = GtopYellow
	}

	line1 := fmt.Sprintf("%s│%s %sCPU%s %s%s%s  %s%3d%%%s",
		GtopDark, Reset,
		GtopCyan, Reset,
		cpuColor, cpuSparkline, Reset,
		GtopWhite, data.ContextPercent, Reset)
	sb.WriteString(gtopPadLine(line1, width, GtopDark+"│"+Reset))

	// Memory sparkline
	memSparkline := t.generateSparkline(data.API5hrPercent)
	line2 := fmt.Sprintf("%s│%s %sMEM%s %s%s%s  %s%3d%%%s  %s%s%s left",
		GtopDark, Reset,
		GtopMagenta, Reset,
		GtopMagenta, memSparkline, Reset,
		GtopWhite, data.API5hrPercent, Reset,
		GtopGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(gtopPadLine(line2, width, GtopDark+"│"+Reset))

	// Network sparkline
	netSparkline := t.generateSparkline(data.API7dayPercent)
	line3 := fmt.Sprintf("%s│%s %sNET%s %s%s%s  %s%3d%%%s  %s%s%s left",
		GtopDark, Reset,
		GtopBlue, Reset,
		GtopBlue, netSparkline, Reset,
		GtopWhite, data.API7dayPercent, Reset,
		GtopGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(gtopPadLine(line3, width, GtopDark+"│"+Reset))

	// Separator
	sb.WriteString(GtopDark + "├" + strings.Repeat("─", width-2) + "┤" + Reset + "\n")

	// Process info style
	line4 := fmt.Sprintf("%s│%s %sPROC%s %s%s%s",
		GtopDark, Reset,
		GtopGray, Reset,
		GtopWhite, ShortenPath(data.ProjectPath, 30), Reset)
	if data.GitBranch != "" {
		line4 += fmt.Sprintf("  %s⎇%s %s%s%s", GtopGray, Reset, GtopCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line4 += fmt.Sprintf(" %s+%d%s", GtopBrightGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line4 += fmt.Sprintf(" %s*%d%s", GtopYellow, data.GitDirty, Reset)
		}
	}
	sb.WriteString(gtopPadLine(line4, width, GtopDark+"│"+Reset))

	// Stats in clean columns
	line5 := fmt.Sprintf("%s│%s %sTOKENS%s %s%s%s  %sMSGS%s %s%d%s  %sTIME%s %s%s%s  %sHIT%s %s%d%%%s",
		GtopDark, Reset,
		GtopGray, Reset, GtopWhite, FormatTokens(data.TokenCount), Reset,
		GtopGray, Reset, GtopWhite, data.MessageCount, Reset,
		GtopGray, Reset, GtopWhite, data.SessionTime, Reset,
		GtopGray, Reset, GtopCyan, data.CacheHitRate, Reset)
	sb.WriteString(gtopPadLine(line5, width, GtopDark+"│"+Reset))

	// Cost row
	line6 := fmt.Sprintf("%s│%s %sSESSION%s %s%s%s  %sRATE%s %s%s/h%s  %sDAY%s %s%s%s",
		GtopDark, Reset,
		GtopGray, Reset, GtopGreen, FormatCostShort(data.SessionCost), Reset,
		GtopGray, Reset, GtopRed, FormatCostShort(data.BurnRate), Reset,
		GtopGray, Reset, GtopYellow, FormatCostShort(data.DayCost), Reset)
	sb.WriteString(gtopPadLine(line6, width, GtopDark+"│"+Reset))

	// Bottom border
	sb.WriteString(GtopDark + "└" + strings.Repeat("─", width-2) + "┘" + Reset + "\n")

	return sb.String()
}

func gtopPadLine(line string, targetWidth int, suffix string) string {
	visible := gtopDisplayWidth(line)
	suffixLen := gtopDisplayWidth(suffix)
	padding := targetWidth - visible - suffixLen
	if padding < 0 {
		padding = 0
	}
	return line + strings.Repeat(" ", padding) + suffix + "\n"
}

// gtopDisplayWidth calculates display width accounting for ANSI codes and emoji width
func gtopDisplayWidth(s string) int {
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
			if gtopIsWideChar(r) {
				width += 2
			} else {
				width += 1
			}
		}
	}
	return width
}

// gtopIsWideChar checks if a rune is a wide character
func gtopIsWideChar(r rune) bool {
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
	// Box Drawing and Block Elements are 1 width
	if r >= 0x2500 && r <= 0x259F {
		return false
	}
	// CJK characters
	if unicode.Is(unicode.Han, r) {
		return true
	}
	// Model icons
	switch r {
	case '💛', '💙', '💚', '⎇':
		return true
	}
	return false
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
