package themes

import (
	"fmt"
	"strings"
)

// MHATheme My Hero Academia quirk analysis style
type MHATheme struct{}

func init() {
	RegisterTheme(&MHATheme{})
}

func (t *MHATheme) Name() string {
	return "mha"
}

func (t *MHATheme) Description() string {
	return "MHA: My Hero Academia quirk analysis style"
}

const (
	MHAGreen   = "\033[38;2;0;180;0m"
	MHARed     = "\033[38;2;220;20;60m"
	MHABlue    = "\033[38;2;30;144;255m"
	MHAYellow  = "\033[38;2;255;215;0m"
	MHAOrange  = "\033[38;2;255;140;0m"
	MHAWhite   = "\033[38;2;255;255;255m"
	MHAGray    = "\033[38;2;100;100;100m"
	MHADark    = "\033[38;2;40;40;40m"
)

func (t *MHATheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(MHAGreen + "╔════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(MHAGreen + "║" + Reset + "  " + MHAYellow + "★" + MHAWhite + " HERO ANALYSIS " + MHAYellow + "★" + Reset + "   " + MHAGreen + "Plus Ultra!" + Reset + "                                            " + MHAGreen + "║" + Reset + "\n")
	sb.WriteString(MHAGreen + "╠════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	heroRank := "#10"
	if data.ModelType == "Opus" {
		heroRank = "#1"
	} else if data.ModelType == "Sonnet" {
		heroRank = "#5"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[NEW QUIRK!]%s", MHAYellow, Reset)
	}

	line1 := fmt.Sprintf("  %sHero:%s %s%s%s  %sRank:%s %s%s%s  %s%s%s%s",
		MHABlue, Reset, modelColor, modelIcon, data.ModelName,
		MHAYellow, Reset, MHAYellow, heroRank, Reset,
		MHAGray, data.Version, Reset, update)

	sb.WriteString(MHAGreen + "║" + Reset)
	sb.WriteString(PadRight(line1, 86))
	sb.WriteString(MHAGreen + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚡%s%s", MHABlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", MHAGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", MHAOrange, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sAgency:%s %s%s",
		MHARed, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)

	sb.WriteString(MHAGreen + "║" + Reset)
	sb.WriteString(PadRight(line2, 86))
	sb.WriteString(MHAGreen + "║" + Reset + "\n")

	sb.WriteString(MHAGreen + "╠════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	powerColor := MHAGreen
	if data.ContextPercent > 75 {
		powerColor = MHARed
	} else if data.ContextPercent > 50 {
		powerColor = MHAOrange
	}

	line3 := fmt.Sprintf("  %sQuirk Power%s  %s  %s%3d%%%s",
		MHAGreen, Reset, t.generateMHABar(data.ContextPercent, 18, powerColor), powerColor, data.ContextPercent, Reset)

	sb.WriteString(MHAGreen + "║" + Reset)
	sb.WriteString(PadRight(line3, 86))
	sb.WriteString(MHAGreen + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sStamina%s      %s  %s%3d%%%s  %s%s%s",
		MHABlue, Reset, t.generateMHABar(100-data.API5hrPercent, 18, MHABlue),
		MHABlue, 100-data.API5hrPercent, Reset, MHAGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(MHAGreen + "║" + Reset)
	sb.WriteString(PadRight(line4, 86))
	sb.WriteString(MHAGreen + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %sResolve%s      %s  %s%3d%%%s  %s%s%s",
		MHAOrange, Reset, t.generateMHABar(100-data.API7dayPercent, 18, MHAOrange),
		MHAOrange, 100-data.API7dayPercent, Reset, MHAGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(MHAGreen + "║" + Reset)
	sb.WriteString(PadRight(line5, 86))
	sb.WriteString(MHAGreen + "║" + Reset + "\n")

	sb.WriteString(MHAGreen + "╠════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sEXP:%s %s%s%s  %sTime:%s %s  %sRescues:%s %s%d%s  %sPay:%s %s%s%s  %sDaily:%s %s%s%s",
		MHAYellow, Reset, MHAYellow, FormatTokens(data.TokenCount), Reset,
		MHAGray, Reset, data.SessionTime,
		MHAGray, Reset, MHAWhite, data.MessageCount, Reset,
		MHAGreen, Reset, MHAGreen, FormatCost(data.SessionCost), Reset,
		MHAOrange, Reset, MHAOrange, FormatCost(data.DayCost), Reset)

	sb.WriteString(MHAGreen + "║" + Reset)
	sb.WriteString(PadRight(line6, 86))
	sb.WriteString(MHAGreen + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sRate:%s %s%s/h%s  %sAccuracy:%s %s%d%%%s",
		MHARed, Reset, MHARed, FormatCost(data.BurnRate), Reset,
		MHABlue, Reset, MHABlue, data.CacheHitRate, Reset)

	sb.WriteString(MHAGreen + "║" + Reset)
	sb.WriteString(PadRight(line7, 86))
	sb.WriteString(MHAGreen + "║" + Reset + "\n")

	sb.WriteString(MHAGreen + "╚════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *MHATheme) generateMHABar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(MHADark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(MHADark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(MHADark + "]" + Reset)
	return bar.String()
}
