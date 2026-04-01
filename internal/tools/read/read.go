package read

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/liao-eli/cc-cli/claude-code-go/internal/tools"
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

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error: file not found: %s", filePath),
			IsError: true,
		}, nil
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error reading file: %v", err),
			IsError: true,
		}, nil
	}

	lines := strings.Split(string(content), "\n")

	start := 0
	if offset > 0 {
		start = int(offset) - 1
		if start >= len(lines) {
			start = len(lines) - 1
		}
	}

	end := len(lines)
	if limit > 0 && start+int(limit) < end {
		end = start + int(limit)
	}

	result := strings.Join(lines[start:end], "\n")

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
