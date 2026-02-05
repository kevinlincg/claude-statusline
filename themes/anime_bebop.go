package themes

import (
	"fmt"
	"strings"
)

// BebopTheme Cowboy Bebop space jazz style
type BebopTheme struct{}

func init() {
	RegisterTheme(&BebopTheme{})
}

func (t *BebopTheme) Name() string {
	return "bebop"
}

func (t *BebopTheme) Description() string {
	return "Bebop: Cowboy Bebop space bounty hunter style"
}

const (
	BebopOrange = "\033[38;2;255;140;0m"
	BebopRed    = "\033[38;2;178;34;34m"
	BebopYellow = "\033[38;2;255;215;0m"
	BebopBlue   = "\033[38;2;70;130;180m"
	BebopGreen  = "\033[38;2;60;179;113m"
	BebopWhite  = "\033[38;2;240;240;240m"
	BebopGray   = "\033[38;2;100;100;100m"
	BebopDark   = "\033[38;2;40;40;40m"
)

func (t *BebopTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(BebopOrange + "╔═══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(BebopOrange + "║" + Reset + "  " + BebopYellow + "★" + BebopWhite + " BEBOP CREW " + BebopYellow + "★" + Reset + "   " + BebopOrange + "Bounty Hunter Database" + Reset + "                                    " + BebopOrange + "║" + Reset + "\n")
	sb.WriteString(BebopOrange + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	crew := "Spike"
	if data.ModelType == "Opus" {
		crew = "Vicious"
	} else if data.ModelType == "Haiku" {
		crew = "Ed"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[NEW BOUNTY]%s", BebopYellow, Reset)
	}

	line1 := fmt.Sprintf("  %sHunter:%s %s%s%s  %sAlias:%s %s%s%s  %s%s%s%s",
		BebopRed, Reset, modelColor, modelIcon, data.ModelName,
		BebopGray, Reset, BebopOrange, crew, Reset,
		BebopGray, data.Version, Reset, update)

	sb.WriteString(BebopOrange + "║" + Reset)
	sb.WriteString(PadRight(line1, 89))
	sb.WriteString(BebopOrange + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚡%s%s", BebopBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", BebopGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", BebopRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sTarget:%s %s%s",
		BebopYellow, Reset, ShortenPath(data.ProjectPath, 45), gitInfo)

	sb.WriteString(BebopOrange + "║" + Reset)
	sb.WriteString(PadRight(line2, 89))
	sb.WriteString(BebopOrange + "║" + Reset + "\n")

	sb.WriteString(BebopOrange + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	fuelColor := BebopBlue
	if data.ContextPercent > 75 {
		fuelColor = BebopRed
	}

	line3 := fmt.Sprintf("  %sFuel%s       %s  %s%3d%%%s",
		BebopBlue, Reset, t.generateBebopBar(data.ContextPercent, 18, fuelColor), fuelColor, data.ContextPercent, Reset)

	sb.WriteString(BebopOrange + "║" + Reset)
	sb.WriteString(PadRight(line3, 89))
	sb.WriteString(BebopOrange + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sAmmo%s       %s  %s%3d%%%s  %s%s%s",
		BebopGreen, Reset, t.generateBebopBar(100-data.API5hrPercent, 18, BebopGreen),
		BebopGreen, 100-data.API5hrPercent, Reset, BebopGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(BebopOrange + "║" + Reset)
	sb.WriteString(PadRight(line4, 89))
	sb.WriteString(BebopOrange + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %sHull%s       %s  %s%3d%%%s  %s%s%s",
		BebopYellow, Reset, t.generateBebopBar(100-data.API7dayPercent, 18, BebopYellow),
		BebopYellow, 100-data.API7dayPercent, Reset, BebopGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(BebopOrange + "║" + Reset)
	sb.WriteString(PadRight(line5, 89))
	sb.WriteString(BebopOrange + "║" + Reset + "\n")

	sb.WriteString(BebopOrange + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sData:%s %s%s%s  %sTime:%s %s  %sHits:%s %s%d%s  %sWoolongs:%s %s%s%s",
		BebopWhite, Reset, BebopWhite, FormatTokens(data.TokenCount), Reset,
		BebopGray, Reset, data.SessionTime,
		BebopGray, Reset, BebopGreen, data.MessageCount, Reset,
		BebopYellow, Reset, BebopYellow, FormatCost(data.SessionCost), Reset)

	sb.WriteString(BebopOrange + "║" + Reset)
	sb.WriteString(PadRight(line6, 89))
	sb.WriteString(BebopOrange + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sAccuracy:%s %s%d%%%s",
		BebopBlue, Reset, BebopBlue, FormatCost(data.DayCost), Reset,
		BebopRed, Reset, BebopRed, FormatCost(data.BurnRate), Reset,
		BebopGreen, Reset, BebopGreen, data.CacheHitRate, Reset)

	sb.WriteString(BebopOrange + "║" + Reset)
	sb.WriteString(PadRight(line7, 89))
	sb.WriteString(BebopOrange + "║" + Reset + "\n")

	sb.WriteString(BebopOrange + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")
	sb.WriteString(BebopOrange + "║" + Reset + "                           " + BebopWhite + "See You Space Cowboy..." + Reset + "                                 " + BebopOrange + "║" + Reset + "\n")
	sb.WriteString(BebopOrange + "╚═══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *BebopTheme) generateBebopBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(BebopDark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(BebopDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(BebopDark + "]" + Reset)
	return bar.String()
}
