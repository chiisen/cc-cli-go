package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/liao-eli/cc-cli-go/internal/types"
)

type Session struct {
	ID        string
	ProjectID string
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []*types.Message
}

type Metadata struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"project_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	MessageCount int       `json:"message_count"`
}

func NewSession(projectID string) *Session {
	return &Session{
		ID:        generateUUID(),
		ProjectID: projectID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Messages:  []*types.Message{},
	}
}

func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getSessionDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".claude", "sessions")
}

func (s *Session) Save() error {
	sessionDir := getSessionDir()
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return err
	}

	transcriptFile := filepath.Join(sessionDir, s.ID+".jsonl")
	f, err := os.Create(transcriptFile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, msg := range s.Messages {
		data, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		f.WriteString(string(data) + "\n")
	}

	metadataFile := filepath.Join(sessionDir, s.ID+".metadata.json")
	metadata := Metadata{
		ID:           s.ID,
		ProjectID:    s.ProjectID,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
		MessageCount: len(s.Messages),
	}

	metadataData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	return os.WriteFile(metadataFile, metadataData, 0644)
}

func LoadSession(sessionID string) (*Session, error) {
	sessionDir := getSessionDir()

	metadataFile := filepath.Join(sessionDir, sessionID+".metadata.json")
	metadataData, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, err
	}

	var metadata Metadata
	if err := json.Unmarshal(metadataData, &metadata); err != nil {
		return nil, err
	}

	session := &Session{
		ID:        metadata.ID,
		ProjectID: metadata.ProjectID,
		CreatedAt: metadata.CreatedAt,
		UpdatedAt: metadata.UpdatedAt,
		Messages:  []*types.Message{},
	}

	transcriptFile := filepath.Join(sessionDir, sessionID+".jsonl")
	f, err := os.Open(transcriptFile)
	if err != nil {
		return session, nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var msg types.Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err == nil {
			session.Messages = append(session.Messages, &msg)
		}
	}

	return session, nil
}

func GetLastSession() (*Session, error) {
	sessionDir := getSessionDir()

	files, err := os.ReadDir(sessionDir)
	if err != nil {
		return nil, err
	}

	var metadataFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".metadata.json") {
			metadataFiles = append(metadataFiles, file.Name())
		}
	}

	if len(metadataFiles) == 0 {
		return nil, fmt.Errorf("no sessions found")
	}

	sort.Slice(metadataFiles, func(i, j int) bool {
		fileI := filepath.Join(sessionDir, metadataFiles[i])
		fileJ := filepath.Join(sessionDir, metadataFiles[j])

		infoI, _ := os.Stat(fileI)
		infoJ, _ := os.Stat(fileJ)

		return infoI.ModTime().After(infoJ.ModTime())
	})

	lastID := strings.TrimSuffix(metadataFiles[0], ".metadata.json")
	return LoadSession(lastID)
}

func (s *Session) AddMessage(msg *types.Message) {
	s.Messages = append(s.Messages, msg)
	s.UpdatedAt = time.Now()
}

func CleanupOldSessions(daysToKeep int) error {
	sessionDir := getSessionDir()

	files, err := os.ReadDir(sessionDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -daysToKeep)

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(sessionDir, file.Name()))
		}
	}

	return nil
}
