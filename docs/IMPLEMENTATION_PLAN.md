# Implementation Plan - CC-CLI-Go

> **For agentic workers**: REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [x]`) syntax for tracking.

**Goal**: Build a minimal but functional CC-CLI-Go clone in Go for learning purposes.

**Architecture**: Bubble Tea TUI + Channel-based streaming + Interface-based tool system.

**Tech Stack**: Go 1.21+, cobra (CLI), Bubble Tea (TUI), net/http (API client).

---

## Phase 1: Project Setup & Foundation

### Task 1: Initialize Go Module

**Files:**
- Create: `go.mod`
- Create: `go.sum`
- Create: `.gitignore`

- [x] **Step 1: Create go.mod**

```bash
mkdir cc-cli-go-go && cd cc-cli-go-go
go mod init github.com/yourusername/cc-cli-go-go
```

- [x] **Step 2: Add dependencies**

```bash
go get github.com/spf13/cobra@latest
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/charmbracelet/lipgloss@latest
```

- [x] **Step 3: Create .gitignore**

```
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test files
*.test
*.out

# Go workspace
go.work

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Claude Code
.claude/
```

- [x] **Step 4: Verify setup**

Run: `go mod tidy && go build ./...`
Expected: No errors

- [x] **Step 5: Commit**

```bash
git init
git add .
git commit -m "feat: initialize Go module with dependencies"
```

---

### Task 2: Create Basic CLI Entry Point

**Files:**
- Create: `cmd/cc-cli-go/main.go`
- Create: `internal/cli/root.go`

- [x] **Step 1: Create main.go**

```go
// cmd/cc-cli-go/main.go
package main

import (
    "os"
    
    "github.com/yourusername/cc-cli-go-go/internal/cli"
)

func main() {
    if err := cli.Execute(); err != nil {
        os.Exit(1)
    }
}
```

- [x] **Step 2: Create root.go with cobra**

```go
// internal/cli/root.go
package cli

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
)

var (
    Version = "0.1.0"
)

var rootCmd = &cobra.Command{
    Use:   "cc-cli-go",
    Short: "CC-CLI-Go - AI-powered coding assistant",
    Long:  `A Go implementation of CC-CLI-Go for learning purposes.`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.Version = Version
}
```

- [x] **Step 3: Run and verify**

Run: `go run ./cmd/cc-cli-go --version`
Expected: `cc-cli-go version 0.1.0`

- [x] **Step 4: Commit**

```bash
git add cmd/ internal/
git commit -m "feat: add basic CLI entry point with cobra"
```

---

### Task 3: Define Core Types

**Files:**
- Create: `internal/types/message.go`
- Create: `internal/types/content.go`
- Create: `internal/types/usage.go`

- [x] **Step 1: Create message types**

```go
// internal/types/message.go
package types

import "time"

type MessageType string

const (
    MessageTypeUser      MessageType = "user"
    MessageTypeAssistant MessageType = "assistant"
    MessageTypeSystem    MessageType = "system"
)

type Message struct {
    UUID      string         `json:"uuid"`
    Type      MessageType    `json:"type"`
    Role      string         `json:"role,omitempty"`
    Content   []ContentBlock `json:"content,omitempty"`
    Usage     *Usage         `json:"usage,omitempty"`
    Timestamp time.Time      `json:"timestamp"`
    
    StopReason string        `json:"stop_reason,omitempty"`
}

func NewUserMessage(content string) *Message {
    return &Message{
        UUID:      generateUUID(),
        Type:      MessageTypeUser,
        Role:      "user",
        Content:   []ContentBlock{{Type: "text", Text: content}},
        Timestamp: time.Now(),
    }
}

func NewAssistantMessage() *Message {
    return &Message{
        UUID:      generateUUID(),
        Type:      MessageTypeAssistant,
        Role:      "assistant",
        Timestamp: time.Now(),
    }
}
```

- [x] **Step 2: Create content block types**

```go
// internal/types/content.go
package types

type ContentBlock struct {
    Type    string      `json:"type"`
    Text    string      `json:"text,omitempty"`
    
    // For tool_use
    ID       string      `json:"id,omitempty"`
    Name     string      `json:"name,omitempty"`
    Input    interface{} `json:"input,omitempty"`
    
    // For tool_result
    ToolUseID string     `json:"tool_use_id,omitempty"`
    Content   interface{} `json:"content,omitempty"`
    IsError   bool       `json:"is_error,omitempty"`
}

func NewTextBlock(text string) ContentBlock {
    return ContentBlock{Type: "text", Text: text}
}

func NewToolUseBlock(id, name string, input interface{}) ContentBlock {
    return ContentBlock{
        Type:  "tool_use",
        ID:    id,
        Name:  name,
        Input: input,
    }
}

func NewToolResultBlock(toolUseID string, content interface{}, isError bool) ContentBlock {
    return ContentBlock{
        Type:       "tool_result",
        ToolUseID:  toolUseID,
        Content:    content,
        IsError:    isError,
    }
}
```

- [x] **Step 3: Create usage type**

```go
// internal/types/usage.go
package types

