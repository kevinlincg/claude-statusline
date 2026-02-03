package themes

import (
	"fmt"
	"strings"
)

// BtopTheme btop++ 現代系統監視器風格
type BtopTheme struct{}

func init() {
	RegisterTheme(&BtopTheme{})
}

func (t *BtopTheme) Name() string {
	return "btop"
}

func (t *BtopTheme) Description() string {
	return "btop：現代系統監視器，漸層色彩與圓角框風格"
}

const (
	// btop uses a modern color scheme with gradients
	BtopBg        = "\033[48;2;16;16;32m"
	BtopFg        = "\033[38;2;200;200;220m"
	BtopDim       = "\033[38;2;80;80;100m"
	BtopBorder    = "\033[38;2;100;100;140m"
	BtopTitle     = "\033[38;2;180;180;220m"
	BtopCyan      = "\033[38;2;0;220;220m"
	BtopMagenta   = "\033[38;2;220;80;220m"
	BtopPink      = "\033[38;2;255;100;150m"
	BtopPurple    = "\033[38;2;160;100;220m"
	BtopBlue      = "\033[38;2;80;140;255m"
	BtopGreen     = "\033[38;2;80;220;120m"
	BtopYellow    = "\033[38;2;255;220;80m"
	BtopOrange    = "\033[38;2;255;150;50m"
	BtopRed       = "\033[38;2;255;80;100m"
	BtopWhite     = "\033[38;2;240;240;255m"
	BtopGrad1     = "\033[38;2;80;200;255m"  // cyan
	BtopGrad2     = "\033[38;2;150;120;255m" // purple
	BtopGrad3     = "\033[38;2;255;100;200m" // pink
)

func (t *BtopTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Rounded top border with title
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	title := fmt.Sprintf(" %s%s%s %s ", modelIcon, data.ModelName, Reset, data.Version)
	if data.UpdateAvailable {
		title += BtopYellow + "⬆ " + Reset
	}

	topLeft := BtopBorder + "╭" + Reset
	topRight := BtopBorder + "╮" + Reset
	titleLen := len(modelIcon) + len(data.ModelName) + len(data.Version) + 4
	if data.UpdateAvailable {
		titleLen += 2
	}
	padding := 78 - titleLen
	leftPad := padding / 2
	rightPad := padding - leftPad

	sb.WriteString(topLeft + BtopBorder + strings.Repeat("─", leftPad) + Reset)
	sb.WriteString(modelColor + title)
	sb.WriteString(BtopBorder + strings.Repeat("─", rightPad) + Reset + topRight + "\n")

	// CPU section with gradient bar
	cpuPct := data.ContextPercent
	cpuBar := t.generateBtopBar(cpuPct, 30)
	cpuColor := BtopGreen
	if cpuPct > 75 {
		cpuColor = BtopRed
	} else if cpuPct > 50 {
		cpuColor = BtopYellow
	}

	line1 := fmt.Sprintf("%s│%s %sCPU%s  %s %s%3d%%%s  %sThreads:%s %s%d%s  %sLoad:%s %s%s%s",
		BtopBorder, Reset,
		BtopCyan, Reset,
		cpuBar, cpuColor, cpuPct, Reset,
		BtopDim, Reset, BtopFg, data.MessageCount, Reset,
		BtopDim, Reset, BtopFg, FormatTokens(data.TokenCount), Reset)
	sb.WriteString(PadRight(line1, 79))
	sb.WriteString(BtopBorder + "│" + Reset + "\n")

	// Memory section
	memPct := data.API5hrPercent
	memBar := t.generateBtopBar(memPct, 30)
	line2 := fmt.Sprintf("%s│%s %sMEM%s  %s %s%3d%%%s  %s5hr:%s %s%s%s",
		BtopBorder, Reset,
		BtopMagenta, Reset,
		memBar, BtopMagenta, memPct, Reset,
		BtopDim, Reset, BtopFg, data.API5hrTimeLeft, Reset)
	sb.WriteString(PadRight(line2, 79))
	sb.WriteString(BtopBorder + "│" + Reset + "\n")

	// Network/Disk section style
	netPct := data.API7dayPercent
	netBar := t.generateBtopBar(netPct, 30)
	line3 := fmt.Sprintf("%s│%s %sNET%s  %s %s%3d%%%s  %s7day:%s %s%s%s",
		BtopBorder, Reset,
		BtopPurple, Reset,
		netBar, BtopPurple, netPct, Reset,
		BtopDim, Reset, BtopFg, data.API7dayTimeLeft, Reset)
	sb.WriteString(PadRight(line3, 79))
	sb.WriteString(BtopBorder + "│" + Reset + "\n")

	// Middle separator
	sb.WriteString(BtopBorder + "├" + strings.Repeat("─", 78) + "┤" + Reset + "\n")

	// Project info
	line4 := fmt.Sprintf("%s│%s %s⌂%s %s%s%s",
		BtopBorder, Reset,
		BtopBlue, Reset,
		BtopFg, ShortenPath(data.ProjectPath, 40), Reset)
	if data.GitBranch != "" {
		line4 += fmt.Sprintf("  %s%s%s", BtopCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line4 += fmt.Sprintf(" %s+%d%s", BtopGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line4 += fmt.Sprintf(" %s~%d%s", BtopOrange, data.GitDirty, Reset)
		}
	}
	sb.WriteString(PadRight(line4, 79))
	sb.WriteString(BtopBorder + "│" + Reset + "\n")

	// Stats row with modern icons
	line5 := fmt.Sprintf("%s│%s %s⏱%s %s%s%s  %s💰%s %s%s%s  %s🔥%s %s%s/h%s  %s⚡%s %s%d%%%s  %s📊%s %s%s%s",
		BtopBorder, Reset,
		BtopDim, Reset, BtopWhite, data.SessionTime, Reset,
		BtopYellow, Reset, BtopGreen, FormatCostShort(data.SessionCost), Reset,
		BtopOrange, Reset, BtopRed, FormatCostShort(data.BurnRate), Reset,
		BtopCyan, Reset, BtopCyan, data.CacheHitRate, Reset,
		BtopPink, Reset, BtopYellow, FormatCostShort(data.DayCost), Reset)
	sb.WriteString(PadRight(line5, 79))
	sb.WriteString(BtopBorder + "│" + Reset + "\n")

	// Rounded bottom border
	sb.WriteString(BtopBorder + "╰" + strings.Repeat("─", 78) + "╯" + Reset + "\n")

	return sb.String()
}

func (t *BtopTheme) generateBtopBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(BtopDim + "[" + Reset)

	// Gradient effect: cyan -> purple -> pink
	for i := 0; i < filled; i++ {
		pos := float64(i) / float64(width)
		if pos < 0.33 {
			bar.WriteString(BtopGrad1)
		} else if pos < 0.66 {
			bar.WriteString(BtopGrad2)
		} else {
			bar.WriteString(BtopGrad3)
		}
		bar.WriteString("━")
	}
	if empty > 0 {
		bar.WriteString(BtopDim)
		bar.WriteString(strings.Repeat("╌", empty))
	}
	bar.WriteString(Reset)
	bar.WriteString(BtopDim + "]" + Reset)
	return bar.String()
}
