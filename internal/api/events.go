package api

import (
	"encoding/json"

	"github.com/liao-eli/cc-cli/claude-code-go/internal/types"
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
	Type  string `json:"type"`
	Delta struct {
		StopReason string `json:"stop_reason"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func ParseMessageStart(data json.RawMessage) (*types.Message, error) {
	var raw struct {
		ID      string               `json:"id"`
		Type    string               `json:"type"`
		Role    string               `json:"role"`
		Content []types.ContentBlock `json:"content"`
		Model   string               `json:"model"`
		Usage   *types.Usage         `json:"usage"`
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
