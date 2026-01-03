package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lancekrogers/lazysql/internal/vim/cmdline"
)

type EditCommand struct{}

func (e *EditCommand) Name() string {
	return "e"
}

func (e *EditCommand) Execute(ctx context.Context, inv cmdline.Invocation, env *cmdline.CommandContext) error {
	if env == nil || env.Buffers == nil {
		return errors.New("buffer manager not configured")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if len(inv.Args) == 0 {
		return errors.New("file name required")
	}

	resolved, err := resolvePath(ctx, inv.Args[0])
	if err != nil {
		return err
	}

	existing := env.Buffers.FindByPath(resolved)
	fileExists := true
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := os.Stat(resolved); err != nil {
		if os.IsNotExist(err) {
			fileExists = false
		} else {
			return fmt.Errorf("stat file %s: %w", resolved, err)
		}
	}

	if existing != nil && !inv.Force {
		if err := env.Buffers.SetActive(existing.ID); err != nil {
			return err
		}
		if env.Status != nil {
			env.Status.Info(fmt.Sprintf("\"%s\" (already open)", resolved))
		}
		return nil
	}

	content := ""
	if fileExists {
		data, err := readFile(ctx, resolved)
		if err != nil {
			return err
		}
		content = string(data)
	}

	buf := existing
	if buf == nil {
		buf = env.Buffers.Create(filepath.Base(resolved))
	}
	if err := env.Buffers.UpdateMetadata(buf.ID, filepath.Base(resolved), resolved); err != nil {
		return err
	}
	shouldLoad := fileExists || inv.Force
	if shouldLoad {
		if err := env.Buffers.UpdateContent(buf.ID, content); err != nil {
			return err
		}
		if err := env.Buffers.MarkSaved(buf.ID); err != nil {
			return err
		}
	}

	if err := env.Buffers.SetActive(buf.ID); err != nil {
		return err
	}

	if env.Status != nil {
		if !fileExists {
			env.Status.Info(fmt.Sprintf("\"%s\" [New File]", resolved))
		} else {
			lines := 0
			if content != "" {
				lines = strings.Count(content, "\n") + 1
			}
			env.Status.Info(fmt.Sprintf("\"%s\" %dL", resolved, lines))
		}
	}

	return nil
}
