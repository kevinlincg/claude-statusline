package themes

import (
	"fmt"
	"strings"
)

// MononokeTheme Princess Mononoke forest god style
type MononokeTheme struct{}

func init() {
	RegisterTheme(&MononokeTheme{})
}

func (t *MononokeTheme) Name() string {
	return "mononoke"
}

func (t *MononokeTheme) Description() string {
	return "Mononoke: Princess Mononoke forest spirit style"
}

const (
	MNKGreen  = "\033[38;2;34;139;34m"
	MNKRed    = "\033[38;2;139;0;0m"
	MNKBrown  = "\033[38;2;139;90;43m"
	MNKWhite  = "\033[38;2;240;240;240m"
	MNKBlue   = "\033[38;2;70;130;180m"
	MNKGray   = "\033[38;2;100;100;100m"
	MNKDark   = "\033[38;2;30;40;30m"
)

func (t *MononokeTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(MNKGreen + "  ══════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")
	sb.WriteString("    " + MNKWhite + "🦌" + Reset + " " + MNKGreen + "Forest Spirit" + Reset + "   " + MNKBrown + "もののけ姫" + Reset + "\n")
	sb.WriteString(MNKGreen + "  ──────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	spirit := "Kodama"
	if data.ModelType == "Opus" {
		spirit = "Shishigami"
	} else if data.ModelType == "Haiku" {
		spirit = "Wolf"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[Awakened]%s", MNKWhite, Reset)
	}

	line1 := fmt.Sprintf("    %sSpirit:%s %s%s%s  %sForm:%s %s%s%s  %s%s%s%s",
		MNKGreen, Reset, modelColor, modelIcon, data.ModelName,
		MNKGray, Reset, MNKWhite, spirit, Reset,
		MNKGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s🌿%s%s", MNKGreen, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", MNKGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", MNKRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("    %sTerritory:%s %s%s",
		MNKBrown, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(MNKGreen + "  ──────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	lifeColor := MNKGreen
	if data.ContextPercent > 75 {
		lifeColor = MNKRed
	}

	line3 := fmt.Sprintf("    %sLife Force%s  %s  %s%3d%%%s",
		MNKGreen, Reset, t.generateMNKBar(data.ContextPercent, 18, lifeColor), lifeColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("    %sNature%s     %s  %s%3d%%%s  %s%s%s",
		MNKBlue, Reset, t.generateMNKBar(100-data.API5hrPercent, 18, MNKBlue),
		MNKBlue, 100-data.API5hrPercent, Reset, MNKGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("    %sCurse%s      %s  %s%3d%%%s  %s%s%s",
		MNKRed, Reset, t.generateMNKBar(data.API7dayPercent, 18, MNKRed),
		MNKRed, data.API7dayPercent, Reset, MNKGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(MNKGreen + "  ──────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	line6 := fmt.Sprintf("    %s%s%s energy  %s%s%s  %s%d%s acts  %s$%s%s  %s$%s/day%s  %s%d%%%s",
		MNKWhite, FormatTokens(data.TokenCount), Reset,
		MNKGray, data.SessionTime, Reset,
		MNKBlue, data.MessageCount, Reset,
		MNKBrown, FormatCost(data.SessionCost), Reset,
		MNKGreen, FormatCost(data.DayCost), Reset,
		MNKWhite, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString(MNKGreen + "  ══════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	return sb.String()
}

func (t *MononokeTheme) generateMNKBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(MNKDark + "〈" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("●", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(MNKDark)
		bar.WriteString(strings.Repeat("○", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(MNKDark + "〉" + Reset)
	return bar.String()
}