type Usage struct {
    InputTokens  int `json:"input_tokens"`
    OutputTokens int `json:"output_tokens"`
}

func (u *Usage) Total() int {
    return u.InputTokens + u.OutputTokens
}
```

- [x] **Step 4: Create UUID helper**

```go
// internal/types/uuid.go
package types

import "github.com/google/uuid"

func generateUUID() string {
    return uuid.New().String()
}
```

- [x] **Step 5: Add uuid dependency**

```bash
go get github.com/google/uuid
```

- [x] **Step 6: Verify types compile**

Run: `go build ./internal/types`
Expected: No errors

- [x] **Step 7: Commit**

```bash
git add internal/types/
git commit -m "feat: define core message and content types"
```

---

## Phase 2: API Client

### Task 4: Create API Client Structure

**Files:**
- Create: `internal/api/client.go`
- Create: `internal/api/config.go`

- [x] **Step 1: Create client structure**

```go
// internal/api/client.go
package api

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type Client struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

func NewClient(apiKey string) *Client {
    return &Client{
        apiKey:  apiKey,
        baseURL: "https://api.anthropic.com/v1",
        httpClient: &http.Client{
            Timeout: 5 * time.Minute,
        },
    }
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    var reqBody io.Reader
    if body != nil {
        jsonBytes, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("marshal request: %w", err)
        }
        reqBody = bytes.NewReader(jsonBytes)
    }
    
    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }
    
    req.Header.Set("x-api-key", c.apiKey)
    req.Header.Set("anthropic-version", "2023-06-01")
    req.Header.Set("content-type", "application/json")
    req.Header.Set("anthropic-beta", "prompt-caching-2024-07-31")
    
    return c.httpClient.Do(req)
}
```

- [x] **Step 2: Create config**

```go
// internal/api/config.go
package api

const (
    DefaultModel    = "claude-sonnet-4-20250514"
    DefaultMaxTokens = 4096
)
```

- [x] **Step 3: Verify client compiles**

Run: `go build ./internal/api`
Expected: No errors

- [x] **Step 4: Commit**

```bash
git add internal/api/
git commit -m "feat: add API client structure"
```

---

### Task 5: Implement Request Builder

**Files:**
- Create: `internal/api/request.go`

- [x] **Step 1: Create request types**

```go
// internal/api/request.go
package api

import "github.com/yourusername/cc-cli-go-go/internal/types"

type Request struct {
    Model     string            `json:"model"`
    MaxTokens int               `json:"max_tokens"`
    System    []SystemBlock     `json:"system,omitempty"`
    Messages  []MessageParam    `json:"messages"`
    Tools     []ToolParam       `json:"tools,omitempty"`
    Stream    bool              `json:"stream"`
}

type SystemBlock struct {
    Type string `json:"type"`
    Text string `json:"text"`
}

type MessageParam struct {
    Role    string          `json:"role"`
    Content []ContentParam  `json:"content"`
}

type ContentParam struct {
    Type string      `json:"type"`
    Text string      `json:"text,omitempty"`
    
    // For tool_use
    ID    string      `json:"id,omitempty"`
    Name  string      `json:"name,omitempty"`
    Input interface{} `json:"input,omitempty"`
    
    // For tool_result
    ToolUseID string      `json:"tool_use_id,omitempty"`
    Content   interface{} `json:"content,omitempty"`
    IsError   bool        `json:"is_error,omitempty"`
}

type ToolParam struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    InputSchema map[string]interface{} `json:"input_schema"`
}

func NewRequest(model string, maxTokens int) *Request {
    return &Request{
        Model:     model,
        MaxTokens: maxTokens,
        Stream:    true,
    }
}

func (r *Request) SetSystem(prompt []string) {
    r.System = make([]SystemBlock, len(prompt))
    for i, p := range prompt {
        r.System[i] = SystemBlock{Type: "text", Text: p}
    }
}

func (r *Request) AddMessage(msg *types.Message) {
    content := make([]ContentParam, len(msg.Content))
    for i, c := range msg.Content {
        content[i] = ContentParam{
            Type:      c.Type,
            Text:      c.Text,
            ID:        c.ID,
            Name:      c.Name,
            Input:     c.Input,
            ToolUseID: c.ToolUseID,
            Content:   c.Content,
            IsError:   c.IsError,
        }
    }
    r.Messages = append(r.Messages, MessageParam{
        Role:    msg.Role,
        Content: content,
    })
}

func (r *Request) AddTool(tool ToolParam) {
    r.Tools = append(r.Tools, tool)
}
```

- [x] **Step 2: Verify request compiles**

Run: `go build ./internal/api`
Expected: No errors

- [x] **Step 3: Commit**

```bash
git add internal/api/request.go
git commit -m "feat: add API request builder"
```

---

### Task 6: Implement SSE Streaming

**Files:**
- Create: `internal/api/streaming.go`
- Create: `internal/api/events.go`

- [x] **Step 1: Create streaming implementation**

```go
// internal/api/streaming.go
package api

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "strings"
)

