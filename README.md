# Claude Statusline (Go 版本)

這是一個用於 Claude Code 自定義狀態欄的 Go 程式。它可以顯示目前的模型、Git 分支、專案名稱、Session 持續時間以及 API 使用量統計。

## 預覽

狀態欄由三行組成，以下是範例輸出：

```
[💠 Claude Sonnet 4] 📂 my-project ⚡ main | ██████████ 45% 90k | 2h35m
│ ⏱️ Session ████████░░ 32% ↻3:45pm | 📅 Week ██░░░░░░░░ 12% ↻1/25 3pm
│ 🔤45.2k 💰$0.135 ⏱️25m 💬12 | 🔥$0.32/hr
```

### 各區塊說明

#### 第一行：基本資訊
| 區塊 | 範例 | 說明 |
|------|------|------|
| **模型顯示** | `[💠 Claude Sonnet 4]` | 目前使用的 Claude 模型（💛 Opus / 💠 Sonnet / 🌸 Haiku） |
| **專案名稱** | `📂 my-project` | 目前工作目錄的資料夾名稱 |
| **Git 分支** | `⚡ main` | 目前的 Git 分支（若不在 Git 專案則不顯示） |
| **Context 進度條** | `██████████ 45% 90k` | Context Window 使用量（綠色 <60% / 金色 60-80% / 紅色 >80%） |
| **每日工時** | `2h35m` | 今日所有 Session 的累積工作時間 |

#### 第二行：API 用量限制
| 區塊 | 範例 | 說明 |
|------|------|------|
| **Session 用量** | `⏱️ Session ████████░░ 32% ↻3:45pm` | 5 小時內的 API 使用率與重置時間 |
| **Week 用量** | `📅 Week ██░░░░░░░░ 12% ↻1/25 3pm` | 7 天內的 API 使用率與重置時間 |

> 進度條顏色會根據使用率變化：綠色 (<50%) → 黃色 (50-75%) → 橘色 (75-90%) → 紅色 (>90%)

#### 第三行：本次 Session 統計
| 區塊 | 範例 | 說明 |
|------|------|------|
| **Token 使用量** | `🔤45.2k` | 本次 Session 累積使用的 Token 數量 |
| **成本估算** | `💰$0.135` | 本次 Session 的預估成本 (USD) |
| **Session 時長** | `⏱️25m` | 本次 Session 的持續時間 |
| **訊息數量** | `💬12` | 本次 Session 的對話訊息數量 |
| **燒錢速度** | `🔥$0.32/hr` | 今日平均每小時花費 |

## 安裝與編譯

### 前置需求

- Go 1.18 或更高版本
- macOS（使用 Keychain 存取 OAuth Token）

### 步驟 1：下載專案

```bash
# 使用 git clone
git clone https://github.com/kevinlincg/claude-statusline.git ~/.claude/statusline-go

# 或手動下載後放到 ~/.claude/statusline-go
```

### 步驟 2：編譯程式

```bash
cd ~/.claude/statusline-go
go build -o statusline statusline.go
```

### 步驟 3：配置 Claude Code

編輯你的 `~/.claude/settings.json`，加入以下設定：

```json
{
  "statusLine": {
    "type": "command",
    "command": "/Users/your-username/.claude/statusline-go/statusline"
  }
}
```

> [!IMPORTANT]
> 請確保將 `/Users/your-username/` 替換為你實際的使用者目錄路徑。
> 可使用 `echo $HOME` 或 `whoami` 確認你的使用者名稱。

### 步驟 4：重新啟動 Claude Code

設定完成後，請重新啟動 Claude Code 以載入自定義狀態欄。

## API 使用量監控

此程式會自動從 Anthropic API 獲取你的使用量資訊。它透過 macOS Keychain 中儲存的 OAuth Token 進行認證，這個 Token 是 Claude Code 自動儲存的，你不需要額外設定。

## 本地資料儲存

本程式會在以下路徑儲存統計資料：

| 路徑 | 說明 |
|------|------|
| `~/.claude/session-tracker/sessions/` | 個別 Session 的時間與資訊 |
| `~/.claude/session-tracker/stats/` | 每日與每週的 Token 使用統計 |

## 可考慮新增的功能

以下是一些可以考慮加入的額外功能：

### 🌐 網路與系統狀態
- **網路延遲**：顯示與 Anthropic API 的 ping 延遲
- **CPU/Memory 使用率**：顯示系統資源使用狀況
- **電池狀態**：顯示筆電電量與充電狀態

### 📊 更多統計資訊
- **今日成本**：本日累積的 API 成本
- **本週成本**：本週累積的 API 成本
- **Cache 命中率**：Cache read vs input tokens 的比例
- **平均回應 Token**：每次回應的平均 Token 數

### 🔔 通知與警告
- **用量警告**：當 API 用量超過閾值時變色或閃爍
- **長對話警告**：當 Context 使用量超過 80% 時提醒
- **Session 超時**：當 Session 持續過久時提醒休息

### 🎨 視覺增強
- **主題切換**：支援亮色/暗色主題
- **自定義顏色**：允許使用者自訂各區塊顏色
- **動態圖示**：根據狀態變換圖示

### 💡 其他實用功能
- **天氣資訊**：顯示當地天氣
- **Pomodoro 計時器**：內建番茄工作法計時
- **Git 狀態**：顯示未提交的檔案數量
- **快捷鍵提示**：顯示常用快捷鍵

## 授權

MIT License
