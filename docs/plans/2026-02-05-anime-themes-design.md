# Anime Themes Design Document

Date: 2026-02-05

## Overview

Add 40 Japanese anime-themed status line themes to the project, mixing direct tributes to classic works with general anime aesthetics.

## Theme List

### Shonen / Action (12 themes)

| # | Name | Anime | Visual Concept | Key Elements |
|---|------|-------|----------------|--------------|
| 1 | `eva` | Evangelion | NERV warning interface | Orange/red alerts, "SYNC RATE", "A.T. FIELD", hexagonal patterns |
| 2 | `dragonball` | Dragon Ball | Scouter display | Green scan lines, "POWER LEVEL", numeric readout |
| 3 | `naruto` | Naruto | Ninja scroll | Chakra gauge, leaf symbol, scroll borders |
| 4 | `onepiece` | One Piece | Wanted poster | Berry symbol, bounty style, pirate flag |
| 5 | `bleach` | Bleach | Reiatsu display | Black/white contrast, spiritual pressure meter |
| 6 | `aot` | Attack on Titan | Survey Corps report | Military style, wall status, ODM gear |
| 7 | `demonslayer` | Demon Slayer | Breathing display | Nichirin blade patterns, breath styles |
| 8 | `jujutsu` | Jujutsu Kaisen | Cursed energy | Domain expansion, curse grade |
| 9 | `mha` | My Hero Academia | Quirk analysis | Plus Ultra, hero ranking |
| 10 | `hxh` | Hunter x Hunter | Nen system | Hunter license, Greed Island |
| 11 | `fma` | Fullmetal Alchemist | Transmutation circle | Equivalent exchange, state alchemist |
| 12 | `jojo` | JoJo's Bizarre Adventure | Stand stats | Power/Speed/Range/Durability, onomatopoeia |

### Classic / Legendary (6 themes)

| # | Name | Anime | Visual Concept | Key Elements |
|---|------|-------|----------------|--------------|
| 13 | `deathnote` | Death Note | Notebook page | Gothic font, shinigami text, rule list |
| 14 | `sailormoon` | Sailor Moon | Transformation | Moon symbol, pastel gradients, sparkles |
| 15 | `bebop` | Cowboy Bebop | Bounty file | Space jazz, "See You Space Cowboy" |
| 16 | `gits` | Ghost in the Shell | Cyberbrain | Section 9, data stream, hack green |
| 17 | `akira` | Akira | Neo-Tokyo | Red warning, experiment number, psychic |
| 18 | `gundam` | Gundam | MS cockpit | Federation/Zeon, ammo status, damage |

### Ghibli Studio (7 themes)

| # | Name | Film | Visual Concept | Key Elements |
|---|------|------|----------------|--------------|
| 19 | `totoro` | My Neighbor Totoro | Forest spirit | Soft green, Totoro silhouette, Catbus |
| 20 | `spirited` | Spirited Away | Bathhouse | Mysterious purple/gold, No-Face |
| 21 | `mononoke` | Princess Mononoke | Forest god | Nature red/green, curse marks |
| 22 | `howl` | Howl's Moving Castle | Steam magic | Warm copper, Calcifer flames |
| 23 | `laputa` | Castle in the Sky | Flying stone | Blue glow, ancient civilization |
| 24 | `kiki` | Kiki's Delivery Service | Delivery slip | Girl pink/purple, broomstick |
| 25 | `nausicaa` | Nausicaä | Toxic jungle | Spore pattern, Ohmu blue |

### Modern Popular (5 themes)

| # | Name | Anime | Visual Concept | Key Elements |
|---|------|-------|----------------|--------------|
| 26 | `spyfamily` | Spy x Family | Secret file | WISE intel, elegant pink/green |
| 27 | `chainsaw` | Chainsaw Man | Devil contract | Blood red/black, wild broken |
| 28 | `sao` | Sword Art Online | VRMMO HUD | Full game UI, HP/MP bars |
| 29 | `rezero` | Re:Zero | Return by death | Witch factor, dark fairy tale |
| 30 | `tokyoghoul` | Tokyo Ghoul | Kagune | Red/black, CCG report, RC cells |

### General Anime Aesthetic (10 themes)

