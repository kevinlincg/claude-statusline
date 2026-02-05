package themes

import (
	"fmt"
	"strings"
)

// KikiTheme Kiki's Delivery Service style
type KikiTheme struct{}

func init() {
	RegisterTheme(&KikiTheme{})
}

func (t *KikiTheme) Name() string {
	return "kiki"
}

func (t *KikiTheme) Description() string {
	return "Kiki: Witch delivery service style"
}

const (
	KikiPurple = "\033[38;2;148;0;211m"
	KikiPink   = "\033[38;2;255;182;193m"
	KikiRed    = "\033[38;2;220;20;60m"
	KikiBlue   = "\033[38;2;135;206;250m"
	KikiWhite  = "\033[38;2;255;250;250m"
	KikiGray   = "\033[38;2;128;128;128m"
	KikiYellow = "\033[38;2;255;223;0m"
)

func (t *KikiTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(KikiPurple + "  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" + Reset + "\n")
	sb.WriteString("    " + KikiRed + "🧹" + Reset + " " + KikiPurple + "Witch Delivery Service" + Reset + "   " + KikiPink + "魔女の宅急便" + Reset + "\n")
	sb.WriteString(KikiPurple + "  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	witch := "Kiki"
	if data.ModelType == "Opus" {
		witch = "Ursula"
	} else if data.ModelType == "Haiku" {
		witch = "Jiji"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s✨%s", KikiYellow, Reset)
	}

	line1 := fmt.Sprintf("    %sWitch:%s %s%s%s  %sName:%s %s%s%s  %s%s%s%s",
		KikiPurple, Reset, modelColor, modelIcon, data.ModelName,
		KikiGray, Reset, KikiPink, witch, Reset,
		KikiGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s🌙%s%s", KikiPurple, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", KikiPink, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", KikiRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("    %sDelivery:%s %s%s",
		KikiRed, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(KikiPurple + "  ─────────────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	flyColor := KikiPurple
	if data.ContextPercent > 75 {
		flyColor = KikiRed
	}

	line3 := fmt.Sprintf("    %sFly Power%s   %s  %s%3d%%%s",
		KikiPurple, Reset, t.generateKikiBar(data.ContextPercent, 18, flyColor), flyColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("    %sEnergy%s      %s  %s%3d%%%s  %s%s%s",
		KikiPink, Reset, t.generateKikiBar(100-data.API5hrPercent, 18, KikiPink),
		KikiPink, 100-data.API5hrPercent, Reset, KikiGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("    %sSpirit%s      %s  %s%3d%%%s  %s%s%s",
		KikiBlue, Reset, t.generateKikiBar(100-data.API7dayPercent, 18, KikiBlue),
		KikiBlue, 100-data.API7dayPercent, Reset, KikiGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(KikiPurple + "  ─────────────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	line6 := fmt.Sprintf("    %s%s%s spells  %s%s%s  %s%d%s deliveries  %s$%s%s  %s$%s/day%s  %s%d%%%s",
		KikiWhite, FormatTokens(data.TokenCount), Reset,
		KikiGray, data.SessionTime, Reset,
		KikiPink, data.MessageCount, Reset,
		KikiYellow, FormatCost(data.SessionCost), Reset,
		KikiPurple, FormatCost(data.DayCost), Reset,
		KikiBlue, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString(KikiPurple + "  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" + Reset + "\n")

	return sb.String()
}

func (t *KikiTheme) generateKikiBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(KikiGray + "〔" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("★", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(KikiGray)
		bar.WriteString(strings.Repeat("☆", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(KikiGray + "〕" + Reset)
	return bar.String()
}
