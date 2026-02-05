package themes

import (
	"fmt"
	"strings"
)

// DragonBallTheme Dragon Ball scouter display style
type DragonBallTheme struct{}

func init() {
	RegisterTheme(&DragonBallTheme{})
}

func (t *DragonBallTheme) Name() string {
	return "dragonball"
}

func (t *DragonBallTheme) Description() string {
	return "Dragon Ball: Scouter power level circular display"
}

const (
	DBGreen  = "\033[38;2;0;255;128m"
	DBYellow = "\033[38;2;255;255;0m"
	DBOrange = "\033[38;2;255;165;0m"
	DBRed    = "\033[38;2;255;0;0m"
	DBCyan   = "\033[38;2;0;255;255m"
	DBDark   = "\033[38;2;0;60;30m"
	DBScan   = "\033[38;2;0;100;50m"
)

func (t *DragonBallTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Power level determines color
	powerLevel := data.TokenCount
	powerColor := DBGreen
	warning := ""
	if powerLevel > 9000 {
		powerColor = DBRed
		warning = " IT'S OVER 9000!!!"
	} else if powerLevel > 5000 {
		powerColor = DBOrange
	} else if powerLevel > 1000 {
		powerColor = DBYellow
	}

	// Scouter circular frame
	sb.WriteString(DBGreen + "    ╭──────────────────────────────────────────────────────────────────────────╮" + Reset + "\n")
	sb.WriteString(DBGreen + "   ╱" + DBScan + "░░" + DBGreen + "╲" + Reset + "  " + DBCyan + "◉ SCOUTER ACTIVATED" + Reset + "                                              " + DBGreen + "│" + Reset + "\n")
	sb.WriteString(DBGreen + "  │" + DBScan + "░░░░" + DBGreen + "│" + Reset + "  ════════════════════════════════════════════════════════════  " + DBGreen + "│" + Reset + "\n")

	// Power level display
	powerStr := fmt.Sprintf("POWER LEVEL: %s%d%s%s", powerColor, powerLevel, Reset, warning)
	sb.WriteString(DBGreen + "  │" + DBScan + "░" + powerColor + "◎" + DBScan + "░" + DBGreen + "│" + Reset + "  " + powerStr)
	padding := 62 - len(fmt.Sprintf("POWER LEVEL: %d%s", powerLevel, warning))
	sb.WriteString(strings.Repeat(" ", padding) + DBGreen + "│" + Reset + "\n")

	// Target info
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	sb.WriteString(DBGreen + "  │" + DBScan + "░░░░" + DBGreen + "│" + Reset + "  ")
	targetLine := fmt.Sprintf("TARGET: %s%s%s  BRANCH: %s%s%s", modelColor, modelIcon+data.ModelName, Reset, DBCyan, data.GitBranch, Reset)
	sb.WriteString(targetLine)
	sb.WriteString(strings.Repeat(" ", 62-len(fmt.Sprintf("TARGET: %s  BRANCH: %s", modelIcon+data.ModelName, data.GitBranch))) + DBGreen + "│" + Reset + "\n")

	sb.WriteString(DBGreen + "   ╲" + DBScan + "░░" + DBGreen + "╱" + Reset + "  ════════════════════════════════════════════════════════════  " + DBGreen + "│" + Reset + "\n")

	// Stats bars
	kiColor := DBGreen
	if data.ContextPercent > 75 {
		kiColor = DBRed
	}

	sb.WriteString(DBGreen + "    │" + Reset + "    " + fmt.Sprintf("%sKI%s %s %s%3d%%%s  %sSTM%s %s %s%3d%%%s  %sEND%s %s %s%3d%%%s",
		DBCyan, Reset, t.generateDBBar(data.ContextPercent, 10, kiColor), kiColor, data.ContextPercent, Reset,
		DBYellow, Reset, t.generateDBBar(100-data.API5hrPercent, 8, DBYellow), DBYellow, 100-data.API5hrPercent, Reset,
		DBOrange, Reset, t.generateDBBar(100-data.API7dayPercent, 8, DBOrange), DBOrange, 100-data.API7dayPercent, Reset))
	sb.WriteString("  " + DBGreen + "│" + Reset + "\n")

	// Bottom stats
	sb.WriteString(DBGreen + "    │" + Reset + fmt.Sprintf("    %sTIME%s %s  %sMSG%s %d  %sZENI%s $%s  %sDAY%s $%s  %sEFF%s %d%%",
		DBScan, Reset, data.SessionTime,
		DBCyan, Reset, data.MessageCount,
		DBYellow, Reset, FormatCost(data.SessionCost),
		DBOrange, Reset, FormatCost(data.DayCost),
		DBGreen, Reset, data.CacheHitRate))
	sb.WriteString("       " + DBGreen + "│" + Reset + "\n")

	sb.WriteString(DBGreen + "    ╰──────────────────────────────────────────────────────────────────────────╯" + Reset + "\n")

	return sb.String()
}

func (t *DragonBallTheme) generateDBBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(DBDark + "‹" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▰", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(DBDark)
		bar.WriteString(strings.Repeat("▱", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(DBDark + "›" + Reset)
	return bar.String()
}
