package themes

import (
	"fmt"
	"strings"
)

// RetroCRTTheme 復古 CRT 顯示器風格
type RetroCRTTheme struct{}

func init() {
	RegisterTheme(&RetroCRTTheme{})
}

func (t *RetroCRTTheme) Name() string {
	return "retro_crt"
}

func (t *RetroCRTTheme) Description() string {
	return "復古 CRT：綠色磷光螢幕，掃描線效果"
}

const (
	CRTGreen      = "\033[38;2;51;255;51m"
	CRTDarkGreen  = "\033[38;2;0;180;0m"
	CRTDimGreen   = "\033[38;2;0;100;0m"
	CRTBrightGreen= "\033[38;2;180;255;180m"
	CRTBgGlow     = "\033[48;2;0;40;0m"
)

func (t *RetroCRTTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Scanline top
	sb.WriteString(CRTDimGreen)
	sb.WriteString("▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// Model + Version + Path + Git
	_, modelIcon := GetModelConfig(data.ModelType)
	update := ""
	if data.UpdateAvailable {
		update = CRTBrightGreen + " [UPDATE]" + Reset
	}

	line1 := fmt.Sprintf(" %s>%s %s%s%s%s%s %s%s%s%s  %s|%s  %s%s%s",
		CRTDarkGreen, Reset,
		CRTGreen, Bold, modelIcon, data.ModelName, Reset,
		CRTDimGreen, data.Version, Reset, update,
		CRTDimGreen, Reset,
		CRTGreen, ShortenPath(data.ProjectPath, 25), Reset)
	if data.GitBranch != "" {
		line1 += fmt.Sprintf("  %s<%s>%s", CRTGreen, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line1 += fmt.Sprintf(" %s+%d%s", CRTBrightGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line1 += fmt.Sprintf(" %s*%d%s", CRTGreen, data.GitDirty, Reset)
		}
	}
	sb.WriteString(line1)
	sb.WriteString("\n")

	// Stats line
	line2 := fmt.Sprintf(" %s>%s %sTOK:%s%s%s  %sMSG:%s%d%s  %sTIME:%s%s%s  %s|%s  %sSES:%s%s%s  %sDAY:%s%s%s  %sRATE:%s%s/h%s  %sHIT:%s%d%%%s",
		CRTDarkGreen, Reset,
		CRTDimGreen, CRTGreen, FormatTokens(data.TokenCount), Reset,
		CRTDimGreen, CRTGreen, data.MessageCount, Reset,
		CRTDimGreen, CRTGreen, data.SessionTime, Reset,
		CRTDimGreen, Reset,
		CRTDimGreen, CRTGreen, FormatCostShort(data.SessionCost), Reset,
		CRTDimGreen, CRTGreen, FormatCostShort(data.DayCost), Reset,
		CRTDimGreen, CRTGreen, FormatCostShort(data.BurnRate), Reset,
		CRTDimGreen, CRTGreen, data.CacheHitRate, Reset)
	sb.WriteString(line2)
	sb.WriteString("\n")

	// Progress bars with CRT glow effect
	ctxBar := t.generateCRTBar(data.ContextPercent, 15)
	bar5 := t.generateCRTBar(data.API5hrPercent, 10)
	bar7 := t.generateCRTBar(data.API7dayPercent, 10)

	line3 := fmt.Sprintf(" %s>%s %sCTX%s%s%s%3d%%%s  %s5HR%s%s%s%3d%%%s %s%s%s  %s7DY%s%s%s%3d%%%s %s%s%s",
		CRTDarkGreen, Reset,
		CRTDimGreen, Reset, ctxBar, CRTGreen, data.ContextPercent, Reset,
		CRTDimGreen, Reset, bar5, CRTGreen, data.API5hrPercent, Reset,
		CRTDimGreen, data.API5hrTimeLeft, Reset,
		CRTDimGreen, Reset, bar7, CRTGreen, data.API7dayPercent, Reset,
		CRTDimGreen, data.API7dayTimeLeft, Reset)
	sb.WriteString(line3)
	sb.WriteString("\n")

	// Scanline bottom
	sb.WriteString(CRTDimGreen)
	sb.WriteString("▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	return sb.String()
}

func (t *RetroCRTTheme) generateCRTBar(percent, width int) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	bar.WriteString("[")
	if filled > 0 {
		bar.WriteString(CRTBgGlow)
		bar.WriteString(CRTGreen)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(CRTDimGreen)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString("]")
	return bar.String()
}
