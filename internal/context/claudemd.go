package context

import (
	"os"
	"path/filepath"
)

func findCLAUDEMDFiles(startDir string) []string {
	var files []string

	dir := startDir
	for {
		claudemd := filepath.Join(dir, "CLAUDE.md")
		if _, err := os.Stat(claudemd); err == nil {
			files = append(files, claudemd)
		}

		geminiMd := filepath.Join(dir, "GEMINI.md")
		if _, err := os.Stat(geminiMd); err == nil {
			files = append(files, geminiMd)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return files
}