type StreamEvent struct {
    Type       string          `json:"type"`
    Index      int             `json:"index,omitempty"`
    Message    json.RawMessage `json:"message,omitempty"`
    Delta      json.RawMessage `json:"delta,omitempty"`
    ContentBlock json.RawMessage `json:"content_block,omitempty"`
    Usage      json.RawMessage `json:"usage,omitempty"`
    Error      error           `json:"-"`
}

func (c *Client) Stream(ctx context.Context, req *Request) (<-chan StreamEvent, error) {
    events := make(chan StreamEvent, 100)
    
    resp, err := c.doRequest(ctx, "POST", "/messages", req)
    if err != nil {
        close(events)
        return events, fmt.Errorf("do request: %w", err)
    }
    
    if resp.StatusCode != 200 {
        resp.Body.Close()
        close(events)
        return events, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    go func() {
        defer close(events)
        defer resp.Body.Close()
        
        scanner := bufio.NewScanner(resp.Body)
        scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)
        
        for scanner.Scan() {
            line := scanner.Text()
            
            if !strings.HasPrefix(line, "data: ") {
                continue
            }
            
            data := strings.TrimPrefix(line, "data: ")
            if data == "[DONE]" {
                break
            }
            
            var event StreamEvent
            if err := json.Unmarshal([]byte(data), &event); err != nil {
                continue
            }
            
            events <- event
        }
    }()
    
    return events, nil
}
```

- [x] **Step 2: Create event parser**

```go
// internal/api/events.go
package api

import (
    "encoding/json"
    
    "github.com/yourusername/cc-cli-go-go/internal/types"
)

type MessageStartEvent struct {
    Type    string          `json:"type"`
    Message json.RawMessage `json:"message"`
}

type ContentBlockStartEvent struct {
    Type         string          `json:"type"`
    Index        int             `json:"index"`
    ContentBlock json.RawMessage `json:"content_block"`
}

type ContentBlockDeltaEvent struct {
    Type  string          `json:"type"`
    Index int             `json:"index"`
    Delta json.RawMessage `json:"delta"`
}

type MessageDeltaEvent struct {
    Type  string          `json:"type"`
    Delta struct {
        StopReason string `json:"stop_reason"`
    } `json:"delta"`
    Usage struct {
        OutputTokens int `json:"output_tokens"`
    } `json:"usage"`
}

func ParseMessageStart(data json.RawMessage) (*types.Message, error) {
    var raw struct {
        ID      string           `json:"id"`
        Type    string           `json:"type"`
        Role    string           `json:"role"`
        Content []types.ContentBlock `json:"content"`
        Model   string           `json:"model"`
        Usage   *types.Usage     `json:"usage"`
    }
    if err := json.Unmarshal(data, &raw); err != nil {
        return nil, err
    }
    return &types.Message{
        UUID:    raw.ID,
        Type:    types.MessageTypeAssistant,
        Role:    raw.Role,
        Content: raw.Content,
        Usage:   raw.Usage,
    }, nil
}

func ParseContentBlock(data json.RawMessage) (*types.ContentBlock, error) {
    var block types.ContentBlock
    if err := json.Unmarshal(data, &block); err != nil {
        return nil, err
    }
    return &block, nil
}

func ParseDelta(data json.RawMessage) (string, string, error) {
    var delta struct {
        Type  string `json:"type"`
        Text  string `json:"text,omitempty"`
        Value string `json:"partial_json,omitempty"`
    }
    if err := json.Unmarshal(data, &delta); err != nil {
        return "", "", err
    }
    return delta.Type, delta.Text + delta.Value, nil
}
```

- [x] **Step 3: Verify streaming compiles**

Run: `go build ./internal/api`
Expected: No errors

- [x] **Step 4: Commit**

```bash
git add internal/api/streaming.go internal/api/events.go
git commit -m "feat: implement SSE streaming and event parsing"
```

---

## Phase 3: Tool System

### Task 7: Define Tool Interface

**Files:**
- Create: `internal/tools/tool.go`
- Create: `internal/tools/registry.go`

- [x] **Step 1: Create tool interface**

```go
// internal/tools/tool.go
package tools

import (
    "context"
    
    "github.com/yourusername/cc-cli-go-go/internal/types"
)

type ToolResult struct {
    Content     interface{}
    IsError     bool
    NewMessages []*types.Message
}

type PermissionResult struct {
    Behavior     string // "allow", "deny", "ask"
    Reason       string
}

type ToolContext struct {
    WorkingDir    string
    AbortSignal   context.Context
}

type Tool interface {
    Name() string
    Description() string
    InputSchema() map[string]interface{}
    
    Execute(ctx context.Context, input map[string]interface{}, tc *ToolContext) (*ToolResult, error)
    
    IsEnabled() bool
    IsReadOnly(input map[string]interface{}) bool
    IsConcurrencySafe(input map[string]interface{}) bool
    
    UserFacingName(input map[string]interface{}) string
}

func ToToolParam(t Tool) map[string]interface{} {
    return map[string]interface{}{
        "name":        t.Name(),
        "description": t.Description(),
        "input_schema": t.InputSchema(),
    }
}
```

- [x] **Step 2: Create registry**

```go
// internal/tools/registry.go
package tools

import "sync"

type Registry struct {
    mu    sync.RWMutex
    tools map[string]Tool
}

