package themes

import (
	"fmt"
	"strings"
)

// ZenTheme 極簡禪風
type ZenTheme struct{}

func init() {
	RegisterTheme(&ZenTheme{})
}

func (t *ZenTheme) Name() string {
	return "zen"
}

func (t *ZenTheme) Description() string {
	return "禪風：極簡留白，寧靜淡雅"
}

const (
	ZenWhite    = "\033[38;2;240;240;240m"
	ZenGray     = "\033[38;2;120;120;120m"
	ZenDimGray  = "\033[38;2;80;80;80m"
	ZenSoftGreen= "\033[38;2;144;180;148m"
	ZenSoftGold = "\033[38;2;200;180;140m"
	ZenSoftRed  = "\033[38;2;180;120;120m"
)

func (t *ZenTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Top border - subtle dots
	sb.WriteString(ZenDimGray)
	sb.WriteString("  · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · ·")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// Line 1: Model and path with generous spacing
	modelColor, _ := GetModelConfig(data.ModelType)
	update := ""
	if data.UpdateAvailable {
		update = ZenSoftGold + " ↑" + Reset
	}

	line1 := fmt.Sprintf("    %s%s%s  %s%s%s%s      %s%s%s",
		modelColor, data.ModelName, Reset,
		ZenDimGray, data.Version, Reset, update,
		ZenWhite, ShortenPath(data.ProjectPath, 25), Reset)
	if data.GitBranch != "" {
		line1 += fmt.Sprintf("  %s· %s%s", ZenDimGray, data.GitBranch, Reset)
		if data.GitStaged > 0 || data.GitDirty > 0 {
			line1 += fmt.Sprintf(" %s(%s+%d%s %s~%d%s)%s",
				ZenDimGray,
				ZenSoftGreen, data.GitStaged, ZenDimGray,
				ZenSoftGold, data.GitDirty, ZenDimGray,
				Reset)
		}
	}
	sb.WriteString(line1)
	sb.WriteString("\n")

	// Empty breathing space
	sb.WriteString("\n")

	// Line 2: Key metrics with lots of space
	line2 := fmt.Sprintf("    %s%s%s     %s%d%s msg     %s%s%s          %s%s%s  %s%s%s  %s%s/h%s",
		ZenGray, FormatTokens(data.TokenCount), Reset,
		ZenGray, data.MessageCount, Reset,
		ZenDimGray, data.SessionTime, Reset,
		ZenSoftGreen, FormatCostShort(data.SessionCost), Reset,
		ZenSoftGold, FormatCostShort(data.DayCost), Reset,
		ZenSoftRed, FormatCostShort(data.BurnRate), Reset)
	sb.WriteString(line2)
	sb.WriteString("\n")

	// Empty breathing space
	sb.WriteString("\n")

	// Line 3: Progress with minimal bars
	ctxBar := t.generateZenBar(data.ContextPercent, 20)
	bar5 := t.generateZenBar(data.API5hrPercent, 12)
	bar7 := t.generateZenBar(data.API7dayPercent, 12)

	ctxColor := ZenSoftGreen
	if data.ContextPercent >= 80 {
		ctxColor = ZenSoftRed
	} else if data.ContextPercent >= 60 {
		ctxColor = ZenSoftGold
	}

	line3 := fmt.Sprintf("    %sctx%s %s %s%d%s      %s5h%s %s %s%d%s      %s7d%s %s %s%d%s",
		ZenDimGray, Reset, ctxBar, ctxColor, data.ContextPercent, Reset,
		ZenDimGray, Reset, bar5, ZenGray, data.API5hrPercent, Reset,
		ZenDimGray, Reset, bar7, ZenGray, data.API7dayPercent, Reset)
	sb.WriteString(line3)
	sb.WriteString("\n")

	// Bottom border
	sb.WriteString(ZenDimGray)
	sb.WriteString("  · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · ·")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	return sb.String()
}

func (t *ZenTheme) generateZenBar(percent, width int) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(ZenGray)
		bar.WriteString(strings.Repeat("─", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(ZenDimGray)
		bar.WriteString(strings.Repeat("·", empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}
