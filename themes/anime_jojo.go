package themes

import (
	"fmt"
	"strings"
)

// JoJoTheme JoJo's Bizarre Adventure stand stats style
type JoJoTheme struct{}

func init() {
	RegisterTheme(&JoJoTheme{})
}

func (t *JoJoTheme) Name() string {
	return "jojo"
}

func (t *JoJoTheme) Description() string {
	return "JoJo: Stand stats and bizarre style"
}

const (
	JoJoGold   = "\033[38;2;255;215;0m"
	JoJoPurple = "\033[38;2;148;0;211m"
	JoJoBlue   = "\033[38;2;65;105;225m"
	JoJoGreen  = "\033[38;2;0;200;100m"
	JoJoRed    = "\033[38;2;220;20;60m"
	JoJoPink   = "\033[38;2;255;105;180m"
	JoJoWhite  = "\033[38;2;255;255;255m"
	JoJoGray   = "\033[38;2;100;100;100m"
	JoJoDark   = "\033[38;2;40;40;40m"
)

func (t *JoJoTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(JoJoGold + "╔═══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(JoJoGold + "║" + Reset + "  " + JoJoPurple + "『" + JoJoWhite + "STAND ANALYSIS" + JoJoPurple + "』" + Reset + "   " + JoJoGold + "ゴゴゴゴゴ" + Reset + "                                           " + JoJoGold + "║" + Reset + "\n")
	sb.WriteString(JoJoGold + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	standType := "Close-Range"
	if data.ModelType == "Opus" {
		standType = "Requiem"
	} else if data.ModelType == "Haiku" {
		standType = "Automatic"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %sメメタァ!%s", JoJoPink, Reset)
	}

	line1 := fmt.Sprintf("  %sStand:%s %s%s%s  %sType:%s %s%s%s  %s%s%s%s",
		JoJoPurple, Reset, modelColor, modelIcon, data.ModelName,
		JoJoGray, Reset, JoJoGold, standType, Reset,
		JoJoGray, data.Version, Reset, update)

	sb.WriteString(JoJoGold + "║" + Reset)
	sb.WriteString(PadRight(line1, 89))
	sb.WriteString(JoJoGold + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s★%s%s", JoJoBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", JoJoGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", JoJoRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sUser:%s %s%s",
		JoJoBlue, Reset, ShortenPath(data.ProjectPath, 45), gitInfo)

	sb.WriteString(JoJoGold + "║" + Reset)
	sb.WriteString(PadRight(line2, 89))
	sb.WriteString(JoJoGold + "║" + Reset + "\n")

	sb.WriteString(JoJoGold + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Stand stats style
	powerRank := t.getRank(100 - data.ContextPercent)
	speedRank := t.getRank(100 - data.API5hrPercent)
	durabilityRank := t.getRank(100 - data.API7dayPercent)

	line3 := fmt.Sprintf("  %sPOWER%s     %s  %s%s%s",
		JoJoRed, Reset, t.generateJoJoBar(100-data.ContextPercent, 16, JoJoRed), JoJoWhite, powerRank, Reset)

	sb.WriteString(JoJoGold + "║" + Reset)
	sb.WriteString(PadRight(line3, 89))
	sb.WriteString(JoJoGold + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sSPEED%s     %s  %s%s%s  %s%s%s",
		JoJoBlue, Reset, t.generateJoJoBar(100-data.API5hrPercent, 16, JoJoBlue),
		JoJoWhite, speedRank, Reset, JoJoGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(JoJoGold + "║" + Reset)
	sb.WriteString(PadRight(line4, 89))
	sb.WriteString(JoJoGold + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %sDURABILITY%s%s  %s%s%s  %s%s%s",
		JoJoGreen, Reset, t.generateJoJoBar(100-data.API7dayPercent, 16, JoJoGreen),
		JoJoWhite, durabilityRank, Reset, JoJoGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(JoJoGold + "║" + Reset)
	sb.WriteString(PadRight(line5, 89))
	sb.WriteString(JoJoGold + "║" + Reset + "\n")

	sb.WriteString(JoJoGold + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sEXP:%s %s%s%s  %sTime:%s %s  %sORA:%s %s%d%s  %sYen:%s %s%s%s  %sDaily:%s %s%s%s",
		JoJoPurple, Reset, JoJoPurple, FormatTokens(data.TokenCount), Reset,
		JoJoGray, Reset, data.SessionTime,
		JoJoGray, Reset, JoJoWhite, data.MessageCount, Reset,
		JoJoGold, Reset, JoJoGold, FormatCost(data.SessionCost), Reset,
		JoJoPink, Reset, JoJoPink, FormatCost(data.DayCost), Reset)

	sb.WriteString(JoJoGold + "║" + Reset)
	sb.WriteString(PadRight(line6, 89))
	sb.WriteString(JoJoGold + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sRate:%s %s%s/h%s  %sPrecision:%s %s%d%%%s %s%s%s",
		JoJoRed, Reset, JoJoRed, FormatCost(data.BurnRate), Reset,
		JoJoBlue, Reset, JoJoBlue, data.CacheHitRate, Reset,
		JoJoGold, t.getRank(data.CacheHitRate), Reset)

	sb.WriteString(JoJoGold + "║" + Reset)
	sb.WriteString(PadRight(line7, 89))
	sb.WriteString(JoJoGold + "║" + Reset + "\n")

	sb.WriteString(JoJoGold + "╚═══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *JoJoTheme) getRank(percent int) string {
	if percent >= 90 {
		return "A"
	} else if percent >= 70 {
		return "B"
	} else if percent >= 50 {
		return "C"
	} else if percent >= 30 {
		return "D"
	}
	return "E"
}

func (t *JoJoTheme) generateJoJoBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(JoJoDark + "「" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("■", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(JoJoDark)
		bar.WriteString(strings.Repeat("□", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(JoJoDark + "」" + Reset)
	return bar.String()
}
