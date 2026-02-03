package themes

import (
	"fmt"
	"strings"
)

// TradeWarsTheme Trade Wars 2002 太空風格
type TradeWarsTheme struct{}

func init() {
	RegisterTheme(&TradeWarsTheme{})
}

func (t *TradeWarsTheme) Name() string {
	return "tradewars"
}

func (t *TradeWarsTheme) Description() string {
	return "Trade Wars：太空貿易遊戲，星艦控制台風格"
}

const (
	TWBlack   = "\033[38;2;0;0;0m"
	TWBlue    = "\033[38;2;0;0;170m"
	TWGreen   = "\033[38;2;0;170;0m"
	TWCyan    = "\033[38;2;0;170;170m"
	TWRed     = "\033[38;2;170;0;0m"
	TWMagenta = "\033[38;2;170;0;170m"
	TWBrown   = "\033[38;2;170;85;0m"
	TWGray    = "\033[38;2;170;170;170m"
	TWDark    = "\033[38;2;85;85;85m"
	TWBrightBlue = "\033[38;2;85;85;255m"
	TWBrightGreen = "\033[38;2;85;255;85m"
	TWBrightCyan = "\033[38;2;85;255;255m"
	TWBrightRed = "\033[38;2;255;85;85m"
	TWBrightMagenta = "\033[38;2;255;85;255m"
	TWYellow  = "\033[38;2;255;255;85m"
	TWWhite   = "\033[38;2;255;255;255m"
	TWBgBlue  = "\033[48;2;0;0;170m"
)

func (t *TradeWarsTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Starfield top border
	sb.WriteString(TWDark + "·  ˚  ·    ✦    ·  ˚ " + TWBrightCyan + "╔═══════════════════════════════════════════════════════╗" + TWDark + " ˚  ·    ✦  ˚  ·" + Reset + "\n")

	// Ship computer header
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	shipClass := "Scout"
	if data.ModelType == "Opus" {
		shipClass = "Imperial StarShip"
	} else if data.ModelType == "Haiku" {
		shipClass = "Merchant Freighter"
	}

	line1 := fmt.Sprintf("%s ✦  · %s║%s %s◆ %s%s%s%s %s[%s]%s  %sSector: %s%s%s",
		TWDark, TWBrightCyan, Reset,
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		TWDark, shipClass, Reset,
		TWYellow, TWBrightGreen, data.Version, Reset)
	if data.UpdateAvailable {
		line1 += TWBrightRed + " !!ALERT!!" + Reset
	}
	sb.WriteString(PadRight(line1, 82))
	sb.WriteString(TWBrightCyan + "║" + TWDark + " ·  ✦" + Reset + "\n")

	// Navigation display
	line2 := fmt.Sprintf("%s˚    ·%s║%s %sNav:%s %s%s%s",
		TWDark, TWBrightCyan, Reset,
		TWCyan, Reset,
		TWWhite, ShortenPath(data.ProjectPath, 30), Reset)
	if data.GitBranch != "" {
		line2 += fmt.Sprintf("  %s►%s%s", TWBrightMagenta, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line2 += fmt.Sprintf(" %s↑%d%s", TWBrightGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line2 += fmt.Sprintf(" %s*%d%s", TWYellow, data.GitDirty, Reset)
		}
	}
	sb.WriteString(PadRight(line2, 82))
	sb.WriteString(TWBrightCyan + "║" + TWDark + "·   ˚" + Reset + "\n")

	// Separator with stars
	sb.WriteString(TWDark + "    ✦  " + TWBrightCyan + "╠═══════════════════════════════════════════════════════╣" + TWDark + "  ✦    " + Reset + "\n")

	// Ship status bars
	shields := 100 - data.ContextPercent
	shieldColor := TWBrightGreen
	if shields <= 20 {
		shieldColor = TWBrightRed
	} else if shields <= 50 {
		shieldColor = TWYellow
	}

	shieldBar := t.generateTWBar(shields, 15, shieldColor)
	fuelBar := t.generateTWBar(100-data.API5hrPercent, 12, TWBrightCyan)
	holdBar := t.generateTWBar(100-data.API7dayPercent, 12, TWBrightMagenta)

	line3 := fmt.Sprintf("%s·  ˚  %s║%s %sShields:%s%s%s%d%%%s  %sFuel:%s%s%s%d%%%s",
		TWDark, TWBrightCyan, Reset,
		TWCyan, Reset, shieldBar, shieldColor, shields, Reset,
		TWCyan, Reset, fuelBar, TWBrightCyan, 100-data.API5hrPercent, Reset)
	sb.WriteString(PadRight(line3, 82))
	sb.WriteString(TWBrightCyan + "║" + TWDark + "  ˚  ·" + Reset + "\n")

	line4 := fmt.Sprintf("%s ✦    %s║%s %sCargo:%s %s%s  %sHolds:%s%s%s%d%%%s  %sTurns:%s %s%s%s",
		TWDark, TWBrightCyan, Reset,
		TWCyan, Reset, TWYellow, FormatTokens(data.TokenCount),
		TWCyan, Reset, holdBar, TWBrightMagenta, 100-data.API7dayPercent, Reset,
		TWCyan, Reset, TWGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(PadRight(line4, 82))
	sb.WriteString(TWBrightCyan + "║" + TWDark + "    ✦ " + Reset + "\n")

	// Credits and stats
	line5 := fmt.Sprintf("%s˚     %s║%s %sCredits:%s %s%s%s  %sFighters:%s %s%d%s  %sRate:%s %s%s/h%s",
		TWDark, TWBrightCyan, Reset,
		TWYellow, Reset, TWBrightGreen, FormatCostShort(data.DayCost), Reset,
		TWRed, Reset, TWBrightRed, data.MessageCount, Reset,
		TWCyan, Reset, TWBrightRed, FormatCostShort(data.BurnRate), Reset)
	sb.WriteString(PadRight(line5, 82))
	sb.WriteString(TWBrightCyan + "║" + TWDark + "     ˚" + Reset + "\n")

	line6 := fmt.Sprintf("%s  ·   %s║%s %sSession:%s %s%s%s  %sTime:%s %s%s%s  %sExp:%s %s%d%%%s",
		TWDark, TWBrightCyan, Reset,
		TWGreen, Reset, TWBrightGreen, FormatCostShort(data.SessionCost), Reset,
		TWGray, Reset, TWWhite, data.SessionTime, Reset,
		TWMagenta, Reset, TWBrightMagenta, data.CacheHitRate, Reset)
	sb.WriteString(PadRight(line6, 82))
	sb.WriteString(TWBrightCyan + "║" + TWDark + "   ·  " + Reset + "\n")

	// Starfield bottom border
	sb.WriteString(TWDark + "·  ˚  ·    ✦    ·  ˚ " + TWBrightCyan + "╚═══════════════════════════════════════════════════════╝" + TWDark + " ˚  ·    ✦  ˚  ·" + Reset + "\n")

	return sb.String()
}

func (t *TradeWarsTheme) generateTWBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(TWDark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(TWDark)
		bar.WriteString(strings.Repeat("·", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(TWDark + "]" + Reset)
	return bar.String()
}
