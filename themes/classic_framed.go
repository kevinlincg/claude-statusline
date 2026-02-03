package themes

import (
	"fmt"
	"strings"
)

// ClassicFramedTheme D 版：經典樹狀+框線
type ClassicFramedTheme struct{}

func init() {
	RegisterTheme(&ClassicFramedTheme{})
}

func (t *ClassicFramedTheme) Name() string {
	return "classic_framed"
}

func (t *ClassicFramedTheme) Description() string {
	return "經典樹狀+框線：左側文字資訊，右側光棒垂直對齊"
}

func (t *ClassicFramedTheme) Render(data StatusData) string {
	var sb strings.Builder

	// 常數
	const leftWidth = 53
	const rightWidth = 45
	const fullWidth = leftWidth + rightWidth + 1 // +1 for middle border

	// 框線字元
	topLeft := "┌"
	topRight := "┐"
	topMid := "┬"
	midLeft := "├"
	midRight := "┤"
	botLeft := "└"
	botRight := "┘"
	botMid := "┴"
	hLine := "─"
	vLine := "│"

	// 第一行：路徑 + Git + 模型（橫跨整個寬度）
	headerContent := t.formatPathGitLine(data, fullWidth-2)
	sb.WriteString(ColorFrame)
	sb.WriteString(topLeft)
	sb.WriteString(strings.Repeat(hLine, fullWidth))
	sb.WriteString(topRight)
	sb.WriteString(Reset)
	sb.WriteString("\n")

	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString(" ")
	sb.WriteString(headerContent)
	sb.WriteString(" ")
	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 分隔線（開始左右分欄）
	sb.WriteString(ColorFrame)
	sb.WriteString(midLeft)
	sb.WriteString(strings.Repeat(hLine, leftWidth))
	sb.WriteString(topMid)
	sb.WriteString(strings.Repeat(hLine, rightWidth))
	sb.WriteString(midRight)
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第二行：Session | Context bar + Cache
	leftContent := t.formatSessionLine(data)
	rightContent := t.formatContextBar(data)

	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(ColorTreeDim)
	sb.WriteString("├─")
	sb.WriteString(Reset)
	sb.WriteString(" ")
	sb.WriteString(PadRight(leftContent, leftWidth-5))
	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(PadRight(rightContent, rightWidth-2))
	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第三行：Cost 1 | 5hr bar
	leftContent = t.formatCostLine1(data)
	rightContent = t.format5hrBar(data)

	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(ColorTreeDim)
	sb.WriteString("├─")
	sb.WriteString(Reset)
	sb.WriteString(" ")
	sb.WriteString(PadRight(leftContent, leftWidth-5))
	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(PadRight(rightContent, rightWidth-2))
	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第四行：Cost 2 | 7day bar
	leftContent = t.formatCostLine2(data)
	rightContent = t.format7dayBar(data)

	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(ColorTreeDim)
	sb.WriteString("└─")
	sb.WriteString(Reset)
	sb.WriteString(" ")
	sb.WriteString(PadRight(leftContent, leftWidth-5))
	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString("  ")
	sb.WriteString(PadRight(rightContent, rightWidth-2))
	sb.WriteString(ColorFrame)
	sb.WriteString(vLine)
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 底部框線
	sb.WriteString(ColorFrame)
	sb.WriteString(botLeft)
	sb.WriteString(strings.Repeat(hLine, leftWidth))
	sb.WriteString(botMid)
	sb.WriteString(strings.Repeat(hLine, rightWidth))
	sb.WriteString(botRight)
	sb.WriteString(Reset)
	sb.WriteString("\n")

	return sb.String()
}

func (t *ClassicFramedTheme) formatPathGitLine(data StatusData, width int) string {
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

	left := path + git

	// 右側：模型 + 版本
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s⬆%s", ColorNeonOrange, Reset)
	}
	model := fmt.Sprintf("%s%s%s%s %s%s%s%s", modelColor, modelIcon, data.ModelName, Reset, ColorNeonGreen, data.Version, Reset, update)

	// 計算填充
	leftVisible := VisibleWidth(left)
	modelVisible := VisibleWidth(model)
	padding := width - leftVisible - modelVisible
	if padding < 1 {
		padding = 1
	}

	return left + strings.Repeat(" ", padding) + model
}

func (t *ClassicFramedTheme) formatSessionLine(data StatusData) string {
	// 使用 │ 分隔符讓數據更清晰，對齊下方 cost 欄位
	return fmt.Sprintf("%s%5s%s tok %s│%s %s%5d%s msg %s│%s %s%7s%s",
		ColorPurple, FormatTokens(data.TokenCount), Reset,
		ColorTreeDim, Reset,
		ColorCyan, data.MessageCount, Reset,
		ColorTreeDim, Reset,
		ColorSilver, data.SessionTime, Reset)
}

func (t *ClassicFramedTheme) formatCostLine1(data StatusData) string {
	return fmt.Sprintf("%s%5s%s ses %s│%s %s%5s%s day %s│%s %s%7s%s",
		ColorGreen, FormatCost(data.SessionCost), Reset,
		ColorTreeDim, Reset,
		ColorYellow, FormatCost(data.DayCost), Reset,
		ColorTreeDim, Reset,
		ColorRed, FormatCost(data.BurnRate)+"/h", Reset)
}

func (t *ClassicFramedTheme) formatCostLine2(data StatusData) string {
	return fmt.Sprintf("%s%5s%s mon %s│%s %s%5s%s wk",
		ColorPurple, FormatCost(data.MonthCost), Reset,
		ColorTreeDim, Reset,
		ColorBlue, FormatCost(data.WeekCost), Reset)
}

func (t *ClassicFramedTheme) formatContextBar(data StatusData) string {
	color, bgColor := GetBarColor(data.ContextPercent)
	bar := GenerateGlowBar(data.ContextPercent, 20, color, bgColor)
	pctColor := GetContextColor(data.ContextPercent)

	return fmt.Sprintf("%sCtx%s %s %s%s%4d%%%s %s%3d%%hit%s",
		ColorLabelDim, Reset,
		bar,
		Bold, pctColor, data.ContextPercent, Reset,
		ColorDim, data.CacheHitRate, Reset)
}

func (t *ClassicFramedTheme) format5hrBar(data StatusData) string {
	color, bgColor := GetBarColor(data.API5hrPercent)
	bar := GenerateGlowBar(data.API5hrPercent, 20, color, bgColor)

	return fmt.Sprintf("%s5hr%s %s %s%s%4d%%%s %s%6s%s",
		ColorLabelDim, Reset,
		bar,
		Bold, color, data.API5hrPercent, Reset,
		ColorDim, data.API5hrTimeLeft, Reset)
}

func (t *ClassicFramedTheme) format7dayBar(data StatusData) string {
	color, bgColor := GetBarColor(data.API7dayPercent)
	bar := GenerateGlowBar(data.API7dayPercent, 20, color, bgColor)

	return fmt.Sprintf("%s7dy%s %s %s%s%4d%%%s %s%6s%s",
		ColorLabelDim, Reset,
		bar,
		Bold, color, data.API7dayPercent, Reset,
		ColorDim, data.API7dayTimeLeft, Reset)
}

