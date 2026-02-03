package themes

import (
	"fmt"
	"strings"
)

// NetHackTheme NetHack/Roguelike 風格
type NetHackTheme struct{}

func init() {
	RegisterTheme(&NetHackTheme{})
}

func (t *NetHackTheme) Name() string {
	return "nethack"
}

func (t *NetHackTheme) Description() string {
	return "NetHack：經典 Roguelike 地牢探索風格"
}

const (
	NHWhite   = "\033[38;2;255;255;255m"
	NHGray    = "\033[38;2;170;170;170m"
	NHDark    = "\033[38;2;85;85;85m"
	NHRed     = "\033[38;2;255;85;85m"
	NHGreen   = "\033[38;2;85;255;85m"
	NHYellow  = "\033[38;2;255;255;85m"
	NHBlue    = "\033[38;2;85;85;255m"
	NHMagenta = "\033[38;2;255;85;255m"
	NHCyan    = "\033[38;2;85;255;255m"
	NHBrown   = "\033[38;2;170;85;0m"
	NHOrange  = "\033[38;2;255;170;0m"
)

func (t *NetHackTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Dungeon map style header with walls
	sb.WriteString(NHDark + "─────┬─────────────────────────────────────────────────────────────────────┬─────" + Reset + "\n")

	// Player @ symbol with class
	modelColor, _ := GetModelConfig(data.ModelType)
	playerClass := "@"
	className := "Tourist"
	if data.ModelType == "Opus" {
		playerClass = "@"
		className = "Wizard"
	} else if data.ModelType == "Haiku" {
		playerClass = "@"
		className = "Monk"
	}

	// Status line 1: Character info (NetHack style)
	line1 := fmt.Sprintf("%s│%s %s%s%s%s %sthe %s%s %s%s%s",
		NHDark, Reset,
		modelColor, Bold, playerClass, Reset,
		NHWhite, className, Reset,
		NHGray, data.Version, Reset)
	if data.UpdateAvailable {
		line1 += NHYellow + " (Dlvl↑)" + Reset
	}
	line1 += fmt.Sprintf("  %sSt:%s18  %sDx:%s%d  %sCo:%s%d",
		NHWhite, NHGreen,
		NHWhite, NHCyan, data.MessageCount,
		NHWhite, NHYellow, data.CacheHitRate)
	sb.WriteString(PadRight(line1, 77))
	sb.WriteString(NHDark + "│" + Reset + "\n")

	// Dungeon level (project path)
	line2 := fmt.Sprintf("%s│%s %sDlvl:%s%s%s",
		NHDark, Reset,
		NHWhite, NHBrown, ShortenPath(data.ProjectPath, 25), Reset)
	if data.GitBranch != "" {
		line2 += fmt.Sprintf("  %s<%s>%s", NHMagenta, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line2 += fmt.Sprintf(" %s+%d%s", NHGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line2 += fmt.Sprintf(" %s~%d%s", NHOrange, data.GitDirty, Reset)
		}
	}
	line2 += fmt.Sprintf("  %s$:%s%s%s  %sT:%s%s",
		NHYellow, NHYellow, FormatCostShort(data.DayCost), Reset,
		NHGray, NHGray, data.SessionTime)
	sb.WriteString(PadRight(line2, 77))
	sb.WriteString(NHDark + "│" + Reset + "\n")

	// Separator (dungeon floor)
	sb.WriteString(NHDark + "─────┼─────────────────────────────────────────────────────────────────────┼─────" + Reset + "\n")

	// Status bars (HP/Pw/AC style)
	hp := 100 - data.ContextPercent
	hpMax := 100
	pw := 100 - data.API5hrPercent
	pwMax := 100

	hpColor := NHGreen
	if hp <= 20 {
		hpColor = NHRed
	} else if hp <= 50 {
		hpColor = NHYellow
	}

	// HP and Pw bars
	hpBar := t.generateNHBar(hp, 15)
	pwBar := t.generateNHBar(pw, 12)
	acBar := t.generateNHBar(100-data.API7dayPercent, 12)

	line3 := fmt.Sprintf("%s│%s %sHP:%s%s%s%d%s(%s%d%s)  %sPw:%s%s%s%d%s(%s%d%s)  %sAC:%s%s%s%d%s  %sXp:%s%s%s",
		NHDark, Reset,
		NHWhite, hpColor, hpBar, hpColor, hp, NHDark, NHGray, hpMax, Reset,
		NHWhite, NHBlue, pwBar, NHBlue, pw, NHDark, NHGray, pwMax, Reset,
		NHWhite, NHCyan, acBar, NHCyan, 100-data.API7dayPercent, Reset,
		NHWhite, NHMagenta, FormatTokens(data.TokenCount), Reset)
	sb.WriteString(PadRight(line3, 77))
	sb.WriteString(NHDark + "│" + Reset + "\n")

	// Bottom status (Hunger/Encumbrance style)
	hungerStatus := "Satiated"
	if data.ContextPercent >= 80 {
		hungerStatus = "Fainting"
	} else if data.ContextPercent >= 60 {
		hungerStatus = "Hungry"
	}

	line4 := fmt.Sprintf("%s│%s %s%s%s  %sBurdened%s  %s%s%s ses  %s%s%s/h rate  %s%s%s left",
		NHDark, Reset,
		NHOrange, hungerStatus, Reset,
		NHYellow, Reset,
		NHGreen, FormatCostShort(data.SessionCost), Reset,
		NHRed, FormatCostShort(data.BurnRate), Reset,
		NHGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(PadRight(line4, 77))
	sb.WriteString(NHDark + "│" + Reset + "\n")

	// Bottom wall
	sb.WriteString(NHDark + "─────┴─────────────────────────────────────────────────────────────────────┴─────" + Reset + "\n")

	return sb.String()
}

func (t *NetHackTheme) generateNHBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(strings.Repeat("█", filled))
	}
	if empty > 0 {
		bar.WriteString(NHDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}