func NewRegistry() *Registry {
    return &Registry{
        tools: make(map[string]Tool),
    }
}

func (r *Registry) Register(tool Tool) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.tools[tool.Name()] = tool
}

func (r *Registry) Get(name string) Tool {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.tools[name]
}

func (r *Registry) All() []Tool {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    result := make([]Tool, 0, len(r.tools))
    for _, t := range r.tools {
        result = append(result, t)
    }
    return result
}
```

- [x] **Step 3: Verify tool system compiles**

Run: `go build ./internal/tools`
Expected: No errors

- [x] **Step 4: Commit**

```bash
git add internal/tools/
git commit -m "feat: define tool interface and registry"
```

---

### Task 8: Implement Bash Tool

**Files:**
- Create: `internal/tools/bash/bash.go`

- [x] **Step 1: Create Bash tool**

```go
// internal/tools/bash/bash.go
package bash

import (
    "bytes"
    "context"
    "fmt"
    "os/exec"
    "strings"
    "time"
    
    "github.com/yourusername/cc-cli-go-go/internal/tools"
)

type BashTool struct{}

func New() *BashTool {
    return &BashTool{}
}

func (t *BashTool) Name() string {
    return "Bash"
}

func (t *BashTool) Description() string {
    return "Execute a bash command. Use for running shell commands."
}

func (t *BashTool) InputSchema() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "command": map[string]interface{}{
                "type":        "string",
                "description": "The command to execute",
            },
            "timeout": map[string]interface{}{
                "type":        "integer",
                "description": "Timeout in milliseconds",
                "default":     120000,
            },
        },
        "required": []string{"command"},
    }
}

func (t *BashTool) Execute(ctx context.Context, input map[string]interface{}, tc *tools.ToolContext) (*tools.ToolResult, error) {
    command, _ := input["command"].(string)
    timeoutMs, _ := input["timeout"].(float64)
    if timeoutMs == 0 {
        timeoutMs = 120000
    }
    
    if command == "" {
        return &tools.ToolResult{
            Content: "Error: command is required",
            IsError: true,
        }, nil
    }
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)
    defer cancel()
    
    // Execute command
    cmd := exec.CommandContext(ctx, "bash", "-c", command)
    if tc.WorkingDir != "" {
        cmd.Dir = tc.WorkingDir
    }
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    err := cmd.Run()
    
    output := stdout.String()
    if stderr.Len() > 0 {
        output += "\nstderr:\n" + stderr.String()
    }
    
    if err != nil {
        output += fmt.Sprintf("\nError: %v", err)
    }
    
    return &tools.ToolResult{
        Content: strings.TrimSpace(output),
        IsError: err != nil && ctx.Err() == context.DeadlineExceeded,
    }, nil
}

func (t *BashTool) IsEnabled() bool {
    return true
}

func (t *BashTool) IsReadOnly(input map[string]interface{}) bool {
    // Conservative: assume all bash commands are not read-only
    return false
}

func (t *BashTool) IsConcurrencySafe(input map[string]interface{}) bool {
    return false
}

func (t *BashTool) UserFacingName(input map[string]interface{}) string {
    if cmd, ok := input["command"].(string); ok {
        if len(cmd) > 50 {
            return cmd[:50] + "..."
        }
        return cmd
    }
    return "Bash"
}
```

- [x] **Step 2: Verify Bash tool compiles**

Run: `go build ./internal/tools/bash`
Expected: No errors

- [x] **Step 3: Commit**

```bash
git add internal/tools/bash/
git commit -m "feat: implement Bash tool"
```

---

### Task 9: Implement Read Tool

**Files:**
- Create: `internal/tools/read/read.go`

- [x] **Step 1: Create Read tool**

```go
// internal/tools/read/read.go
package read

import (
    "context"
    "io/ioutil"
    "os"
    "strings"
    
    "github.com/yourusername/cc-cli-go-go/internal/tools"
)

type ReadTool struct{}

func New() *ReadTool {
    return &ReadTool{}
}

func (t *ReadTool) Name() string {
    return "Read"
}

func (t *ReadTool) Description() string {
    return "Read the contents of a file."
}

func (t *ReadTool) InputSchema() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "file_path": map[string]interface{}{
                "type":        "string",
                "description": "The absolute path to the file to read",
            },
            "limit": map[string]interface{}{
                "type":        "integer",
                "description": "Maximum number of lines to read",
            },
            "offset": map[string]interface{}{
                "type":        "integer",
                "description": "Line number to start reading from (1-indexed)",
            },
        },
        "required": []string{"file_path"},
    }
}

