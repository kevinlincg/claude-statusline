package themes

import (
	"fmt"
	"strings"
)

// CompactTheme compact three-line style
type CompactTheme struct{}

func init() {
	RegisterTheme(&CompactTheme{})
}

func (t *CompactTheme) Name() string {
	return "compact"
}

func (t *CompactTheme) Description() string {
	return "Compact three-line: minimum height, complete information"
}

func (t *CompactTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Model + Version
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s⬆%s", ColorNeonOrange, Reset)
	}

	// Line 1: Model | Path + Git | Time
	sb.WriteString(fmt.Sprintf(" %s%s%s%s %s%s%s%s",
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		ColorNeonGreen, data.Version, Reset))
	sb.WriteString(update)
	sb.WriteString(fmt.Sprintf("  %s│%s  ", ColorFrame, Reset))
	sb.WriteString(fmt.Sprintf("%s📂 %s%s", ColorYellow, data.ProjectPath, Reset))
	if data.GitBranch != "" {
		sb.WriteString(fmt.Sprintf("  %s⚡%s%s", ColorCyan, data.GitBranch, Reset))
		if data.GitStaged > 0 {
			sb.WriteString(fmt.Sprintf(" %s+%d%s", ColorGreen, data.GitStaged, Reset))
		}
		if data.GitDirty > 0 {
			sb.WriteString(fmt.Sprintf(" %s~%d%s", ColorOrange, data.GitDirty, Reset))
		}
		sb.WriteString(FormatGitExtras(data, ColorGreen, ColorOrange, Dim))
	}
	sb.WriteString("\n")

	// Line 2: Session Stats | Cost
	sb.WriteString(fmt.Sprintf(" %s%5s%s tok  %s%3d%s msg  %s%6s%s",
		ColorPurple, FormatTokens(data.TokenCount), Reset,
		ColorCyan, data.MessageCount, Reset,
		ColorSilver, data.SessionTime, Reset))
	if tps := FormatTokensPerSec(data.TokensPerSec); tps != "" {
		sb.WriteString(fmt.Sprintf("  %s%s%s", ColorPurple, tps, Reset))
	}
	if lines := FormatLinesChanged(data.LinesAdded, data.LinesRemoved, ColorGreen, ColorRed); lines != "" {
		sb.WriteString("  " + lines)
	}
	sb.WriteString(fmt.Sprintf("  %s│%s  ", ColorFrame, Reset))
	sb.WriteString(fmt.Sprintf("%s%s%s ses  %s%s%s day  %s%s%s mon  %s%s/h%s  %s%d%%hit%s",
		ColorGreen, FormatCost(data.SessionCost), Reset,
		ColorYellow, FormatCost(data.DayCost), Reset,
		ColorPurple, FormatCost(data.MonthCost), Reset,
		ColorRed, FormatCost(data.BurnRate), Reset,
		ColorGreen, data.CacheHitRate, Reset))
	sb.WriteString("\n")

	// Line 3: Three progress bars
	color1, bg1 := GetBarColor(data.ContextPercent)
	color5, bg5 := GetBarColor(data.API5hrPercent)
	color7, bg7 := GetBarColor(data.API7dayPercent)

	sb.WriteString(fmt.Sprintf(" %sCtx%s %s %s%3d%%%s",
		ColorLabelDim, Reset,
		GenerateGlowBar(data.ContextPercent, 12, color1, bg1),
		color1, data.ContextPercent, Reset))
	sb.WriteString(fmt.Sprintf("  %s│%s  ", ColorFrame, Reset))
	sb.WriteString(fmt.Sprintf("%s5hr%s %s %s%3d%%%s %s%s%s",
		ColorLabelDim, Reset,
		GenerateGlowBar(data.API5hrPercent, 12, color5, bg5),
		color5, data.API5hrPercent, Reset,
		ColorDim, data.API5hrTimeLeft, Reset))
	sb.WriteString(fmt.Sprintf("  %s│%s  ", ColorFrame, Reset))
	sb.WriteString(fmt.Sprintf("%s7dy%s %s %s%3d%%%s %s%s%s",
		ColorLabelDim, Reset,
		GenerateGlowBar(data.API7dayPercent, 12, color7, bg7),
		color7, data.API7dayPercent, Reset,
		ColorDim, data.API7dayTimeLeft, Reset))
	sb.WriteString("\n")

	return sb.String()
}
