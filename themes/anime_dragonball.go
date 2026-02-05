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
	return "Dragon Ball: Scouter power level display"
}

const (
	DBGreen       = "\033[38;2;0;255;128m"
	DBBrightGreen = "\033[38;2;128;255;128m"
	DBYellow      = "\033[38;2;255;255;0m"
	DBOrange      = "\033[38;2;255;165;0m"
	DBRed         = "\033[38;2;255;0;0m"
	DBCyan        = "\033[38;2;0;255;255m"
	DBDark        = "\033[38;2;0;80;40m"
	DBScanLine    = "\033[38;2;0;100;50m"
)

func (t *DragonBallTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Scouter frame
	sb.WriteString(DBGreen + "┌─────────────────────────────────────────────────────────────────────────────────┐" + Reset + "\n")

	// Model and scan line effect
	modelColor, modelIcon := GetModelConfig(data.ModelType)

	// Power level calculation (tokens as power level)
	powerLevel := data.TokenCount
	powerColor := DBGreen
	if powerLevel > 100000 {
		powerColor = DBRed // It's over 9000!!!
	} else if powerLevel > 50000 {
		powerColor = DBOrange
	} else if powerLevel > 10000 {
		powerColor = DBYellow
	}

	line1 := fmt.Sprintf(" %s▓▓%s SCOUTER v%s %s▓▓%s  %s%s%s %s",
		DBScanLine, Reset, data.Version, DBScanLine, Reset,
		modelColor, modelIcon, data.ModelName, Reset)
	if data.UpdateAvailable {
		line1 += fmt.Sprintf(" %s[NEW]%s", DBYellow, Reset)
	}

	sb.WriteString(DBGreen + "│" + Reset)
	sb.WriteString(PadRight(line1, 83))
	sb.WriteString(DBGreen + "│" + Reset + "\n")

	sb.WriteString(DBGreen + "├" + DBScanLine + "─────────────────────────────────────────────────────────────────────────────────" + DBGreen + "┤" + Reset + "\n")

	// Power Level display (main focus)
	powerStr := fmt.Sprintf("%d", powerLevel)
	line2 := fmt.Sprintf(" %s>>> POWER LEVEL:%s %s%s%s",
		DBGreen, Reset, powerColor, powerStr, Reset)

	sb.WriteString(DBGreen + "│" + Reset)
	sb.WriteString(PadRight(line2, 83))
	sb.WriteString(DBGreen + "│" + Reset + "\n")

	// Target info
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf(" %s⚡%s%s", DBCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", DBBrightGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", DBYellow, data.GitDirty, Reset)
		}
	}

	line3 := fmt.Sprintf(" %sTARGET:%s %s%s",
		DBGreen, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)

	sb.WriteString(DBGreen + "│" + Reset)
	sb.WriteString(PadRight(line3, 83))
	sb.WriteString(DBGreen + "│" + Reset + "\n")

	sb.WriteString(DBGreen + "├" + DBScanLine + "─────────────────────────────────────────────────────────────────────────────────" + DBGreen + "┤" + Reset + "\n")

	// Ki gauge (Context)
	kiColor := DBGreen
	if data.ContextPercent > 75 {
		kiColor = DBRed
	} else if data.ContextPercent > 50 {
		kiColor = DBYellow
	}

	line4 := fmt.Sprintf(" %sKI%s %s %s%3d%%%s  %sSTAMINA%s %s %s%3d%%%s %s%s%s  %sENDURE%s %s %s%3d%%%s %s%s%s",
		DBCyan, Reset,
		t.generateDBBar(data.ContextPercent, 12, kiColor),
		kiColor, data.ContextPercent, Reset,
		DBYellow, Reset,
		t.generateDBBar(100-data.API5hrPercent, 10, DBYellow),
		DBYellow, 100-data.API5hrPercent, Reset,
		DBDark, data.API5hrTimeLeft, Reset,
		DBOrange, Reset,
		t.generateDBBar(100-data.API7dayPercent, 10, DBOrange),
		DBOrange, 100-data.API7dayPercent, Reset,
		DBDark, data.API7dayTimeLeft, Reset)

	sb.WriteString(DBGreen + "│" + Reset)
	sb.WriteString(PadRight(line4, 83))
	sb.WriteString(DBGreen + "│" + Reset + "\n")

	// Stats row
	line5 := fmt.Sprintf(" %sTIME%s %s  %sMSG%s %s%d%s  %sZENI%s %s%s%s  %sDAY%s %s%s%s  %sRATE%s %s%s/h%s  %sFOCUS%s %s%d%%%s",
		DBGreen, Reset, data.SessionTime,
		DBCyan, Reset, DBCyan, data.MessageCount, Reset,
		DBYellow, Reset, DBYellow, FormatCost(data.SessionCost), Reset,
		DBOrange, Reset, DBOrange, FormatCost(data.DayCost), Reset,
		DBRed, Reset, DBRed, FormatCost(data.BurnRate), Reset,
		DBBrightGreen, Reset, DBBrightGreen, data.CacheHitRate, Reset)

	sb.WriteString(DBGreen + "│" + Reset)
	sb.WriteString(PadRight(line5, 83))
	sb.WriteString(DBGreen + "│" + Reset + "\n")

	sb.WriteString(DBGreen + "└─────────────────────────────────────────────────────────────────────────────────┘" + Reset + "\n")

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
