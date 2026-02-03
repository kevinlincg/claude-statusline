package themes

import (
	"fmt"
	"strings"
	"unicode"
)

// MinimalTheme A 版：簡潔樹狀
type MinimalTheme struct{}

func init() {
	RegisterTheme(&MinimalTheme{})
}

func (t *MinimalTheme) Name() string {
	return "minimal"
}

func (t *MinimalTheme) Description() string {
	return "簡潔樹狀：無外框，樹狀結構顯示資訊"
}

func (t *MinimalTheme) Render(data StatusData) string {
	var sb strings.Builder
	width := 80

	// 第一行：路徑 + Git + 模型 + 版本
	line1 := t.formatHeader(data)
	sb.WriteString(minimalPadLine(" "+line1, width, ""))

	// 第二行：Session | Context bar
	leftSide2 := fmt.Sprintf(" %s├─%s %s",
		ColorTreeDim, Reset,
		t.formatSessionLine(data))
	rightSide2 := t.formatContextBar(data)
	sb.WriteString(minimalTwoColumn(leftSide2, rightSide2, width))

	// 第三行：Cost | 5hr bar
	leftSide3 := fmt.Sprintf(" %s├─%s %s",
		ColorTreeDim, Reset,
		t.formatCostLine(data))
	rightSide3 := t.format5hrBar(data)
	sb.WriteString(minimalTwoColumn(leftSide3, rightSide3, width))

	// 第四行：| 7day bar
	leftSide4 := fmt.Sprintf(" %s└─%s %s",
		ColorTreeDim, Reset,
		t.formatCostLine2(data))
	rightSide4 := t.format7dayBar(data)
	sb.WriteString(minimalTwoColumn(leftSide4, rightSide4, width))

	return sb.String()
}

func minimalTwoColumn(left, right string, totalWidth int) string {
	leftWidth := minimalDisplayWidth(left)
	rightWidth := minimalDisplayWidth(right)
	sep := fmt.Sprintf("  %s│%s  ", ColorFrame, Reset)
	sepWidth := minimalDisplayWidth(sep)

	padding := totalWidth - leftWidth - sepWidth - rightWidth
	if padding < 0 {
		padding = 0
	}
	return left + strings.Repeat(" ", padding) + sep + right + "\n"
}

func minimalPadLine(line string, targetWidth int, suffix string) string {
	visible := minimalDisplayWidth(line)
	suffixLen := minimalDisplayWidth(suffix)
	padding := targetWidth - visible - suffixLen
	if padding < 0 {
		padding = 0
	}
	return line + strings.Repeat(" ", padding) + suffix + "\n"
}

// minimalDisplayWidth calculates display width accounting for:
// - ANSI escape codes (0 width)
// - Emojis and wide characters (2 width)
// - Regular ASCII and box drawing (1 width)
func minimalDisplayWidth(s string) int {
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
			if minimalIsWideChar(r) {
				width += 2
			} else {
				width += 1
			}
		}
	}
	return width
}

// minimalIsWideChar checks if a rune is a wide character (emoji or CJK)
func minimalIsWideChar(r rune) bool {
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
	// Specific emojis
	switch r {
	case '📂', '⚡', '💛', '💙', '💚', '⬆':
		return true
	}
	return false
}

func (t *MinimalTheme) formatHeader(data StatusData) string {
	modelColor, modelIcon := GetModelConfig(data.ModelType)

	path := fmt.Sprintf("%s📂%s %s", ColorYellow, Reset, ShortenPath(data.ProjectPath, 25))

	git := ""
	if data.GitBranch != "" {
		git = fmt.Sprintf("  %s⚡%s %s", ColorCyan, Reset, data.GitBranch)
		if data.GitStaged > 0 {
			git += fmt.Sprintf(" %s+%d%s", ColorGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			git += fmt.Sprintf(" %s~%d%s", ColorOrange, data.GitDirty, Reset)
		}
	}

	model := fmt.Sprintf("%s[%s%s%s %s%s]%s", ColorFrame, modelColor, modelIcon, Reset, data.ModelName, ColorFrame, Reset)
	version := fmt.Sprintf(" %s%s%s", ColorNeonGreen, data.Version, Reset)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s⬆%s", ColorNeonOrange, Reset)
	}

	// Calculate spacing to right-align the model info
	leftPart := path + git
	rightPart := model + version + update
	leftWidth := minimalDisplayWidth(leftPart)
	rightWidth := minimalDisplayWidth(rightPart)
	spacing := 78 - leftWidth - rightWidth // 78 = 80 - 2 for margins
	if spacing < 2 {
		spacing = 2
	}

	return leftPart + strings.Repeat(" ", spacing) + rightPart
}

func (t *MinimalTheme) formatSessionLine(data StatusData) string {
	return fmt.Sprintf("%sSession%s  %s%s%s tok  %s%d%s msg  %s%s%s  %s%d%%%s hit",
		ColorLabel, Reset,
		ColorPurple, FormatTokens(data.TokenCount), Reset,
		ColorCyan, data.MessageCount, Reset,
		ColorSilver, data.SessionTime, Reset,
		ColorGreen, data.CacheHitRate, Reset)
}

func (t *MinimalTheme) formatCostLine(data StatusData) string {
	return fmt.Sprintf("%sCost%s     ses %s%s%s  day %s%s%s  %s%s/h%s",
		ColorLabel, Reset,
		ColorGreen, FormatCostShort(data.SessionCost), Reset,
		ColorYellow, FormatCostShort(data.DayCost), Reset,
		ColorRed, FormatCostShort(data.BurnRate), Reset)
}

func (t *MinimalTheme) formatCostLine2(data StatusData) string {
	return fmt.Sprintf("         mon %s%s%s  wk %s%s%s",
		ColorPurple, FormatCostShort(data.MonthCost), Reset,
		ColorBlue, FormatCostShort(data.WeekCost), Reset)
}

func (t *MinimalTheme) formatContextBar(data StatusData) string {
	color, bgColor := GetBarColor(data.ContextPercent)
	bar := GenerateGlowBar(data.ContextPercent, 18, color, bgColor)
	pctColor := GetContextColor(data.ContextPercent)

	return fmt.Sprintf("%sCtx%s %s %s%d%%%s %s%s%s",
		ColorLabelDim, Reset,
		bar,
		pctColor, data.ContextPercent, Reset,
		ColorDim, FormatNumber(data.ContextUsed), Reset)
}

func (t *MinimalTheme) format5hrBar(data StatusData) string {
	color, bgColor := GetBarColor(data.API5hrPercent)
	bar := GenerateGlowBar(data.API5hrPercent, 18, color, bgColor)

	return fmt.Sprintf("%s5hr%s %s %s%d%%%s %s%s%s",
		ColorLabelDim, Reset,
		bar,
		color, data.API5hrPercent, Reset,
		ColorDim, data.API5hrTimeLeft, Reset)
}

func (t *MinimalTheme) format7dayBar(data StatusData) string {
	color, bgColor := GetBarColor(data.API7dayPercent)
	bar := GenerateGlowBar(data.API7dayPercent, 18, color, bgColor)

	return fmt.Sprintf("%s7dy%s %s %s%d%%%s %s%s%s",
		ColorLabelDim, Reset,
		bar,
		color, data.API7dayPercent, Reset,
		ColorDim, data.API7dayTimeLeft, Reset)
}
