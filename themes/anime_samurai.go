package themes

import (
	"fmt"
	"strings"
)

// SamuraiTheme Traditional Japanese samurai brush style
type SamuraiTheme struct{}

func init() {
	RegisterTheme(&SamuraiTheme{})
}

func (t *SamuraiTheme) Name() string {
	return "samurai"
}

func (t *SamuraiTheme) Description() string {
	return "Samurai: Traditional Japanese brush calligraphy style"
}

const (
	SMRRed   = "\033[38;2;180;50;50m"
	SMRGold  = "\033[38;2;200;160;80m"
	SMRBlack = "\033[38;2;30;30;30m"
	SMRWhite = "\033[38;2;245;240;230m"
	SMRGray  = "\033[38;2;120;110;100m"
	SMRInk   = "\033[38;2;50;50;60m"
)

func (t *SamuraiTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Compact scroll header
	sb.WriteString(SMRGold + "    ╭━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╮" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┃" + Reset + "        " + SMRRed + "武" + SMRInk + " 士 " + SMRRed + "道" + Reset + "   " + SMRGray + "━━ SAMURAI ━━" + Reset + "   " + SMRInk + "侍" + Reset + "                            " + SMRGold + "┃" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	rank := "浪人"
	if data.ModelType == "Opus" {
		rank = "将軍"
	} else if data.ModelType == "Haiku" {
		rank = "足軽"
	}

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf(" %s⚔%s%s", SMRInk, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", SMRGold, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", SMRRed, data.GitDirty, Reset)
		}
	}

	line1 := fmt.Sprintf("  %s刀:%s %s%s%s  %s位:%s %s%s%s  %s道:%s %s%s",
		SMRInk, Reset, modelColor, modelIcon, data.ModelName,
		SMRInk, Reset, SMRGold, rank, Reset,
		SMRInk, Reset, ShortenPath(data.ProjectPath, 28), gitInfo)

	sb.WriteString(SMRGold + "    ┃" + Reset)
	sb.WriteString(PadRight(line1, 75))
	sb.WriteString(SMRGold + "┃" + Reset + "\n")

	sb.WriteString(SMRGold + "    ┣" + SMRGray + "┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄" + SMRGold + "┫" + Reset + "\n")

	// Stats with brush-style bars
	kiColor := SMRGold
	if data.ContextPercent > 75 {
		kiColor = SMRRed
	}

	line2 := fmt.Sprintf("  %s気%s %s %s%3d%%%s  %s力%s %s %s%3d%%%s %s%s%s  %s魂%s %s %s%3d%%%s %s%s%s",
		SMRRed, Reset, t.generateSMRBar(data.ContextPercent, 12, kiColor), kiColor, data.ContextPercent, Reset,
		SMRGold, Reset, t.generateSMRBar(100-data.API5hrPercent, 10, SMRGold), SMRGold, 100-data.API5hrPercent, Reset, SMRGray, data.API5hrTimeLeft, Reset,
		SMRInk, Reset, t.generateSMRBar(100-data.API7dayPercent, 10, SMRInk), SMRInk, 100-data.API7dayPercent, Reset, SMRGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(SMRGold + "    ┃" + Reset)
	sb.WriteString(PadRight(line2, 75))
	sb.WriteString(SMRGold + "┃" + Reset + "\n")

	sb.WriteString(SMRGold + "    ┣" + SMRGray + "┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄" + SMRGold + "┫" + Reset + "\n")

	line3 := fmt.Sprintf("  %s%s%s 文字  %s%s%s  %s%d%s 斬  %s金%s%s  %s$%s/日%s  %s%d%%%s 効  %s%s%s",
		SMRWhite, FormatTokens(data.TokenCount), Reset,
		SMRGray, data.SessionTime, Reset,
		SMRInk, data.MessageCount, Reset,
		SMRGold, FormatCost(data.SessionCost), Reset,
		SMRRed, FormatCost(data.DayCost), Reset,
		SMRGold, data.CacheHitRate, Reset,
		SMRGray, data.Version, Reset)

	sb.WriteString(SMRGold + "    ┃" + Reset)
	sb.WriteString(PadRight(line3, 75))
	sb.WriteString(SMRGold + "┃" + Reset + "\n")

	sb.WriteString(SMRGold + "    ╰━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╯" + Reset + "\n")

	return sb.String()
}

func (t *SamuraiTheme) generateSMRBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(SMRGray + "〘" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("━", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(SMRGray)
		bar.WriteString(strings.Repeat("─", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(SMRGray + "〙" + Reset)
	return bar.String()
}
