package sqlstore

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Interface interface {
	GetQuery(queryName string) (string, error)
}

type SqlStore struct {
	queries map[string]string
}

func New(sqlPath string) (*SqlStore, error) {
	dirs, err := os.ReadDir(sqlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL directory: %w", err)
	}

	queries := make(map[string]string)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		fullPath := filepath.Join(sqlPath, dir.Name())

		files, err := os.ReadDir(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read subdirectory %s: %w", fullPath, err)
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
				continue
			}

			content, err := os.ReadFile(filepath.Join(fullPath, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", file.Name(), err)
			}

			queries[file.Name()] = string(bytes.Join(bytes.Fields(content), []byte(" ")))
		}
	}

	if len(queries) == 0 {
		return nil, fmt.Errorf("no .sql files loaded from: %s", sqlPath)
	}

	return &SqlStore{
		queries: queries,
	}, nil
}

func (s *SqlStore) GetQuery(queryName string) (string, error) {
	query, exists := s.queries[queryName]
	if !exists {
		return "", fmt.Errorf("SQL query %s not found", queryName)
	}
	return query, nil
}
