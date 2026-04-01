package bash

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/liao-eli/cc-cli/claude-code-go/internal/tools"
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

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()

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
