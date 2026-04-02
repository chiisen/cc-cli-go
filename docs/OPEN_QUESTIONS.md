# Open Questions - CC-CLI-Go

> **Purpose**: Track discussion topics worth deeper analysis before making new technical decisions.
> 
> **Last Updated**: 2026-04-02

---

## 📊 議題狀態總覽 / Status Overview

| # | 議題 | 狀態 | 相關文件 |
|---|------|------|----------|
| 1 | 串流模型 | ✅ 已決策 | TECH_DECISIONS.md |
| 2 | 工具並行策略 | ✅ 已決策 | TECH_DECISIONS.md |
| 3 | Session 儲存格式 | ✅ 已決策 | TECH_DECISIONS.md |
| 4 | 權限系統 trust boundary | ✅ 已實作 | CORE_FEATURES.md |
| 5 | Context compaction 品質驗證 | ✅ 已實作 | CORE_FEATURES.md |
| 6 | 錯誤模型 | ✅ 已實作 | TECH_DECISIONS.md |
| 7 | Provider 抽象 | ✅ 已決策 | TECH_DECISIONS.md |
| 8 | 文件治理自動化 | 🔄 部分完成 | README.md |
| 9 | `internal/context` 獨立性 | ❓ 待討論 | GO_MODULE_STRUCTURE.md |
| 10 | `cmd`/`internal` 邊界 | ✅ 已規範 | GO_MODULE_STRUCTURE.md |
| 11 | 架構文件對齊 | 🔄 部分完成 | README.md |
| 12 | Phase 驗收門檻 | ✅ 已有標準 | CORE_FEATURES.md |

**狀態說明**:
- ✅ 已決策/已實作/已規範：已有明確決策或實作
- 🔄 部分完成：已部分實作，尚需優化
- ❓ 待討論：尚未定案，需要進一步分析

## 待完成事項 / Unfinished Items

### 1. 文件治理自動化

- 狀態：`🔄 部分完成`
- 尚未完成：CI / 檢查腳本
- 下一步：驗證文件連結與索引一致性

### 2. `internal/context` 邊界定義

- 狀態：`❓ 待討論`
- 尚未完成：是否合併到 `query` 或保留獨立層的決策
- 下一步：確認是否有多入口重用需求

### 3. 架構文件與實作對齊

- 狀態：`🔄 部分完成`
- 尚未完成：完整同步所有架構段落
- 下一步：在 release 流程加入文件檢查

### 4. Context compaction 回歸測試

- 狀態：`🔄 待優化`
- 尚未完成：固定 prompt 測試集
- 下一步：建立摘要品質驗證資料集

### 5. 錯誤分類規範

- 狀態：`🔄 待優化`
- 尚未完成：明確錯誤分類文件
- 下一步：定義 UI、重試、日誌共用的分類規則

---

## 使用方式 / How to Use

- 先確認題目是否已被現有 ADR 解答。
- 若尚未定案，將題目維持在此檔案。
- 若已形成明確決策，移入 [TECH_DECISIONS.md](TECH_DECISIONS.md)。

---

## 議題清單 / Discussion Topics

### 1. 串流模型是否要維持 `channel` 實作？ ✅ 已決策

**問題**: 現在採用 `Channel-based Streaming`，是否已足夠，還是需要抽象成 event bus 或 callback pipeline？

**值得分析的點**:
- 可讀性
- 可測試性
- 背壓處理
- cancellation 行為

**決策結果**: 維持 channel-based streaming，符合 Go 慣例。

