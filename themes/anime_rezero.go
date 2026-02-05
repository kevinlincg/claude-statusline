package themes

import (
	"fmt"
	"strings"
)

// ReZeroTheme Re:Zero return by death style
type ReZeroTheme struct{}

func init() {
	RegisterTheme(&ReZeroTheme{})
}

func (t *ReZeroTheme) Name() string {
	return "rezero"
}

func (t *ReZeroTheme) Description() string {
	return "Re:Zero: Return by Death checkpoint style"
}

const (
	RZPurple = "\033[38;2;128;0;128m"
	RZBlue   = "\033[38;2;100;149;237m"
	RZSilver = "\033[38;2;192;192;192m"
	RZRed    = "\033[38;2;178;34;34m"
	RZWhite  = "\033[38;2;245;245;245m"
	RZGray   = "\033[38;2;105;105;105m"
	RZDark   = "\033[38;2;30;30;40m"
)

func (t *ReZeroTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(RZPurple + "════════════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")
	sb.WriteString("  " + RZRed + "☠" + RZWhite + " Re:ZERO " + RZRed + "☠" + Reset + "   " + RZPurple + "リゼロ" + Reset + "   " + RZSilver + "「Return by Death」" + Reset + "\n")
	sb.WriteString(RZPurple + "════════════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	character := "Rem"
	if data.ModelType == "Opus" {
		character = "Subaru"
	} else if data.ModelType == "Haiku" {
		character = "Puck"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[Checkpoint]%s", RZPurple, Reset)
	}

	line1 := fmt.Sprintf("  %sPlayer:%s %s%s%s  %sSpirit:%s %s%s%s  %s%s%s%s",
		RZPurple, Reset, modelColor, modelIcon, data.ModelName,
		RZGray, Reset, RZSilver, character, Reset,
		RZGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s◈%s%s", RZPurple, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", RZBlue, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", RZRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sSave Point:%s %s%s",
		RZBlue, Reset, ShortenPath(data.ProjectPath, 38), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(RZPurple + "────────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	witchColor := RZPurple
	if data.ContextPercent > 75 {
		witchColor = RZRed
	}

	line3 := fmt.Sprintf("  %sWitch Miasma%s  %s  %s%3d%%%s",
		RZPurple, Reset, t.generateRZBar(data.ContextPercent, 18, witchColor), witchColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("  %sMana%s          %s  %s%3d%%%s  %s%s%s",
		RZBlue, Reset, t.generateRZBar(100-data.API5hrPercent, 18, RZBlue),
		RZBlue, 100-data.API5hrPercent, Reset, RZGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("  %sDeaths%s        %s  %s%3d%%%s  %s%s%s",
		RZRed, Reset, t.generateRZBar(data.API7dayPercent, 18, RZRed),
		RZRed, data.API7dayPercent, Reset, RZGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(RZPurple + "────────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	line6 := fmt.Sprintf("  %s%s%s memories  %s%s%s  %s%d%s loops  %s$%s%s  %s$%s/day%s  %s%d%%%s",
		RZWhite, FormatTokens(data.TokenCount), Reset,
		RZGray, data.SessionTime, Reset,
		RZPurple, data.MessageCount, Reset,
		RZSilver, FormatCost(data.SessionCost), Reset,
		RZBlue, FormatCost(data.DayCost), Reset,
		RZPurple, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString(RZPurple + "════════════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	return sb.String()
}

func (t *ReZeroTheme) generateRZBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(RZDark + "【" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("◆", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(RZDark)
		bar.WriteString(strings.Repeat("◇", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(RZDark + "】" + Reset)
	return bar.String()
}
