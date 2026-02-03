package themes

import (
	"fmt"
	"strings"
)

// CyberpunkTheme 賽博朋克霓虹風格
type CyberpunkTheme struct{}

func init() {
	RegisterTheme(&CyberpunkTheme{})
}

func (t *CyberpunkTheme) Name() string {
	return "cyberpunk"
}

func (t *CyberpunkTheme) Description() string {
	return "賽博朋克霓虹：青色/洋紅雙色框線，霓虹發光效果"
}

const (
	CyberCyan    = "\033[38;2;0;255;255m"
	CyberMagenta = "\033[38;2;255;0;255m"
	CyberPink    = "\033[38;2;255;100;255m"
	CyberGreen   = "\033[38;2;0;255;136m"
	CyberYellow  = "\033[38;2;255;255;0m"
	CyberOrange  = "\033[38;2;255;150;0m"
)

func (t *CyberpunkTheme) Render(data StatusData) string {
	var sb strings.Builder

	const width = 90

	// 頂部框線
	sb.WriteString(CyberCyan)
	sb.WriteString("╔")
	sb.WriteString(strings.Repeat("═", width))
	sb.WriteString("╗")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 標題行：CLAUDE CODE + 版本 + 模型
	titleLine := t.formatTitleLine(data, width)
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString(titleLine)
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 分隔線
	sb.WriteString(CyberCyan)
	sb.WriteString("╠")
	sb.WriteString(strings.Repeat("═", width))
	sb.WriteString("╣")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 路徑 + Git + Session 資訊
	infoLine := t.formatInfoLine(data, width)
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString(infoLine)
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 分隔線
	sb.WriteString(CyberCyan)
	sb.WriteString("╠")
	sb.WriteString(strings.Repeat("═", width))
	sb.WriteString("╣")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 成本行
	costLine := t.formatCostLine(data, width)
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString(costLine)
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 分隔線
	sb.WriteString(CyberCyan)
	sb.WriteString("╠")
	sb.WriteString(strings.Repeat("═", width))
	sb.WriteString("╣")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// Context 行
	ctxLine := t.formatContextLine(data, width)
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString(ctxLine)
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// API 行
	apiLine := t.formatAPILine(data, width)
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString(apiLine)
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 底部框線
	sb.WriteString(CyberCyan)
	sb.WriteString("╚")
	sb.WriteString(strings.Repeat("═", width))
	sb.WriteString("╝")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	return sb.String()
}

func (t *CyberpunkTheme) formatTitleLine(data StatusData, width int) string {
	modelColor, modelIcon := GetModelConfig(data.ModelType)

	title := fmt.Sprintf("  %s%sCLAUDE CODE%s", Bold, CyberMagenta, Reset)
	version := fmt.Sprintf("  %s%s%s", CyberGreen, data.Version, Reset)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s%s⬆ NEW%s", Bold, CyberOrange, Reset)
	}
	model := fmt.Sprintf("  %s│%s  %s%s %s%s", ColorDim, Reset, modelColor, modelIcon, data.ModelName, Reset)

	content := title + version + update + model
	return PadRight(content, width)
}

func (t *CyberpunkTheme) formatInfoLine(data StatusData, width int) string {
	path := fmt.Sprintf(" %s📂 %s%s", ColorYellow, data.ProjectPath, Reset)

	git := ""
	if data.GitBranch != "" {
		git = fmt.Sprintf("  %s⚡ %s%s", CyberCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			git += fmt.Sprintf(" %s+%d%s", ColorGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			git += fmt.Sprintf(" %s~%d%s", ColorOrange, data.GitDirty, Reset)
		}
	}

	stats := fmt.Sprintf("  %s│%s  %s⏱ %s%s  %s│%s  %s💬 %d msg%s  %s│%s  %s📊 %s tok%s",
		ColorDim, Reset,
		CyberMagenta, data.SessionTime, Reset,
		ColorDim, Reset,
		CyberCyan, data.MessageCount, Reset,
		ColorDim, Reset,
		ColorPurple, FormatTokens(data.TokenCount), Reset)

	return PadRight(path+git+stats, width)
}

func (t *CyberpunkTheme) formatCostLine(data StatusData, width int) string {
	content := fmt.Sprintf(" %s💰 COST%s  ses %s%s%s  day %s%s%s  mon %s%s%s  %s│%s  wk %s%s%s  avg %s%s/h%s  hit %s%d%%%s",
		Bold, Reset,
		CyberGreen, FormatCost(data.SessionCost), Reset,
		ColorYellow, FormatCost(data.DayCost), Reset,
		ColorPurple, FormatCost(data.MonthCost), Reset,
		ColorDim, Reset,
		ColorBlue, FormatCost(data.WeekCost), Reset,
		ColorRed, FormatCost(data.BurnRate), Reset,
		ColorGreen, data.CacheHitRate, Reset)

	return PadRight(content, width)
}

func (t *CyberpunkTheme) formatContextLine(data StatusData, width int) string {
	color, bgColor := GetBarColor(data.ContextPercent)
	bar := GenerateGlowBar(data.ContextPercent, 20, color, bgColor)
	pctColor := GetContextColor(data.ContextPercent)

	color5, bgColor5 := GetBarColor(data.API5hrPercent)
	bar5 := GenerateGlowBar(data.API5hrPercent, 10, color5, bgColor5)

	content := fmt.Sprintf(" %s📈 CTX%s   %s %s%4d%%%s  %s%s  %s│%s  %s5hr%s %s %s%d%%%s (%s)",
		Bold, Reset,
		bar, pctColor, data.ContextPercent, Reset,
		ColorDim, FormatNumber(data.ContextUsed),
		ColorDim, Reset,
		Bold, Reset,
		bar5, color5, data.API5hrPercent, Reset,
		data.API5hrTimeLeft)

	return PadRight(content, width)
}

func (t *CyberpunkTheme) formatAPILine(data StatusData, width int) string {
	color7, bgColor7 := GetBarColor(data.API7dayPercent)
	bar7 := GenerateGlowBar(data.API7dayPercent, 20, color7, bgColor7)

	content := fmt.Sprintf("          %s└─ context window ─┘%s            %s│%s  %s7dy%s %s %s%d%%%s (%s)",
		ColorDim, Reset,
		ColorDim, Reset,
		Bold, Reset,
		bar7, color7, data.API7dayPercent, Reset,
		data.API7dayTimeLeft)

	return PadRight(content, width)
}
