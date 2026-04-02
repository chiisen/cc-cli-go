package context

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ContextInfo struct {
	WorkingDir    string
	GitBranch     string
	GitStatus     string
	DateTime      string
	CLAUDEMDFiles []string
}

func BuildContext() (*ContextInfo, error) {
	info := &ContextInfo{}

	info.WorkingDir, _ = os.Getwd()

	info.DateTime = time.Now().Format("2006-01-02 15:04:05")

	info.GitBranch = getGitBranch()
	info.GitStatus = getGitStatus()

	info.CLAUDEMDFiles = findCLAUDEMDFiles(info.WorkingDir)

	return info, nil
}

func getGitBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getGitStatus() string {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (info *ContextInfo) ToSystemPrompt() string {
	var prompt strings.Builder

	prompt.WriteString("Environment Information:\n")
	prompt.WriteString(fmt.Sprintf("- Working Directory: %s\n", info.WorkingDir))
	prompt.WriteString(fmt.Sprintf("- Date/Time: %s\n", info.DateTime))

	if info.GitBranch != "" {
		prompt.WriteString(fmt.Sprintf("- Git Branch: %s\n", info.GitBranch))
	}

	if info.GitStatus != "" {
		prompt.WriteString(fmt.Sprintf("- Git Status (modified/untracked files):\n%s\n", info.GitStatus))
	}

	if len(info.CLAUDEMDFiles) > 0 {
		prompt.WriteString("\nCLAUDE.md files found:\n")
		for _, file := range info.CLAUDEMDFiles {
			content, err := os.ReadFile(file)
			if err == nil {
				prompt.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", file, string(content)))
			}
		}
	}

	return prompt.String()
}
