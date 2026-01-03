package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/lancekrogers/lazysql/internal/vim/buffer"
	"github.com/lancekrogers/lazysql/internal/vim/cmdline"
)

type QuitCommand struct{}

func (q *QuitCommand) Name() string {
	return "q"
}

func (q *QuitCommand) Execute(ctx context.Context, inv cmdline.Invocation, env *cmdline.CommandContext) error {
	if env == nil || env.Buffers == nil {
		return errors.New("buffer manager not configured")
	}
	if env.App == nil {
		return errors.New("application not configured")
	}

	buf := env.Buffers.Active()
	if buf == nil {
		env.App.Stop()
		return nil
	}

	var err error
	if inv.Force {
		err = env.Buffers.CloseForce(buf.ID)
	} else {
		err = env.Buffers.Close(buf.ID)
	}

	if errors.Is(err, buffer.ErrBufferDirty) {
		return errors.New("no write since last change (add ! to override)")
	}
	if err != nil {
		return err
	}

	if env.Buffers.Count() == 0 {
		env.App.Stop()
	}
	return nil
}

type WriteQuitCommand struct {
	write *WriteCommand
	quit  *QuitCommand
}

func (wq *WriteQuitCommand) Name() string {
	return "wq"
}

func (wq *WriteQuitCommand) Execute(ctx context.Context, inv cmdline.Invocation, env *cmdline.CommandContext) error {
	if wq.write == nil {
		wq.write = &WriteCommand{}
	}
	if wq.quit == nil {
		wq.quit = &QuitCommand{}
	}
	if err := wq.write.Execute(ctx, inv, env); err != nil {
		return err
	}
	return wq.quit.Execute(ctx, cmdline.Invocation{Name: "q", Force: inv.Force, Raw: inv.Raw}, env)
}

type QuitAllCommand struct{}

func (qa *QuitAllCommand) Name() string {
	return "qa"
}

func (qa *QuitAllCommand) Execute(ctx context.Context, inv cmdline.Invocation, env *cmdline.CommandContext) error {
	if env == nil || env.Buffers == nil {
		return errors.New("buffer manager not configured")
	}
	if env.App == nil {
		return errors.New("application not configured")
	}

	if !inv.Force {
		if dirty := env.Buffers.GetDirtyBuffers(); len(dirty) > 0 {
			return errors.New(fmt.Sprintf("%d file(s) have unsaved changes (add ! to override)", len(dirty)))
		}
	}

	for _, buf := range env.Buffers.ListBuffers() {
		if inv.Force {
			if err := env.Buffers.CloseForce(buf.ID); err != nil {
				return err
			}
		} else {
			if err := env.Buffers.Close(buf.ID); err != nil {
				if errors.Is(err, buffer.ErrBufferDirty) {
					return errors.New("no write since last change (add ! to override)")
				}
				return err
			}
		}
	}

	env.App.Stop()
	return nil
}
