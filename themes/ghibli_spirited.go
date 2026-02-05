package themes

import (
	"fmt"
	"strings"
)

// SpiritedTheme Spirited Away bathhouse style
type SpiritedTheme struct{}

func init() {
	RegisterTheme(&SpiritedTheme{})
}

func (t *SpiritedTheme) Name() string {
	return "spirited"
}

func (t *SpiritedTheme) Description() string {
	return "Spirited Away: Bathhouse mysterious style"
}

const (
	SPPurple = "\033[38;2;128;0;128m"
	SPGold   = "\033[38;2;218;165;32m"
	SPRed    = "\033[38;2;139;0;0m"
	SPBlue   = "\033[38;2;70;130;180m"
	SPCream  = "\033[38;2;255;248;220m"
	SPGray   = "\033[38;2;100;100;100m"
	SPDark   = "\033[38;2;40;30;50m"
)

func (t *SpiritedTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(SPPurple + "  ═══════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")
	sb.WriteString("    " + SPGold + "油屋" + Reset + "  " + SPCream + "Aburaya Bathhouse" + Reset + "   " + SPPurple + "千と千尋の神隠し" + Reset + "\n")
	sb.WriteString(SPPurple + "  ═══════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	worker := "Sen"
	if data.ModelType == "Opus" {
		worker = "Yubaba"
	} else if data.ModelType == "Haiku" {
		worker = "Lin"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s✨%s", SPGold, Reset)
	}

	line1 := fmt.Sprintf("    %sWorker:%s %s%s%s  %sName:%s %s%s%s  %s%s%s%s",
		SPRed, Reset, modelColor, modelIcon, data.ModelName,
		SPGray, Reset, SPPurple, worker, Reset,
		SPGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s🌊%s%s", SPBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", SPGold, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", SPRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("    %sTask:%s %s%s",
		SPBlue, Reset, ShortenPath(data.ProjectPath, 42), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(SPPurple + "  ─────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	spiritColor := SPPurple
	if data.ContextPercent > 75 {
		spiritColor = SPRed
	}

	line3 := fmt.Sprintf("    %sSpirit%s     %s  %s%3d%%%s",
		SPPurple, Reset, t.generateSPBar(data.ContextPercent, 18, spiritColor), spiritColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("    %sBath Water%s %s  %s%3d%%%s  %s%s%s",
		SPBlue, Reset, t.generateSPBar(100-data.API5hrPercent, 18, SPBlue),
		SPBlue, 100-data.API5hrPercent, Reset, SPGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("    %sGold%s       %s  %s%3d%%%s  %s%s%s",
		SPGold, Reset, t.generateSPBar(100-data.API7dayPercent, 18, SPGold),
		SPGold, 100-data.API7dayPercent, Reset, SPGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(SPPurple + "  ─────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	line6 := fmt.Sprintf("    %s%s%s work  %s%s%s  %s%d%s guests  %s$%s%s  %s$%s/day%s  %s%d%%%s luck",
		SPCream, FormatTokens(data.TokenCount), Reset,
		SPGray, data.SessionTime, Reset,
		SPBlue, data.MessageCount, Reset,
		SPGold, FormatCost(data.SessionCost), Reset,
		SPPurple, FormatCost(data.DayCost), Reset,
		SPGold, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString(SPPurple + "  ═══════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	return sb.String()
}

func (t *SpiritedTheme) generateSPBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(SPDark + "〔" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("◆", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(SPDark)
		bar.WriteString(strings.Repeat("◇", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(SPDark + "〕" + Reset)
	return bar.String()
}
