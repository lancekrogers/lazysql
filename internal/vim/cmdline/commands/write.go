package commands

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/lancekrogers/lazysql/internal/vim/cmdline"
)

type WriteCommand struct{}

func (w *WriteCommand) Name() string {
	return "w"
}

func (w *WriteCommand) Execute(ctx context.Context, inv cmdline.Invocation, env *cmdline.CommandContext) error {
	if env == nil || env.Buffers == nil {
		return errors.New("buffer manager not configured")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	buf := env.Buffers.Active()
	if buf == nil {
		return errors.New("no active buffer")
	}

	target := ""
	if len(inv.Args) > 0 {
		target = inv.Args[0]
	}
	if target == "" {
		target = buf.FilePath
	}
	if target == "" {
		return errors.New("no file name")
	}

	resolved, err := resolvePath(ctx, target)
	if err != nil {
		return err
	}

	allowOverwrite := inv.Force || buf.FilePath == resolved
	if err := writeFile(ctx, resolved, []byte(buf.Content), allowOverwrite); err != nil {
		return err
	}
	if err := env.Buffers.UpdateMetadata(buf.ID, filepath.Base(resolved), resolved); err != nil {
		return err
	}
	if err := env.Buffers.MarkSaved(buf.ID); err != nil {
		return err
	}
	if env.Status != nil {
		env.Status.Info(fmt.Sprintf("\"%s\" written", resolved))
	}
	return nil
}
