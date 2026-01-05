package namespace

import (
	"errors"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type ShellController interface {
	ToggleShell() error
	CloseShell() error
	FocusShellInput() error
	ClearShellOutput() error
	ResizeShell() error
	ShowShellHistory() error
}

type ShellNamespace struct {
	controller ShellController
}

func NewShellNamespace(controller ShellController) *ShellNamespace {
	return &ShellNamespace{controller: controller}
}

func (s *ShellNamespace) Prefix() rune {
	return 's'
}

func (s *ShellNamespace) Name() string {
	return "shell"
}

func (s *ShellNamespace) Commands() []leader.Command {
	return []leader.Command{
		{Sequence: []rune{'s'}, Description: "Toggle shell pane", Handler: s.toggleShell},
		{Sequence: []rune{'q'}, Description: "Close shell pane", Handler: s.closeShell},
		{Sequence: []rune{'f'}, Description: "Focus shell input", Handler: s.focusShell},
		{Sequence: []rune{'c'}, Description: "Clear shell output", Handler: s.clearShell},
		{Sequence: []rune{'r'}, Description: "Resize shell pane", Handler: s.resizeShell},
		{Sequence: []rune{'h'}, Description: "Shell history", Handler: s.showHistory},
	}
}

func (s *ShellNamespace) toggleShell() error {
	if s.controller == nil {
		return errors.New("shell controller not configured")
	}
	return s.controller.ToggleShell()
}

func (s *ShellNamespace) closeShell() error {
	if s.controller == nil {
		return errors.New("shell controller not configured")
	}
	return s.controller.CloseShell()
}

func (s *ShellNamespace) focusShell() error {
	if s.controller == nil {
		return errors.New("shell controller not configured")
	}
	return s.controller.FocusShellInput()
}

func (s *ShellNamespace) clearShell() error {
	if s.controller == nil {
		return errors.New("shell controller not configured")
	}
	return s.controller.ClearShellOutput()
}

func (s *ShellNamespace) resizeShell() error {
	if s.controller == nil {
		return errors.New("shell controller not configured")
	}
	return s.controller.ResizeShell()
}

func (s *ShellNamespace) showHistory() error {
	if s.controller == nil {
		return errors.New("shell controller not configured")
	}
	return s.controller.ShowShellHistory()
}
