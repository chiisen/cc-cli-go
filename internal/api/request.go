package api

import "github.com/liao-eli/cc-cli/claude-code-go/internal/types"

type Request struct {
	Model     string         `json:"model"`
	MaxTokens int            `json:"max_tokens"`
	System    []SystemBlock  `json:"system,omitempty"`
	Messages  []MessageParam `json:"messages"`
	Tools     []ToolParam    `json:"tools,omitempty"`
	Stream    bool           `json:"stream"`
}

type SystemBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type MessageParam struct {
	Role    string         `json:"role"`
	Content []ContentParam `json:"content"`
}

type ContentParam struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`

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
