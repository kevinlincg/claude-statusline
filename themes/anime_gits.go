package themes

import (
	"fmt"
	"strings"
)

// GITSTheme Ghost in the Shell cyberbrain style
type GITSTheme struct{}

func init() {
	RegisterTheme(&GITSTheme{})
}

func (t *GITSTheme) Name() string {
	return "gits"
}

func (t *GITSTheme) Description() string {
	return "GITS: Ghost in the Shell cyberbrain interface"
}

const (
	GITSGreen  = "\033[38;2;0;255;128m"
	GITSCyan   = "\033[38;2;0;200;200m"
	GITSBlue   = "\033[38;2;0;150;255m"
	GITSPurple = "\033[38;2;180;0;255m"
	GITSWhite  = "\033[38;2;200;200;200m"
	GITSGray   = "\033[38;2;80;80;80m"
	GITSDark   = "\033[38;2;30;30;30m"
	GITSRed    = "\033[38;2;255;50;50m"
)

func (t *GITSTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(GITSGreen + "╔══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(GITSGreen + "║" + Reset + "  " + GITSCyan + "◈" + GITSWhite + " SECTION 9 " + GITSCyan + "◈" + Reset + "   " + GITSGreen + "公安9課 CYBERBRAIN INTERFACE" + Reset + "                          " + GITSGreen + "║" + Reset + "\n")
	sb.WriteString(GITSGreen + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	clearance := "Level 3"
	if data.ModelType == "Opus" {
		clearance = "Level 9"
	} else if data.ModelType == "Haiku" {
		clearance = "Level 1"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[PATCH READY]%s", GITSCyan, Reset)
	}

	line1 := fmt.Sprintf("  %sAgent:%s %s%s%s  %sClearance:%s %s%s%s  %s%s%s%s",
		GITSCyan, Reset, modelColor, modelIcon, data.ModelName,
		GITSGray, Reset, GITSPurple, clearance, Reset,
		GITSGray, data.Version, Reset, update)

	sb.WriteString(GITSGreen + "║" + Reset)
	sb.WriteString(PadRight(line1, 88))
	sb.WriteString(GITSGreen + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚡%s%s", GITSBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", GITSGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", GITSRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sOperation:%s %s%s",
		GITSBlue, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)

	sb.WriteString(GITSGreen + "║" + Reset)
	sb.WriteString(PadRight(line2, 88))
	sb.WriteString(GITSGreen + "║" + Reset + "\n")

	sb.WriteString(GITSGreen + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	memColor := GITSCyan
	if data.ContextPercent > 75 {
		memColor = GITSRed
	}

	line3 := fmt.Sprintf("  %sMEMORY%s      %s  %s%3d%%%s",
		GITSCyan, Reset, t.generateGITSBar(data.ContextPercent, 18, memColor), memColor, data.ContextPercent, Reset)

	sb.WriteString(GITSGreen + "║" + Reset)
	sb.WriteString(PadRight(line3, 88))
	sb.WriteString(GITSGreen + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sBANDWIDTH%s   %s  %s%3d%%%s  %s%s%s",
		GITSBlue, Reset, t.generateGITSBar(100-data.API5hrPercent, 18, GITSBlue),
		GITSBlue, 100-data.API5hrPercent, Reset, GITSGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(GITSGreen + "║" + Reset)
	sb.WriteString(PadRight(line4, 88))
	sb.WriteString(GITSGreen + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %sGHOST%s       %s  %s%3d%%%s  %s%s%s",
		GITSPurple, Reset, t.generateGITSBar(100-data.API7dayPercent, 18, GITSPurple),
		GITSPurple, 100-data.API7dayPercent, Reset, GITSGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(GITSGreen + "║" + Reset)
	sb.WriteString(PadRight(line5, 88))
	sb.WriteString(GITSGreen + "║" + Reset + "\n")

	sb.WriteString(GITSGreen + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sData:%s %s%s%s  %sUptime:%s %s  %sOps:%s %s%d%s  %sCost:%s %s%s%s  %sDaily:%s %s%s%s",
		GITSWhite, Reset, GITSWhite, FormatTokens(data.TokenCount), Reset,
		GITSGray, Reset, data.SessionTime,
		GITSGray, Reset, GITSCyan, data.MessageCount, Reset,
		GITSGreen, Reset, GITSGreen, FormatCost(data.SessionCost), Reset,
		GITSBlue, Reset, GITSBlue, FormatCost(data.DayCost), Reset)

	sb.WriteString(GITSGreen + "║" + Reset)
	sb.WriteString(PadRight(line6, 88))
	sb.WriteString(GITSGreen + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sRate:%s %s%s/h%s  %sSync:%s %s%d%%%s",
		GITSRed, Reset, GITSRed, FormatCost(data.BurnRate), Reset,
		GITSPurple, Reset, GITSPurple, data.CacheHitRate, Reset)

	sb.WriteString(GITSGreen + "║" + Reset)
	sb.WriteString(PadRight(line7, 88))
	sb.WriteString(GITSGreen + "║" + Reset + "\n")

	sb.WriteString(GITSGreen + "╚══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *GITSTheme) generateGITSBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(GITSDark + "〈" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(GITSDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(GITSDark + "〉" + Reset)
	return bar.String()
}
