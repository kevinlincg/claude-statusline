package themes

import (
	"fmt"
	"strings"
)

// CompactTheme C 版：精簡對稱
type CompactTheme struct{}

func init() {
	RegisterTheme(&CompactTheme{})
}

func (t *CompactTheme) Name() string {
	return "compact"
}

func (t *CompactTheme) Description() string {
	return "精簡對稱：只有上下橫線，兩行內容"
}

func (t *CompactTheme) Render(data StatusData) string {
	var sb strings.Builder

	const width = 97

	// 頂部橫線
	sb.WriteString(" ")
	sb.WriteString(ColorFrameDim)
	sb.WriteString(strings.Repeat("─", width))
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第一行：所有文字資訊
	line1 := t.formatLine1(data)
	sb.WriteString(ColorFrame)
	sb.WriteString(" │")
	sb.WriteString(Reset)
	sb.WriteString(line1)
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第二行：所有光棒
	line2 := t.formatLine2(data)
	sb.WriteString(ColorFrame)
	sb.WriteString(" │")
	sb.WriteString(Reset)
	sb.WriteString(line2)
	sb.WriteString(ColorFrame)
	sb.WriteString("│")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 底部橫線
	sb.WriteString(" ")
	sb.WriteString(ColorFrameDim)
	sb.WriteString(strings.Repeat("─", width))
	sb.WriteString(Reset)
	sb.WriteString("\n")

	return sb.String()
}

func (t *CompactTheme) formatLine1(data StatusData) string {
	modelColor, modelIcon := GetModelConfig(data.ModelType)

	// 模型 + 版本
	model := fmt.Sprintf(" %s%s%s%s%s", modelColor, Bold, modelIcon, data.ModelName, Reset)
	version := fmt.Sprintf(" %s%s%s", ColorNeonGreen, data.Version, Reset)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf("%s⬆%s", ColorNeonOrange, Reset)
	}

	// 分隔
	sep := fmt.Sprintf(" %s│%s ", ColorFrame, Reset)

	// 路徑 + Git
	path := fmt.Sprintf("%s%s%s", ColorYellow, shortenPath(data.ProjectPath), Reset)
	git := ""
	if data.GitBranch != "" {
		git = fmt.Sprintf(" %s⚡%s%s", ColorCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			git += fmt.Sprintf("%s+%d%s", ColorGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			git += fmt.Sprintf("%s~%d%s", ColorOrange, data.GitDirty, Reset)
		}
	}

	// Token + Msg + Time
	stats := fmt.Sprintf("%s%s%stok %s%d%sm %s%s%s",
		ColorPurple, FormatTokens(data.TokenCount), Reset,
		ColorCyan, data.MessageCount, Reset,
		ColorSilver, data.SessionTime, Reset)

	// Cost
	cost := fmt.Sprintf("%s%s%ss %s%s%sd %s%s%sm %s%s%sw %s%s/h%s %s%d%%h%s",
		ColorGreen, FormatCostShort(data.SessionCost), Reset,
		ColorYellow, FormatCostShort(data.DayCost), Reset,
		ColorPurple, FormatCostShort(data.MonthCost), Reset,
		ColorBlue, FormatCostShort(data.WeekCost), Reset,
		ColorRed, FormatCostShort(data.BurnRate), Reset,
		ColorGreen, data.CacheHitRate, Reset)

	return model + version + update + sep + path + git + sep + stats + sep + cost
}

func (t *CompactTheme) formatLine2(data StatusData) string {
	sep := fmt.Sprintf(" %s│%s ", ColorFrame, Reset)

	// Context bar
	color1, bgColor1 := GetBarColor(data.ContextPercent)
	bar1 := GenerateGlowBar(data.ContextPercent, 20, color1, bgColor1)
	ctx := fmt.Sprintf(" %sCtx%s%s %s%d%%%s%s%s%s",
		ColorLabelDim, Reset,
		bar1,
		color1, data.ContextPercent, Reset,
		ColorDim, FormatNumber(data.ContextUsed), Reset)

	// 5hr bar
	color5, bgColor5 := GetBarColor(data.API5hrPercent)
	bar5 := GenerateGlowBar(data.API5hrPercent, 10, color5, bgColor5)
	api5 := fmt.Sprintf("%s5h%s%s %s%d%%%s%s%s%s",
		ColorLabelDim, Reset,
		bar5,
		color5, data.API5hrPercent, Reset,
		ColorDim, data.API5hrTimeLeft, Reset)

	// 7day bar
	color7, bgColor7 := GetBarColor(data.API7dayPercent)
	bar7 := GenerateGlowBar(data.API7dayPercent, 10, color7, bgColor7)
	api7 := fmt.Sprintf("%s7d%s%s %s%d%%%s%s%s%s",
		ColorLabelDim, Reset,
		bar7,
		color7, data.API7dayPercent, Reset,
		ColorDim, data.API7dayTimeLeft, Reset)

	return ctx + sep + api5 + sep + api7
}

func shortenPath(path string) string {
	if len(path) > 20 {
		parts := strings.Split(path, "/")
		if len(parts) > 2 {
			return "~/" + parts[len(parts)-1]
		}
	}
	return path
}
