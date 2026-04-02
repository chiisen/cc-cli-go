package query

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/liao-eli/cc-cli-go/internal/api"
	"github.com/liao-eli/cc-cli-go/internal/permission"
	"github.com/liao-eli/cc-cli-go/internal/tools"
	"github.com/liao-eli/cc-cli-go/internal/types"
)

type Engine struct {
	client  *api.Client
	toolReg *tools.Registry
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

	stream, err := e.client.Stream(ctx, req)
	if err != nil {
		results <- QueryResult{Reason: "error", Error: err}
		return
	}

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

			if len(toolUses) > 0 {
				toolResults := e.executeTools(ctx, toolUses, params, events)
				_ = toolResults
			}

			results <- QueryResult{Reason: "completed"}
			return
		}
	}
}

func (e *Engine) executeTools(ctx context.Context, toolUses []types.ContentBlock, params QueryParams, events chan<- StreamEvent) []*tools.ToolResult {
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

			if params.PermissionChecker != nil {
				decision := params.PermissionChecker.Check(toolUse.Name, input)

				if decision.Behavior == permission.BehaviorDeny {
					results[idx] = &tools.ToolResult{
						Content: fmt.Sprintf("Permission denied: %s", decision.Reason),
						IsError: true,
					}
					return
				}

				if decision.Behavior == permission.BehaviorAsk {
					events <- StreamEvent{
						Type: "permission_request",
						PermissionRequest: &PermissionRequestEvent{
							ToolName: toolUse.Name,
							Input:    input,
							Decision: decision,
							Index:    idx,
						},
					}
				}
			}

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
