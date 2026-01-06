package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func resolvePath(ctx context.Context, path string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	cleaned := strings.TrimSpace(path)
	if cleaned == "" {
		return "", errors.New("file name required")
	}
	if cleaned == "~" || strings.HasPrefix(cleaned, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		if cleaned == "~" {
			cleaned = home
		} else {
			cleaned = filepath.Join(home, cleaned[2:])
		}
	}
	cleaned = filepath.Clean(cleaned)
	if !filepath.IsAbs(cleaned) {
		abs, err := filepath.Abs(cleaned)
		if err != nil {
			return "", fmt.Errorf("resolve absolute path: %w", err)
		}
		cleaned = abs
	}
	return cleaned, nil
}

func readFile(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	// #nosec G304 -- paths are resolved and validated before read.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}
	return data, nil
}

func writeFile(ctx context.Context, path string, content []byte, allowOverwrite bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !allowOverwrite {
		if _, err := os.Stat(path); err == nil {
			return errors.New("file exists (use :w!)")
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("stat file %s: %w", path, err)
		}
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("create directory %s: %w", dir, err)
	}

	tmpFile, err := os.CreateTemp(dir, ".lazysql-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmpFile.Write(content); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}
