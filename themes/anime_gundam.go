package themes

import (
	"fmt"
	"strings"
)

// GundamTheme Gundam mobile suit cockpit style
type GundamTheme struct{}

func init() {
	RegisterTheme(&GundamTheme{})
}

func (t *GundamTheme) Name() string {
	return "gundam"
}

func (t *GundamTheme) Description() string {
	return "Gundam: Mobile Suit cockpit interface"
}

const (
	GundamRed    = "\033[38;2;200;0;0m"
	GundamBlue   = "\033[38;2;0;100;200m"
	GundamWhite  = "\033[38;2;240;240;240m"
	GundamYellow = "\033[38;2;255;200;0m"
	GundamGreen  = "\033[38;2;0;200;100m"
	GundamGray   = "\033[38;2;100;100;100m"
	GundamDark   = "\033[38;2;40;40;40m"
)

func (t *GundamTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(GundamBlue + "╔══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(GundamBlue + "║" + Reset + "  " + GundamWhite + "◆ MOBILE SUIT SYSTEM ◆" + Reset + "   " + GundamYellow + "E.F.S.F." + Reset + "                                          " + GundamBlue + "║" + Reset + "\n")
	sb.WriteString(GundamBlue + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	msType := "GM"
	if data.ModelType == "Opus" {
		msType = "RX-78-2"
	} else if data.ModelType == "Haiku" {
		msType = "Ball"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[UPGRADE]%s", GundamGreen, Reset)
	}

	line1 := fmt.Sprintf("  %sUnit:%s %s%s%s  %sType:%s %s%s%s  %s%s%s%s",
		GundamYellow, Reset, modelColor, modelIcon, data.ModelName,
		GundamGray, Reset, GundamWhite, msType, Reset,
		GundamGray, data.Version, Reset, update)

	sb.WriteString(GundamBlue + "║" + Reset)
	sb.WriteString(PadRight(line1, 88))
	sb.WriteString(GundamBlue + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s◈%s%s", GundamBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", GundamGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s!%d%s", GundamRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sMission:%s %s%s",
		GundamRed, Reset, ShortenPath(data.ProjectPath, 42), gitInfo)

	sb.WriteString(GundamBlue + "║" + Reset)
	sb.WriteString(PadRight(line2, 88))
	sb.WriteString(GundamBlue + "║" + Reset + "\n")

	sb.WriteString(GundamBlue + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	reactorColor := GundamGreen
	if data.ContextPercent > 75 {
		reactorColor = GundamRed
	} else if data.ContextPercent > 50 {
		reactorColor = GundamYellow
	}

	line3 := fmt.Sprintf("  %sREACTOR%s     %s  %s%3d%%%s",
		GundamGreen, Reset, t.generateGundamBar(data.ContextPercent, 18, reactorColor), reactorColor, data.ContextPercent, Reset)

	sb.WriteString(GundamBlue + "║" + Reset)
	sb.WriteString(PadRight(line3, 88))
	sb.WriteString(GundamBlue + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sAMMO%s        %s  %s%3d%%%s  %s%s%s",
		GundamYellow, Reset, t.generateGundamBar(100-data.API5hrPercent, 18, GundamYellow),
		GundamYellow, 100-data.API5hrPercent, Reset, GundamGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(GundamBlue + "║" + Reset)
	sb.WriteString(PadRight(line4, 88))
	sb.WriteString(GundamBlue + "║" + Reset + "\n")

	armorColor := GundamBlue
	if data.API7dayPercent > 75 {
		armorColor = GundamRed
	}

	line5 := fmt.Sprintf("  %sARMOR%s       %s  %s%3d%%%s  %s%s%s",
		GundamBlue, Reset, t.generateGundamBar(100-data.API7dayPercent, 18, armorColor),
		armorColor, 100-data.API7dayPercent, Reset, GundamGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(GundamBlue + "║" + Reset)
	sb.WriteString(PadRight(line5, 88))
	sb.WriteString(GundamBlue + "║" + Reset + "\n")

	sb.WriteString(GundamBlue + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sOutput:%s %s%s%s  %sTime:%s %s  %sSorties:%s %s%d%s  %sCost:%s %s%s%s",
		GundamWhite, Reset, GundamWhite, FormatTokens(data.TokenCount), Reset,
		GundamGray, Reset, data.SessionTime,
		GundamGray, Reset, GundamYellow, data.MessageCount, Reset,
		GundamGreen, Reset, GundamGreen, FormatCost(data.SessionCost), Reset)

	sb.WriteString(GundamBlue + "║" + Reset)
	sb.WriteString(PadRight(line6, 88))
	sb.WriteString(GundamBlue + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sAccuracy:%s %s%d%%%s",
		GundamBlue, Reset, GundamBlue, FormatCost(data.DayCost), Reset,
		GundamRed, Reset, GundamRed, FormatCost(data.BurnRate), Reset,
		GundamGreen, Reset, GundamGreen, data.CacheHitRate, Reset)

	sb.WriteString(GundamBlue + "║" + Reset)
	sb.WriteString(PadRight(line7, 88))
	sb.WriteString(GundamBlue + "║" + Reset + "\n")

	sb.WriteString(GundamBlue + "╚══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *GundamTheme) generateGundamBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(GundamDark + "〔" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▰", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(GundamDark)
		bar.WriteString(strings.Repeat("▱", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(GundamDark + "〕" + Reset)
	return bar.String()
}
