package connector

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Checkmarx/gen-ai-wrapper/pkg/message"

	"github.com/google/uuid"
)

const innerDir = "cx-gpt"

type FileSystemConnector struct {
	BaseDir string
}

func NewFileSystemConnector(baseDir string) Connector {
	if baseDir == "" {
		baseDir = os.TempDir()
	}
	return FileSystemConnector{
		BaseDir: baseDir,
	}
}

func (w FileSystemConnector) HistoryById(id uuid.UUID) ([]message.Message, error) {
	filePath, err := w.getFilePathById(id)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	filePath = filepath.Clean(filePath)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, nil
	}

	if !json.Valid(bytes) {
		return nil, fmt.Errorf("invalid JSON data in history file")
	}
	var history []message.Message
	err = json.Unmarshal(bytes, &history)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (w FileSystemConnector) DeleteHistory(id uuid.UUID) error {
	filePath, err := w.getFilePathById(id)
	if err != nil {
		return err
	}

	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	filePath = filepath.Clean(filePath)
	err = os.Remove(filePath)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied when deleting history file: %w", err)
		}
		return err
	}
	return nil
}

func (w FileSystemConnector) SaveHistory(id uuid.UUID, history []message.Message) error {
	filePath, err := w.getFilePathById(id)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(history)
	if err != nil {
		return err
	}

	return w.writeHistory(filePath, bytes)
}

func (w FileSystemConnector) writeHistory(fp string, bytes []byte) error {
	basePath, err := w.safeBasePath()
	if err != nil {
		return err
	}

	_, err = os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(basePath, 0700)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	fp = filepath.Clean(fp)
	return os.WriteFile(fp, bytes, 0600)
}

func (w FileSystemConnector) getFilePathById(id uuid.UUID) (string, error) {
	basePath, err := w.safeBasePath()
	if err != nil {
		return "", err
	}
	p := filepath.Join(basePath, id.String())
	return w.validatePath(p, basePath)
}

func (w FileSystemConnector) safeBasePath() (string, error) {
	p := filepath.Clean(filepath.Join(w.BaseDir, innerDir))
	base := filepath.Clean(w.BaseDir)
	if !strings.HasPrefix(p, base+string(filepath.Separator)) && p != base {
		return "", fmt.Errorf("computed base path escapes allowed directory")
	}
	return p, nil
}

func (w FileSystemConnector) validatePath(p, basePath string) (string, error) {
	clean := filepath.Clean(p)
	if !strings.HasPrefix(clean, basePath+string(filepath.Separator)) && clean != basePath {
		return "", fmt.Errorf("path traversal detected: path escapes allowed directory")
	}
	return clean, nil
}
