package cmdline

import (
	"context"
	"errors"
	"fmt"
)

type CommandExecutor struct {
	registry *Registry
	context  *CommandContext
}

func NewCommandExecutor(registry *Registry, context *CommandContext) *CommandExecutor {
	return &CommandExecutor{
		registry: registry,
		context:  context,
	}
}

func (e *CommandExecutor) Execute(ctx context.Context, command string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if e.registry == nil {
		return errors.New("command registry not configured")
	}
	invocation, err := ParseInvocation(command)
	if err != nil {
		return err
	}
	if invocation.Name == "" {
		return nil
	}

	cmd := e.registry.Lookup(invocation.Name)
	if cmd == nil {
		return fmt.Errorf("not an editor command: %s", invocation.Raw)
	}
	return cmd.Execute(ctx, invocation, e.context)
}
