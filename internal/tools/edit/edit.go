package edit

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/liao-eli/cc-cli/claude-code-go/internal/tools"
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

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error reading file: %v", err),
			IsError: true,
		}, nil
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, oldString) {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error: old_string not found in file. Make sure you're using the EXACT string."),
			IsError: true,
		}, nil
	}

	count := strings.Count(contentStr, oldString)
	if count > 1 {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error: old_string appears %d times in file. Please provide more context to make it unique.", count),
			IsError: true,
		}, nil
	}

	newContent := strings.Replace(contentStr, oldString, newString, 1)

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
