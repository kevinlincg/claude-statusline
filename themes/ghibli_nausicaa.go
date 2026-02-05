package themes

import (
	"fmt"
	"strings"
)

// NausicaaTheme Nausicaa Valley of the Wind style
type NausicaaTheme struct{}

func init() {
	RegisterTheme(&NausicaaTheme{})
}

func (t *NausicaaTheme) Name() string {
	return "nausicaa"
}

func (t *NausicaaTheme) Description() string {
	return "Nausicaa: Valley of the Wind toxic jungle style"
}

const (
	NausBlue   = "\033[38;2;70;130;180m"
	NausCyan   = "\033[38;2;0;180;180m"
	NausGreen  = "\033[38;2;85;170;127m"
	NausPurple = "\033[38;2;147;112;219m"
	NausGold   = "\033[38;2;218;165;32m"
	NausWhite  = "\033[38;2;245;245;245m"
	NausGray   = "\033[38;2;119;136;153m"
	NausDark   = "\033[38;2;46;64;83m"
)

func (t *NausicaaTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(NausBlue + "══════════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")
	sb.WriteString("  " + NausCyan + "🌬" + Reset + " " + NausBlue + "Valley of the Wind" + Reset + "   " + NausGreen + "風の谷のナウシカ" + Reset + "\n")
	sb.WriteString(NausBlue + "══════════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	role := "Mehve Pilot"
	if data.ModelType == "Opus" {
		role = "Nausicaa"
	} else if data.ModelType == "Haiku" {
		role = "Ohmu Friend"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[Wind Rising]%s", NausCyan, Reset)
	}

	line1 := fmt.Sprintf("  %sPilot:%s %s%s%s  %sRole:%s %s%s%s  %s%s%s%s",
		NausBlue, Reset, modelColor, modelIcon, data.ModelName,
		NausGray, Reset, NausGold, role, Reset,
		NausGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s🍃%s%s", NausGreen, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", NausCyan, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", NausPurple, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sValley:%s %s%s",
		NausGreen, Reset, ShortenPath(data.ProjectPath, 42), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(NausBlue + "──────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	windColor := NausCyan
	if data.ContextPercent > 75 {
		windColor = NausPurple
	}

	line3 := fmt.Sprintf("  %sWind%s       %s  %s%3d%%%s",
		NausCyan, Reset, t.generateNausBar(data.ContextPercent, 18, windColor), windColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("  %sPurity%s     %s  %s%3d%%%s  %s%s%s",
		NausGreen, Reset, t.generateNausBar(100-data.API5hrPercent, 18, NausGreen),
		NausGreen, 100-data.API5hrPercent, Reset, NausGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("  %sOhmu Bond%s  %s  %s%3d%%%s  %s%s%s",
		NausBlue, Reset, t.generateNausBar(100-data.API7dayPercent, 18, NausBlue),
		NausBlue, 100-data.API7dayPercent, Reset, NausGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(NausBlue + "──────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	line6 := fmt.Sprintf("  %s%s%s spores  %s%s%s  %s%d%s flights  %s$%s%s  %s$%s/day%s  %s%d%%%s sync",
		NausWhite, FormatTokens(data.TokenCount), Reset,
		NausGray, data.SessionTime, Reset,
		NausCyan, data.MessageCount, Reset,
		NausGold, FormatCost(data.SessionCost), Reset,
		NausGreen, FormatCost(data.DayCost), Reset,
		NausBlue, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString(NausBlue + "══════════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	return sb.String()
}

func (t *NausicaaTheme) generateNausBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(NausDark + "〈" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("●", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(NausDark)
		bar.WriteString(strings.Repeat("○", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(NausDark + "〉" + Reset)
	return bar.String()
}
