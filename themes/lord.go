package themes

import (
	"fmt"
	"strings"
)

// LORDTheme Legend of the Red Dragon 風格
type LORDTheme struct{}

func init() {
	RegisterTheme(&LORDTheme{})
}

func (t *LORDTheme) Name() string {
	return "lord"
}

func (t *LORDTheme) Description() string {
	return "LORD：紅龍傳說 BBS 經典文字遊戲風格"
}

const (
	LORDRed      = "\033[38;2;170;0;0m"
	LORDBrightRed= "\033[38;2;255;85;85m"
	LORDGreen    = "\033[38;2;0;170;0m"
	LORDBrightGreen = "\033[38;2;85;255;85m"
	LORDYellow   = "\033[38;2;170;170;0m"
	LORDBrightYellow = "\033[38;2;255;255;85m"
	LORDBlue     = "\033[38;2;0;0;170m"
	LORDBrightBlue = "\033[38;2;85;85;255m"
	LORDMagenta  = "\033[38;2;170;0;170m"
	LORDBrightMagenta = "\033[38;2;255;85;255m"
	LORDCyan     = "\033[38;2;0;170;170m"
	LORDBrightCyan = "\033[38;2;85;255;255m"
	LORDWhite    = "\033[38;2;170;170;170m"
	LORDBrightWhite = "\033[38;2;255;255;255m"
	LORDDark     = "\033[38;2;85;85;85m"
)

func (t *LORDTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Forest clearing style header
	sb.WriteString(LORDGreen + "  ══════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	// Location banner
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	location := "The Forest"
	if data.ModelType == "Opus" {
		location = "The Dragon's Lair"
	} else if data.ModelType == "Haiku" {
		location = "The Village Square"
	}

	sb.WriteString(fmt.Sprintf("  %s%s %s%s%s%s %s- %s%s%s",
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		LORDBrightYellow, location, LORDCyan, data.Version, Reset))
	if data.UpdateAvailable {
		sb.WriteString(LORDBrightRed + " (News at Inn!)" + Reset)
	}
	sb.WriteString("\n")

	sb.WriteString(LORDGreen + "  ══════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	// Your Quest info
	sb.WriteString(fmt.Sprintf("  %sYour Quest:%s %s%s%s",
		LORDYellow, Reset,
		LORDBrightWhite, ShortenPath(data.ProjectPath, 35), Reset))
	if data.GitBranch != "" {
		sb.WriteString(fmt.Sprintf("  %s[%s%s%s]%s",
			LORDDark, LORDBrightCyan, data.GitBranch, LORDDark, Reset))
		if data.GitStaged > 0 {
			sb.WriteString(fmt.Sprintf(" %s+%d%s", LORDBrightGreen, data.GitStaged, Reset))
		}
		if data.GitDirty > 0 {
			sb.WriteString(fmt.Sprintf(" %s*%d%s", LORDBrightYellow, data.GitDirty, Reset))
		}
	}
	sb.WriteString("\n\n")

	// Stats in classic LORD style
	hitPoints := 100 - data.ContextPercent
	hpColor := LORDBrightGreen
	if hitPoints <= 20 {
		hpColor = LORDBrightRed
	} else if hitPoints <= 50 {
		hpColor = LORDBrightYellow
	}

	sb.WriteString(fmt.Sprintf("  %sHit Points:%s %s%d%s/100    %sForest Fights:%s %s%d%s/100    %sGold:%s %s%s%s\n",
		LORDCyan, Reset, hpColor, hitPoints, Reset,
		LORDCyan, Reset, LORDBrightGreen, 100-data.API5hrPercent, Reset,
		LORDCyan, Reset, LORDBrightYellow, FormatCostShort(data.DayCost), Reset))

	sb.WriteString(fmt.Sprintf("  %sExperience:%s %s%s%s       %sCharm:%s %s%d%s           %sGems:%s %s%s%s\n",
		LORDCyan, Reset, LORDBrightMagenta, FormatTokens(data.TokenCount), Reset,
		LORDCyan, Reset, LORDBrightCyan, data.CacheHitRate, Reset,
		LORDCyan, Reset, LORDBrightBlue, FormatCostShort(data.SessionCost), Reset))

	sb.WriteString("\n")

	// Progress bars in tavern menu style
	sb.WriteString(fmt.Sprintf("  %s(%s1%s)%s Vitality   %s  %s%d%%%s remaining\n",
		LORDDark, LORDBrightWhite, LORDDark, Reset,
		t.generateLORDBar(100-data.ContextPercent, 20, LORDBrightGreen),
		LORDBrightGreen, 100-data.ContextPercent, Reset))

	sb.WriteString(fmt.Sprintf("  %s(%s2%s)%s Daily Limit%s  %s%d%%%s remaining  %s%s%s\n",
		LORDDark, LORDBrightWhite, LORDDark, Reset,
		t.generateLORDBar(100-data.API5hrPercent, 20, LORDBrightCyan),
		LORDBrightCyan, 100-data.API5hrPercent, Reset,
		LORDDark, data.API5hrTimeLeft, Reset))

	sb.WriteString(fmt.Sprintf("  %s(%s3%s)%s Weekly Limit%s  %s%d%%%s remaining  %s%s%s\n",
		LORDDark, LORDBrightWhite, LORDDark, Reset,
		t.generateLORDBar(100-data.API7dayPercent, 19, LORDBrightYellow),
		LORDBrightYellow, 100-data.API7dayPercent, Reset,
		LORDDark, data.API7dayTimeLeft, Reset))

	// Footer
	sb.WriteString(LORDGreen + "  ══════════════════════════════════════════════════════════════════════════════" + Reset + "\n")
	sb.WriteString(fmt.Sprintf("  %s(%sR%s)%seturn to Town  %s(%sQ%s)%suit Game  %sTime:%s %s  %sRate:%s %s%s/h%s\n",
		LORDDark, LORDBrightWhite, LORDDark, Reset,
		LORDDark, LORDBrightWhite, LORDDark, Reset,
		LORDYellow, Reset, data.SessionTime,
		LORDYellow, Reset, LORDBrightRed, FormatCostShort(data.BurnRate), Reset))

	return sb.String()
}

func (t *LORDTheme) generateLORDBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(LORDDark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(LORDDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(LORDDark + "]" + Reset)
	return bar.String()
}
