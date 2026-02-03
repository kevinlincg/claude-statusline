package themes

import (
	"fmt"
	"strings"
)

// BBSTheme 經典 BBS 佈告欄風格
type BBSTheme struct{}

func init() {
	RegisterTheme(&BBSTheme{})
}

func (t *BBSTheme) Name() string {
	return "bbs"
}

func (t *BBSTheme) Description() string {
	return "BBS：經典電子佈告欄 ANSI 藝術風格"
}

const (
	BBSBlue      = "\033[38;2;0;0;170m"
	BBSBrightBlue= "\033[38;2;85;85;255m"
	BBSCyan      = "\033[38;2;0;170;170m"
	BBSBrightCyan= "\033[38;2;85;255;255m"
	BBSWhite     = "\033[38;2;170;170;170m"
	BBSBrightWhite = "\033[38;2;255;255;255m"
	BBSYellow    = "\033[38;2;170;170;0m"
	BBSBrightYellow = "\033[38;2;255;255;85m"
	BBSRed       = "\033[38;2;170;0;0m"
	BBSBrightRed = "\033[38;2;255;85;85m"
	BBSGreen     = "\033[38;2;0;170;0m"
	BBSBrightGreen = "\033[38;2;85;255;85m"
	BBSMagenta   = "\033[38;2;170;0;170m"
	BBSBrightMagenta = "\033[38;2;255;85;255m"
	BBSDark      = "\033[38;2;85;85;85m"
	BBSBgBlue    = "\033[48;2;0;0;170m"
)

