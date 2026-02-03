package themes

import (
	"fmt"
	"strings"
	"unicode"
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
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	className := "Artificer"
	if data.ModelType == "Opus" {
		className = "Archmage"
	} else if data.ModelType == "Haiku" {
		className = "Apprentice"
	}

	// Top border with double-line box drawing
	sb.WriteString(MUDDark + "╔" + strings.Repeat("═", width-2) + "╗" + Reset + "\n")

	// Character name and class
	update := ""
	if data.UpdateAvailable {
		update = MUDGold + " *UP*" + Reset
	}

	gitStr := ""
	if data.GitBranch != "" {
		gitStr = fmt.Sprintf(" %s<%s>%s", MUDBrown, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitStr += fmt.Sprintf("%s+%d%s", MUDGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitStr += fmt.Sprintf("%s~%d%s", MUDRed, data.GitDirty, Reset)
		}
	}

	line1 := fmt.Sprintf("%s║%s %s%s%s%s Lv.%s%s%s%s  %s%s%s%s",
		MUDDark, Reset,
		modelColor, modelIcon, data.ModelName, Reset,
		MUDCyan, data.Version, Reset, update,
		MUDGray, ShortenPath(data.ProjectPath, 22), Reset, gitStr)
	sb.WriteString(mudPadLine(line1, width, MUDDark+"║"+Reset))

	// Separator
	sb.WriteString(MUDDark + "╠" + strings.Repeat("═", width-2) + "╣" + Reset + "\n")

	// Stats row 1: HP / MP / XP / Class
	hpBar := t.generateMUDBar(100-data.ContextPercent, 8, MUDRed)
	mpBar := t.generateMUDBar(100-data.API5hrPercent, 6, MUDBlue)
	xpBar := t.generateMUDBar(100-data.API7dayPercent, 6, MUDCyan)

	hpColor := MUDGreen
	if data.ContextPercent >= 80 {
		hpColor = MUDRed
	} else if data.ContextPercent >= 60 {
		hpColor = MUDGold
	}

	line2 := fmt.Sprintf("%s║%s %sHP%s%s%s%3d%%%s %sMP%s%s%s%3d%%%s %sXP%s%s%s%3d%%%s %s%s%s %sGP%s%s",
		MUDDark, Reset,
		MUDRed, Reset, hpBar, hpColor, 100-data.ContextPercent, Reset,
		MUDBlue, Reset, mpBar, MUDBlue, 100-data.API5hrPercent, Reset,
		MUDCyan, Reset, xpBar, MUDCyan, 100-data.API7dayPercent, Reset,
		MUDMagenta, className, Reset,
		MUDGold, MUDGold, FormatCostShort(data.DayCost))
	sb.WriteString(mudPadLine(line2, width, MUDDark+"║"+Reset))

	// Stats row 2: Combat stats
	line3 := fmt.Sprintf("%s║%s %sATK%s%-6s %sDEF%s%-3d %sSPD%s%-6s %sLUK%s%d%% %sSes%s%s %sRate%s%s/h",
		MUDDark, Reset,
		MUDRed, MUDWhite, FormatTokens(data.TokenCount),
		MUDBlue, MUDWhite, data.MessageCount,
		MUDGreen, MUDWhite, data.SessionTime,
		MUDMagenta, MUDWhite, data.CacheHitRate,
		MUDBrown, MUDGold, FormatCostShort(data.SessionCost),
		MUDBrown, MUDRed, FormatCostShort(data.BurnRate))
	sb.WriteString(mudPadLine(line3, width, MUDDark+"║"+Reset))

	// Stats row 3: Time left
	line4 := fmt.Sprintf("%s║%s %s5hr%s %-6s  %s7day%s %-6s",
		MUDDark, Reset,
		MUDGray, MUDCyan, data.API5hrTimeLeft,
		MUDGray, MUDCyan, data.API7dayTimeLeft)
	sb.WriteString(mudPadLine(line4, width, MUDDark+"║"+Reset))

	// Bottom border
	sb.WriteString(MUDDark + "╚" + strings.Repeat("═", width-2) + "╝" + Reset + "\n")

	return sb.String()
}

func mudPadLine(line string, targetWidth int, suffix string) string {
	visible := mudDisplayWidth(line)
	suffixLen := mudDisplayWidth(suffix)
	padding := targetWidth - visible - suffixLen
	if padding < 0 {
		padding = 0
	}
	return line + strings.Repeat(" ", padding) + suffix + "\n"
}

// mudDisplayWidth calculates display width accounting for:
// - ANSI escape codes (0 width)
// - Emojis and wide characters (2 width)
// - Regular ASCII (1 width)
func mudDisplayWidth(s string) int {
	inEscape := false
	width := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
		} else if inEscape {
			if r == 'm' {
				inEscape = false
			}
		} else {
			if mudIsWideChar(r) {
				width += 2
			} else {
				width += 1
			}
		}
	}
	return width
}

// mudIsWideChar checks if a rune is a wide character (emoji or CJK)
func mudIsWideChar(r rune) bool {
	// Emojis and symbols that are typically 2 cells wide
	if r >= 0x1F300 && r <= 0x1F9FF { // Misc Symbols, Emoticons, etc.
		return true
	}
	if r >= 0x2600 && r <= 0x26FF { // Misc Symbols
		return true
	}
	if r >= 0x2700 && r <= 0x27BF { // Dingbats
		return true
	}
	// Box Drawing characters are 1 width
	if r >= 0x2500 && r <= 0x257F {
		return false
	}
	// CJK characters
	if unicode.Is(unicode.Han, r) {
		return true
	}
	// Full-width characters
	if r >= 0xFF00 && r <= 0xFFEF {
		return true
	}
	// Model icons (emojis)
	switch r {
	case '💛', '💙', '💚', '⚔', '🛡', '⏱', '★':
		return true
	}
	return false
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
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(MUDDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(MUDDark + "]" + Reset)
	return bar.String()
}
