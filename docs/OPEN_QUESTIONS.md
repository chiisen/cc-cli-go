# Open Questions - CC-CLI-Go

> **Purpose**: Track discussion topics worth deeper analysis before making new technical decisions.

---

## 使用方式 / How to Use

- 先確認題目是否已被現有 ADR 解答。
- 若尚未定案，將題目維持在此檔案。
- 若已形成明確決策，移入 [TECH_DECISIONS.md](TECH_DECISIONS.md)。

---

## 議題清單 / Discussion Topics

### 1. 串流模型是否要維持 `channel` 實作？

**問題**: 現在採用 `Channel-based Streaming`，是否已足夠，還是需要抽象成 event bus 或 callback pipeline？

**值得分析的點**:
- 可讀性
- 可測試性
- 背壓處理
- cancellation 行為

**目前傾向**: 維持 channel，除非出現明確的擴充需求。

---

### 2. 工具並行策略是否需要 worker pool？

**問題**: 現在以 `goroutine + WaitGroup` 進行並行工具執行，是否需要加上 worker pool、併發上限與排程策略？

**值得分析的點**:
- 併發安全性
- 資源控制
- 重試與 timeout
- 實作複雜度

**目前傾向**: 先保持簡單，等有明確壓力再升級。

---

### 3. Session 儲存格式是否要從 JSONL 演進？

**問題**: 現在使用 JSONL，未來是否需要換成 SQLite 或其他結構化儲存？

**值得分析的點**:
- append-only 寫入效率
- 查詢能力
- 資料一致性
- migrate 成本

**目前傾向**: JSONL 足夠，除非出現大量歷史查詢或索引需求。

---

### 4. 權限系統的 trust boundary 要多嚴格？

**問題**: 哪些工具輸入必須永遠詢問？危險命令檢測要做到多保守？

**值得分析的點**:
- false positive / false negative
- 危險命令規則維護成本
- 使用者體驗
- 安全性下限

**目前傾向**: 保守優先，先防漏判。

---

### 5. Context compaction 如何驗證品質？

**問題**: 自動摘要舊對話後，如何確保資訊沒有失真到影響後續推理？

**值得分析的點**:
- 摘要完整度
- token 節省效果
- 回歸測試資料集
- 人工可讀性

**目前傾向**: 需要建立固定 prompt 集合做回歸。

---

### 6. 錯誤模型是否足以支撐 UI 與重試？

**問題**: `errors.Is/As` + 自訂錯誤型別，是否已足夠支援 UI 顯示、重試分類與日誌分析？

**值得分析的點**:
- 錯誤分層
- 使用者訊息
- 程式可判斷性
- 錯誤來源追蹤

**目前傾向**: 夠用，但需要明確錯誤分類規範。

---

### 7. Provider 是否要提前抽象？

**問題**: 現在只支援 Anthropic，是否要先設計 `Provider` interface 以便未來擴充？

**值得分析的點**:
- 過早抽象風險
- 未來擴充成本
- API 層的耦合程度
- 測試替身需求

**目前傾向**: 先不抽象，等有第二供應商需求再切。

---

### 8. 文件治理要不要自動化？

**問題**: README、TODO、CHANGELOG、docs 索引的同步，是否要用 CI 或檢查腳本保證？

**值得分析的點**:
- 維護成本
- 文件失真風險
- CI 複雜度
- 對團隊協作的幫助

**目前傾向**: 先手動規範，之後再加自動檢查。

---

## 從現有文件延伸 / Topics Derived from Existing Docs

### 9. `internal/context` 是否應該獨立成更清楚的上下文層？

**來源**: [GO_MODULE_STRUCTURE.md](GO_MODULE_STRUCTURE.md), [GO_ARCHITECTURE.md](GO_ARCHITECTURE.md)

**問題**: 目前文件中同時出現 `query` 的 context builder 與獨立 `context` package，這兩者邊界是否需要再切清楚？

**值得分析的點**:
- package 責任切分
- 減少重複實作
- future refactor 成本
- 測試邊界

**目前傾向**: 若 context 只服務 query，可合併；若要擴展到多入口，再保留獨立層。

---

### 10. `cmd` / `internal` 的邊界是否足夠嚴謹？

**來源**: [GO_MODULE_STRUCTURE.md](GO_MODULE_STRUCTURE.md)

**問題**: CLI entrypoint、工具實作、API client、TUI 各自應該放在哪一層，才能維持長期可維護性？

**值得分析的點**:
- 對外 API 是否暴露
- 測試與重用
- package coupling
- 探索性功能的放置位置

**目前傾向**: `cmd` 保持極薄，核心能力收斂在 `internal`。

---

### 11. 架構文件與實作狀態是否要維持一一對齊？

**來源**: [ARCHITECTURE_ANALYSIS.md](ARCHITECTURE_ANALYSIS.md), [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)

**問題**: 原始架構分析、Go 架構設計、實作計畫三份文件，要不要要求每次實作後都同步修正？

**值得分析的點**:
- 文件一致性
- 決策追蹤
- 過期資訊風險
- 維護成本

**目前傾向**: 至少關鍵段落要對齊，完整對齊可用 release 檢查。

---

### 12. Phase 規劃是否需要更細的驗收門檻？

**來源**: [CORE_FEATURES.md](CORE_FEATURES.md), [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)

**問題**: 每個 phase 結束時要用什麼標準判定「真的可以進下一階段」？

**值得分析的點**:
- 測試覆蓋率
- 功能驗收
- 效能門檻
- 文件完成度

**目前傾向**: 以測試 + 功能演示 + 文件更新三者共同作為門檻。