func (t *BBSTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Classic BBS header with shadow box
	sb.WriteString(BBSBrightBlue + "▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄" + Reset + "\n")

	// BBS Name banner
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	bbsName := "Claude Code BBS"

	line1 := fmt.Sprintf("%s█%s %s%s  %s«« %s%s%s %s»»%s  %sSysOp: %s%s%s%s",
		BBSBrightBlue, Reset,
		BBSBgBlue+BBSBrightWhite, bbsName, Reset,
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		BBSBrightCyan, BBSBrightWhite, data.Version, Reset,
		BBSBrightBlue+"█"+Reset)
	if data.UpdateAvailable {
		line1 = fmt.Sprintf("%s█%s %s%s%s  «« %s%s%s%s%s »»  %s★ NEW FILES! ★%s",
			BBSBrightBlue, Reset,
			BBSBgBlue+BBSBrightWhite, bbsName, Reset,
			modelColor, Bold, modelIcon, data.ModelName, Reset,
			BBSBrightYellow, Reset)
	}
	sb.WriteString(PadRight(line1, 80))
	sb.WriteString("\n")

	// Menu bar
	sb.WriteString(BBSBrightBlue + "█" + BBSBrightCyan + "════════════════════════════════════════════════════════════════════════════" + BBSBrightBlue + "█" + Reset + "\n")

	// File area (project path)
	line2 := fmt.Sprintf("%s█%s  %sFile Area:%s %s%s%s",
		BBSBrightBlue, Reset,
		BBSYellow, Reset,
		BBSBrightWhite, ShortenPath(data.ProjectPath, 30), Reset)
	if data.GitBranch != "" {
		line2 += fmt.Sprintf("  %s[%s%s%s]%s", BBSDark, BBSBrightCyan, data.GitBranch, BBSDark, Reset)
		if data.GitStaged > 0 {
			line2 += fmt.Sprintf(" %s↑%d%s", BBSBrightGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line2 += fmt.Sprintf(" %s*%d%s", BBSBrightYellow, data.GitDirty, Reset)
		}
	}
	sb.WriteString(PadRight(line2, 79))
	sb.WriteString(BBSBrightBlue + "█" + Reset + "\n")

	// Stats display
	line3 := fmt.Sprintf("%s█%s  %sCalls:%s %s%s%s   %sMsgs:%s %s%d%s   %sTime:%s %s%s%s   %sCredits:%s %s%s%s",
		BBSBrightBlue, Reset,
		BBSCyan, Reset, BBSBrightWhite, FormatTokens(data.TokenCount), Reset,
		BBSCyan, Reset, BBSBrightWhite, data.MessageCount, Reset,
		BBSCyan, Reset, BBSBrightWhite, data.SessionTime, Reset,
		BBSCyan, Reset, BBSBrightYellow, FormatCostShort(data.DayCost), Reset)
	sb.WriteString(PadRight(line3, 79))
	sb.WriteString(BBSBrightBlue + "█" + Reset + "\n")

	// Separator
	sb.WriteString(BBSBrightBlue + "█" + BBSDark + "────────────────────────────────────────────────────────────────────────────" + BBSBrightBlue + "█" + Reset + "\n")

	// Status bars (BBS style ratio bars)
	ctxPct := 100 - data.ContextPercent
	ctxBar := t.generateBBSBar(ctxPct, 20)
	ctxColor := BBSBrightGreen
	if ctxPct <= 20 {
		ctxColor = BBSBrightRed
	} else if ctxPct <= 50 {
		ctxColor = BBSBrightYellow
	}

	line4 := fmt.Sprintf("%s█%s  %sSystem Load:%s  %s %s%3d%%%s",
		BBSBrightBlue, Reset,
		BBSYellow, Reset,
		ctxBar, ctxColor, ctxPct, Reset)
	sb.WriteString(PadRight(line4, 79))
	sb.WriteString(BBSBrightBlue + "█" + Reset + "\n")

	dlBar := t.generateBBSBar(100-data.API5hrPercent, 20)
	line5 := fmt.Sprintf("%s█%s  %sD/L Ratio:%s    %s %s%3d%%%s  %s(%s)%s",
		BBSBrightBlue, Reset,
		BBSYellow, Reset,
		dlBar, BBSBrightCyan, 100-data.API5hrPercent, Reset,
		BBSDark, data.API5hrTimeLeft, Reset)
	sb.WriteString(PadRight(line5, 79))
	sb.WriteString(BBSBrightBlue + "█" + Reset + "\n")

	ulBar := t.generateBBSBar(100-data.API7dayPercent, 20)
	line6 := fmt.Sprintf("%s█%s  %sU/L Ratio:%s    %s %s%3d%%%s  %s(%s)%s",
		BBSBrightBlue, Reset,
		BBSYellow, Reset,
		ulBar, BBSBrightMagenta, 100-data.API7dayPercent, Reset,
		BBSDark, data.API7dayTimeLeft, Reset)
	sb.WriteString(PadRight(line6, 79))
	sb.WriteString(BBSBrightBlue + "█" + Reset + "\n")

	// Bottom info bar
	sb.WriteString(BBSBrightBlue + "█" + BBSDark + "────────────────────────────────────────────────────────────────────────────" + BBSBrightBlue + "█" + Reset + "\n")

	line7 := fmt.Sprintf("%s█%s  %sSession:%s %s%s%s  %sRate:%s %s%s/h%s  %sHit:%s %s%d%%%s  %s[Press any key...]%s",
		BBSBrightBlue, Reset,
		BBSGreen, Reset, BBSBrightGreen, FormatCostShort(data.SessionCost), Reset,
		BBSRed, Reset, BBSBrightRed, FormatCostShort(data.BurnRate), Reset,
		BBSCyan, Reset, BBSBrightCyan, data.CacheHitRate, Reset,
		BBSDark, Reset)
	sb.WriteString(PadRight(line7, 79))
	sb.WriteString(BBSBrightBlue + "█" + Reset + "\n")

	// Footer
	sb.WriteString(BBSBrightBlue + "▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀" + Reset + "\n")

	return sb.String()
}

func (t *BBSTheme) generateBBSBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(BBSDark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(BBSBrightCyan)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(BBSDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(BBSDark + "]" + Reset)
	return bar.String()
}
