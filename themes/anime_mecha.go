package themes

import (
	"fmt"
	"strings"
)

// MechaTheme Mecha anime cockpit HUD style
type MechaTheme struct{}

func init() {
	RegisterTheme(&MechaTheme{})
}

func (t *MechaTheme) Name() string {
	return "mecha"
}

func (t *MechaTheme) Description() string {
	return "Mecha: Giant robot cockpit HUD targeting system"
}

const (
	MCHGreen  = "\033[38;2;0;255;100m"
	MCHCyan   = "\033[38;2;0;255;255m"
	MCHYellow = "\033[38;2;255;255;0m"
	MCHRed    = "\033[38;2;255;50;50m"
	MCHWhite  = "\033[38;2;200;255;200m"
	MCHDark   = "\033[38;2;20;40;20m"
)

func (t *MechaTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Top targeting frame
	sb.WriteString(MCHGreen + "╔═══╗" + Reset + "                                                                           " + MCHGreen + "╔═══╗" + Reset + "\n")
	sb.WriteString(MCHGreen + "║" + MCHCyan + "HUD" + MCHGreen + "║" + Reset + MCHGreen + "═══════════════════════════════════════════════════════════════════════════" + MCHGreen + "║" + MCHCyan + "SYS" + MCHGreen + "║" + Reset + "\n")
	sb.WriteString(MCHGreen + "╠═══╩═══════════════════════════════════════════════════════════════════════════════╩═══╣" + Reset + "\n")

	// Title with angular brackets
	sb.WriteString(MCHGreen + "║" + Reset + "  " + MCHCyan + "《《《" + MCHWhite + " MECHA SYSTEM ONLINE " + MCHCyan + "》》》" + Reset + "   " + MCHYellow + "ロボットアニメ" + Reset + "                         " + MCHGreen + "║" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	pilot := "Newtype"
	if data.ModelType == "Opus" {
		pilot = "Ace Pilot"
	} else if data.ModelType == "Haiku" {
		pilot = "Trainee"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s<UPGRADE>%s", MCHYellow, Reset)
	}

	line1 := fmt.Sprintf("  %s▶ PILOT:%s %s%s%s  %s▶ CLASS:%s %s%s%s  %s%s%s%s",
		MCHCyan, Reset, modelColor, modelIcon, data.ModelName,
		MCHGreen, Reset, MCHYellow, pilot, Reset,
		MCHDark, data.Version, Reset, update)

	sb.WriteString(MCHGreen + "║" + Reset)
	sb.WriteString(PadRight(line1, 87))
	sb.WriteString(MCHGreen + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s◇%s%s", MCHCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", MCHGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", MCHYellow, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %s▶ MISSION:%s %s%s",
		MCHGreen, Reset, ShortenPath(data.ProjectPath, 38), gitInfo)

	sb.WriteString(MCHGreen + "║" + Reset)
	sb.WriteString(PadRight(line2, 87))
	sb.WriteString(MCHGreen + "║" + Reset + "\n")

	sb.WriteString(MCHGreen + "╠═══╦═══════════════════════════════════════════════════════════════════════════════╦═══╣" + Reset + "\n")
	sb.WriteString(MCHGreen + "║" + MCHRed + "WRN" + MCHGreen + "║" + Reset + "  " + MCHCyan + "▼ SYSTEM STATUS ▼" + Reset + "                                                          " + MCHGreen + "║" + MCHCyan + "PWR" + MCHGreen + "║" + Reset + "\n")
	sb.WriteString(MCHGreen + "╠═══╩═══════════════════════════════════════════════════════════════════════════════╩═══╣" + Reset + "\n")

	// Power levels with targeting style
	reactorColor := MCHGreen
	if data.ContextPercent > 75 {
		reactorColor = MCHRed
	}

	line3 := fmt.Sprintf("  %s├─ REACTOR%s  %s  %s%3d%%%s %s◀◀◀%s",
		MCHCyan, Reset, t.generateMCHBar(data.ContextPercent, 16, reactorColor), reactorColor, data.ContextPercent, Reset, reactorColor, Reset)

	sb.WriteString(MCHGreen + "║" + Reset)
	sb.WriteString(PadRight(line3, 87))
	sb.WriteString(MCHGreen + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %s├─ ENERGY%s   %s  %s%3d%%%s  %s%s%s",
		MCHGreen, Reset, t.generateMCHBar(100-data.API5hrPercent, 16, MCHGreen),
		MCHGreen, 100-data.API5hrPercent, Reset, MCHDark, data.API5hrTimeLeft, Reset)

	sb.WriteString(MCHGreen + "║" + Reset)
	sb.WriteString(PadRight(line4, 87))
	sb.WriteString(MCHGreen + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %s└─ ARMOR%s    %s  %s%3d%%%s  %s%s%s",
		MCHYellow, Reset, t.generateMCHBar(100-data.API7dayPercent, 16, MCHYellow),
		MCHYellow, 100-data.API7dayPercent, Reset, MCHDark, data.API7dayTimeLeft, Reset)

	sb.WriteString(MCHGreen + "║" + Reset)
	sb.WriteString(PadRight(line5, 87))
	sb.WriteString(MCHGreen + "║" + Reset + "\n")

	sb.WriteString(MCHGreen + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sDATA:%s%s%s  %sTIME:%s%s  %sSORT:%s%s%d%s  %sFUEL:%s%s$%s%s  %sRATE:%s%s$%s/h%s",
		MCHWhite, Reset, FormatTokens(data.TokenCount), Reset,
		MCHDark, Reset, data.SessionTime,
		MCHDark, Reset, MCHCyan, data.MessageCount, Reset,
		MCHYellow, Reset, MCHYellow, FormatCost(data.SessionCost), Reset,
		MCHGreen, Reset, MCHGreen, FormatCost(data.BurnRate), Reset)

	sb.WriteString(MCHGreen + "║" + Reset)
	sb.WriteString(PadRight(line6, 87))
	sb.WriteString(MCHGreen + "║" + Reset + "\n")

	sb.WriteString(MCHGreen + "╚═══╗" + Reset + "                                                                           " + MCHGreen + "╔═══╝" + Reset + "\n")
	sb.WriteString(MCHGreen + "    ╚═══════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *MechaTheme) generateMCHBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(MCHDark + "‹" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("■", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(MCHDark)
		bar.WriteString(strings.Repeat("□", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(MCHDark + "›" + Reset)
	return bar.String()
}
