package themes

import (
	"fmt"
	"strings"
)

// SpyFamilyTheme Spy x Family WISE operation style
type SpyFamilyTheme struct{}

func init() {
	RegisterTheme(&SpyFamilyTheme{})
}

func (t *SpyFamilyTheme) Name() string {
	return "spyfamily"
}

func (t *SpyFamilyTheme) Description() string {
	return "Spy x Family: WISE operation mission briefing style"
}

const (
	SPYBlack  = "\033[38;2;30;30;30m"
	SPYRed    = "\033[38;2;200;50;70m"
	SPYPink   = "\033[38;2;255;182;193m"
	SPYGold   = "\033[38;2;218;165;32m"
	SPYGreen  = "\033[38;2;100;180;100m"
	SPYWhite  = "\033[38;2;245;245;245m"
	SPYGray   = "\033[38;2;128;128;128m"
)

func (t *SpyFamilyTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(SPYRed + "┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓" + Reset + "\n")
	sb.WriteString(SPYRed + "┃" + Reset + "  " + SPYRed + "🎯" + SPYWhite + " OPERATION STRIX " + SPYRed + "🎯" + Reset + "   " + SPYPink + "スパイファミリー" + Reset + "   " + SPYGold + "[CLASSIFIED]" + Reset + "                  " + SPYRed + "┃" + Reset + "\n")
	sb.WriteString(SPYRed + "┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	agent := "Anya"
	if data.ModelType == "Opus" {
		agent = "Loid"
	} else if data.ModelType == "Haiku" {
		agent = "Bond"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[New Intel]%s", SPYGold, Reset)
	}

	line1 := fmt.Sprintf("  %sAgent:%s %s%s%s  %sCodename:%s %s%s%s  %s%s%s%s",
		SPYRed, Reset, modelColor, modelIcon, data.ModelName,
		SPYGray, Reset, SPYPink, agent, Reset,
		SPYGray, data.Version, Reset, update)

	sb.WriteString(SPYRed + "┃" + Reset)
	sb.WriteString(PadRight(line1, 87))
	sb.WriteString(SPYRed + "┃" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s📁%s%s", SPYGold, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", SPYGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", SPYRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sMission:%s %s%s",
		SPYGold, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)

	sb.WriteString(SPYRed + "┃" + Reset)
	sb.WriteString(PadRight(line2, 87))
	sb.WriteString(SPYRed + "┃" + Reset + "\n")

	sb.WriteString(SPYRed + "┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫" + Reset + "\n")

	telepathyColor := SPYPink
	if data.ContextPercent > 75 {
		telepathyColor = SPYRed
	}

	line3 := fmt.Sprintf("  %sTelepathy%s   %s  %s%3d%%%s",
		SPYPink, Reset, t.generateSPYBar(data.ContextPercent, 18, telepathyColor), telepathyColor, data.ContextPercent, Reset)

	sb.WriteString(SPYRed + "┃" + Reset)
	sb.WriteString(PadRight(line3, 87))
	sb.WriteString(SPYRed + "┃" + Reset + "\n")

	line4 := fmt.Sprintf("  %sCover%s       %s  %s%3d%%%s  %s%s%s",
		SPYGreen, Reset, t.generateSPYBar(100-data.API5hrPercent, 18, SPYGreen),
		SPYGreen, 100-data.API5hrPercent, Reset, SPYGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(SPYRed + "┃" + Reset)
	sb.WriteString(PadRight(line4, 87))
	sb.WriteString(SPYRed + "┃" + Reset + "\n")

	line5 := fmt.Sprintf("  %sNetwork%s     %s  %s%3d%%%s  %s%s%s",
		SPYGold, Reset, t.generateSPYBar(100-data.API7dayPercent, 18, SPYGold),
		SPYGold, 100-data.API7dayPercent, Reset, SPYGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(SPYRed + "┃" + Reset)
	sb.WriteString(PadRight(line5, 87))
	sb.WriteString(SPYRed + "┃" + Reset + "\n")

	sb.WriteString(SPYRed + "┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫" + Reset + "\n")

	line6 := fmt.Sprintf("  %sIntel:%s %s%s%s  %sTime:%s %s  %sOps:%s %s%d%s  %sBudget:%s %s$%s%s",
		SPYWhite, Reset, SPYWhite, FormatTokens(data.TokenCount), Reset,
		SPYGray, Reset, data.SessionTime,
		SPYGray, Reset, SPYPink, data.MessageCount, Reset,
		SPYGold, Reset, SPYGold, FormatCost(data.SessionCost), Reset)

	sb.WriteString(SPYRed + "┃" + Reset)
	sb.WriteString(PadRight(line6, 87))
	sb.WriteString(SPYRed + "┃" + Reset + "\n")

	line7 := fmt.Sprintf("  %sDaily:%s %s$%s%s  %sRate:%s %s$%s/h%s  %sWaku:%s %s%d%%%s",
		SPYGreen, Reset, SPYGreen, FormatCost(data.DayCost), Reset,
		SPYRed, Reset, SPYRed, FormatCost(data.BurnRate), Reset,
		SPYPink, Reset, SPYPink, data.CacheHitRate, Reset)

	sb.WriteString(SPYRed + "┃" + Reset)
	sb.WriteString(PadRight(line7, 87))
	sb.WriteString(SPYRed + "┃" + Reset + "\n")

	sb.WriteString(SPYRed + "┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛" + Reset + "\n")

	return sb.String()
}

func (t *SpyFamilyTheme) generateSPYBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(SPYGray + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(SPYBlack)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(SPYGray + "]" + Reset)
	return bar.String()
}
