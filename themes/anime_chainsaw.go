package themes

import (
	"fmt"
	"strings"
)

// ChainsawTheme Chainsaw Man devil contract style
type ChainsawTheme struct{}

func init() {
	RegisterTheme(&ChainsawTheme{})
}

func (t *ChainsawTheme) Name() string {
	return "chainsaw"
}

func (t *ChainsawTheme) Description() string {
	return "Chainsaw Man: Devil contract blood price style"
}

const (
	CSMRed     = "\033[38;2;180;30;30m"
	CSMOrange  = "\033[38;2;255;100;50m"
	CSMYellow  = "\033[38;2;255;200;50m"
	CSMBlack   = "\033[38;2;20;20;20m"
	CSMWhite   = "\033[38;2;240;240;240m"
	CSMGray    = "\033[38;2;100;100;100m"
	CSMBlood   = "\033[38;2;139;0;0m"
)

func (t *ChainsawTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(CSMRed + "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓" + Reset + "\n")
	sb.WriteString("  " + CSMOrange + "⛓" + CSMWhite + " CHAINSAW MAN " + CSMOrange + "⛓" + Reset + "   " + CSMRed + "チェンソーマン" + Reset + "   " + CSMBlood + "// DEVIL CONTRACT //" + Reset + "\n")
	sb.WriteString(CSMRed + "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	devil := "Pochita"
	if data.ModelType == "Opus" {
		devil = "Chainsaw"
	} else if data.ModelType == "Haiku" {
		devil = "Fox"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s⛓NEW⛓%s", CSMOrange, Reset)
	}

	line1 := fmt.Sprintf("  %sHunter:%s %s%s%s  %sDevil:%s %s%s%s  %s%s%s%s",
		CSMRed, Reset, modelColor, modelIcon, data.ModelName,
		CSMGray, Reset, CSMOrange, devil, Reset,
		CSMGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⛓%s%s", CSMOrange, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", CSMYellow, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", CSMBlood, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sContract:%s %s%s",
		CSMBlood, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(CSMRed + "───────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	bloodColor := CSMRed
	if data.ContextPercent > 75 {
		bloodColor = CSMBlood
	}

	line3 := fmt.Sprintf("  %sBlood%s      %s  %s%3d%%%s",
		CSMRed, Reset, t.generateCSMBar(data.ContextPercent, 18, bloodColor), bloodColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("  %sContract%s   %s  %s%3d%%%s  %s%s%s",
		CSMOrange, Reset, t.generateCSMBar(100-data.API5hrPercent, 18, CSMOrange),
		CSMOrange, 100-data.API5hrPercent, Reset, CSMGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("  %sPrice%s      %s  %s%3d%%%s  %s%s%s",
		CSMBlood, Reset, t.generateCSMBar(data.API7dayPercent, 18, CSMBlood),
		CSMBlood, data.API7dayPercent, Reset, CSMGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(CSMRed + "───────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	line6 := fmt.Sprintf("  %s%s%s kills  %s%s%s  %s%d%s hunts  %s$%s%s  %s$%s/day%s  %s%d%%%s",
		CSMWhite, FormatTokens(data.TokenCount), Reset,
		CSMGray, data.SessionTime, Reset,
		CSMOrange, data.MessageCount, Reset,
		CSMYellow, FormatCost(data.SessionCost), Reset,
		CSMRed, FormatCost(data.DayCost), Reset,
		CSMBlood, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString(CSMRed + "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓" + Reset + "\n")

	return sb.String()
}

func (t *ChainsawTheme) generateCSMBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(CSMBlack + "⟨" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(CSMBlack)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(CSMBlack + "⟩" + Reset)
	return bar.String()
}
