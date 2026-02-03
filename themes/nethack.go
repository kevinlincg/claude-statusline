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
	width := 80

	// Dungeon map style header with walls
	sb.WriteString(NHDark + "-----+" + strings.Repeat("-", width-12) + "+-----" + Reset + "\n")

	// Player @ symbol with class
	modelColor, _ := GetModelConfig(data.ModelType)
	className := "Tourist"
	if data.ModelType == "Opus" {
		className = "Wizard"
	} else if data.ModelType == "Haiku" {
		className = "Monk"
	}

	// Status line 1: Character info (NetHack style)
	update := ""
	if data.UpdateAvailable {
		update = NHYellow + " (Dlvl^)" + Reset
	}

	line1 := fmt.Sprintf("%s|%s %s%s@%s %sthe %s%s %s%s%s%s  St:%s18%s Dx:%s%d%s Co:%s%d%s",
		NHDark, Reset,
		modelColor, Bold, Reset,
		NHWhite, className, Reset,
		NHGray, data.Version, Reset, update,
		NHGreen, Reset,
		NHCyan, data.MessageCount, Reset,
		NHYellow, data.CacheHitRate, Reset)
	sb.WriteString(nhPadLine(line1, width, NHDark+"|"+Reset))

	// Dungeon level (project path)
	gitStr := ""
	if data.GitBranch != "" {
		gitStr = fmt.Sprintf("  %s<%s>%s", NHMagenta, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitStr += fmt.Sprintf(" %s+%d%s", NHGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitStr += fmt.Sprintf(" %s~%d%s", NHOrange, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("%s|%s %sDlvl:%s%s%s%s  %s$:%s%s  %sT:%s%s",
		NHDark, Reset,
		NHWhite, NHBrown, ShortenPath(data.ProjectPath, 25), Reset, gitStr,
		NHYellow, FormatCostShort(data.DayCost), Reset,
		NHGray, data.SessionTime, Reset)
	sb.WriteString(nhPadLine(line2, width, NHDark+"|"+Reset))

	// Separator (dungeon floor)
	sb.WriteString(NHDark + "-----+" + strings.Repeat("-", width-12) + "+-----" + Reset + "\n")

	// Status bars (HP/Pw/AC style)
	hp := 100 - data.ContextPercent
	hpMax := 100
	pw := 100 - data.API5hrPercent
	pwMax := 100
	ac := 100 - data.API7dayPercent

	hpColor := NHGreen
	if hp <= 20 {
		hpColor = NHRed
	} else if hp <= 50 {
		hpColor = NHYellow
	}

	// HP and Pw bars
	hpBar := t.generateNHBar(hp, 12)
	pwBar := t.generateNHBar(pw, 10)
	acBar := t.generateNHBar(ac, 10)

	line3 := fmt.Sprintf("%s|%s %sHP:%s%s%s%d%s(%s%d%s) %sPw:%s%s%s%d%s(%s%d%s) %sAC:%s%s%s%d%s %sXp:%s%s%s",
		NHDark, Reset,
		NHWhite, hpColor, hpBar, hpColor, hp, NHDark, NHGray, hpMax, Reset,
		NHWhite, NHBlue, pwBar, NHBlue, pw, NHDark, NHGray, pwMax, Reset,
		NHWhite, NHCyan, acBar, NHCyan, ac, Reset,
		NHWhite, NHMagenta, FormatTokens(data.TokenCount), Reset)
	sb.WriteString(nhPadLine(line3, width, NHDark+"|"+Reset))

	// Bottom status (Hunger/Encumbrance style)
	hungerStatus := "Satiated"
	if data.ContextPercent >= 80 {
		hungerStatus = "Fainting"
	} else if data.ContextPercent >= 60 {
		hungerStatus = "Hungry"
	}

	line4 := fmt.Sprintf("%s|%s %s%s%s  %sBurdened%s  %s%s%s ses  %s%s/h%s rate  %s%s%s left",
		NHDark, Reset,
		NHOrange, hungerStatus, Reset,
		NHYellow, Reset,
		NHGreen, FormatCostShort(data.SessionCost), Reset,
		NHRed, FormatCostShort(data.BurnRate), Reset,
		NHGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(nhPadLine(line4, width, NHDark+"|"+Reset))

	// Bottom wall
	sb.WriteString(NHDark + "-----+" + strings.Repeat("-", width-12) + "+-----" + Reset + "\n")

	return sb.String()
}

func nhPadLine(line string, targetWidth int, suffix string) string {
	visible := nhVisibleLen(line)
	suffixLen := nhVisibleLen(suffix)
	padding := targetWidth - visible - suffixLen
	if padding < 0 {
		padding = 0
	}
	return line + strings.Repeat(" ", padding) + suffix + "\n"
}

func nhVisibleLen(s string) int {
	inEscape := false
	count := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
		} else if inEscape {
			if r == 'm' {
				inEscape = false
			}
		} else {
			count++
		}
	}
	return count
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
		bar.WriteString(strings.Repeat("#", filled))
	}
	if empty > 0 {
		bar.WriteString(NHDark)
		bar.WriteString(strings.Repeat("-", empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}
