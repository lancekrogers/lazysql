package cmdline

import (
	"context"

	"github.com/lancekrogers/lazysql/internal/vim/buffer"
)

type StatusReporter interface {
	Info(message string)
	Error(message string)
}

type AppController interface {
	Stop()
}

type CommandContext struct {
	Buffers *buffer.Manager
	Status  StatusReporter
	App     AppController
}

type Invocation struct {
	Name  string
	Args  []string
	Force bool
	Raw   string
}

type Command interface {
	Name() string
	Execute(ctx context.Context, inv Invocation, env *CommandContext) error
}
