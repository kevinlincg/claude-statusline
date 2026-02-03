package themes

import (
	"fmt"
	"strings"
)

// MatrixTheme 矩陣駭客風格
type MatrixTheme struct{}

func init() {
	RegisterTheme(&MatrixTheme{})
}

func (t *MatrixTheme) Name() string {
	return "matrix"
}

func (t *MatrixTheme) Description() string {
	return "矩陣駭客風：綠色主調，終端機命令風格"
}

const (
	MatrixGreen     = "\033[38;2;0;255;0m"
	MatrixDarkGreen = "\033[38;2;0;200;0m"
	MatrixDim       = "\033[38;2;0;100;0m"
	MatrixGray      = "\033[38;2;80;80;80m"
)

func (t *MatrixTheme) Render(data StatusData) string {
	var sb strings.Builder

	const width = 89

	// 頂部邊框
	sb.WriteString(MatrixGreen)
	sb.WriteString("░▒▓")
	sb.WriteString(Bold)
	sb.WriteString(strings.Repeat("█", width-6))
	sb.WriteString(Reset)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓▒░")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 標題行
	titleLine := t.formatTitleLine(data, width)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString(titleLine)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 路徑行
	pathLine := t.formatPathLine(data, width)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString(pathLine)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 分隔線
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(strings.Repeat("▄", width))
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 成本行
	costLine := t.formatCostLine(data, width)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString(costLine)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 資料行
	dataLine := t.formatDataLine(data, width)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString(dataLine)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 分隔線
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(strings.Repeat("▀", width))
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// API 行
	apiLine := t.formatAPILine(data, width)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString(apiLine)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 底部邊框
	sb.WriteString(MatrixGreen)
	sb.WriteString("░▒▓")
	sb.WriteString(Bold)
	sb.WriteString(strings.Repeat("█", width-6))
	sb.WriteString(Reset)
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓▒░")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	return sb.String()
}

func (t *MatrixTheme) formatTitleLine(data StatusData, width int) string {
	_, modelIcon := GetModelConfig(data.ModelType)

	title := fmt.Sprintf("  %s%s⟦ CLAUDE CODE ⟧%s", Bold, MatrixGreen, Reset)
	version := fmt.Sprintf("  %s%s%s", MatrixDarkGreen, data.Version, Reset)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s⬆%s", ColorYellow, Reset)
	}

	filler := fmt.Sprintf("  %s════════════════════════════%s", MatrixGray, Reset)
	model := fmt.Sprintf("  %s⟦ %s %s ⟧%s", ColorGold, modelIcon, data.ModelName, Reset)

	content := title + version + update + filler + model
	return PadRight(content, width)
}

func (t *MatrixTheme) formatPathLine(data StatusData, width int) string {
	prompt := fmt.Sprintf("  %s├──▶%s", MatrixDarkGreen, Reset)
	path := fmt.Sprintf(" %s%s%s", ColorYellow, data.ProjectPath, Reset)

	git := ""
	if data.GitBranch != "" {
		git = fmt.Sprintf("  %s├──▶%s %s⚡ %s%s", MatrixDarkGreen, Reset, ColorCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			git += fmt.Sprintf(" %s+%d%s", ColorGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			git += fmt.Sprintf(" %s~%d%s", ColorOrange, data.GitDirty, Reset)
		}
	}

	time := fmt.Sprintf("  %s├──▶%s %s⏱ %s%s", MatrixDarkGreen, Reset, MatrixGreen, data.SessionTime, Reset)

	content := prompt + path + git + time
	return PadRight(content, width)
}

func (t *MatrixTheme) formatCostLine(data StatusData, width int) string {
	prompt := fmt.Sprintf("  %s$>%s", MatrixDarkGreen, Reset)
	content := fmt.Sprintf(" %sCOST%s ses:%s%s%s  day:%s%s%s  mon:%s%s%s  wk:%s%s%s  %s│%s  rate:%s%s/h%s  cache:%s%d%%%s",
		Bold, Reset,
		MatrixGreen, FormatCost(data.SessionCost), Reset,
		ColorYellow, FormatCost(data.DayCost), Reset,
		ColorPurple, FormatCost(data.MonthCost), Reset,
		ColorBlue, FormatCost(data.WeekCost), Reset,
		MatrixGray, Reset,
		ColorRed, FormatCost(data.BurnRate), Reset,
		ColorGreen, data.CacheHitRate, Reset)

	return PadRight(prompt+content, width)
}

func (t *MatrixTheme) formatDataLine(data StatusData, width int) string {
	prompt := fmt.Sprintf("  %s$>%s", MatrixDarkGreen, Reset)

	bar := GenerateGlowBar(data.ContextPercent, 20, MatrixGreen, "\033[48;2;0;40;0m")
	pctColor := GetContextColor(data.ContextPercent)

	content := fmt.Sprintf(" %sDATA%s tok:%s%s%s  msg:%s%d%s  %s│%s  ctx:%s %s%d%%%s %s",
		Bold, Reset,
		ColorPurple, FormatTokens(data.TokenCount), Reset,
		ColorCyan, data.MessageCount, Reset,
		MatrixGray, Reset,
		bar, pctColor, data.ContextPercent, Reset,
		FormatNumber(data.ContextUsed))

	return PadRight(prompt+content, width)
}

func (t *MatrixTheme) formatAPILine(data StatusData, width int) string {
	prompt := fmt.Sprintf("  %s$>%s", MatrixDarkGreen, Reset)

	color5, _ := GetBarColor(data.API5hrPercent)
	bar5 := GenerateGlowBar(data.API5hrPercent, 10, color5, "\033[48;2;0;40;0m")

	color7, _ := GetBarColor(data.API7dayPercent)
	bar7 := GenerateGlowBar(data.API7dayPercent, 10, color7, "\033[48;2;0;40;0m")

	content := fmt.Sprintf(" %sAPI%s  5hr:%s %s%d%%%s (%s)  %s│%s  7dy:%s %s%d%%%s (%s)",
		Bold, Reset,
		bar5, color5, data.API5hrPercent, Reset, data.API5hrTimeLeft,
		MatrixGray, Reset,
		bar7, color7, data.API7dayPercent, Reset, data.API7dayTimeLeft)

	return PadRight(prompt+content, width)
}
