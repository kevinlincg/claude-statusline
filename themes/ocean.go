package themes

import (
	"fmt"
	"strings"
)

// OceanTheme 海洋深海風格
type OceanTheme struct{}

func init() {
	RegisterTheme(&OceanTheme{})
}

func (t *OceanTheme) Name() string {
	return "ocean"
}

func (t *OceanTheme) Description() string {
	return "深海：海洋波浪漸層，寧靜藍調"
}

const (
	OceanDeep    = "\033[38;2;0;40;80m"
	OceanMid     = "\033[38;2;0;80;140m"
	OceanLight   = "\033[38;2;0;150;200m"
	OceanSurf    = "\033[38;2;100;200;255m"
	OceanFoam    = "\033[38;2;200;240;255m"
	OceanSand    = "\033[38;2;240;220;180m"
	OceanCoral   = "\033[38;2;255;127;80m"
	OceanGreen   = "\033[38;2;32;178;170m"
	OceanGold    = "\033[38;2;255;215;0m"
	OceanDim     = "\033[38;2;60;80;100m"
	OceanBgDeep  = "\033[48;2;0;30;60m"
)

func (t *OceanTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Wave top border
	sb.WriteString(OceanDeep + "~" + OceanMid + "~" + OceanLight + "≈" + OceanSurf + "≋")
	sb.WriteString(OceanFoam + "⌇⌇⌇")
	sb.WriteString(OceanSurf + "≋" + OceanLight + "≈" + OceanMid + "~" + OceanDeep + "~")
	sb.WriteString(OceanMid + strings.Repeat("~", 50))
	sb.WriteString(OceanDeep + "~" + OceanMid + "~" + OceanLight + "≈" + OceanSurf + "≋")
	sb.WriteString(OceanFoam + "⌇⌇⌇")
	sb.WriteString(OceanSurf + "≋" + OceanLight + "≈" + OceanMid + "~" + OceanDeep + "~")
	sb.WriteString(Reset + "\n")

	// Model + Version + Path (fish swimming)
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s⚡%s", OceanGold, Reset)
	}

	line1 := fmt.Sprintf(" %s><>%s %s%s%s%s%s %s%s%s%s  %s~%s  %s◈%s %s%s",
		OceanGreen, Reset,
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		OceanDim, data.Version, Reset, update,
		OceanMid, Reset,
		OceanSand, Reset, ShortenPath(data.ProjectPath, 22), Reset)
	if data.GitBranch != "" {
		line1 += fmt.Sprintf("  %s⚓%s%s", OceanLight, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line1 += fmt.Sprintf(" %s+%d%s", OceanGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line1 += fmt.Sprintf(" %s~%d%s", OceanCoral, data.GitDirty, Reset)
		}
	}
	sb.WriteString(line1)
	sb.WriteString("\n")

	// Stats line
	line2 := fmt.Sprintf(" %s><>%s %s%s%s tok  %s%d%s msg  %s%s%s  %s~%s  %s%s%s  %s%s%s  %s%s/h%s  %s%d%%hit%s",
		OceanGreen, Reset,
		OceanSurf, FormatTokens(data.TokenCount), Reset,
		OceanLight, data.MessageCount, Reset,
		OceanDim, data.SessionTime, Reset,
		OceanMid, Reset,
		OceanGreen, FormatCostShort(data.SessionCost), Reset,
		OceanSand, FormatCostShort(data.DayCost), Reset,
		OceanCoral, FormatCostShort(data.BurnRate), Reset,
		OceanGreen, data.CacheHitRate, Reset)
	sb.WriteString(line2)
	sb.WriteString("\n")

	// Depth gauges (progress bars)
	ctxBar := t.generateOceanBar(data.ContextPercent, 14)
	bar5 := t.generateOceanBar(data.API5hrPercent, 10)
	bar7 := t.generateOceanBar(data.API7dayPercent, 10)

	ctxColor := OceanGreen
	if data.ContextPercent >= 80 {
		ctxColor = OceanCoral
	} else if data.ContextPercent >= 60 {
		ctxColor = OceanGold
	}

	line3 := fmt.Sprintf(" %s><>%s %sCtx%s%s%s%3d%%%s  %s5hr%s%s%s%3d%%%s %s%s%s  %s7dy%s%s%s%3d%%%s %s%s%s",
		OceanGreen, Reset,
		OceanDim, Reset, ctxBar, ctxColor, data.ContextPercent, Reset,
		OceanDim, Reset, bar5, OceanLight, data.API5hrPercent, Reset,
		OceanDim, data.API5hrTimeLeft, Reset,
		OceanDim, Reset, bar7, OceanSurf, data.API7dayPercent, Reset,
		OceanDim, data.API7dayTimeLeft, Reset)
	sb.WriteString(line3)
	sb.WriteString("\n")

	// Wave bottom border
	sb.WriteString(OceanDeep + "~" + OceanMid + "~" + OceanLight + "≈" + OceanSurf + "≋")
	sb.WriteString(OceanFoam + "⌇⌇⌇")
	sb.WriteString(OceanSurf + "≋" + OceanLight + "≈" + OceanMid + "~" + OceanDeep + "~")
	sb.WriteString(OceanMid + strings.Repeat("~", 50))
	sb.WriteString(OceanDeep + "~" + OceanMid + "~" + OceanLight + "≈" + OceanSurf + "≋")
	sb.WriteString(OceanFoam + "⌇⌇⌇")
	sb.WriteString(OceanSurf + "≋" + OceanLight + "≈" + OceanMid + "~" + OceanDeep + "~")
	sb.WriteString(Reset + "\n")

	return sb.String()
}

func (t *OceanTheme) generateOceanBar(percent, width int) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(OceanDeep + "〔" + Reset)

	if filled > 0 {
		bar.WriteString(OceanBgDeep)
		bar.WriteString(OceanSurf)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(OceanDeep)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(OceanDeep + "〕" + Reset)
	return bar.String()
}
