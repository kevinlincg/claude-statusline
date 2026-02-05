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
	return "One Piece: Wanted poster and pirate log style"
}

const (
	OPBrown       = "\033[38;2;139;90;43m"
	OPGold        = "\033[38;2;255;215;0m"
	OPRed         = "\033[38;2;180;30;30m"
	OPCream       = "\033[38;2;255;248;220m"
	OPBlue        = "\033[38;2;0;105;148m"
	OPDarkBrown   = "\033[38;2;101;67;33m"
	OPBlack       = "\033[38;2;30;30;30m"
)

func (t *OnePieceTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Wanted poster style header
	sb.WriteString(OPBrown + "╔═══════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")

	// WANTED banner
	modelColor, modelIcon := GetModelConfig(data.ModelType)

	line1 := fmt.Sprintf("                         %s%sW A N T E D%s", OPRed, Bold, Reset)
	sb.WriteString(OPBrown + "║" + Reset)
	sb.WriteString(PadRight(line1, 85))
	sb.WriteString(OPBrown + "║" + Reset + "\n")

	// Model as pirate name
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s★NEW★%s", OPGold, Reset)
	}
	line2 := fmt.Sprintf("                    %s%s%s %s%s%s",
		modelColor, modelIcon, data.ModelName, OPDarkBrown, data.Version, Reset) + update

	sb.WriteString(OPBrown + "║" + Reset)
	sb.WriteString(PadRight(line2, 85))
	sb.WriteString(OPBrown + "║" + Reset + "\n")

	sb.WriteString(OPBrown + "╠═══════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Bounty (tokens)
	bounty := data.TokenCount * 1000 // Make it look impressive
	line3 := fmt.Sprintf("  %sBOUNTY:%s %s฿ %d%s",
		OPRed, Reset, OPGold, bounty, Reset)

	sb.WriteString(OPBrown + "║" + Reset)
	sb.WriteString(PadRight(line3, 85))
	sb.WriteString(OPBrown + "║" + Reset + "\n")

	// Crew (project) and ship (branch)
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %sShip:%s %s%s%s", OPBlue, Reset, OPBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", OPGold, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", OPRed, data.GitDirty, Reset)
		}
	}

	line4 := fmt.Sprintf("  %sCrew:%s %s%s", OPDarkBrown, Reset, ShortenPath(data.ProjectPath, 35), gitInfo)

	sb.WriteString(OPBrown + "║" + Reset)
	sb.WriteString(PadRight(line4, 85))
	sb.WriteString(OPBrown + "║" + Reset + "\n")

	sb.WriteString(OPBrown + "╠═══════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Log Pose (progress bars)
	line5 := fmt.Sprintf("  %sLog Pose%s %s %s%3d%%%s  %sMorale%s %s %s%3d%%%s %s%s%s",
		OPBlue, Reset,
		t.generateOPBar(data.ContextPercent, 12),
		OPBlue, data.ContextPercent, Reset,
		OPGold, Reset,
		t.generateOPBar(100-data.API5hrPercent, 12),
		OPGold, 100-data.API5hrPercent, Reset,
		OPDarkBrown, data.API5hrTimeLeft, Reset)

	sb.WriteString(OPBrown + "║" + Reset)
	sb.WriteString(PadRight(line5, 85))
	sb.WriteString(OPBrown + "║" + Reset + "\n")

	// Provisions (weekly limit)
	line6 := fmt.Sprintf("  %sProvisions%s %s %s%3d%%%s %s%s%s",
		OPRed, Reset,
		t.generateOPBar(100-data.API7dayPercent, 12),
		OPRed, 100-data.API7dayPercent, Reset,
		OPDarkBrown, data.API7dayTimeLeft, Reset)

	sb.WriteString(OPBrown + "║" + Reset)
	sb.WriteString(PadRight(line6, 85))
	sb.WriteString(OPBrown + "║" + Reset + "\n")

	sb.WriteString(OPBrown + "╠═══════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Treasure stats
	line7 := fmt.Sprintf("  %sVoyage%s %s  %sBattles%s %s%d%s  %sBerry%s %s฿%s%s  %sDaily%s %s฿%s%s  %sRate%s %s฿%s/h%s  %sLuck%s %s%d%%%s",
		OPBlue, Reset, data.SessionTime,
		OPRed, Reset, OPCream, data.MessageCount, Reset,
		OPGold, Reset, OPGold, FormatCost(data.SessionCost), Reset,
		OPBrown, Reset, OPBrown, FormatCost(data.DayCost), Reset,
		OPRed, Reset, OPRed, FormatCost(data.BurnRate), Reset,
		OPGold, Reset, OPGold, data.CacheHitRate, Reset)

	sb.WriteString(OPBrown + "║" + Reset)
	sb.WriteString(PadRight(line7, 85))
	sb.WriteString(OPBrown + "║" + Reset + "\n")

	sb.WriteString(OPBrown + "╚═══════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

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
	if percent < 30 {
		color = OPRed
	} else if percent < 60 {
		color = OPBrown
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
