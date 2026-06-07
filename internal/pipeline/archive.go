// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// SPDX-License-Identifier: AGPL-3.0-only

package pipeline

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ── PromptScaffold ──

// PromptScaffold represents the full prompt scaffold with all dimension fields.
type PromptScaffold struct {
	Scene        string            `json:"scene"`
	Platform     string            `json:"platform"`
	Fields       map[string]string `json:"fields"`
	Negative     []string          `json:"negative"`
	OriginalScene string           `json:"original_scene,omitempty"`
	RawPrompt    string            `json:"raw_prompt,omitempty"`
}

// ArchiveRecord is a single entry in the JSONL archive.
type ArchiveRecord struct {
	ID          string            `json:"id"`
	Scene       string            `json:"scene"`
	Platform    string            `json:"platform"`
	Fields      map[string]string `json:"fields"`
	Negative    []string          `json:"negative"`
	FinalPrompt string            `json:"final_prompt"`
	CreatedAt   string            `json:"created_at"`
	Tags        []string          `json:"tags"`
}

// ── Archive ──

// Archive manages JSONL-based prompt storage.
type Archive struct {
	BaseDir string // Directory for archive.jsonl (default: ~/.via54)
}

// NewArchive creates a new Archive with the given base directory.
// If baseDir is empty, uses $HOME/.via54.
func NewArchive(baseDir string) *Archive {
	if baseDir == "" {
		home, _ := os.UserHomeDir()
		baseDir = filepath.Join(home, ".via54")
	}
	return &Archive{BaseDir: baseDir}
}

// archivePath returns the path to the archive JSONL file.
func (a *Archive) archivePath() string {
	return filepath.Join(a.BaseDir, "archive.jsonl")
}

// ensureDir creates the archive directory if needed.
func (a *Archive) ensureDir() error {
	return os.MkdirAll(a.BaseDir, 0755)
}

// generateID generates a short random hex ID.
func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Save stores a prompt scaffold in the archive.
func (a *Archive) Save(scaffold *PromptScaffold, tags []string) (string, error) {
	if err := a.ensureDir(); err != nil {
		return "", fmt.Errorf("create archive dir: %w", err)
	}

	recordID := generateID()
	record := ArchiveRecord{
		ID:          recordID,
		Scene:       scaffold.Scene,
		Platform:    scaffold.Platform,
		Fields:      scaffold.Fields,
		Negative:    scaffold.Negative,
		FinalPrompt: scaffold.RawPrompt,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		Tags:        tags,
	}

	f, err := os.OpenFile(a.archivePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("open archive: %w", err)
	}
	defer f.Close()

	line, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("marshal record: %w", err)
	}

	if _, err := f.Write(append(line, '\n')); err != nil {
		return "", fmt.Errorf("write record: %w", err)
	}

	return recordID, nil
}

// Search finds records matching a query in scene or tags.
func (a *Archive) Search(query string, limit int) ([]ArchiveRecord, error) {
	if limit <= 0 {
		limit = 10
	}

	var results []ArchiveRecord
	queryLower := strings.ToLower(query)

	f, err := os.Open(a.archivePath())
	if err != nil {
		if os.IsNotExist(err) {
			return results, nil
		}
		return nil, fmt.Errorf("open archive: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var record ArchiveRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			continue
		}

		// Check scene match
		if strings.Contains(strings.ToLower(record.Scene), queryLower) {
			results = append(results, record)
			if len(results) >= limit {
				break
			}
			continue
		}

		// Check tags match
		for _, tag := range record.Tags {
			if strings.Contains(strings.ToLower(tag), queryLower) {
				results = append(results, record)
				if len(results) >= limit {
					break
				}
				break
			}
		}
	}

	return results, scanner.Err()
}

// List returns recent archive entries (newest first).
func (a *Archive) List(recent int) ([]ArchiveRecord, error) {
	if recent <= 0 {
		recent = 20
	}

	var records []ArchiveRecord

	f, err := os.Open(a.archivePath())
	if err != nil {
		if os.IsNotExist(err) {
			return records, nil
		}
		return nil, fmt.Errorf("open archive: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var record ArchiveRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			continue
		}
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan archive: %w", err)
	}

	// Return last 'recent' entries, newest first
	if len(records) > recent {
		records = records[len(records)-recent:]
	}

	// Reverse for newest first
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}

	return records, nil
}

// Delete removes an archive entry by its ID.
func (a *Archive) Delete(recordID string) (bool, error) {
	path := a.archivePath()

	// Read all records
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("open archive: %w", err)
	}

	var kept []ArchiveRecord
	found := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var record ArchiveRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			continue
		}

		if record.ID == recordID {
			found = true
			continue
		}
		kept = append(kept, record)
	}

	scanErr := scanner.Err()
	f.Close() // Close before writing

	if scanErr != nil {
		return false, fmt.Errorf("scan archive: %w", scanErr)
	}

	if !found {
		return false, nil
	}

	// Rewrite the archive without the deleted record
	tmpPath := path + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return false, fmt.Errorf("create temp file: %w", err)
	}

	for _, record := range kept {
		line, err := json.Marshal(record)
		if err != nil {
			out.Close()
			os.Remove(tmpPath)
			return false, fmt.Errorf("marshal record: %w", err)
		}
		if _, err := out.Write(append(line, '\n')); err != nil {
			out.Close()
			os.Remove(tmpPath)
			return false, fmt.Errorf("write record: %w", err)
		}
	}

	out.Close()

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return false, fmt.Errorf("replace archive: %w", err)
	}

	return true, nil
}
