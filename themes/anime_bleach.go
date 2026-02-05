package themes

import (
	"fmt"
	"strings"
)

// BleachTheme Bleach reiatsu/spiritual pressure style
type BleachTheme struct{}

func init() {
	RegisterTheme(&BleachTheme{})
}

func (t *BleachTheme) Name() string {
	return "bleach"
}

func (t *BleachTheme) Description() string {
	return "Bleach: Reiatsu spiritual pressure display"
}

const (
	BleachWhite    = "\033[38;2;255;255;255m"
	BleachBlack    = "\033[38;2;20;20;20m"
	BleachBlue     = "\033[38;2;0;150;255m"
	BleachRed      = "\033[38;2;200;0;0m"
	BleachPurple   = "\033[38;2;128;0;128m"
	BleachGold     = "\033[38;2;255;215;0m"
	BleachGray     = "\033[38;2;100;100;100m"
	BleachCyan     = "\033[38;2;0;200;200m"
	BleachBgBlack  = "\033[48;2;10;10;10m"
)

func (t *BleachTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Black and white contrast header
	sb.WriteString(BleachWhite + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" + Reset + "\n")

	// Soul Reaper status
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	division := "13th"
	if data.ModelType == "Opus" {
		division = "1st"
	} else if data.ModelType == "Haiku" {
		division = "4th"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s卍 BANKAI%s", BleachRed, Reset)
	}

	line1 := fmt.Sprintf(" %s死神%s %s%s%s %s  %sDivision:%s %s%s%s  %s%s%s%s",
		BleachWhite, Reset,
		modelColor, modelIcon, data.ModelName,
		Reset,
		BleachGray, Reset, BleachGold, division, Reset,
		BleachGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	// Target
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚔%s%s", BleachCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", BleachBlue, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", BleachRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf(" %sTarget:%s %s%s",
		BleachGray, Reset, ShortenPath(data.ProjectPath, 45), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(BleachGray + "────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	// Reiatsu (spiritual pressure) - context
	reiatsuColor := BleachBlue
	if data.ContextPercent > 75 {
		reiatsuColor = BleachRed
	} else if data.ContextPercent > 50 {
		reiatsuColor = BleachPurple
	}

	line3 := fmt.Sprintf(" %sREIATSU%s    %s  %s%3d%%%s",
		BleachBlue, Reset,
		t.generateBleachBar(data.ContextPercent, 22, reiatsuColor),
		reiatsuColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	// Reiryoku (spiritual power) - 5hr
	line4 := fmt.Sprintf(" %sREIRYOKU%s   %s  %s%3d%%%s  %s%s%s",
		BleachCyan, Reset,
		t.generateBleachBar(100-data.API5hrPercent, 22, BleachCyan),
		BleachCyan, 100-data.API5hrPercent, Reset,
		BleachGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	// Endurance - 7day
	line5 := fmt.Sprintf(" %sENDURANCE%s  %s  %s%3d%%%s  %s%s%s",
		BleachPurple, Reset,
		t.generateBleachBar(100-data.API7dayPercent, 22, BleachPurple),
		BleachPurple, 100-data.API7dayPercent, Reset,
		BleachGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(BleachGray + "────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	// Stats
	line6 := fmt.Sprintf(" %sPower:%s %s%s%s  %sTime:%s %s  %sStrikes:%s %s%d%s  %sCost:%s %s%s%s  %sDaily:%s %s%s%s",
		BleachWhite, Reset, BleachWhite, FormatTokens(data.TokenCount), Reset,
		BleachGray, Reset, data.SessionTime,
		BleachGray, Reset, BleachCyan, data.MessageCount, Reset,
		BleachGold, Reset, BleachGold, FormatCost(data.SessionCost), Reset,
		BleachPurple, Reset, BleachPurple, FormatCost(data.DayCost), Reset)
	sb.WriteString(line6 + "\n")

	line7 := fmt.Sprintf(" %sRate:%s %s%s/h%s  %sAccuracy:%s %s%d%%%s",
		BleachRed, Reset, BleachRed, FormatCost(data.BurnRate), Reset,
		BleachBlue, Reset, BleachBlue, data.CacheHitRate, Reset)
	sb.WriteString(line7 + "\n")

	sb.WriteString(BleachWhite + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" + Reset + "\n")

	return sb.String()
}

func (t *BleachTheme) generateBleachBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(BleachGray + "【" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("■", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(BleachGray)
		bar.WriteString(strings.Repeat("□", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(BleachGray + "】" + Reset)
	return bar.String()
}