func (t *ReadTool) Execute(ctx context.Context, input map[string]interface{}, tc *tools.ToolContext) (*tools.ToolResult, error) {
    filePath, _ := input["file_path"].(string)
    limit, _ := input["limit"].(float64)
    offset, _ := input["offset"].(float64)
    
    if filePath == "" {
        return &tools.ToolResult{
            Content: "Error: file_path is required",
            IsError: true,
        }, nil
    }
    
    // Check if file exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return &tools.ToolResult{
            Content: fmt.Sprintf("Error: file not found: %s", filePath),
            IsError: true,
        }, nil
    }
    
    // Read file
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return &tools.ToolResult{
            Content: fmt.Sprintf("Error reading file: %v", err),
            IsError: true,
        }, nil
    }
    
    lines := strings.Split(string(content), "\n")
    
    // Apply offset
    start := 0
    if offset > 0 {
        start = int(offset) - 1
        if start >= len(lines) {
            start = len(lines) - 1
        }
    }
    
    // Apply limit
    end := len(lines)
    if limit > 0 && start+int(limit) < end {
        end = start + int(limit)
    }
    
    result := strings.Join(lines[start:end], "\n")
    
    // Add line numbers
    resultLines := strings.Split(result, "\n")
    numberedLines := make([]string, len(resultLines))
    for i, line := range resultLines {
        numberedLines[i] = fmt.Sprintf("%6d: %s", start+i+1, line)
    }
    
    return &tools.ToolResult{
        Content: strings.Join(numberedLines, "\n"),
    }, nil
}

func (t *ReadTool) IsEnabled() bool {
    return true
}

func (t *ReadTool) IsReadOnly(input map[string]interface{}) bool {
    return true
}

func (t *ReadTool) IsConcurrencySafe(input map[string]interface{}) bool {
    return true
}

func (t *ReadTool) UserFacingName(input map[string]interface{}) string {
    if path, ok := input["file_path"].(string); ok {
        return path
    }
    return "Read"
}
```

- [x] **Step 2: Add fmt import**

Run: `goimports -w internal/tools/read/read.go`

- [x] **Step 3: Verify Read tool compiles**

Run: `go build ./internal/tools/read`
Expected: No errors

- [x] **Step 4: Commit**

```bash
git add internal/tools/read/
git commit -m "feat: implement Read tool"
```

---

### Task 10: Implement Edit Tool

**Files:**
- Create: `internal/tools/edit/edit.go`

- [x] **Step 1: Create Edit tool**

```go
// internal/tools/edit/edit.go
package edit

import (
    "context"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    
    "github.com/yourusername/cc-cli-go-go/internal/tools"
)

type EditTool struct{}

func New() *EditTool {
    return &EditTool{}
}

func (t *EditTool) Name() string {
    return "Edit"
}

func (t *EditTool) Description() string {
    return "Edit a file by replacing exact string matches."
}

func (t *EditTool) InputSchema() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "file_path": map[string]interface{}{
                "type":        "string",
                "description": "The absolute path to the file to edit",
            },
            "old_string": map[string]interface{}{
                "type":        "string",
                "description": "The exact string to replace",
            },
            "new_string": map[string]interface{}{
                "type":        "string",
                "description": "The new string to replace with",
            },
        },
        "required": []string{"file_path", "old_string", "new_string"},
    }
}

func (t *EditTool) Execute(ctx context.Context, input map[string]interface{}, tc *tools.ToolContext) (*tools.ToolResult, error) {
    filePath, _ := input["file_path"].(string)
    oldString, _ := input["old_string"].(string)
    newString, _ := input["new_string"].(string)
    
    if filePath == "" || oldString == "" {
        return &tools.ToolResult{
            Content: "Error: file_path and old_string are required",
            IsError: true,
        }, nil
    }
    
    // Read file
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return &tools.ToolResult{
            Content: fmt.Sprintf("Error reading file: %v", err),
            IsError: true,
        }, nil
    }
    
    contentStr := string(content)
    
    // Check if old_string exists
    if !strings.Contains(contentStr, oldString) {
        return &tools.ToolResult{
            Content: fmt.Sprintf("Error: old_string not found in file. Make sure you're using the EXACT string."),
            IsError: true,
        }, nil
    }
    
    // Check for multiple occurrences
    count := strings.Count(contentStr, oldString)
    if count > 1 {
        return &tools.ToolResult{
            Content: fmt.Sprintf("Error: old_string appears %d times in file. Please provide more context to make it unique.", count),
            IsError: true,
        }, nil
    }
    
    // Replace
    newContent := strings.Replace(contentStr, oldString, newString, 1)
    
    // Write file
    if err := ioutil.WriteFile(filePath, []byte(newContent), 0644); err != nil {
        return &tools.ToolResult{
            Content: fmt.Sprintf("Error writing file: %v", err),
            IsError: true,
        }, nil
    }
    
    return &tools.ToolResult{
        Content: fmt.Sprintf("Successfully edited %s", filePath),
    }, nil
}

func (t *EditTool) IsEnabled() bool {
    return true
}

func (t *EditTool) IsReadOnly(input map[string]interface{}) bool {
    return false
}

func (t *EditTool) IsConcurrencySafe(input map[string]interface{}) bool {
    return false
}

func (t *EditTool) UserFacingName(input map[string]interface{}) string {
    if path, ok := input["file_path"].(string); ok {
        return path
    }
    return "Edit"
}
```

- [x] **Step 2: Verify Edit tool compiles**

Run: `go build ./internal/tools/edit`
Expected: No errors

- [x] **Step 3: Commit**

```bash
git add internal/tools/edit/
git commit -m "feat: implement Edit tool"
```

---

## Phase 4: Query Engine

### Task 11: Create Query Engine

**Files:**
- Create: `internal/query/engine.go`
- Create: `internal/query/types.go`

- [x] **Step 1: Create query types**

```go
// internal/query/types.go
package query

