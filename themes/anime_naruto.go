package themes

import (
	"fmt"
	"strings"
)

// NarutoTheme Naruto ninja scroll style
type NarutoTheme struct{}

func init() {
	RegisterTheme(&NarutoTheme{})
}

func (t *NarutoTheme) Name() string {
	return "naruto"
}

func (t *NarutoTheme) Description() string {
	return "Naruto: Ninja scroll and chakra gauge style"
}

const (
	NarutoOrange    = "\033[38;2;255;140;0m"
	NarutoBlue      = "\033[38;2;65;105;225m"
	NarutoRed       = "\033[38;2;178;34;34m"
	NarutoGreen     = "\033[38;2;0;128;0m"
	NarutoCream     = "\033[38;2;255;248;220m"
	NarutoBrown     = "\033[38;2;139;90;43m"
	NarutoPurple    = "\033[38;2;148;0;211m"
	NarutoYellow    = "\033[38;2;255;215;0m"
	NarutoDark      = "\033[38;2;60;60;60m"
)

func (t *NarutoTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Scroll top
	sb.WriteString(NarutoBrown + "  ╭─────────────────────────────────────────────────────────────────────────────────╮" + Reset + "\n")
	sb.WriteString(NarutoBrown + "══╡" + NarutoCream + "                              " + NarutoOrange + "忍" + NarutoCream + " NINJA STATUS " + NarutoOrange + "忍" + NarutoCream + "                              " + NarutoBrown + "╞══" + Reset + "\n")
	sb.WriteString(NarutoBrown + "  ├─────────────────────────────────────────────────────────────────────────────────┤" + Reset + "\n")

	// Ninja info
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	rank := "Genin"
	if data.ModelType == "Opus" {
		rank = "Hokage"
	} else if data.ModelType == "Sonnet" {
		rank = "Jonin"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[New Jutsu!]%s", NarutoYellow, Reset)
	}

	line1 := fmt.Sprintf("  %s│%s  %s%s%s %s  %sRank:%s %s%s%s  %sv%s%s%s",
		NarutoBrown, Reset,
		modelColor, modelIcon, data.ModelName,
		Reset,
		NarutoOrange, Reset, NarutoYellow, rank, Reset,
		NarutoDark, data.Version, Reset, update)

	sb.WriteString(line1 + "\n")

	// Mission (project)
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚡%s%s", NarutoBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", NarutoGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", NarutoOrange, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %s│%s  %sMission:%s %s%s",
		NarutoBrown, Reset,
		NarutoRed, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(NarutoBrown + "  ├─────────────────────────────────────────────────────────────────────────────────┤" + Reset + "\n")

	// Chakra gauge
	chakraColor := NarutoBlue
	if data.ContextPercent > 75 {
		chakraColor = NarutoRed
	} else if data.ContextPercent > 50 {
		chakraColor = NarutoOrange
	}

	line3 := fmt.Sprintf("  %s│%s  %sChakra%s     %s  %s%3d%%%s",
		NarutoBrown, Reset,
		NarutoBlue, Reset,
		t.generateNarutoBar(data.ContextPercent, 20, chakraColor),
		chakraColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	// Stamina (5hr)
	line4 := fmt.Sprintf("  %s│%s  %sStamina%s    %s  %s%3d%%%s  %s%s%s",
		NarutoBrown, Reset,
		NarutoGreen, Reset,
		t.generateNarutoBar(100-data.API5hrPercent, 20, NarutoGreen),
		NarutoGreen, 100-data.API5hrPercent, Reset,
		NarutoDark, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	// Will of Fire (7day)
	line5 := fmt.Sprintf("  %s│%s  %sWill%s       %s  %s%3d%%%s  %s%s%s",
		NarutoBrown, Reset,
		NarutoOrange, Reset,
		t.generateNarutoBar(100-data.API7dayPercent, 20, NarutoOrange),
		NarutoOrange, 100-data.API7dayPercent, Reset,
		NarutoDark, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(NarutoBrown + "  ├─────────────────────────────────────────────────────────────────────────────────┤" + Reset + "\n")

	// Stats
	line6 := fmt.Sprintf("  %s│%s  %sJutsu:%s %s%s%s  %sTime:%s %s  %sMissions:%s %s%d%s  %sRyo:%s %s%s%s  %sDaily:%s %s%s%s",
		NarutoBrown, Reset,
		NarutoPurple, Reset, NarutoPurple, FormatTokens(data.TokenCount), Reset,
		NarutoDark, Reset, data.SessionTime,
		NarutoDark, Reset, NarutoCream, data.MessageCount, Reset,
		NarutoYellow, Reset, NarutoYellow, FormatCost(data.SessionCost), Reset,
		NarutoOrange, Reset, NarutoOrange, FormatCost(data.DayCost), Reset)
	sb.WriteString(line6 + "\n")

	line7 := fmt.Sprintf("  %s│%s  %sRate:%s %s%s/h%s  %sAccuracy:%s %s%d%%%s",
		NarutoBrown, Reset,
		NarutoRed, Reset, NarutoRed, FormatCost(data.BurnRate), Reset,
		NarutoGreen, Reset, NarutoGreen, data.CacheHitRate, Reset)
	sb.WriteString(line7 + "\n")

	// Scroll bottom
	sb.WriteString(NarutoBrown + "  ├─────────────────────────────────────────────────────────────────────────────────┤" + Reset + "\n")
	sb.WriteString(NarutoBrown + "══╡" + NarutoCream + "                                    " + NarutoOrange + "木ノ葉" + NarutoCream + "                                    " + NarutoBrown + "╞══" + Reset + "\n")
	sb.WriteString(NarutoBrown + "  ╰─────────────────────────────────────────────────────────────────────────────────╯" + Reset + "\n")

	return sb.String()
}

func (t *NarutoTheme) generateNarutoBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(NarutoBrown + "〔" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▰", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(NarutoDark)
		bar.WriteString(strings.Repeat("▱", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(NarutoBrown + "〕" + Reset)
	return bar.String()
}
