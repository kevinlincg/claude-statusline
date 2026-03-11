package themes

import (
	"fmt"
	"strings"
)

// TwolinePillsTheme two-line with rounded pill badges
type TwolinePillsTheme struct{}

func init() {
	RegisterTheme(&TwolinePillsTheme{})
}

func (t *TwolinePillsTheme) Name() string {
	return "twoline_pills"
}

func (t *TwolinePillsTheme) Description() string {
	return "Two-line pills: line 1 identity + workspace, line 2 API limits + session stats"
}

func (t *TwolinePillsTheme) Render(data StatusData) string {
	var sb strings.Builder

	// ── Line 1: Model | Path | Git ──
	modelColor, _ := GetModelConfig(data.ModelType)
	sb.WriteString(t.pill(
		fmt.Sprintf("%s%s%s %s%s%s", modelColor, Bold, data.ModelName, Reset+ColorDim, data.Version, Reset),
		PillBorder))

	sb.WriteString(" ")

	sb.WriteString(t.pill(
		fmt.Sprintf("%s%s%s", ColorBlue, ShortenPath(data.ProjectPath, 20), Reset),
		PillBorder))

	if data.GitBranch != "" {
		sb.WriteString(" ")
		gitContent := fmt.Sprintf("%s%s%s", ColorGreen, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitContent += fmt.Sprintf(" %s+%d%s", ColorGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitContent += fmt.Sprintf(" %s~%d%s", ColorOrange, data.GitDirty, Reset)
		}
		sb.WriteString(t.pill(gitContent, PillBorder))
	}

	sb.WriteString("\n")

	// ── Line 2: 5hr | 7day | Ctx | Session cost ──

	// 5hr
	color5, _ := GetBarColor(data.API5hrPercent)
	bar5 := t.miniBar(data.API5hrPercent, 8, color5)
	time5 := ""
	if data.API5hrTimeLeft != "" {
		time5 = fmt.Sprintf(" %s%s%s", ColorDim, data.API5hrTimeLeft, Reset)
	}
	sb.WriteString(t.pill(
		fmt.Sprintf("%s5h%s %s %s%d%%%s%s", ColorDim, Reset, bar5, color5, data.API5hrPercent, Reset, time5),
		PillBorder))

	sb.WriteString(" ")

	// 7day
	color7, _ := GetBarColor(data.API7dayPercent)
	bar7 := t.miniBar(data.API7dayPercent, 8, color7)
	time7 := ""
	if data.API7dayTimeLeft != "" {
		time7 = fmt.Sprintf(" %s%s%s", ColorDim, data.API7dayTimeLeft, Reset)
	}
	sb.WriteString(t.pill(
		fmt.Sprintf("%s7d%s %s %s%d%%%s%s", ColorDim, Reset, bar7, color7, data.API7dayPercent, Reset, time7),
		PillBorder))

	sb.WriteString(" ")

	// Context
	ctxColor := GetContextColor(data.ContextPercent)
	ctxBar := t.miniBar(data.ContextPercent, 6, ctxColor)
	sb.WriteString(t.pill(
		fmt.Sprintf("%sctx%s %s %s%d%%%s", ColorDim, Reset, ctxBar, ctxColor, data.ContextPercent, Reset),
		PillBorder))

	sb.WriteString(" ")

	// Session cost + day cost
	sb.WriteString(t.pill(
		fmt.Sprintf("%s%s%s %sses%s %s·%s %s%s%s %sday%s",
			ColorGreen, FormatCostShort(data.SessionCost), Reset, ColorDim, Reset,
			ColorDim, Reset,
			ColorYellow, FormatCostShort(data.DayCost), Reset, ColorDim, Reset),
		PillBorder))

	sb.WriteString("\n")

	return sb.String()
}

func (t *TwolinePillsTheme) pill(content, borderColor string) string {
	return fmt.Sprintf("%s(%s %s %s)%s", borderColor, Reset, content, borderColor, Reset)
}

func (t *TwolinePillsTheme) miniBar(percent, width int, color string) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▮", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(PillDim)
		bar.WriteString(strings.Repeat("▯", empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}