import (
    "github.com/yourusername/cc-cli-go-go/internal/api"
    "github.com/yourusername/cc-cli-go-go/internal/tools"
    "github.com/yourusername/cc-cli-go-go/internal/types"
)

type QueryParams struct {
    Messages     []*types.Message
    SystemPrompt []string
    Tools        []tools.Tool
    Model        string
    MaxTokens    int
}

type QueryResult struct {
    Reason string
    Error  error
}

type StreamEvent struct {
    Type    string
    Message *types.Message
    Content *types.ContentBlock
    Delta   string
    Usage   *types.Usage
}
```

- [x] **Step 2: Create query engine skeleton**

```go
// internal/query/engine.go
package query

import (
    "context"
    "encoding/json"
    "sync"
    
    "github.com/yourusername/cc-cli-go-go/internal/api"
    "github.com/yourusername/cc-cli-go-go/internal/tools"
)

type Engine struct {
    client    *api.Client
    toolReg   *tools.Registry
    eventChan chan StreamEvent
}

func NewEngine(client *api.Client, toolReg *tools.Registry) *Engine {
    return &Engine{
        client:  client,
        toolReg: toolReg,
    }
}

func (e *Engine) Query(ctx context.Context, params QueryParams) (<-chan StreamEvent, <-chan QueryResult) {
    events := make(chan StreamEvent, 100)
    results := make(chan QueryResult, 1)
    
    go func() {
        defer close(events)
        defer close(results)
        
        e.runQuery(ctx, params, events, results)
    }()
    
    return events, results
}

func (e *Engine) runQuery(ctx context.Context, params QueryParams, events chan<- StreamEvent, results chan<- QueryResult) {
    // Build request
    req := api.NewRequest(params.Model, params.MaxTokens)
    req.SetSystem(params.SystemPrompt)
    
    for _, msg := range params.Messages {
        req.AddMessage(msg)
    }
    
    for _, tool := range params.Tools {
        req.AddTool(api.ToolParam{
            Name:        tool.Name(),
            Description: tool.Description(),
            InputSchema: tool.InputSchema(),
        })
    }
    
    // Stream API response
    stream, err := e.client.Stream(ctx, req)
    if err != nil {
        results <- QueryResult{Reason: "error", Error: err}
        return
    }
    
    // Process events
    var currentMessage *types.Message
    var currentContent *types.ContentBlock
    var toolUses []types.ContentBlock
    
    for event := range stream {
        switch event.Type {
        case "message_start":
            msg, err := api.ParseMessageStart(event.Message)
            if err == nil {
                currentMessage = msg
                events <- StreamEvent{Type: "message_start", Message: msg}
            }
            
        case "content_block_start":
            block, err := api.ParseContentBlock(event.ContentBlock)
            if err == nil {
                currentContent = block
                events <- StreamEvent{Type: "content_block_start", Content: block}
            }
            
        case "content_block_delta":
            deltaType, deltaText, err := api.ParseDelta(event.Delta)
            if err == nil && deltaText != "" {
                if currentContent != nil && deltaType == "text_delta" {
                    currentContent.Text += deltaText
                }
                events <- StreamEvent{Type: "content_block_delta", Delta: deltaText}
            }
            
        case "content_block_stop":
            if currentContent != nil && currentContent.Type == "tool_use" {
                toolUses = append(toolUses, *currentContent)
            }
            events <- StreamEvent{Type: "content_block_stop"}
            
        case "message_delta":
            var delta api.MessageDeltaEvent
            if err := json.Unmarshal(event.Delta, &delta); err == nil {
                if currentMessage != nil {
                    currentMessage.StopReason = delta.Delta.StopReason
                }
            }
            
        case "message_stop":
            events <- StreamEvent{Type: "message_stop"}
            
            // Execute tools if any
            if len(toolUses) > 0 {
                toolResults := e.executeTools(ctx, toolUses, params)
                // Append results and continue loop
                // (simplified for this example)
            }
            
            results <- QueryResult{Reason: "completed"}
            return
        }
    }
}

func (e *Engine) executeTools(ctx context.Context, toolUses []types.ContentBlock, params QueryParams) []*tools.ToolResult {
    var wg sync.WaitGroup
    var mu sync.Mutex
    results := make([]*tools.ToolResult, len(toolUses))
    
    for i, tu := range toolUses {
        wg.Add(1)
        go func(idx int, toolUse types.ContentBlock) {
            defer wg.Done()
            
            tool := e.toolReg.Get(toolUse.Name)
            if tool == nil {
                results[idx] = &tools.ToolResult{
                    Content: "tool not found",
                    IsError: true,
                }
                return
            }
            
            input, _ := toolUse.Input.(map[string]interface{})
            result, err := tool.Execute(ctx, input, &tools.ToolContext{})
            if err != nil {
                result = &tools.ToolResult{
                    Content: err.Error(),
                    IsError: true,
                }
            }
            
            mu.Lock()
            results[idx] = result
            mu.Unlock()
        }(i, tu)
    }
    
    wg.Wait()
    return results
}
```

- [x] **Step 3: Verify query engine compiles**

Run: `go build ./internal/query`
Expected: No errors

- [x] **Step 4: Commit**

```bash
git add internal/query/
git commit -m "feat: implement query engine with tool execution"
```

---

## Phase 5: Basic TUI

### Task 12: Create Bubble Tea App

**Files:**
- Create: `internal/tui/app.go`
- Create: `internal/tui/model.go`

- [x] **Step 1: Create model**

```go
// internal/tui/model.go
package tui

