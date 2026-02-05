package themes

import (
	"fmt"
	"strings"
)

// AOTTheme Attack on Titan Survey Corps style
type AOTTheme struct{}

func init() {
	RegisterTheme(&AOTTheme{})
}

func (t *AOTTheme) Name() string {
	return "aot"
}

func (t *AOTTheme) Description() string {
	return "AOT: Survey Corps military report style"
}

const (
	AOTGreen     = "\033[38;2;0;100;0m"
	AOTBrown     = "\033[38;2;101;67;33m"
	AOTWhite     = "\033[38;2;230;230;230m"
	AOTRed       = "\033[38;2;139;0;0m"
	AOTBlue      = "\033[38;2;70;130;180m"
	AOTGold      = "\033[38;2;184;134;11m"
	AOTGray      = "\033[38;2;105;105;105m"
	AOTDark      = "\033[38;2;50;50;50m"
)

func (t *AOTTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Military report header
	sb.WriteString(AOTGreen + "╔══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(AOTGreen + "║" + Reset + "  " + AOTWhite + "█▀▀ █ █ █▀▄ █ █ █▀▀ █ █   █▀▀ █▀█ █▀▄ █▀█ █▀▀" + Reset + "   " + AOTGreen + "◇ SURVEY CORPS ◇" + Reset + "              " + AOTGreen + "║" + Reset + "\n")
	sb.WriteString(AOTGreen + "║" + Reset + "  " + AOTWhite + "▀▀█ █ █ █▀▄ ▀▄▀ █▀▀  █    █   █ █ █▀▄ █▀▀ ▀▀█" + Reset + "   " + AOTGold + "自由の翼" + Reset + "                       " + AOTGreen + "║" + Reset + "\n")
	sb.WriteString(AOTGreen + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Unit info
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	rank := "Cadet"
	if data.ModelType == "Opus" {
		rank = "Commander"
	} else if data.ModelType == "Sonnet" {
		rank = "Captain"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[NEW ORDERS]%s", AOTRed, Reset)
	}

	line1 := fmt.Sprintf("  %sUNIT:%s %s%s%s  %sRANK:%s %s%s%s  %sVER:%s %s%s",
		AOTBrown, Reset, modelColor, modelIcon, data.ModelName,
		AOTBrown, Reset, AOTGold, rank, Reset,
		AOTGray, Reset, data.Version, update)

	sb.WriteString(AOTGreen + "║" + Reset)
	sb.WriteString(PadRight(line1, 88))
	sb.WriteString(AOTGreen + "║" + Reset + "\n")

	// Mission
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %sBranch:%s %s", AOTBlue, Reset, data.GitBranch)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", AOTGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s!%d%s", AOTRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sMISSION:%s %s%s",
		AOTRed, Reset, ShortenPath(data.ProjectPath, 35), gitInfo)

	sb.WriteString(AOTGreen + "║" + Reset)
	sb.WriteString(PadRight(line2, 88))
	sb.WriteString(AOTGreen + "║" + Reset + "\n")

	sb.WriteString(AOTGreen + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// ODM Gas (Context)
	gasColor := AOTBlue
	if data.ContextPercent > 75 {
		gasColor = AOTRed
	} else if data.ContextPercent > 50 {
		gasColor = AOTGold
	}

	line3 := fmt.Sprintf("  %sODM GAS%s    %s  %s%3d%%%s",
		AOTBlue, Reset,
		t.generateAOTBar(100-data.ContextPercent, 20, gasColor),
		gasColor, 100-data.ContextPercent, Reset)

	sb.WriteString(AOTGreen + "║" + Reset)
	sb.WriteString(PadRight(line3, 88))
	sb.WriteString(AOTGreen + "║" + Reset + "\n")

	// Blade durability (5hr)
	line4 := fmt.Sprintf("  %sBLADES%s     %s  %s%3d%%%s  %sResupply: %s%s",
		AOTWhite, Reset,
		t.generateAOTBar(100-data.API5hrPercent, 20, AOTWhite),
		AOTWhite, 100-data.API5hrPercent, Reset,
		AOTGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(AOTGreen + "║" + Reset)
	sb.WriteString(PadRight(line4, 88))
	sb.WriteString(AOTGreen + "║" + Reset + "\n")

	// Wall integrity (7day)
	wallColor := AOTGreen
	if data.API7dayPercent > 75 {
		wallColor = AOTRed
	}

	line5 := fmt.Sprintf("  %sWALL HP%s    %s  %s%3d%%%s  %sRepair: %s%s",
		AOTGreen, Reset,
		t.generateAOTBar(100-data.API7dayPercent, 20, wallColor),
		wallColor, 100-data.API7dayPercent, Reset,
		AOTGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(AOTGreen + "║" + Reset)
	sb.WriteString(PadRight(line5, 88))
	sb.WriteString(AOTGreen + "║" + Reset + "\n")

	sb.WriteString(AOTGreen + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Combat stats
	line6 := fmt.Sprintf("  %sKills:%s %s%s%s  %sTime:%s %s  %sEngagements:%s %s%d%s  %sCost:%s %s%s%s",
		AOTRed, Reset, AOTRed, FormatTokens(data.TokenCount), Reset,
		AOTGray, Reset, data.SessionTime,
		AOTGray, Reset, AOTWhite, data.MessageCount, Reset,
		AOTGold, Reset, AOTGold, FormatCost(data.SessionCost), Reset)

	sb.WriteString(AOTGreen + "║" + Reset)
	sb.WriteString(PadRight(line6, 88))
	sb.WriteString(AOTGreen + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sHit%%:%s %s%d%%%s",
		AOTBrown, Reset, AOTBrown, FormatCost(data.DayCost), Reset,
		AOTRed, Reset, AOTRed, FormatCost(data.BurnRate), Reset,
		AOTBlue, Reset, AOTBlue, data.CacheHitRate, Reset)

	sb.WriteString(AOTGreen + "║" + Reset)
	sb.WriteString(PadRight(line7, 88))
	sb.WriteString(AOTGreen + "║" + Reset + "\n")

	sb.WriteString(AOTGreen + "╚══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *AOTTheme) generateAOTBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(AOTDark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(AOTDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(AOTDark + "]" + Reset)
	return bar.String()
}
