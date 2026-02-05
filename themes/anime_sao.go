package themes

import (
	"fmt"
	"strings"
)

// SAOTheme Sword Art Online game interface style
type SAOTheme struct{}

func init() {
	RegisterTheme(&SAOTheme{})
}

func (t *SAOTheme) Name() string {
	return "sao"
}

func (t *SAOTheme) Description() string {
	return "SAO: Sword Art Online game interface style"
}

const (
	SAOBlue   = "\033[38;2;0;150;200m"
	SAOCyan   = "\033[38;2;0;200;200m"
	SAOGreen  = "\033[38;2;100;200;100m"
	SAOYellow = "\033[38;2;255;200;50m"
	SAORed    = "\033[38;2;200;50;50m"
	SAOWhite  = "\033[38;2;240;240;240m"
	SAOGray   = "\033[38;2;100;120;140m"
	SAODark   = "\033[38;2;30;40;50m"
)

func (t *SAOTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(SAOBlue + "╔══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(SAOBlue + "║" + Reset + "  " + SAOCyan + "⚔" + SAOWhite + " SWORD ART ONLINE " + SAOCyan + "⚔" + Reset + "   " + SAOBlue + "ソードアート・オンライン" + Reset + "                      " + SAOBlue + "║" + Reset + "\n")
	sb.WriteString(SAOBlue + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	player := "Silica"
	if data.ModelType == "Opus" {
		player = "Kirito"
	} else if data.ModelType == "Haiku" {
		player = "Pina"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[Level Up!]%s", SAOYellow, Reset)
	}

	line1 := fmt.Sprintf("  %sPlayer:%s %s%s%s  %sName:%s %s%s%s  %s%s%s%s",
		SAOCyan, Reset, modelColor, modelIcon, data.ModelName,
		SAOGray, Reset, SAOGreen, player, Reset,
		SAOGray, data.Version, Reset, update)

	sb.WriteString(SAOBlue + "║" + Reset)
	sb.WriteString(PadRight(line1, 88))
	sb.WriteString(SAOBlue + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚔%s%s", SAOCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", SAOGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", SAOYellow, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sFloor:%s %s%s",
		SAOGreen, Reset, ShortenPath(data.ProjectPath, 42), gitInfo)

	sb.WriteString(SAOBlue + "║" + Reset)
	sb.WriteString(PadRight(line2, 88))
	sb.WriteString(SAOBlue + "║" + Reset + "\n")

	sb.WriteString(SAOBlue + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	hpColor := SAOGreen
	if data.ContextPercent > 50 {
		hpColor = SAOYellow
	}
	if data.ContextPercent > 75 {
		hpColor = SAORed
	}

	line3 := fmt.Sprintf("  %sHP%s          %s  %s%3d%%%s",
		SAOGreen, Reset, t.generateSAOBar(data.ContextPercent, 18, hpColor), hpColor, data.ContextPercent, Reset)

	sb.WriteString(SAOBlue + "║" + Reset)
	sb.WriteString(PadRight(line3, 88))
	sb.WriteString(SAOBlue + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sMP%s          %s  %s%3d%%%s  %s%s%s",
		SAOBlue, Reset, t.generateSAOBar(100-data.API5hrPercent, 18, SAOBlue),
		SAOBlue, 100-data.API5hrPercent, Reset, SAOGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(SAOBlue + "║" + Reset)
	sb.WriteString(PadRight(line4, 88))
	sb.WriteString(SAOBlue + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %sSTAMINA%s     %s  %s%3d%%%s  %s%s%s",
		SAOCyan, Reset, t.generateSAOBar(100-data.API7dayPercent, 18, SAOCyan),
		SAOCyan, 100-data.API7dayPercent, Reset, SAOGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(SAOBlue + "║" + Reset)
	sb.WriteString(PadRight(line5, 88))
	sb.WriteString(SAOBlue + "║" + Reset + "\n")

	sb.WriteString(SAOBlue + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sEXP:%s %s%s%s  %sTime:%s %s  %sQuests:%s %s%d%s  %sCol:%s %s%s%s",
		SAOWhite, Reset, SAOWhite, FormatTokens(data.TokenCount), Reset,
		SAOGray, Reset, data.SessionTime,
		SAOGray, Reset, SAOCyan, data.MessageCount, Reset,
		SAOYellow, Reset, SAOYellow, FormatCost(data.SessionCost), Reset)

	sb.WriteString(SAOBlue + "║" + Reset)
	sb.WriteString(PadRight(line6, 88))
	sb.WriteString(SAOBlue + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sCrit:%s %s%d%%%s",
		SAOGreen, Reset, SAOGreen, FormatCost(data.DayCost), Reset,
		SAOYellow, Reset, SAOYellow, FormatCost(data.BurnRate), Reset,
		SAOCyan, Reset, SAOCyan, data.CacheHitRate, Reset)

	sb.WriteString(SAOBlue + "║" + Reset)
	sb.WriteString(PadRight(line7, 88))
	sb.WriteString(SAOBlue + "║" + Reset + "\n")

	sb.WriteString(SAOBlue + "╚══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *SAOTheme) generateSAOBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(SAODark + "〈" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(SAODark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(SAODark + "〉" + Reset)
	return bar.String()
}
