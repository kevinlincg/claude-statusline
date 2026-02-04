package themes

import (
	"fmt"
	"strings"
)

// SynthwaveTheme 合成波霓虹日落風格
type SynthwaveTheme struct{}

func init() {
	RegisterTheme(&SynthwaveTheme{})
}

func (t *SynthwaveTheme) Name() string {
	return "synthwave"
}

func (t *SynthwaveTheme) Description() string {
	return "合成波：霓虹日落漸層，80年代復古未來"
}

const (
	SynthPink    = "\033[38;2;255;41;117m"
	SynthCyan    = "\033[38;2;0;255;255m"
	SynthPurple  = "\033[38;2;140;30;255m"
	SynthOrange  = "\033[38;2;255;144;31m"
	SynthYellow  = "\033[38;2;255;211;25m"
	SynthMagenta = "\033[38;2;242;34;255m"
	SynthDim     = "\033[38;2;100;60;120m"
	SynthBgPink  = "\033[48;2;60;10;30m"
	SynthBgCyan  = "\033[48;2;0;40;50m"
)

func (t *SynthwaveTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Top gradient border (sunset colors)
	sb.WriteString(SynthYellow + "▄" + SynthOrange + "▄▄" + SynthPink + "▄▄▄▄▄" + SynthMagenta + "▄▄▄▄▄▄▄▄" + SynthPurple + "▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄" + Reset)
	sb.WriteString("\n")

	// Model + Version + Path + Git
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s⬆%s", SynthOrange, Reset)
	}

	line1 := fmt.Sprintf(" %s%s%s%s%s %s%s%s%s  %s░%s  %s📂 %s%s",
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		SynthCyan, data.Version, Reset, update,
		SynthDim, Reset,
		SynthYellow, ShortenPath(data.ProjectPath, 20), Reset)
	if data.GitBranch != "" {
		line1 += fmt.Sprintf("  %s⚡%s%s", SynthPink, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line1 += fmt.Sprintf(" %s+%d%s", SynthCyan, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line1 += fmt.Sprintf(" %s~%d%s", SynthOrange, data.GitDirty, Reset)
		}
	}
	sb.WriteString(line1)
	sb.WriteString("\n")

	// Stats with neon colors
	line2 := fmt.Sprintf(" %s%s%s tok  %s%d%s msg  %s%s%s  %s░%s  %s%s%s ses  %s%s%s day  %s%s/h%s  %s%d%%hit%s",
		SynthPurple, FormatTokens(data.TokenCount), Reset,
		SynthCyan, data.MessageCount, Reset,
		SynthDim, data.SessionTime, Reset,
		SynthDim, Reset,
		SynthCyan, FormatCostShort(data.SessionCost), Reset,
		SynthYellow, FormatCostShort(data.DayCost), Reset,
		SynthPink, FormatCostShort(data.BurnRate), Reset,
		SynthCyan, data.CacheHitRate, Reset)
	sb.WriteString(line2)
	sb.WriteString("\n")

	// Neon progress bars
	ctxBar := t.generateNeonBar(data.ContextPercent, 14, SynthCyan, SynthBgCyan)
	bar5 := t.generateNeonBar(data.API5hrPercent, 10, SynthPink, SynthBgPink)
	bar7 := t.generateNeonBar(data.API7dayPercent, 10, SynthMagenta, SynthBgPink)

	line3 := fmt.Sprintf(" %sCtx%s%s%s%3d%%%s  %s░%s  %s5hr%s%s%s%3d%%%s %s%s%s  %s░%s  %s7dy%s%s%s%3d%%%s %s%s%s",
		SynthDim, Reset, ctxBar, SynthCyan, data.ContextPercent, Reset,
		SynthDim, Reset,
		SynthDim, Reset, bar5, SynthPink, data.API5hrPercent, Reset,
		SynthDim, data.API5hrTimeLeft, Reset,
		SynthDim, Reset,
		SynthDim, Reset, bar7, SynthMagenta, data.API7dayPercent, Reset,
		SynthDim, data.API7dayTimeLeft, Reset)
	sb.WriteString(line3)
	sb.WriteString("\n")

	// Bottom gradient border (reverse sunset)
	sb.WriteString(SynthPurple + "▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀" + SynthMagenta + "▀▀▀▀▀▀▀▀" + SynthPink + "▀▀▀▀▀" + SynthOrange + "▀▀" + SynthYellow + "▀" + Reset)
	sb.WriteString("\n")

	return sb.String()
}

func (t *SynthwaveTheme) generateNeonBar(percent, width int, color, bgColor string) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(bgColor)
		bar.WriteString(Bold)
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(SynthDim)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}
