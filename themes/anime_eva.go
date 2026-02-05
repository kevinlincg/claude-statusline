package themes

import (
	"fmt"
	"strings"
)

// EVATheme Evangelion NERV terminal style
type EVATheme struct{}

func init() {
	RegisterTheme(&EVATheme{})
}

func (t *EVATheme) Name() string {
	return "eva"
}

func (t *EVATheme) Description() string {
	return "EVA: NERV terminal interface with sync rate display"
}

const (
	EVAOrange = "\033[38;2;255;103;0m"
	EVAPurple = "\033[38;2;128;0;128m"
	EVAGreen  = "\033[38;2;0;255;0m"
	EVARed    = "\033[38;2;255;0;0m"
	EVABlue   = "\033[38;2;0;150;255m"
	EVAWhite  = "\033[38;2;255;255;255m"
	EVAGray   = "\033[38;2;100;100;100m"
	EVADark   = "\033[38;2;30;30;30m"
)

func (t *EVATheme) Render(data StatusData) string {
	var sb strings.Builder

	// NERV Logo Header
	sb.WriteString("\n")
	sb.WriteString("        " + EVARed + "███╗   ██╗" + EVAOrange + "███████╗" + EVARed + "██████╗ " + EVAOrange + "██╗   ██╗" + Reset + "\n")
	sb.WriteString("        " + EVARed + "████╗  ██║" + EVAOrange + "██╔════╝" + EVARed + "██╔══██╗" + EVAOrange + "██║   ██║" + Reset + "\n")
	sb.WriteString("        " + EVARed + "██╔██╗ ██║" + EVAOrange + "█████╗  " + EVARed + "██████╔╝" + EVAOrange + "██║   ██║" + Reset + "\n")
	sb.WriteString("        " + EVARed + "██║╚██╗██║" + EVAOrange + "██╔══╝  " + EVARed + "██╔══██╗" + EVAOrange + "╚██╗ ██╔╝" + Reset + "\n")
	sb.WriteString("        " + EVARed + "██║ ╚████║" + EVAOrange + "███████╗" + EVARed + "██║  ██║" + EVAOrange + " ╚████╔╝ " + Reset + "\n")
	sb.WriteString("        " + EVAGray + "╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝  ╚═══╝ " + Reset + "\n")
	sb.WriteString("      " + EVAGray + "God's in his heaven. All's right with the world." + Reset + "\n")
	sb.WriteString("\n")

	// Warning bar if context high
	if data.ContextPercent > 75 {
		sb.WriteString("  " + EVARed + "▓▓▓ WARNING ▓▓▓ PATTERN BLUE ▓▓▓ ANGEL DETECTED ▓▓▓ WARNING ▓▓▓" + Reset + "\n")
	} else {
		sb.WriteString("  " + EVAGreen + "─────────────────── SYSTEM STATUS: NOMINAL ───────────────────" + Reset + "\n")
	}
	sb.WriteString("\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	pilot := "REI"
	unit := "EVA-00"
	if data.ModelType == "Opus" {
		pilot = "SHINJI"
		unit = "EVA-01"
	} else if data.ModelType == "Haiku" {
		pilot = "ASUKA"
		unit = "EVA-02"
	}

	// Pilot Info Block
	sb.WriteString(fmt.Sprintf("  %s┌─ PILOT ─────────────────────────────────────────────────────┐%s\n", EVAOrange, Reset))
	sb.WriteString(fmt.Sprintf("  %s│%s  NAME: %s%-10s%s  UNIT: %s%-8s%s  MODEL: %s%s%-10s%s  %s│%s\n",
		EVAOrange, Reset,
		EVAWhite, pilot, Reset,
		EVAWhite, unit, Reset,
		modelColor, modelIcon, data.ModelName, Reset,
		EVAOrange, Reset))
	sb.WriteString(fmt.Sprintf("  %s└─────────────────────────────────────────────────────────────┘%s\n", EVAOrange, Reset))

	// Sync Rate Display (Context)
	syncColor := EVAGreen
	syncStatus := "STABLE"
	if data.ContextPercent > 75 {
		syncColor = EVARed
		syncStatus = "CRITICAL"
	} else if data.ContextPercent > 50 {
		syncColor = EVAOrange
		syncStatus = "ELEVATED"
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s╔═══════════════════════════════════════╗%s\n", EVAPurple, Reset))
	sb.WriteString(fmt.Sprintf("  %s║%s    %sSYNCHRONIZATION  RATE%s              %s║%s\n", EVAPurple, Reset, EVAWhite, Reset, EVAPurple, Reset))
	sb.WriteString(fmt.Sprintf("  %s║%s               %s", EVAPurple, Reset, syncColor))
	sb.WriteString(fmt.Sprintf("  %3d.%02d %%  ", data.ContextPercent, (data.ContextPercent*7)%100))
	sb.WriteString(fmt.Sprintf("%s              %s║%s\n", Reset, EVAPurple, Reset))
	sb.WriteString(fmt.Sprintf("  %s║%s  %s  %s%s%s  %s║%s\n",
		EVAPurple, Reset,
		t.generateEVABar(data.ContextPercent, 28, syncColor),
		syncColor, syncStatus, Reset, EVAPurple, Reset))
	sb.WriteString(fmt.Sprintf("  %s╚═══════════════════════════════════════╝%s\n", EVAPurple, Reset))

	// A.T. Field and Umbilical Status
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s▸ A.T. FIELD%s    ", EVAOrange, Reset))
	sb.WriteString(t.generateEVABar(100-data.API5hrPercent, 18, EVAOrange))
	sb.WriteString(fmt.Sprintf("  %s%3d%%%s  %s%s%s\n", EVAOrange, 100-data.API5hrPercent, Reset, EVAGray, data.API5hrTimeLeft, Reset))

	sb.WriteString(fmt.Sprintf("  %s▸ UMBILICAL%s     ", EVABlue, Reset))
	sb.WriteString(t.generateEVABar(100-data.API7dayPercent, 18, EVABlue))
	sb.WriteString(fmt.Sprintf("  %s%3d%%%s  %s%s%s\n", EVABlue, 100-data.API7dayPercent, Reset, EVAGray, data.API7dayTimeLeft, Reset))

	// Bottom stats
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s─────────────────────────────────────────────────────────────%s\n", EVAGray, Reset))
	sb.WriteString(fmt.Sprintf("  %sDATA:%s %s  %sTIME:%s %s  %sOPS:%s %d  %sCOST:%s $%s  %sVER:%s %s\n",
		EVAGray, Reset, FormatTokens(data.TokenCount),
		EVAGray, Reset, data.SessionTime,
		EVAGray, Reset, data.MessageCount,
		EVAGray, Reset, FormatCost(data.SessionCost),
		EVAGray, Reset, data.Version))

	return sb.String()
}

func (t *EVATheme) generateEVABar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(EVAGray + "│" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(EVADark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(EVAGray + "│" + Reset)
	return bar.String()
}
