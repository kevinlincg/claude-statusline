package themes

import (
	"fmt"
	"strings"
)

// ShonenTheme Shonen manga action panel style
type ShonenTheme struct{}

func init() {
	RegisterTheme(&ShonenTheme{})
}

func (t *ShonenTheme) Name() string {
	return "shonen"
}

func (t *ShonenTheme) Description() string {
	return "Shonen: Action manga panel with speed lines"
}

const (
	SHNRed    = "\033[38;2;255;50;50m"
	SHNOrange = "\033[38;2;255;150;0m"
	SHNYellow = "\033[38;2;255;255;0m"
	SHNBlue   = "\033[38;2;50;150;255m"
	SHNWhite  = "\033[38;2;255;255;255m"
	SHNBlack  = "\033[38;2;30;30;30m"
)

func (t *ShonenTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Action lines header
	sb.WriteString(SHNBlack + "█" + SHNRed + "///" + SHNOrange + "///" + SHNYellow + "///" + SHNWhite + "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓" + SHNYellow + "\\\\\\" + SHNOrange + "\\\\\\" + SHNRed + "\\\\\\" + SHNBlack + "█" + Reset + "\n")
	sb.WriteString(SHNBlack + "█" + Reset + "                                                                                     " + SHNBlack + "█" + Reset + "\n")
	sb.WriteString(SHNBlack + "█" + Reset + "       " + SHNRed + "【" + SHNWhite + " Ｓ Ｈ Ｏ Ｎ Ｅ Ｎ  " + SHNYellow + "少年マンガ" + SHNRed + " 】" + Reset + "                                  " + SHNBlack + "█" + Reset + "\n")
	sb.WriteString(SHNBlack + "█" + Reset + "                                                                                     " + SHNBlack + "█" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	fighter := "Rival"
	if data.ModelType == "Opus" {
		fighter = "Protagonist"
	} else if data.ModelType == "Haiku" {
		fighter = "Sidekick"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s⚡POWER UP!⚡%s", SHNYellow, Reset)
	}

	line1 := fmt.Sprintf("  %s▶▶%s %sFIGHTER:%s %s%s%s  %sCLASS:%s %s%s%s  %s%s%s%s",
		SHNRed, Reset, SHNOrange, Reset, modelColor, modelIcon, data.ModelName,
		SHNBlack, Reset, SHNRed, fighter, Reset,
		SHNBlack, data.Version, Reset, update)

	sb.WriteString(SHNBlack + "█" + Reset)
	sb.WriteString(PadRight(line1, 84))
	sb.WriteString(SHNBlack + "█" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚔%s%s", SHNOrange, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", SHNYellow, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", SHNRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %s▶▶%s %sARENa:%s %s%s",
		SHNOrange, Reset, SHNBlue, Reset, ShortenPath(data.ProjectPath, 42), gitInfo)

	sb.WriteString(SHNBlack + "█" + Reset)
	sb.WriteString(PadRight(line2, 84))
	sb.WriteString(SHNBlack + "█" + Reset + "\n")

	sb.WriteString(SHNBlack + "█" + Reset + "                                                                                     " + SHNBlack + "█" + Reset + "\n")
	sb.WriteString(SHNBlack + "█" + SHNRed + "═══════════════════════════════════════════════════════════════════════════════════" + SHNBlack + "█" + Reset + "\n")
	sb.WriteString(SHNBlack + "█" + Reset + "                                                                                     " + SHNBlack + "█" + Reset + "\n")

	// Power levels with action style
	powerColor := SHNOrange
	if data.ContextPercent > 75 {
		powerColor = SHNRed
	}

	line3 := fmt.Sprintf("    %s>>> POWER LEVEL%s   %s %s%3d%%%s %s<<<<<<<%s",
		SHNRed, Reset, t.generateSHNBar(data.ContextPercent, 14, powerColor), powerColor, data.ContextPercent, Reset, powerColor, Reset)

	sb.WriteString(SHNBlack + "█" + Reset)
	sb.WriteString(PadRight(line3, 84))
	sb.WriteString(SHNBlack + "█" + Reset + "\n")

	line4 := fmt.Sprintf("    %s>>> SPIRIT%s        %s %s%3d%%%s  %s%s%s",
		SHNOrange, Reset, t.generateSHNBar(100-data.API5hrPercent, 14, SHNOrange),
		SHNOrange, 100-data.API5hrPercent, Reset, SHNBlack, data.API5hrTimeLeft, Reset)

	sb.WriteString(SHNBlack + "█" + Reset)
	sb.WriteString(PadRight(line4, 84))
	sb.WriteString(SHNBlack + "█" + Reset + "\n")

	line5 := fmt.Sprintf("    %s>>> STAMINA%s       %s %s%3d%%%s  %s%s%s",
		SHNYellow, Reset, t.generateSHNBar(100-data.API7dayPercent, 14, SHNYellow),
		SHNYellow, 100-data.API7dayPercent, Reset, SHNBlack, data.API7dayTimeLeft, Reset)

	sb.WriteString(SHNBlack + "█" + Reset)
	sb.WriteString(PadRight(line5, 84))
	sb.WriteString(SHNBlack + "█" + Reset + "\n")

	sb.WriteString(SHNBlack + "█" + Reset + "                                                                                     " + SHNBlack + "█" + Reset + "\n")
	sb.WriteString(SHNBlack + "█" + SHNOrange + "═══════════════════════════════════════════════════════════════════════════════════" + SHNBlack + "█" + Reset + "\n")
	sb.WriteString(SHNBlack + "█" + Reset + "                                                                                     " + SHNBlack + "█" + Reset + "\n")

	line6 := fmt.Sprintf("    %s%s%s exp  %s%s%s  %s%d%s battles  %s$%s%s cost  %s%d%%%s crit",
		SHNWhite, FormatTokens(data.TokenCount), Reset,
		SHNBlack, data.SessionTime, Reset,
		SHNBlue, data.MessageCount, Reset,
		SHNYellow, FormatCost(data.SessionCost), Reset,
		SHNRed, data.CacheHitRate, Reset)

	sb.WriteString(SHNBlack + "█" + Reset)
	sb.WriteString(PadRight(line6, 84))
	sb.WriteString(SHNBlack + "█" + Reset + "\n")

	sb.WriteString(SHNBlack + "█" + Reset + "                                                                                     " + SHNBlack + "█" + Reset + "\n")
	sb.WriteString(SHNBlack + "█" + SHNRed + "\\\\\\" + SHNOrange + "\\\\\\" + SHNYellow + "\\\\\\" + SHNWhite + "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓" + SHNYellow + "///" + SHNOrange + "///" + SHNRed + "///" + SHNBlack + "█" + Reset + "\n")

	return sb.String()
}

func (t *ShonenTheme) generateSHNBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(SHNBlack + "〖" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▰", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(SHNBlack)
		bar.WriteString(strings.Repeat("▱", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(SHNBlack + "〗" + Reset)
	return bar.String()
}
