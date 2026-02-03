package themes

import (
	"fmt"
	"strings"
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

	// 第一行：路徑 + Git + 模型 + 版本
	sb.WriteString(" ")
	sb.WriteString(t.formatHeader(data))
	sb.WriteString("\n")

	// 第二行：Session | Context bar
	sb.WriteString(fmt.Sprintf(" %s├─%s %s  %s│%s  %s\n",
		ColorTreeDim, Reset,
		t.formatSessionLine(data),
		ColorFrame, Reset,
		t.formatContextBar(data)))

	// 第三行：Cost | 5hr bar
	sb.WriteString(fmt.Sprintf(" %s├─%s %s  %s│%s  %s\n",
		ColorTreeDim, Reset,
		t.formatCostLine(data),
		ColorFrame, Reset,
		t.format5hrBar(data)))

	// 第四行：| 7day bar
	sb.WriteString(fmt.Sprintf(" %s└─%s %s  %s│%s  %s\n",
		ColorTreeDim, Reset,
		t.formatCostLine2(data),
		ColorFrame, Reset,
		t.format7dayBar(data)))

	return sb.String()
}

func (t *MinimalTheme) formatHeader(data StatusData) string {
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
		update = fmt.Sprintf(" %s%s⬆%s", Bold, ColorNeonOrange, Reset)
	}

	return path + git + strings.Repeat(" ", 30) + model + version + update
}

func (t *MinimalTheme) formatSessionLine(data StatusData) string {
	return fmt.Sprintf("%sSession%s   %s%s%s tok   %s%d%s msg   %s%s%s   %s%d%%%s hit",
		ColorLabel, Reset,
		ColorPurple, FormatTokens(data.TokenCount), Reset,
		ColorCyan, data.MessageCount, Reset,
		ColorSilver, data.SessionTime, Reset,
		ColorGreen, data.CacheHitRate, Reset)
}

func (t *MinimalTheme) formatCostLine(data StatusData) string {
	return fmt.Sprintf("%sCost%s      ses %s%s%s  day %s%s%s  %s%s/h%s",
		ColorLabel, Reset,
		ColorGreen, FormatCost(data.SessionCost), Reset,
		ColorYellow, FormatCost(data.DayCost), Reset,
		ColorRed, FormatCost(data.BurnRate), Reset)
}

func (t *MinimalTheme) formatCostLine2(data StatusData) string {
	return fmt.Sprintf("          mon %s%s%s  wk %s%s%s",
		ColorPurple, FormatCost(data.MonthCost), Reset,
		ColorBlue, FormatCost(data.WeekCost), Reset)
}

func (t *MinimalTheme) formatContextBar(data StatusData) string {
	color, bgColor := GetBarColor(data.ContextPercent)
	bar := GenerateGlowBar(data.ContextPercent, 20, color, bgColor)
	pctColor := GetContextColor(data.ContextPercent)

	return fmt.Sprintf("%sCtx%s  %s %s%d%%%s %s%s%s",
		ColorLabelDim, Reset,
		bar,
		pctColor, data.ContextPercent, Reset,
		ColorDim, FormatNumber(data.ContextUsed), Reset)
}

func (t *MinimalTheme) format5hrBar(data StatusData) string {
	color, bgColor := GetBarColor(data.API5hrPercent)
	bar := GenerateGlowBar(data.API5hrPercent, 20, color, bgColor)

	return fmt.Sprintf("%s5hr%s  %s %s%d%%%s %s%s%s",
		ColorLabelDim, Reset,
		bar,
		color, data.API5hrPercent, Reset,
		ColorDim, data.API5hrTimeLeft, Reset)
}

func (t *MinimalTheme) format7dayBar(data StatusData) string {
	color, bgColor := GetBarColor(data.API7dayPercent)
	bar := GenerateGlowBar(data.API7dayPercent, 20, color, bgColor)

	return fmt.Sprintf("%s7dy%s  %s %s%d%%%s %s%s%s",
		ColorLabelDim, Reset,
		bar,
		color, data.API7dayPercent, Reset,
		ColorDim, data.API7dayTimeLeft, Reset)
}
