package commands

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestResolvePath(t *testing.T) {
	abs, err := resolvePath(context.Background(), "testdata/file.sql")
	if err != nil {
		t.Fatalf("resolve path: %v", err)
	}
	if !filepath.IsAbs(abs) {
		t.Fatalf("expected absolute path, got %q", abs)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("user home dir: %v", err)
	}
	homePath, err := resolvePath(context.Background(), "~")
	if err != nil {
		t.Fatalf("resolve home: %v", err)
	}
	if homePath != filepath.Clean(home) {
		t.Fatalf("expected home path %q, got %q", filepath.Clean(home), homePath)
	}

	tildePath, err := resolvePath(context.Background(), "~/file.sql")
	if err != nil {
		t.Fatalf("resolve tilde file: %v", err)
	}
	expected := filepath.Join(home, "file.sql")
	if tildePath != expected {
		t.Fatalf("expected tilde path %q, got %q", expected, tildePath)
	}
}

func TestResolvePathErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := resolvePath(ctx, "file.sql"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled context error, got %v", err)
	}

	if _, err := resolvePath(context.Background(), "   "); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestReadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.sql")
	if err := os.WriteFile(path, []byte("select 1;"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	data, err := readFile(context.Background(), path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "select 1;" {
		t.Fatalf("unexpected content: %q", string(data))
	}
}

func TestReadFileErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := readFile(ctx, "file.sql"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled context error, got %v", err)
	}

	missing := filepath.Join(t.TempDir(), "missing.sql")
	if _, err := readFile(context.Background(), missing); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "data.sql")
	if err := writeFile(context.Background(), path, []byte("select 2;"), false); err != nil {
		t.Fatalf("write file: %v", err)
	}
	// #nosec G304 -- test reads temp file content.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "select 2;" {
		t.Fatalf("unexpected content: %q", string(data))
	}
}

func TestWriteFileOverwriteAndCancel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.sql")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if err := writeFile(context.Background(), path, []byte("new"), false); err == nil {
		t.Fatal("expected overwrite error without force")
	}

	if err := writeFile(context.Background(), path, []byte("new"), true); err != nil {
		t.Fatalf("write file force: %v", err)
	}
	// #nosec G304 -- test reads temp file content.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "new" {
		t.Fatalf("unexpected content: %q", string(data))
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := writeFile(ctx, filepath.Join(dir, "other.sql"), []byte("x"), true); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled context error, got %v", err)
	}
}

func TestWriteFileStatError(t *testing.T) {
	dir := t.TempDir()
	parent := filepath.Join(dir, "parent")
	if err := os.WriteFile(parent, []byte("x"), 0o644); err != nil {
		t.Fatalf("write parent: %v", err)
	}

	path := filepath.Join(parent, "child.sql")
	if err := writeFile(context.Background(), path, []byte("x"), false); err == nil {
		t.Fatal("expected stat error for non-directory path")
	}
}

func TestWriteFileMkdirAllError(t *testing.T) {
	dir := t.TempDir()
	parent := filepath.Join(dir, "parent")
	if err := os.WriteFile(parent, []byte("x"), 0o644); err != nil {
		t.Fatalf("write parent: %v", err)
	}

	path := filepath.Join(parent, "child.sql")
	if err := writeFile(context.Background(), path, []byte("x"), true); err == nil {
		t.Fatal("expected mkdir error for non-directory path")
	}
}

func TestWriteFileRenameError(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.sql")
	if err := os.MkdirAll(target, 0o750); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}

	if err := writeFile(context.Background(), target, []byte("x"), true); err == nil {
		t.Fatal("expected rename error when target is directory")
	}
}

func TestWriteFileCreateTempError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod permission test not reliable on windows")
	}

	dir := t.TempDir()
	// #nosec G302 -- test relies on read-only directory permissions.
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Skipf("chmod failed: %v", err)
	}
	defer func() {
		// #nosec G302 -- restore directory permissions for cleanup.
		_ = os.Chmod(dir, 0o700)
	}()

	path := filepath.Join(dir, "data.sql")
	if err := writeFile(context.Background(), path, []byte("x"), true); err == nil {
		t.Fatal("expected error creating temp file in read-only dir")
	}
}
