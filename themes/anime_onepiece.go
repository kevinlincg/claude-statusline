package themes

import (
	"fmt"
	"strings"
)

// OnePieceTheme One Piece wanted poster / pirate style
type OnePieceTheme struct{}

func init() {
	RegisterTheme(&OnePieceTheme{})
}

func (t *OnePieceTheme) Name() string {
	return "onepiece"
}

func (t *OnePieceTheme) Description() string {
	return "One Piece: Wanted poster bounty style"
}

const (
	OPBrown     = "\033[38;2;139;90;43m"
	OPGold      = "\033[38;2;255;215;0m"
	OPRed       = "\033[38;2;180;30;30m"
	OPCream     = "\033[38;2;255;248;220m"
	OPBlue      = "\033[38;2;0;105;148m"
	OPDarkBrown = "\033[38;2;101;67;33m"
	OPBlack     = "\033[38;2;30;30;30m"
)

func (t *OnePieceTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Weathered poster border
	sb.WriteString(OPDarkBrown + "▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄" + Reset + "\n")
	sb.WriteString(OPBrown + "█" + OPCream + "                                                                           " + OPBrown + "█" + Reset + "\n")

	// WANTED banner
	sb.WriteString(OPBrown + "█" + OPCream + "              " + OPRed + "██╗    ██╗ █████╗ ███╗   ██╗████████╗███████╗██████╗" + OPCream + "               " + OPBrown + "█" + Reset + "\n")
	sb.WriteString(OPBrown + "█" + OPCream + "              " + OPRed + "██║    ██║██╔══██╗████╗  ██║╚══██╔══╝██╔════╝██╔══██╗" + OPCream + "              " + OPBrown + "█" + Reset + "\n")
	sb.WriteString(OPBrown + "█" + OPCream + "              " + OPRed + "██║ █╗ ██║███████║██╔██╗ ██║   ██║   █████╗  ██║  ██║" + OPCream + "              " + OPBrown + "█" + Reset + "\n")
	sb.WriteString(OPBrown + "█" + OPCream + "              " + OPRed + "╚██╗██╔╝██╔══██║██║╚██╗██║   ██║   ██╔══╝  ██║  ██║" + OPCream + "              " + OPBrown + "█" + Reset + "\n")
	sb.WriteString(OPBrown + "█" + OPCream + "              " + OPRed + " ╚███╔╝ ██║  ██║██║ ╚████║   ██║   ███████╗██████╔╝" + OPCream + "              " + OPBrown + "█" + Reset + "\n")
	sb.WriteString(OPBrown + "█" + OPCream + "                                                                           " + OPBrown + "█" + Reset + "\n")

	// Model as pirate name
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	crewName := "Straw Hat Pirates"
	if data.ModelType == "Opus" {
		crewName = "Pirate King"
	} else if data.ModelType == "Haiku" {
		crewName = "East Blue Rookie"
	}

	sb.WriteString(OPBrown + "█" + OPCream + "                    " + modelColor + modelIcon + data.ModelName + Reset + "  " + OPDarkBrown + "「" + crewName + "」" + Reset)
	sb.WriteString(strings.Repeat(" ", 30-len(data.ModelName)-len(crewName)) + OPBrown + "█" + Reset + "\n")

	sb.WriteString(OPBrown + "█" + OPCream + "                                                                           " + OPBrown + "█" + Reset + "\n")

	// Bounty
	bounty := data.TokenCount * 1000
	sb.WriteString(OPBrown + "█" + OPCream + "                      " + OPGold + fmt.Sprintf("฿ %d", bounty) + Reset)
	sb.WriteString(strings.Repeat(" ", 43-len(fmt.Sprintf("%d", bounty))) + OPBrown + "█" + Reset + "\n")

	sb.WriteString(OPBrown + "█" + OPCream + "                         " + OPBlack + "DEAD OR ALIVE" + OPCream + "                                  " + OPBrown + "█" + Reset + "\n")
	sb.WriteString(OPBrown + "█" + OPDarkBrown + "─────────────────────────────────────────────────────────────────────────" + OPBrown + "█" + Reset + "\n")

	// Stats
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("%s⚓%s", OPBlue, data.GitBranch)
	}

	sb.WriteString(OPBrown + "█" + Reset + fmt.Sprintf(" %sShip:%s %-20s %sLog:%s %s %2d%%  %sMorale:%s %s %2d%% %s",
		OPBlue, Reset, gitInfo,
		OPGold, Reset, t.generateOPBar(data.ContextPercent, 8), data.ContextPercent,
		OPRed, Reset, t.generateOPBar(100-data.API5hrPercent, 8), 100-data.API5hrPercent,
		data.API5hrTimeLeft) + OPBrown + "█" + Reset + "\n")

	sb.WriteString(OPBrown + "█" + Reset + fmt.Sprintf(" %sTime:%s %s  %sMsg:%s %d  %sBerry:%s $%s  %sRate:%s $%s/h  %sLuck:%s %d%%",
		OPDarkBrown, Reset, data.SessionTime,
		OPBlue, Reset, data.MessageCount,
		OPGold, Reset, FormatCost(data.SessionCost),
		OPRed, Reset, FormatCost(data.BurnRate),
		OPGold, Reset, data.CacheHitRate) + "      " + OPBrown + "█" + Reset + "\n")

	sb.WriteString(OPDarkBrown + "▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀" + Reset + "\n")

	return sb.String()
}

func (t *OnePieceTheme) generateOPBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	color := OPGold
	if percent > 75 {
		color = OPRed
	}

	var bar strings.Builder
	bar.WriteString(OPDarkBrown + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(OPDarkBrown)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(OPDarkBrown + "]" + Reset)
	return bar.String()
}
