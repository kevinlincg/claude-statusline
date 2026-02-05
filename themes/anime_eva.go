package themes

import (
	"fmt"
	"strings"
)

// EVATheme Neon Genesis Evangelion NERV interface style
type EVATheme struct{}

func init() {
	RegisterTheme(&EVATheme{})
}

func (t *EVATheme) Name() string {
	return "eva"
}

func (t *EVATheme) Description() string {
	return "EVA: NERV system interface, sync rate and A.T. Field"
}

const (
	EVAOrange     = "\033[38;2;255;102;0m"
	EVARed        = "\033[38;2;204;0;0m"
	EVAPurple     = "\033[38;2;128;0;128m"
	EVAGreen      = "\033[38;2;0;255;0m"
	EVAYellow     = "\033[38;2;255;255;0m"
	EVABlue       = "\033[38;2;0;128;255m"
	EVADark       = "\033[38;2;40;40;40m"
	EVAWhite      = "\033[38;2;255;255;255m"
	EVABgOrange   = "\033[48;2;255;102;0m"
	EVABgRed      = "\033[48;2;80;0;0m"
)

func (t *EVATheme) Render(data StatusData) string {
	var sb strings.Builder

	// NERV Header
	sb.WriteString(EVAOrange + "╔══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")

	// Title line with NERV branding
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	warning := ""
	if data.UpdateAvailable {
		warning = EVARed + " [UPDATE AVAILABLE]" + Reset
	}

	title := fmt.Sprintf(" %s%s NERV SYSTEM%s  %s%s%s %s%s%s%s",
		EVABgOrange, EVADark, Reset,
		modelColor, modelIcon, data.ModelName,
		EVAGreen, data.Version, Reset, warning)

	sb.WriteString(EVAOrange + "║" + Reset)
	sb.WriteString(PadRight(title, 88))
	sb.WriteString(EVAOrange + "║" + Reset + "\n")

	sb.WriteString(EVAOrange + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Project and Git info
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("%s[%s]%s", EVABlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", EVAGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s*%d%s", EVAYellow, data.GitDirty, Reset)
		}
	}

	line1 := fmt.Sprintf(" %sPILOT:%s %s  %sUNIT:%s %s",
		EVAOrange, Reset, ShortenPath(data.ProjectPath, 30),
		EVAOrange, Reset, gitInfo)

	sb.WriteString(EVAOrange + "║" + Reset)
	sb.WriteString(PadRight(line1, 88))
	sb.WriteString(EVAOrange + "║" + Reset + "\n")

	sb.WriteString(EVAOrange + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Sync Rate (Context) and A.T. Field (API limits)
	syncColor := EVAGreen
	if data.ContextPercent > 75 {
		syncColor = EVARed
	} else if data.ContextPercent > 50 {
		syncColor = EVAYellow
	}

	line2 := fmt.Sprintf(" %sSYNC RATE%s    %s  %s%3d%%%s   %sA.T. FIELD%s  %s  %s%3d%%%s  %s%s%s",
		EVAOrange, Reset,
		t.generateEVABar(data.ContextPercent, 15, syncColor),
		syncColor, data.ContextPercent, Reset,
		EVAPurple, Reset,
		t.generateEVABar(100-data.API5hrPercent, 12, EVAPurple),
		EVAPurple, 100-data.API5hrPercent, Reset,
		EVADark, data.API5hrTimeLeft, Reset)

	sb.WriteString(EVAOrange + "║" + Reset)
	sb.WriteString(PadRight(line2, 88))
	sb.WriteString(EVAOrange + "║" + Reset + "\n")

	// Energy and damage stats
	line3 := fmt.Sprintf(" %sENERGY%s %s%s%s tok  %sTIME%s %s  %sDAMAGE%s %s%s%s  %sWEEKLY%s %s  %s%3d%%%s  %s%s%s",
		EVAGreen, Reset, EVAWhite, FormatTokens(data.TokenCount), Reset,
		EVABlue, Reset, data.SessionTime,
		EVARed, Reset, EVAYellow, FormatCost(data.SessionCost), Reset,
		EVAOrange, Reset,
		t.generateEVABar(100-data.API7dayPercent, 10, EVAOrange),
		EVAOrange, 100-data.API7dayPercent, Reset,
		EVADark, data.API7dayTimeLeft, Reset)

	sb.WriteString(EVAOrange + "║" + Reset)
	sb.WriteString(PadRight(line3, 88))
	sb.WriteString(EVAOrange + "║" + Reset + "\n")

	// Bottom stats
	line4 := fmt.Sprintf(" %sMSG%s %s%d%s  %sDAY%s %s%s%s  %sMON%s %s%s%s  %sRATE%s %s%s/h%s  %sEFF%s %s%d%%%s",
		EVABlue, Reset, EVAWhite, data.MessageCount, Reset,
		EVAYellow, Reset, EVAYellow, FormatCost(data.DayCost), Reset,
		EVAPurple, Reset, EVAPurple, FormatCost(data.MonthCost), Reset,
		EVARed, Reset, EVARed, FormatCost(data.BurnRate), Reset,
		EVAGreen, Reset, EVAGreen, data.CacheHitRate, Reset)

	sb.WriteString(EVAOrange + "║" + Reset)
	sb.WriteString(PadRight(line4, 88))
	sb.WriteString(EVAOrange + "║" + Reset + "\n")

	sb.WriteString(EVAOrange + "╚══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *EVATheme) generateEVABar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(EVADark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(EVADark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(EVADark + "]" + Reset)
	return bar.String()
}