**相關文件**:
- [TECH_DECISIONS.md - ADR-005](TECH_DECISIONS.md#adr-005-streaming-architecture)
- [GO_ARCHITECTURE.md - 並發模型](GO_ARCHITECTURE.md#5-並發模型)

---

### 2. 工具並行策略是否需要 worker pool？ ✅ 已決策

**問題**: 現在以 `goroutine + WaitGroup` 進行並行工具執行，是否需要加上 worker pool、併發上限與排程策略？

**值得分析的點**:
- 併發安全性
- 資源控制
- 重試與 timeout
- 實作複雜度

**決策結果**: 使用 `goroutine + WaitGroup + Mutex`，符合 Go 慣例，保持簡單。

**相關文件**:
- [TECH_DECISIONS.md - ADR-008](TECH_DECISIONS.md#adr-008-tool-execution-concurrency)
- [GO_ARCHITECTURE.md - Query Loop](GO_ARCHITECTURE.md#41-query-loop-go-版本)
- [TODO.md - 已完成功能](../TODO.md#已完成功能--completed-features)

---

### 3. Session 儲存格式是否要從 JSONL 演進？ ✅ 已決策

**問題**: 現在使用 JSONL，未來是否需要換成 SQLite 或其他結構化儲存？

**值得分析的點**:
- append-only 寫入效率
- 查詢能力
- 資料一致性
- migrate 成本

**決策結果**: 維持 JSONL 格式，與原始專案相同，簡單且足夠。

**相關文件**:
- [TECH_DECISIONS.md - ADR-007](TECH_DECISIONS.md#adr-007-session-storage-format)
- [TODO.md - Session Storage](../TODO.md#phase-16-會話管理--session-management)
- [CHANGELOG.md - Session Storage](../CHANGELOG.md#新增功能--added)

---

### 4. 權限系統的 trust boundary 要多嚴格？ ✅ 已實作

**問題**: 哪些工具輸入必須永遠詢問？危險命令檢測要做到多保守？

**值得分析的點**:
- false positive / false negative
- 危險命令規則維護成本
- 使用者體驗
- 安全性下限

**實作狀態**: 已實作完整的 Permission System，包含：
- Permission modes (default/accept/plan/auto)
- Dangerous command detection
- Rule-based permission control

**相關文件**:
- [CORE_FEATURES.md - Permission System](CORE_FEATURES.md#15-permission-system)
- [TODO.md - Phase 1.7](../TODO.md#phase-17-權限系統--permission-system)
- [CHANGELOG.md - Permission System](../CHANGELOG.md#權限系統--permission-system)

---

### 5. Context compaction 如何驗證品質？ ✅ 已實作

**問題**: 自動摘要舊對話後，如何確保資訊沒有失真到影響後續推理？

**值得分析的點**:
- 摘要完整度
- token 節省效果
- 回歸測試資料集
- 人工可讀性

**實作狀態**: 已實作 Context Compaction System，包含：
- Auto-compact trigger (80% threshold)
- Summary generation
- Token estimation
- Manual compaction support

**相關文件**:
- [CORE_FEATURES.md - Context Compaction](CORE_FEATURES.md#31-context-compaction)
- [TODO.md - Context Compaction](../TODO.md#context-compaction-環境壓縮)
- [CHANGELOG.md - Context Compaction](../CHANGELOG.md#環境壓縮--context-compaction)

**待優化**: 建立固定 prompt 集合做回歸測試，驗證摘要品質。

---

### 6. 錯誤模型是否足以支撐 UI 與重試？ ✅ 已實作

**問題**: `errors.Is/As` + 自訂錯誤型別，是否已足夠支援 UI 顯示、重試分類與日誌分析？

**值得分析的點**:
- 錯誤分層
- 使用者訊息
- 程式可判斷性
- 錯誤來源追蹤

**實作狀態**: 已實作統一錯誤類型系統，包含：
- Unified Error Types (API/Tool/Permission/Config/Session/Internal)
- Error Wrapping with context
- User-friendly error messages
- API & Tool error handling

**相關文件**:
- [TECH_DECISIONS.md - ADR-010](TECH_DECISIONS.md#adr-010-error-handling)
- [TODO.md - Error Handling](../TODO.md#quality-assurance-品質保證)
- [CHANGELOG.md - Error Handling](../CHANGELOG.md#錯誤處理--error-handling)

**待優化**: 需要明確錯誤分類規範文件。

---

### 7. Provider 是否要提前抽象？ ✅ 已決策

**問題**: 現在只支援 Anthropic，是否要先設計 `Provider` interface 以便未來擴充？

**值得分析的點**:
- 過早抽象風險
- 未來擴充成本
- API 層的耦合程度
- 測試替身需求

**決策結果**: 先不抽象，等有第二供應商需求再切。遵循 YAGNI 原則。

**相關文件**:
- [TECH_DECISIONS.md - ADR-003](TECH_DECISIONS.md#adr-003-api-provider-support)
- [CORE_FEATURES.md - 決策摘要](CORE_FEATURES.md#決策摘要)

---

### 8. 文件治理要不要自動化？ 🔄 部分完成

**問題**: README、TODO、CHANGELOG、docs 索引的同步，是否要用 CI 或檢查腳本保證？

**值得分析的點**:
- 維護成本
- 文件失真風險
- CI 複雜度
- 對團隊協作的幫助

**實作狀態**: 
- ✅ 已建立 Documentation Maintenance Rules
- ✅ README.md 有完整的文件索引
- ✅ CHANGELOG.md 使用 Keep a Changelog 格式
- ❌ 尚未自動化檢查

**相關文件**:
- [README.md - Documentation Index](../README.md#documentation-index--文件索引)
- [README.md - Maintenance Rules](../README.md#documentation-maintenance-rules--文件維護規則)
- [CHANGELOG.md](../CHANGELOG.md)

**下一步**: 建立 CI 檢查腳本，驗證文件連結有效性與一致性。

---

## 從現有文件延伸 / Topics Derived from Existing Docs

### 9. `internal/context` 是否應該獨立成更清楚的上下文層？ ❓ 待討論

**來源**: [GO_MODULE_STRUCTURE.md](GO_MODULE_STRUCTURE.md), [GO_ARCHITECTURE.md](GO_ARCHITECTURE.md)

**問題**: 目前文件中同時出現 `query` 的 context builder 與獨立 `context` package，這兩者邊界是否需要再切清楚？

**值得分析的點**:
- package 責任切分
- 減少重複實作
- future refactor 成本
- 測試邊界

**目前狀態**: 
- `internal/context` 負責 Git 狀態、CLAUDE.md 發現、系統上下文
- `internal/query` 的 context builder 可能與 `internal/context` 有重疊

**建議**: 若 context 只服務 query，可合併；若要擴展到多入口，再保留獨立層。

**相關文件**:
- [GO_MODULE_STRUCTURE.md - internal/context](GO_MODULE_STRUCTURE.md#internal-context)
- [GO_ARCHITECTURE.md - Context Building](GO_ARCHITECTURE.md#context--system-prompt)

---

### 10. `cmd` / `internal` 的邊界是否足夠嚴謹？ ✅ 已規範

**來源**: [GO_MODULE_STRUCTURE.md](GO_MODULE_STRUCTURE.md)

**問題**: CLI entrypoint、工具實作、API client、TUI 各自應該放在哪一層，才能維持長期可維護性？

**值得分析的點**:
- 對外 API 是否暴露
- 測試與重用
- package coupling
- 探索性功能的放置位置

**規範結果**: 
- `cmd` 保持極薄，只負責入口與初始化
- 核心能力收斂在 `internal`
- `internal/*` packages 有清楚的職責邊界
- 使用 dependency injection 降低耦合

**相關文件**:
- [GO_MODULE_STRUCTURE.md - 套件依賴關係圖](GO_MODULE_STRUCTURE.md#套件依賴關係圖)
- [GO_ARCHITECTURE.md - Package 結構](GO_ARCHITECTURE.md#3-package-結構)

---

### 11. 架構文件與實作狀態是否要維持一一對齊？ 🔄 部分完成

**來源**: [ARCHITECTURE_ANALYSIS.md](ARCHITECTURE_ANALYSIS.md), [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)

**問題**: 原始架構分析、Go 架構設計、實作計畫三份文件，要不要要求每次實作後都同步修正？

**值得分析的點**:
- 文件一致性
- 決策追蹤
- 過期資訊風險
- 維護成本

**實作狀態**:
- ✅ README.md 有 Documentation Maintenance Rules
- ✅ CHANGELOG.md 記錄所有重要變更
- ✅ TODO.md 反映最新實作狀態
- ⚠️ 架構文件可能與實作有落差

**相關文件**:
- [README.md - Maintenance Rules](../README.md#documentation-maintenance-rules--文件維護規則)
- [CHANGELOG.md](../CHANGELOG.md)
- [TODO.md](../TODO.md)

**建議**: 至少關鍵段落要對齊，完整對齊可用 release 檢查。

---

### 12. Phase 規劃是否需要更細的驗收門檻？ ✅ 已有標準

**來源**: [CORE_FEATURES.md](CORE_FEATURES.md), [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)

**問題**: 每個 phase 結束時要用什麼標準判定「真的可以進下一階段」？

**值得分析的點**:
- 測試覆蓋率
- 功能驗收
- 效能門檻
- 文件完成度

**驗收標準**:
- ✅ 測試：每個工具都有完整的單元測試（平均覆蓋率 81.9%）
- ✅ 功能演示：所有核心功能已實作並通過測試
- ✅ 文件更新：README、TODO、CHANGELOG 已同步更新

**相關文件**:
- [CORE_FEATURES.md - 驗收標準](CORE_FEATURES.md#驗收標準)
- [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)
- [TODO.md - 進度統計](../TODO.md#進度統計--progress-statistics)
- [TESTING.md](../TESTING.md)

**實作結果**: 所有 P0 任務已 100% 完成，測試覆蓋率達標。
