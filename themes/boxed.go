package themes

import (
	"fmt"
	"strings"
)

// BoxedTheme B 版：框線整齊
type BoxedTheme struct{}

func init() {
	RegisterTheme(&BoxedTheme{})
}

func (t *BoxedTheme) Name() string {
	return "boxed"
}

func (t *BoxedTheme) Description() string {
	return "框線整齊：完整框線包圍，左右對稱分區"
}

func (t *BoxedTheme) Render(data StatusData) string {
	var sb strings.Builder

	const leftWidth = 49
	const rightWidth = 43
	const fullWidth = leftWidth + rightWidth + 1

	// 頂部框線
	sb.WriteString(ColorFrame)
	sb.WriteString(" ┌")
	sb.WriteString(strings.Repeat("─", fullWidth))
	sb.WriteString("┐")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 標題行：路徑 + Git + 模型
	headerContent := t.formatHeader(data)
	sb.WriteString(ColorFrame)
	sb.WriteString(" │")
	sb.WriteString(Reset)
	sb.WriteString(" ")
	sb.WriteString(PadRight(headerContent, fullWidth-1))
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 分隔線
	sb.WriteString(ColorFrame)
	sb.WriteString(" ├")
	sb.WriteString(strings.Repeat("─", leftWidth))
	sb.WriteString("┼")
	sb.WriteString(strings.Repeat("─", rightWidth))
	sb.WriteString("┤")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第一行：Session 資訊 | Context bar
	leftContent := t.formatSessionInfo(data)
	rightContent := t.formatContextBar(data)

	sb.WriteString(ColorFrame)
	sb.WriteString(" │")
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(PadRight(leftContent, leftWidth-2))
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(PadRight(rightContent, rightWidth-2))
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第二行：Cache | 5hr bar
	leftContent = t.formatCacheInfo(data)
	rightContent = t.format5hrBar(data)

	sb.WriteString(ColorFrame)
	sb.WriteString(" │")
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(PadRight(leftContent, leftWidth-2))
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(PadRight(rightContent, rightWidth-2))
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第三行：Cost 1 | 7day bar
	leftContent = t.formatCostInfo1(data)
	rightContent = t.format7dayBar(data)

	sb.WriteString(ColorFrame)
	sb.WriteString(" │")
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(PadRight(leftContent, leftWidth-2))
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(PadRight(rightContent, rightWidth-2))
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第四行：Cost 2 | 空白
	leftContent = t.formatCostInfo2(data)

	sb.WriteString(ColorFrame)
	sb.WriteString(" │")
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(PadRight(leftContent, leftWidth-2))
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString(strings.Repeat(" ", rightWidth))
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 底部框線
	sb.WriteString(ColorFrame)
	sb.WriteString(" └")
	sb.WriteString(strings.Repeat("─", leftWidth))
	sb.WriteString("┴")
	sb.WriteString(strings.Repeat("─", rightWidth))
	sb.WriteString("┘")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	return sb.String()
}

func (t *BoxedTheme) formatHeader(data StatusData) string {
	modelColor, modelIcon := GetModelConfig(data.ModelType)

	path := fmt.Sprintf("%s📂 %s%s", ColorYellow, data.ProjectPath, Reset)

	git := ""
	if data.GitBranch != "" {
		git = fmt.Sprintf("  %s⚡ %s%s", ColorCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			git += fmt.Sprintf(" %s+%d%s", ColorGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			git += fmt.Sprintf(" %s~%d%s", ColorOrange, data.GitDirty, Reset)
		}
	}

	model := fmt.Sprintf("%s[%s%s%s %s%s]%s", ColorFrame, modelColor, Bold, modelIcon, data.ModelName, Reset+ColorFrame, Reset)
	version := fmt.Sprintf("  %s%s%s", ColorNeonGreen, data.Version, Reset)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s%s⬆ UPDATE%s", Bold, ColorNeonOrange, Reset)
	}

	return path + git + "  " + model + version + update
}

func (t *BoxedTheme) formatSessionInfo(data StatusData) string {
	return fmt.Sprintf("%sSession%s   %s%s%s tok   %s%d%s msg   %s%s%s",
		ColorLabel, Reset,
		ColorPurple, FormatTokens(data.TokenCount), Reset,
		ColorCyan, data.MessageCount, Reset,
		ColorSilver, data.SessionTime, Reset)
}

func (t *BoxedTheme) formatCacheInfo(data StatusData) string {
	color := ColorGreen
	if data.CacheHitRate < 40 {
		color = ColorOrange
	} else if data.CacheHitRate < 70 {
		color = ColorYellow
	}
	return fmt.Sprintf("%sCache%s     %s%d%%%s hit",
		ColorLabel, Reset,
		color, data.CacheHitRate, Reset)
}

func (t *BoxedTheme) formatCostInfo1(data StatusData) string {
	return fmt.Sprintf("%sCost%s      ses %s%s%s   day %s%s%s   %s%s/h%s",
		ColorLabel, Reset,
		ColorGreen, FormatCost(data.SessionCost), Reset,
		ColorYellow, FormatCost(data.DayCost), Reset,
		ColorRed, FormatCost(data.BurnRate), Reset)
}

func (t *BoxedTheme) formatCostInfo2(data StatusData) string {
	return fmt.Sprintf("          mon %s%s%s   wk %s%s%s",
		ColorPurple, FormatCost(data.MonthCost), Reset,
		ColorBlue, FormatCost(data.WeekCost), Reset)
}

func (t *BoxedTheme) formatContextBar(data StatusData) string {
	color, bgColor := GetBarColor(data.ContextPercent)
	bar := GenerateGlowBar(data.ContextPercent, 20, color, bgColor)
	pctColor := GetContextColor(data.ContextPercent)

	return fmt.Sprintf("%sCtx%s  %s %s%d%%%s %s%s%s",
		ColorLabelDim, Reset,
		bar,
		pctColor, data.ContextPercent, Reset,
		ColorDim, FormatNumber(data.ContextUsed), Reset)
}

func (t *BoxedTheme) format5hrBar(data StatusData) string {
	color, bgColor := GetBarColor(data.API5hrPercent)
	bar := GenerateGlowBar(data.API5hrPercent, 20, color, bgColor)

	return fmt.Sprintf("%s5hr%s  %s %s%d%%%s %s%s%s",
		ColorLabelDim, Reset,
		bar,
		color, data.API5hrPercent, Reset,
		ColorDim, data.API5hrTimeLeft, Reset)
}

func (t *BoxedTheme) format7dayBar(data StatusData) string {
	color, bgColor := GetBarColor(data.API7dayPercent)
	bar := GenerateGlowBar(data.API7dayPercent, 20, color, bgColor)

	return fmt.Sprintf("%s7dy%s  %s %s%d%%%s %s%s%s",
		ColorLabelDim, Reset,
		bar,
		color, data.API7dayPercent, Reset,
		ColorDim, data.API7dayTimeLeft, Reset)
}
