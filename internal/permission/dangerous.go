package permission

import (
	"strings"
)

func isDangerousCommand(toolName string, input map[string]interface{}) bool {
	if toolName != "Bash" {
		return false
	}

	command, _ := input["command"].(string)
	if command == "" {
		return false
	}

	dangerousPatterns := []string{
		"rm -rf",
		"rm -fr",
		"DROP TABLE",
		"DROP DATABASE",
		"DELETE FROM",
		"git push --force",
		"git push -f",
		"git reset --hard",
		"rm -rf /",
		":(){ :|:& };:",
		"chmod -R 777",
		"mkfs",
		"dd if=",
		"> /dev/sda",
		"curl | bash",
		"wget | bash",
		"apt-get remove",
		"yum remove",
		"pacman -R",
		"brew uninstall",
		"npm uninstall -g",
		"pip uninstall",
		"go clean -modcache",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(command, pattern) {
			return true
		}
	}

	return false
}

func getDangerousReason(command string) string {
	dangerousReasons := map[string]string{
		"rm -rf":           "destructive file deletion",
		"rm -fr":           "destructive file deletion",
		"DROP TABLE":       "destructive SQL operation",
		"DROP DATABASE":    "destructive SQL operation",
		"DELETE FROM":      "destructive SQL operation",
		"git push --force": "force push to remote",
		"git push -f":      "force push to remote",
		"git reset --hard": "hard reset (destructive)",
	}

	for pattern, reason := range dangerousReasons {
		if strings.Contains(command, pattern) {
			return reason
		}
	}

	return "potentially dangerous command"
}
