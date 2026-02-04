package themes

import (
	"fmt"
	"strings"
	"unicode"
)

// DungeonTheme 地牢火把風格
type DungeonTheme struct{}

func init() {
	RegisterTheme(&DungeonTheme{})
}

func (t *DungeonTheme) Name() string {
	return "dungeon"
}

func (t *DungeonTheme) Description() string {
	return "地牢：石牆火把照明，黑暗冒險氛圍"
}

const (
	DunStone     = "\033[38;2;105;105;105m"
	DunDarkStone = "\033[38;2;64;64;64m"
	DunTorch     = "\033[38;2;255;147;41m"
	DunFlame     = "\033[38;2;255;100;0m"
	DunGold      = "\033[38;2;255;215;0m"
	DunRed       = "\033[38;2;178;34;34m"
	DunGreen     = "\033[38;2;34;139;34m"
	DunBlue      = "\033[38;2;70;130;180m"
	DunPurple    = "\033[38;2;138;43;226m"
	DunBone      = "\033[38;2;255;250;240m"
	DunShadow    = "\033[38;2;25;25;25m"
	DunMoss      = "\033[38;2;85;107;47m"
)

func (t *DungeonTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Fixed width: 78 characters
	width := 78

	// Stone wall top with torches
	sb.WriteString(DunDarkStone + "▓▓▓" + DunTorch + "╔" + DunFlame + "҈" + DunTorch + "╗" + DunDarkStone + strings.Repeat("▓", width-12) + DunTorch + "╔" + DunFlame + "҈" + DunTorch + "╗" + DunDarkStone + "▓▓▓" + Reset + "\n")

	// Chamber name
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	chamberName := "The Dark Chamber"
	if data.ModelType == "Opus" {
		chamberName = "The Arcane Sanctum"
	} else if data.ModelType == "Haiku" {
		chamberName = "The Monk's Cell"
	}

	update := ""
	if data.UpdateAvailable {
		update = DunGold + " ★" + Reset
	}

	line1 := fmt.Sprintf("%s▓%s%s║%s %s%s%s%s  %s~ %s ~%s%s  %s%s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		modelColor, modelIcon, data.ModelName, Reset,
		DunTorch, chamberName, Reset, update,
		DunStone, data.Version)
	sb.WriteString(dunPadLine(line1, width, DunTorch+"║"+DunDarkStone+"▓"+Reset))

	// Quest scroll
	gitStr := ""
	if data.GitBranch != "" {
		gitStr = fmt.Sprintf("  %s⚔%s%s", DunMoss, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitStr += fmt.Sprintf(" %s+%d%s", DunGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitStr += fmt.Sprintf(" %s*%d%s", DunRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("%s▓%s%s║%s %s📜 %s%s%s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunBone, ShortenPath(data.ProjectPath, 30), Reset, gitStr)
	sb.WriteString(dunPadLine(line2, width, DunTorch+"║"+DunDarkStone+"▓"+Reset))

	// Stone separator
	sb.WriteString(DunDarkStone + "▓" + DunTorch + "╠" + DunStone + strings.Repeat("═", width-4) + DunTorch + "╣" + DunDarkStone + "▓" + Reset + "\n")

	// Stats as dungeon items
	line3 := fmt.Sprintf("%s▓%s%s║%s %s🗡%s %s  %s🛡%s %d  %s⏳%s %s  %s💀%s %s  %s💎%s %s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunRed, Reset, FormatTokens(data.TokenCount),
		DunBlue, Reset, data.MessageCount,
		DunStone, Reset, data.SessionTime,
		DunPurple, Reset, FormatCostShort(data.BurnRate),
		DunGold, Reset, FormatCostShort(data.DayCost))
	sb.WriteString(dunPadLine(line3, width, DunTorch+"║"+DunDarkStone+"▓"+Reset))

	// Health/Mana pools
	hp := 100 - data.ContextPercent
	hpColor := DunGreen
	if hp <= 20 {
		hpColor = DunRed
	} else if hp <= 50 {
		hpColor = DunTorch
	}

	hpBar := t.generateDungeonBar(hp, 15, hpColor)
	mpBar := t.generateDungeonBar(100-data.API5hrPercent, 12, DunBlue)
	xpBar := t.generateDungeonBar(100-data.API7dayPercent, 12, DunPurple)

	line4 := fmt.Sprintf("%s▓%s%s║%s %s❤%s%s%s%d%s  %s✦%s%s%s%d%s  %s⚡%s%s%s%d%s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunRed, Reset, hpBar, hpColor, hp, Reset,
		DunBlue, Reset, mpBar, DunBlue, 100-data.API5hrPercent, Reset,
		DunPurple, Reset, xpBar, DunPurple, 100-data.API7dayPercent, Reset)
	sb.WriteString(dunPadLine(line4, width, DunTorch+"║"+DunDarkStone+"▓"+Reset))

	// Treasure info
	line5 := fmt.Sprintf("%s▓%s%s║%s %s💰%s %s ses  %s⚗%s %d%% hit  %s⌛%s %s  %s⌛%s %s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunGold, Reset, FormatCostShort(data.SessionCost),
		DunGreen, Reset, data.CacheHitRate,
		DunStone, Reset, data.API5hrTimeLeft,
		DunStone, Reset, data.API7dayTimeLeft)
	sb.WriteString(dunPadLine(line5, width, DunTorch+"║"+DunDarkStone+"▓"+Reset))

	// Stone wall bottom with torches
	sb.WriteString(DunDarkStone + "▓▓▓" + DunTorch + "╚" + DunFlame + "҈" + DunTorch + "╝" + DunDarkStone + strings.Repeat("▓", width-12) + DunTorch + "╚" + DunFlame + "҈" + DunTorch + "╝" + DunDarkStone + "▓▓▓" + Reset + "\n")

	return sb.String()
}

func dunPadLine(line string, targetWidth int, suffix string) string {
	visible := dunDisplayWidth(line)
	suffixLen := dunDisplayWidth(suffix)
	padding := targetWidth - visible - suffixLen
	if padding < 0 {
		padding = 0
	}
	return line + strings.Repeat(" ", padding) + suffix + "\n"
}

// dunDisplayWidth calculates display width accounting for:
// - ANSI escape codes (0 width)
// - Emojis and wide characters (2 width)
// - Regular ASCII (1 width)
func dunDisplayWidth(s string) int {
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
			// Check if it's an emoji or wide character
			if isWideChar(r) {
				width += 2
			} else {
				width += 1
			}
		}
	}
	return width
}

// isWideChar checks if a rune is a wide character (emoji or CJK)
func isWideChar(r rune) bool {
	// Box Drawing characters are 1 width (check first)
	if r >= 0x2500 && r <= 0x257F {
		return false
	}
	// Block elements (▓▒░█) are 1 width
	if r >= 0x2580 && r <= 0x259F {
		return false
	}
	// Geometric shapes - most are 1 width
	if r >= 0x25A0 && r <= 0x25FF {
		return false
	}
	// Special 1-width characters
	switch r {
	case '▓', '▰', '▱', '҈', '★', '✦':
		return false
	}
	// CJK brackets - these ARE wide (2 cells)
	if r >= 0x3000 && r <= 0x303F {
		return true
	}
	// CJK punctuation brackets
	switch r {
	case '〔', '〕', '「', '」', '『', '』':
		return true
	}
	// Emojis
	if r >= 0x1F300 && r <= 0x1F9FF {
		return true
	}
	if r >= 0x2600 && r <= 0x26FF {
		return true
	}
	if r >= 0x2700 && r <= 0x27BF {
		return true
	}
	// CJK characters
	if unicode.Is(unicode.Han, r) {
		return true
	}
	// Full-width characters
	if r >= 0xFF00 && r <= 0xFFEF {
		return true
	}
	// Specific emojis
	switch r {
	case '❤', '⚡', '💰', '💎', '💀', '🗡', '🛡', '⏳', '⚗', '⌛', '⚔', '📜', '💛', '💙', '💚':
		return true
	}
	return false
}

func (t *DungeonTheme) generateDungeonBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(DunShadow + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▰", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(DunShadow)
		bar.WriteString(strings.Repeat("▱", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(DunShadow + "]" + Reset)
	return bar.String()
}
