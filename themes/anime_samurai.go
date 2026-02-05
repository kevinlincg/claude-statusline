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
	SMRRed    = "\033[38;2;180;50;50m"
	SMRGold   = "\033[38;2;200;160;80m"
	SMRBlack  = "\033[38;2;30;30;30m"
	SMRWhite  = "\033[38;2;245;240;230m"
	SMRGray   = "\033[38;2;120;110;100m"
	SMRInk    = "\033[38;2;50;50;60m"
)

func (t *SamuraiTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Traditional scroll top
	sb.WriteString(SMRGold + "                    ╭━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╮" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┏━━━━━━━━━━━━━━━┫" + Reset + "                                           " + SMRGold + "┣━━━━━━━━━━━━━━━┓" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┃" + Reset + "                                                                              " + SMRGold + "┃" + Reset + "\n")

	// Title with brush style
	sb.WriteString(SMRGold + "    ┃" + Reset + "                        " + SMRRed + "武" + SMRBlack + " 士 " + SMRRed + "道" + Reset + "   " + SMRInk + "侍" + Reset + "                               " + SMRGold + "┃" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┃" + Reset + "                      " + SMRGray + "━━ SAMURAI ━━" + Reset + "                                " + SMRGold + "┃" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┃" + Reset + "                                                                              " + SMRGold + "┃" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	rank := "Ronin"
	if data.ModelType == "Opus" {
		rank = "Shogun"
	} else if data.ModelType == "Haiku" {
		rank = "Ashigaru"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s【新】%s", SMRRed, Reset)
	}

	line1 := fmt.Sprintf("    %s刀:%s %s%s%s    %s位:%s %s%s%s    %s%s%s%s",
		SMRInk, Reset, modelColor, modelIcon, data.ModelName,
		SMRInk, Reset, SMRGold, rank, Reset,
		SMRGray, data.Version, Reset, update)

	sb.WriteString(SMRGold + "    ┃" + Reset + "  ")
	sb.WriteString(PadRight(line1, 72))
	sb.WriteString(SMRGold + "┃" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚔%s%s", SMRInk, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", SMRGold, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", SMRRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("    %s道:%s %s%s",
		SMRInk, Reset, ShortenPath(data.ProjectPath, 45), gitInfo)

	sb.WriteString(SMRGold + "    ┃" + Reset + "  ")
	sb.WriteString(PadRight(line2, 72))
	sb.WriteString(SMRGold + "┃" + Reset + "\n")

	sb.WriteString(SMRGold + "    ┃" + Reset + "                                                                              " + SMRGold + "┃" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┃" + Reset + "  " + SMRGray + "┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄" + Reset + "  " + SMRGold + "┃" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┃" + Reset + "                                                                              " + SMRGold + "┃" + Reset + "\n")

	// Stats with brush-style bars
	kiColor := SMRGold
	if data.ContextPercent > 75 {
		kiColor = SMRRed
	}

	line3 := fmt.Sprintf("        %s気%s  %s  %s%3d%%%s",
		SMRRed, Reset, t.generateSMRBar(data.ContextPercent, 20, kiColor), kiColor, data.ContextPercent, Reset)

	sb.WriteString(SMRGold + "    ┃" + Reset)
	sb.WriteString(PadRight(line3, 74))
	sb.WriteString(SMRGold + "┃" + Reset + "\n")

	line4 := fmt.Sprintf("        %s力%s  %s  %s%3d%%%s  %s%s%s",
		SMRGold, Reset, t.generateSMRBar(100-data.API5hrPercent, 20, SMRGold),
		SMRGold, 100-data.API5hrPercent, Reset, SMRGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(SMRGold + "    ┃" + Reset)
	sb.WriteString(PadRight(line4, 74))
	sb.WriteString(SMRGold + "┃" + Reset + "\n")

	line5 := fmt.Sprintf("        %s魂%s  %s  %s%3d%%%s  %s%s%s",
		SMRInk, Reset, t.generateSMRBar(100-data.API7dayPercent, 20, SMRInk),
		SMRInk, 100-data.API7dayPercent, Reset, SMRGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(SMRGold + "    ┃" + Reset)
	sb.WriteString(PadRight(line5, 74))
	sb.WriteString(SMRGold + "┃" + Reset + "\n")

	sb.WriteString(SMRGold + "    ┃" + Reset + "                                                                              " + SMRGold + "┃" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┃" + Reset + "  " + SMRGray + "┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄" + Reset + "  " + SMRGold + "┃" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┃" + Reset + "                                                                              " + SMRGold + "┃" + Reset + "\n")

	line6 := fmt.Sprintf("        %s%s%s 文字  %s%s%s  %s%d%s 斬  %s金%s%s  %s%d%%%s 効",
		SMRWhite, FormatTokens(data.TokenCount), Reset,
		SMRGray, data.SessionTime, Reset,
		SMRInk, data.MessageCount, Reset,
		SMRGold, FormatCost(data.SessionCost), Reset,
		SMRRed, data.CacheHitRate, Reset)

	sb.WriteString(SMRGold + "    ┃" + Reset)
	sb.WriteString(PadRight(line6, 74))
	sb.WriteString(SMRGold + "┃" + Reset + "\n")

	sb.WriteString(SMRGold + "    ┃" + Reset + "                                                                              " + SMRGold + "┃" + Reset + "\n")
	sb.WriteString(SMRGold + "    ┗━━━━━━━━━━━━━━━┫" + Reset + "                                           " + SMRGold + "┣━━━━━━━━━━━━━━━┛" + Reset + "\n")
	sb.WriteString(SMRGold + "                    ╰━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╯" + Reset + "\n")

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