| # | Name | Style | Visual Concept | Key Elements |
|---|------|-------|----------------|--------------|
| 31 | `isekai` | Isekai RPG | Status window | HP/MP/EXP, skill tree, reincarnation bonus |
| 32 | `mecha` | Mecha cockpit | Pilot HUD | Damage diagram, ammo count, warnings |
| 33 | `mahou_shoujo` | Magical Girl | Transformation | Stars/hearts, dreamy pink/purple |
| 34 | `shonen` | Shonen battle | Power burst | Breaking limits, red/orange flames |
| 35 | `visual_novel` | Dating sim | Dialog box | Affection hearts, choice branches |
| 36 | `chibi` | Kawaii | Super cute | Kaomoji everywhere, candy colors |
| 37 | `samurai` | Bushido | Japanese ink | Sword patterns, Edo colors |
| 38 | `idol` | Idol stage | Concert lights | Glow sticks, center position |
| 39 | `school` | School life | Blackboard | Chalk text, notebook grid, class schedule |
| 40 | `yokai` | Supernatural | Ghost story | Ghost fire purple, Hyakki Yagyo |

## Color Palettes

### EVA Theme
```go
EVAOrange    = "\033[38;2;255;102;0m"
EVARed       = "\033[38;2;204;0;0m"
EVAPurple    = "\033[38;2;128;0;128m"
EVAGreen     = "\033[38;2;0;255;0m"
```

### Dragon Ball Theme
```go
DBGreen      = "\033[38;2;0;255;128m"
DBYellow     = "\033[38;2;255;255;0m"
DBOrange     = "\033[38;2;255;165;0m"
```

### Ghibli Totoro Theme
```go
TotoroGreen  = "\033[38;2;144;238;144m"
TotoroBrown  = "\033[38;2;139;90;43m"
TotoroSky    = "\033[38;2;135;206;235m"
```

(Full color palettes will be defined per theme during implementation)

## Data Field Mappings

Each theme can creatively rename the status fields:

| Standard Field | EVA | Dragon Ball | One Piece | Isekai RPG |
|----------------|-----|-------------|-----------|------------|
| TokenCount | Sync Rate | Power Level | Bounty | EXP |
| ContextPercent | A.T. Field | Ki | Log Pose | MP |
| SessionCost | Damage | Senzu Beans | Berry | Gold |
| API5hrPercent | Battery | Stamina | Crew Morale | Stamina |
| CacheHitRate | Efficiency | Focus | Navigation | Luck |

## Implementation Approach

1. **Phase 1**: Create base templates for each category
   - Shonen action template (bordered, stats-focused)
   - Ghibli template (soft, nature-inspired)
   - Generic anime template (clean, modern)

2. **Phase 2**: Implement themes in batches of 5
   - Each batch: 5 themes → screenshots → test → commit

3. **Phase 3**: Update documentation
   - Add anime section to all THEMES.md files
   - Generate new intro GIF with anime themes

## File Structure

```
themes/
├── anime_eva.go
├── anime_dragonball.go
├── anime_naruto.go
├── anime_onepiece.go
├── anime_bleach.go
├── anime_aot.go
├── anime_demonslayer.go
├── anime_jujutsu.go
├── anime_mha.go
├── anime_hxh.go
├── anime_fma.go
├── anime_jojo.go
├── anime_deathnote.go
├── anime_sailormoon.go
├── anime_bebop.go
├── anime_gits.go
├── anime_akira.go
├── anime_gundam.go
├── ghibli_totoro.go
├── ghibli_spirited.go
├── ghibli_mononoke.go
├── ghibli_howl.go
├── ghibli_laputa.go
├── ghibli_kiki.go
├── ghibli_nausicaa.go
├── anime_spyfamily.go
├── anime_chainsaw.go
├── anime_sao.go
├── anime_rezero.go
├── anime_tokyoghoul.go
├── anime_isekai.go
├── anime_mecha.go
├── anime_mahou_shoujo.go
├── anime_shonen.go
├── anime_visual_novel.go
├── anime_chibi.go
├── anime_samurai.go
├── anime_idol.go
├── anime_school.go
└── anime_yokai.go
```

## Testing

- Each theme must render correctly with test data
- Preview command must work: `./statusline --preview <theme>`
- No ANSI escape code leaks
- Consistent width across themes

## Timeline

Estimated: 40 themes × ~30 lines each = ~1200 lines of Go code