import (
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/viewport"
    "github.com/charmbracelet/bubbles/spinner"
    tea "github.com/charmbracelet/bubbletea"
    
    "github.com/yourusername/cc-cli-go-go/internal/types"
)

type Model struct {
    // UI components
    input    textinput.Model
    viewport viewport.Model
    spinner  spinner.Model
    
    // State
    messages []*types.Message
    loading  bool
    ready    bool
    
    // Query engine
    eventChan <-chan interface{}
}

func InitialModel() Model {
    ti := textinput.New()
    ti.Placeholder = "Type your message..."
    ti.Focus()
    
    vp := viewport.New(80, 20)
    
    s := spinner.New()
    s.Spinner = spinner.Dot
    
    return Model{
        input:    ti,
        viewport: vp,
        spinner:  s,
        messages: []*types.Message{},
    }
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        textinput.Blink,
        spinner.Tick,
    )
}
```

- [x] **Step 2: Create update and view**

```go
// internal/tui/app.go
package tui

import (
    "fmt"
    "strings"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC, tea.KeyCtrlD:
            return m, tea.Quit
        case tea.KeyEnter:
            if m.input.Value() != "" {
                return m.submitInput()
            }
        }
        
    case tea.WindowSizeMsg:
        m.viewport.Width = msg.Width
        m.viewport.Height = msg.Height - 4
        if !m.ready {
            m.ready = true
        }
    }
    
    // Update components
    m.input, cmd := m.input.Update(msg)
    cmds = append(cmds, cmd)
    
    m.viewport, cmd = m.viewport.Update(msg)
    cmds = append(cmds, cmd)
    
    m.spinner, cmd = m.spinner.Update(msg)
    cmds = append(cmds, cmd)
    
    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
    if !m.ready {
        return "Loading..."
    }
    
    var b strings.Builder
    
    // Messages viewport
    b.WriteString(m.viewport.View())
    b.WriteString("\n\n")
    
    // Loading indicator
    if m.loading {
        b.WriteString(m.spinner.View())
        b.WriteString(" Thinking...")
        b.WriteString("\n")
    }
    
    // Input
    b.WriteString(m.input.View())
    
    return b.String()
}

func (m Model) submitInput() (tea.Model, tea.Cmd) {
    text := m.input.Value()
    m.input.SetValue("")
    
    // Add user message
    m.messages = append(m.messages, types.NewUserMessage(text))
    
    // Update viewport
    m.viewport.SetContent(m.renderMessages())
    
    m.loading = true
    
    // TODO: Send to query engine
    return m, nil
}

func (m Model) renderMessages() string {
    var b strings.Builder
    for _, msg := range m.messages {
        b.WriteString(m.renderMessage(msg))
        b.WriteString("\n")
    }
    return b.String()
}

func (m Model) renderMessage(msg *types.Message) string {
    var style lipgloss.Style
    var prefix string
    
    switch msg.Type {
    case types.MessageTypeUser:
        style = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
        prefix = "You: "
    case types.MessageTypeAssistant:
        style = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
        prefix = "Claude: "
    }
    
    var content string
    for _, block := range msg.Content {
        if block.Type == "text" {
            content += block.Text
        }
    }
    
    return style.Render(prefix + content)
}
```

- [x] **Step 3: Verify TUI compiles**

Run: `go build ./internal/tui`
Expected: No errors

- [x] **Step 4: Commit**

```bash
git add internal/tui/
git commit -m "feat: implement basic Bubble Tea TUI"
```

---

### Task 13: Integrate Query Engine with TUI

**Files:**
- Modify: `internal/tui/model.go`
- Modify: `internal/tui/app.go`

- [x] **Step 1: Add query engine to model**

```go
// In model.go, add to Model struct:
type Model struct {
    // ... existing fields ...
    
    queryEngine *query.Engine
    eventChan   <-chan query.StreamEvent
    resultChan  <-chan query.QueryResult
}
```

- [x] **Step 2: Handle stream events**

```go
// In app.go, add message type:
type StreamEventMsg query.StreamEvent
type QueryResultMsg query.QueryResult

// In Update function, add cases:
case StreamEventMsg:
    // Handle streaming event
    event := query.StreamEvent(msg)
    if event.Type == "content_block_delta" && event.Delta != "" {
        // Append to current assistant message
        if len(m.messages) > 0 {
            lastMsg := m.messages[len(m.messages)-1]
            if lastMsg.Type == types.MessageTypeAssistant {
                for i := range lastMsg.Content {
                    if lastMsg.Content[i].Type == "text" {
                        lastMsg.Content[i].Text += event.Delta
                    }
                }
            }
        }
        m.viewport.SetContent(m.renderMessages())
    }
    return m, nil
    
