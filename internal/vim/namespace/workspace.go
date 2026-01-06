package namespace

import (
	"errors"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type WorkspaceController interface {
	Connect() error
	SwitchConnection() error
	ListConnections() error
	Disconnect() error
}

type WorkspaceNamespace struct {
	controller WorkspaceController
}

func NewWorkspaceNamespace(controller WorkspaceController) *WorkspaceNamespace {
	return &WorkspaceNamespace{controller: controller}
}

func (w *WorkspaceNamespace) Prefix() rune {
	return 'w'
}

func (w *WorkspaceNamespace) Name() string {
	return "workspace"
}

func (w *WorkspaceNamespace) Commands() []leader.Command {
	return []leader.Command{
		{Sequence: []rune{'c'}, Description: "Connect", Handler: w.connect},
		{Sequence: []rune{'s'}, Description: "Switch connection", Handler: w.switchConnection},
		{Sequence: []rune{'l'}, Description: "List connections", Handler: w.listConnections},
		{Sequence: []rune{'d'}, Description: "Disconnect", Handler: w.disconnect},
	}
}

func (w *WorkspaceNamespace) connect() error {
	if w.controller == nil {
		return errors.New("workspace controller not configured")
	}
	return w.controller.Connect()
}

func (w *WorkspaceNamespace) switchConnection() error {
	if w.controller == nil {
		return errors.New("workspace controller not configured")
	}
	return w.controller.SwitchConnection()
}

func (w *WorkspaceNamespace) listConnections() error {
	if w.controller == nil {
		return errors.New("workspace controller not configured")
	}
	return w.controller.ListConnections()
}

func (w *WorkspaceNamespace) disconnect() error {
	if w.controller == nil {
		return errors.New("workspace controller not configured")
	}
	return w.controller.Disconnect()
}
