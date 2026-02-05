package themes

import (
	"fmt"
	"strings"
)

// VisualNovelTheme Visual Novel dialog box style
type VisualNovelTheme struct{}

func init() {
	RegisterTheme(&VisualNovelTheme{})
}

func (t *VisualNovelTheme) Name() string {
	return "visualnovel"
}

func (t *VisualNovelTheme) Description() string {
	return "Visual Novel: Dialog box with character portrait style"
}

const (
	VNBlue    = "\033[38;2;100;150;200m"
	VNPink    = "\033[38;2;255;182;193m"
	VNPurple  = "\033[38;2;180;150;200m"
	VNGold    = "\033[38;2;255;215;150m"
	VNWhite   = "\033[38;2;255;255;255m"
	VNGray    = "\033[38;2;150;150;150m"
	VNDark    = "\033[38;2;40;50;60m"
)

func (t *VisualNovelTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Background frame
	sb.WriteString(VNDark + "▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄" + Reset + "\n")

	// Character name box
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	character := "Mysterious Girl"
	if data.ModelType == "Opus" {
		character = "Protagonist"
	} else if data.ModelType == "Haiku" {
		character = "Childhood Friend"
	}

	sb.WriteString(VNDark + "█" + Reset + "  " + VNBlue + "┌─────────────────┐" + Reset + "                                                           " + VNDark + "█" + Reset + "\n")
	nameBox := fmt.Sprintf("%s│ %s%-15s%s │%s", VNBlue, VNGold, character, VNBlue, Reset)
	sb.WriteString(VNDark + "█" + Reset + "  " + nameBox + "   " + VNPurple + "ビジュアルノベル" + Reset + "                                       " + VNDark + "█" + Reset + "\n")
	sb.WriteString(VNDark + "█" + Reset + "  " + VNBlue + "└─────────────────┘" + Reset + "                                                           " + VNDark + "█" + Reset + "\n")

	// Dialog box
	sb.WriteString(VNDark + "█" + Reset + "  " + VNPurple + "╔════════════════════════════════════════════════════════════════════════════════╗" + Reset + "  " + VNDark + "█" + Reset + "\n")

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[New Route!]%s", VNPink, Reset)
	}

	line1 := fmt.Sprintf("%s「%sModel: %s%s%s%s  %sVersion: %s%s%s%s」%s",
		VNWhite, Reset, modelColor, modelIcon, data.ModelName, Reset,
		VNGray, Reset, data.Version, update, VNWhite, Reset)
	sb.WriteString(VNDark + "█" + Reset + "  " + VNPurple + "║" + Reset + " " + PadRight(line1, 78) + " " + VNPurple + "║" + Reset + "  " + VNDark + "█" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf(" %s[%s]%s", VNBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", VNPink, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", VNPurple, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("%s「%sScene: %s%s%s」%s",
		VNWhite, Reset, ShortenPath(data.ProjectPath, 45), gitInfo, VNWhite, Reset)
	sb.WriteString(VNDark + "█" + Reset + "  " + VNPurple + "║" + Reset + " " + PadRight(line2, 78) + " " + VNPurple + "║" + Reset + "  " + VNDark + "█" + Reset + "\n")

	sb.WriteString(VNDark + "█" + Reset + "  " + VNPurple + "╠════════════════════════════════════════════════════════════════════════════════╣" + Reset + "  " + VNDark + "█" + Reset + "\n")

	// Stats as dialog choices
	affectionColor := VNPink
	if data.ContextPercent > 75 {
		affectionColor = VNPurple
	}

	line3 := fmt.Sprintf("%s▸ Affection%s   %s  %s%3d%%%s",
		VNPink, Reset, t.generateVNBar(data.ContextPercent, 16, affectionColor), affectionColor, data.ContextPercent, Reset)
	sb.WriteString(VNDark + "█" + Reset + "  " + VNPurple + "║" + Reset + "   " + PadRight(line3, 75) + " " + VNPurple + "║" + Reset + "  " + VNDark + "█" + Reset + "\n")

	line4 := fmt.Sprintf("%s▸ Trust%s       %s  %s%3d%%%s  %s%s%s",
		VNBlue, Reset, t.generateVNBar(100-data.API5hrPercent, 16, VNBlue),
		VNBlue, 100-data.API5hrPercent, Reset, VNGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(VNDark + "█" + Reset + "  " + VNPurple + "║" + Reset + "   " + PadRight(line4, 75) + " " + VNPurple + "║" + Reset + "  " + VNDark + "█" + Reset + "\n")

	line5 := fmt.Sprintf("%s▸ Destiny%s     %s  %s%3d%%%s  %s%s%s",
		VNGold, Reset, t.generateVNBar(100-data.API7dayPercent, 16, VNGold),
		VNGold, 100-data.API7dayPercent, Reset, VNGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(VNDark + "█" + Reset + "  " + VNPurple + "║" + Reset + "   " + PadRight(line5, 75) + " " + VNPurple + "║" + Reset + "  " + VNDark + "█" + Reset + "\n")

	sb.WriteString(VNDark + "█" + Reset + "  " + VNPurple + "╠════════════════════════════════════════════════════════════════════════════════╣" + Reset + "  " + VNDark + "█" + Reset + "\n")

	line6 := fmt.Sprintf("%sWords: %s%s%s  Time: %s  Choices: %s%d%s  Route: %s$%s%s",
		VNGray, VNWhite, FormatTokens(data.TokenCount), Reset,
		data.SessionTime,
		VNPink, data.MessageCount, Reset,
		VNGold, FormatCost(data.SessionCost), Reset)
	sb.WriteString(VNDark + "█" + Reset + "  " + VNPurple + "║" + Reset + "   " + PadRight(line6, 75) + " " + VNPurple + "║" + Reset + "  " + VNDark + "█" + Reset + "\n")

	sb.WriteString(VNDark + "█" + Reset + "  " + VNPurple + "╚════════════════════════════════════════════════════════════════════════════════╝" + Reset + "  " + VNDark + "█" + Reset + "\n")

	// Choice indicator
	sb.WriteString(VNDark + "█" + Reset + "                                                   " + VNPink + "▼ Press SPACE to continue ▼" + Reset + "       " + VNDark + "█" + Reset + "\n")
	sb.WriteString(VNDark + "▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀" + Reset + "\n")

	return sb.String()
}

func (t *VisualNovelTheme) generateVNBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(VNGray + "「" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("●", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(VNGray)
		bar.WriteString(strings.Repeat("○", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(VNGray + "」" + Reset)
	return bar.String()
}
