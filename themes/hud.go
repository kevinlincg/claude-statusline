package themes

import (
	"fmt"
	"strings"
)

// HUDTheme 科幻 HUD 風格
type HUDTheme struct{}

func init() {
	RegisterTheme(&HUDTheme{})
}

func (t *HUDTheme) Name() string {
	return "hud"
}

func (t *HUDTheme) Description() string {
	return "科幻 HUD：未來感介面，角括號標籤"
}

const (
	HUDCyan      = "\033[38;2;0;200;200m"
	HUDBrightCyan = "\033[38;2;0;255;255m"
	HUDGreen     = "\033[38;2;0;255;180m"
	HUDYellow    = "\033[38;2;255;220;0m"
)

func (t *HUDTheme) Render(data StatusData) string {
	var sb strings.Builder

	// 第一行：模型 + 版本 + 路徑 + Git + Cost
	line1 := t.formatLine1(data)
	sb.WriteString(line1)
	sb.WriteString("\n")

	// 第二行：所有光棒
	line2 := t.formatLine2(data)
	sb.WriteString(line2)
	sb.WriteString("\n")

	return sb.String()
}

func (t *HUDTheme) formatLine1(data StatusData) string {
	modelColor, modelIcon := GetModelConfig(data.ModelType)

	// ⟨◆Opus4.5⟩
	model := fmt.Sprintf(" %s⟨%s%s%s%s%s⟩%s",
		HUDCyan,
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		HUDCyan+Reset)

	// v1.75↑
	version := fmt.Sprintf(" %s%s%s", HUDGreen, data.Version, Reset)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf("%s↑%s", ColorRed, Reset)
	}

	// ⟨~/proj⟩
	path := fmt.Sprintf(" %s⟨%s%s%s%s⟩%s",
		HUDCyan,
		ColorYellow, shortenPath(data.ProjectPath), Reset,
		HUDCyan, Reset)

	// ⟨⚡main+3~5⟩
	git := ""
	if data.GitBranch != "" {
		gitContent := fmt.Sprintf("%s⚡%s%s", ColorCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitContent += fmt.Sprintf("%s+%d%s", ColorGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitContent += fmt.Sprintf("%s~%d%s", ColorOrange, data.GitDirty, Reset)
		}
		git = fmt.Sprintf(" %s⟨%s%s⟩%s", HUDCyan, gitContent, HUDCyan, Reset)
	}

	// ⟨$.12/$3/$68⟩
	cost := fmt.Sprintf(" %s⟨%s%s%s/%s%s%s/%s%s%s%s⟩%s",
		HUDCyan,
		ColorGreen, FormatCostShort(data.SessionCost), Reset,
		ColorYellow, FormatCostShort(data.DayCost), Reset,
		ColorPurple, FormatCostShort(data.MonthCost), Reset,
		HUDCyan, Reset)

	// $5/h
	rate := fmt.Sprintf(" %s%s/h%s", ColorRed, FormatCostShort(data.BurnRate), Reset)

	return model + version + update + path + git + cost + rate
}

func (t *HUDTheme) formatLine2(data StatusData) string {
	// CTX bar
	ctxBar := t.generateHUDBar(data.ContextPercent, 20, HUDBrightCyan, BgCyanGlow)
	ctxColor := GetContextColor(data.ContextPercent)
	ctx := fmt.Sprintf("  %sCTX%s%s%s%d%s",
		ColorDim, Reset,
		ctxBar,
		ctxColor, data.ContextPercent, Reset)

	// 5H bar
	color5, bgColor5 := GetBarColor(data.API5hrPercent)
	bar5 := t.generateHUDBar(data.API5hrPercent, 10, color5, bgColor5)
	api5 := fmt.Sprintf("  %s5H%s%s%s%d%s",
		ColorDim, Reset,
		bar5,
		color5, data.API5hrPercent, Reset)

	// 7D bar
	color7, bgColor7 := GetBarColor(data.API7dayPercent)
	bar7 := t.generateHUDBar(data.API7dayPercent, 10, color7, bgColor7)
	api7 := fmt.Sprintf("  %s7D%s%s%s%d%s",
		ColorDim, Reset,
		bar7,
		color7, data.API7dayPercent, Reset)

	return ctx + api5 + api7
}

func (t *HUDTheme) generateHUDBar(percent, width int, color, bgColor string) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(bgColor)
		bar.WriteString(Bold)
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("━", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString("\033[38;2;30;50;50m")
		bar.WriteString(strings.Repeat("╌", empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}
