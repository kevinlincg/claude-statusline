package themes

import (
	"fmt"
	"strings"
)

// HxHTheme Hunter x Hunter nen system style
type HxHTheme struct{}

func init() {
	RegisterTheme(&HxHTheme{})
}

func (t *HxHTheme) Name() string {
	return "hxh"
}

func (t *HxHTheme) Description() string {
	return "HxH: Hunter x Hunter nen system style"
}

const (
	HxHGreen   = "\033[38;2;0;200;100m"
	HxHBlue    = "\033[38;2;50;150;255m"
	HxHYellow  = "\033[38;2;255;220;0m"
	HxHRed     = "\033[38;2;255;50;50m"
	HxHPurple  = "\033[38;2;180;0;255m"
	HxHOrange  = "\033[38;2;255;165;0m"
	HxHWhite   = "\033[38;2;240;240;240m"
	HxHGray    = "\033[38;2;100;100;100m"
	HxHDark    = "\033[38;2;50;50;50m"
)

func (t *HxHTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(HxHGreen + "┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓" + Reset + "\n")
	sb.WriteString(HxHGreen + "┃" + Reset + "  " + HxHYellow + "◆" + HxHWhite + " HUNTER LICENSE " + HxHYellow + "◆" + Reset + "   " + HxHGreen + "ハンター協会" + Reset + "                                         " + HxHGreen + "┃" + Reset + "\n")
	sb.WriteString(HxHGreen + "┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	nenType := "Enhancer"
	nenColor := HxHYellow
	if data.ModelType == "Opus" {
		nenType = "Specialist"
		nenColor = HxHPurple
	} else if data.ModelType == "Haiku" {
		nenType = "Conjurer"
		nenColor = HxHBlue
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[NEW ABILITY]%s", HxHPurple, Reset)
	}

	line1 := fmt.Sprintf("  %sHunter:%s %s%s%s  %sNen:%s %s%s%s  %s%s%s%s",
		HxHGreen, Reset, modelColor, modelIcon, data.ModelName,
		HxHGray, Reset, nenColor, nenType, Reset,
		HxHGray, data.Version, Reset, update)

	sb.WriteString(HxHGreen + "┃" + Reset)
	sb.WriteString(PadRight(line1, 84))
	sb.WriteString(HxHGreen + "┃" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s♦%s%s", HxHBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", HxHGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", HxHOrange, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sQuest:%s %s%s",
		HxHYellow, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)

	sb.WriteString(HxHGreen + "┃" + Reset)
	sb.WriteString(PadRight(line2, 84))
	sb.WriteString(HxHGreen + "┃" + Reset + "\n")

	sb.WriteString(HxHGreen + "┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫" + Reset + "\n")

	auraColor := nenColor
	if data.ContextPercent > 75 {
		auraColor = HxHRed
	}

	line3 := fmt.Sprintf("  %sAura%s      %s  %s%3d%%%s",
		nenColor, Reset, t.generateHxHBar(data.ContextPercent, 18, auraColor), auraColor, data.ContextPercent, Reset)

	sb.WriteString(HxHGreen + "┃" + Reset)
	sb.WriteString(PadRight(line3, 84))
	sb.WriteString(HxHGreen + "┃" + Reset + "\n")

	line4 := fmt.Sprintf("  %sStamina%s   %s  %s%3d%%%s  %s%s%s",
		HxHGreen, Reset, t.generateHxHBar(100-data.API5hrPercent, 18, HxHGreen),
		HxHGreen, 100-data.API5hrPercent, Reset, HxHGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(HxHGreen + "┃" + Reset)
	sb.WriteString(PadRight(line4, 84))
	sb.WriteString(HxHGreen + "┃" + Reset + "\n")

	line5 := fmt.Sprintf("  %sResolve%s   %s  %s%3d%%%s  %s%s%s",
		HxHOrange, Reset, t.generateHxHBar(100-data.API7dayPercent, 18, HxHOrange),
		HxHOrange, 100-data.API7dayPercent, Reset, HxHGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(HxHGreen + "┃" + Reset)
	sb.WriteString(PadRight(line5, 84))
	sb.WriteString(HxHGreen + "┃" + Reset + "\n")

	sb.WriteString(HxHGreen + "┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫" + Reset + "\n")

	line6 := fmt.Sprintf("  %sEXP:%s %s%s%s  %sTime:%s %s  %sHunts:%s %s%d%s  %sJenny:%s %s%s%s  %sDaily:%s %s%s%s",
		HxHPurple, Reset, HxHPurple, FormatTokens(data.TokenCount), Reset,
		HxHGray, Reset, data.SessionTime,
		HxHGray, Reset, HxHWhite, data.MessageCount, Reset,
		HxHYellow, Reset, HxHYellow, FormatCost(data.SessionCost), Reset,
		HxHOrange, Reset, HxHOrange, FormatCost(data.DayCost), Reset)

	sb.WriteString(HxHGreen + "┃" + Reset)
	sb.WriteString(PadRight(line6, 84))
	sb.WriteString(HxHGreen + "┃" + Reset + "\n")

	line7 := fmt.Sprintf("  %sRate:%s %s%s/h%s  %sFocus:%s %s%d%%%s",
		HxHRed, Reset, HxHRed, FormatCost(data.BurnRate), Reset,
		HxHBlue, Reset, HxHBlue, data.CacheHitRate, Reset)

	sb.WriteString(HxHGreen + "┃" + Reset)
	sb.WriteString(PadRight(line7, 84))
	sb.WriteString(HxHGreen + "┃" + Reset + "\n")

	sb.WriteString(HxHGreen + "┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛" + Reset + "\n")

	return sb.String()
}

func (t *HxHTheme) generateHxHBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(HxHDark + "〈" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("●", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(HxHDark)
		bar.WriteString(strings.Repeat("○", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(HxHDark + "〉" + Reset)
	return bar.String()
}