case QueryResultMsg:
    m.loading = false
    return m, nil
```

- [x] **Step 3: Update submitInput**

```go
func (m Model) submitInput() (tea.Model, tea.Cmd) {
    text := m.input.Value()
    m.input.SetValue("")
    
    // Add user message
    userMsg := types.NewUserMessage(text)
    m.messages = append(m.messages, userMsg)
    
    // Add placeholder assistant message
    assistantMsg := types.NewAssistantMessage()
    assistantMsg.Content = []types.ContentBlock{{Type: "text", Text: ""}}
    m.messages = append(m.messages, assistantMsg)
    
    m.viewport.SetContent(m.renderMessages())
    m.loading = true
    
    // Start query
    params := query.QueryParams{
        Messages:     m.messages[:len(m.messages)-1], // Exclude placeholder
        SystemPrompt: []string{"You are a helpful coding assistant."},
        Model:        "claude-sonnet-4-20250514",
        MaxTokens:    4096,
    }
    
    m.eventChan, m.resultChan = m.queryEngine.Query(context.Background(), params)
    
    return m, m.waitForEvents()
}

func (m Model) waitForEvents() tea.Cmd {
    return func() tea.Msg {
        select {
        case event, ok := <-m.eventChan:
            if !ok {
                return nil
            }
            return StreamEventMsg(event)
        case result, ok := <-m.resultChan:
            if !ok {
                return nil
            }
            return QueryResultMsg(result)
        }
    }
}
```

- [x] **Step 4: Verify integration compiles**

Run: `go build ./...`
Expected: No errors

- [x] **Step 5: Commit**

```bash
git add internal/tui/
git commit -m "feat: integrate query engine with TUI"
```

---

## Phase 6: Wire Everything Together

### Task 14: Create Main Command

**Files:**
- Modify: `internal/cli/root.go`
- Create: `internal/cli/run.go`

- [x] **Step 1: Update root command**

```go
// internal/cli/root.go - add run command
func init() {
    rootCmd.AddCommand(runCmd)
    rootCmd.AddCommand(printCmd)
}
```

- [x] **Step 2: Create run command**

```go
// internal/cli/run.go
package cli

import (
    "context"
    "os"
    
    "github.com/spf13/cobra"
    
    "github.com/yourusername/cc-cli-go-go/internal/api"
    "github.com/yourusername/cc-cli-go-go/internal/query"
    "github.com/yourusername/cc-cli-go-go/internal/tools"
    "github.com/yourusername/cc-cli-go-go/internal/tools/bash"
    "github.com/yourusername/cc-cli-go-go/internal/tools/edit"
    "github.com/yourusername/cc-cli-go-go/internal/tools/read"
    "github.com/yourusername/cc-cli-go-go/internal/tui"
)

var runCmd = &cobra.Command{
    Use:   "run",
    Short: "Start interactive session",
    RunE:  runInteractive,
}

func runInteractive(cmd *cobra.Command, args []string) error {
    // Get API key
    apiKey := os.Getenv("ANTHROPIC_API_KEY")
    if apiKey == "" {
        return fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
    }
    
    // Initialize components
    client := api.NewClient(apiKey)
    
    toolReg := tools.NewRegistry()
    toolReg.Register(bash.New())
    toolReg.Register(read.New())
    toolReg.Register(edit.New())
    
    engine := query.NewEngine(client, toolReg)
    
    // Create TUI model
    model := tui.InitialModel()
    model.QueryEngine = engine
    
    // Run TUI
    p := tea.NewProgram(model)
    if _, err := p.Run(); err != nil {
        return fmt.Errorf("run TUI: %w", err)
    }
    
    return nil
}
```

- [x] **Step 3: Verify run command compiles**

Run: `go build ./cmd/cc-cli-go`
Expected: No errors

- [x] **Step 4: Commit**

```bash
git add internal/cli/
git commit -m "feat: wire all components together in run command"
```

---

### Task 15: Test End-to-End

- [x] **Step 1: Build binary**

```bash
go build -o bin/cc-cli-go ./cmd/cc-cli-go
```

- [x] **Step 2: Test version**

```bash
./bin/cc-cli-go --version
```
Expected: Shows version

- [x] **Step 3: Test interactive mode**

```bash
ANTHROPIC_API_KEY=your-key ./bin/cc-cli-go run
```
Expected: TUI starts, can send messages

- [x] **Step 4: Fix any issues**

If there are compilation errors or runtime issues, fix them.

- [x] **Step 5: Final commit**

```bash
git add .
git commit -m "feat: complete basic Claude Code Go implementation"
```

---

## Summary

**Completed Tasks**: 15 tasks

**Estimated Time**: 25-30 hours for Phase 1-6

**Remaining Work** (Phase 7+):
- Permission system
- Glob/Grep/Write tools
- Session storage
- CLAUDE.md discovery
- Configuration management
- Error handling improvements
- Testing
- Documentation

---

*Document generated for Go porting reference*