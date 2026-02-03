package themes

import (
	"fmt"
	"strings"
)

// MUDRPGTheme MUD RPG 角色狀態風格
type MUDRPGTheme struct{}

func init() {
	RegisterTheme(&MUDRPGTheme{})
}

func (t *MUDRPGTheme) Name() string {
	return "mud_rpg"
}

func (t *MUDRPGTheme) Description() string {
	return "MUD RPG：經典文字冒險遊戲角色狀態介面"
}

const (
	MUDGold    = "\033[38;2;255;215;0m"
	MUDRed     = "\033[38;2;255;80;80m"
	MUDBlue    = "\033[38;2;100;149;237m"
	MUDGreen   = "\033[38;2;50;205;50m"
	MUDCyan    = "\033[38;2;0;206;209m"
	MUDMagenta = "\033[38;2;218;112;214m"
	MUDWhite   = "\033[38;2;245;245;245m"
	MUDGray    = "\033[38;2;128;128;128m"
	MUDDark    = "\033[38;2;64;64;64m"
	MUDBrown   = "\033[38;2;139;90;43m"
)

func (t *MUDRPGTheme) Render(data StatusData) string {
	var sb strings.Builder
	width := 80

	// Header with character class style
	modelColor, _ := GetModelConfig(data.ModelType)
	className := "Artificer"
	if data.ModelType == "Opus" {
		className = "Archmage"
	} else if data.ModelType == "Haiku" {
		className = "Apprentice"
	}

	// Top border
	sb.WriteString(MUDDark + "+" + strings.Repeat("=", width-2) + "+" + Reset + "\n")

	// Character name and class
	update := ""
	if data.UpdateAvailable {
		update = MUDGold + " *LEVEL UP*" + Reset
	}

	gitStr := ""
	if data.GitBranch != "" {
		gitStr = fmt.Sprintf("  %s[%s]%s", MUDBrown, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitStr += fmt.Sprintf("%s+%d%s", MUDGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitStr += fmt.Sprintf("%s~%d%s", MUDRed, data.GitDirty, Reset)
		}
	}

	line1 := fmt.Sprintf("%s|%s %s[%s]%s Lv.%s%s%s%s  %s%s%s%s",
		MUDDark, Reset,
		modelColor, data.ModelName, Reset,
		MUDCyan, data.Version, Reset, update,
		MUDGray, ShortenPath(data.ProjectPath, 25), Reset, gitStr)
	sb.WriteString(mudPadLine(line1, width, MUDDark+"|"+Reset))

	// Separator
	sb.WriteString(MUDDark + "+" + strings.Repeat("-", width-2) + "+" + Reset + "\n")

	// Stats row 1: HP / MP / Class
	hpBar := t.generateMUDBar(100-data.ContextPercent, 10, MUDRed)
	mpBar := t.generateMUDBar(100-data.API5hrPercent, 8, MUDBlue)
	xpBar := t.generateMUDBar(100-data.API7dayPercent, 8, MUDCyan)

	hpColor := MUDGreen
	if data.ContextPercent >= 80 {
		hpColor = MUDRed
	} else if data.ContextPercent >= 60 {
		hpColor = MUDGold
	}

	line2 := fmt.Sprintf("%s|%s %sHP%s%s%s%3d%%%s  %sMP%s%s%s%3d%%%s  %sXP%s%s%s%3d%%%s  %s%s%s",
		MUDDark, Reset,
		MUDRed, Reset, hpBar, hpColor, 100-data.ContextPercent, Reset,
		MUDBlue, Reset, mpBar, MUDBlue, 100-data.API5hrPercent, Reset,
		MUDCyan, Reset, xpBar, MUDCyan, 100-data.API7dayPercent, Reset,
		MUDMagenta, className, Reset)
	sb.WriteString(mudPadLine(line2, width, MUDDark+"|"+Reset))

	// Stats row 2: Combat stats
	line3 := fmt.Sprintf("%s|%s %sATK%s %-6s  %sDEF%s %-3d  %sSPD%s %-6s  %sLUK%s %d%%  %sGP%s %s",
		MUDDark, Reset,
		MUDRed, MUDWhite, FormatTokens(data.TokenCount),
		MUDBlue, MUDWhite, data.MessageCount,
		MUDGreen, MUDWhite, data.SessionTime,
		MUDMagenta, MUDWhite, data.CacheHitRate,
		MUDGold, MUDGold, FormatCostShort(data.DayCost))
	sb.WriteString(mudPadLine(line3, width, MUDDark+"|"+Reset))

	// Stats row 3: Equipment
	line4 := fmt.Sprintf("%s|%s %sSes%s %s  %sRate%s %s/h  %s5hr%s %s  %s7day%s %s",
		MUDDark, Reset,
		MUDBrown, MUDGold, FormatCostShort(data.SessionCost),
		MUDBrown, MUDRed, FormatCostShort(data.BurnRate),
		MUDGray, MUDGray, data.API5hrTimeLeft,
		MUDGray, MUDGray, data.API7dayTimeLeft)
	sb.WriteString(mudPadLine(line4, width, MUDDark+"|"+Reset))

	// Bottom border
	sb.WriteString(MUDDark + "+" + strings.Repeat("=", width-2) + "+" + Reset + "\n")

	return sb.String()
}

func mudPadLine(line string, targetWidth int, suffix string) string {
	visible := mudVisibleLen(line)
	suffixLen := mudVisibleLen(suffix)
	padding := targetWidth - visible - suffixLen
	if padding < 0 {
		padding = 0
	}
	return line + strings.Repeat(" ", padding) + suffix + "\n"
}

func mudVisibleLen(s string) int {
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

func (t *MUDRPGTheme) generateMUDBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(MUDDark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("=", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(MUDDark)
		bar.WriteString(strings.Repeat("-", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(MUDDark + "]" + Reset)
	return bar.String()
}
