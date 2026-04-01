package query

import (
	"github.com/liao-eli/cc-cli/claude-code-go/internal/tools"
	"github.com/liao-eli/cc-cli/claude-code-go/internal/types"
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
