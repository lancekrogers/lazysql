package namespace

import (
	"errors"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type CancelController interface {
	CancelQuery() error
	ClearResults() error
}

type CancelNamespace struct {
	controller CancelController
}

func NewCancelNamespace(controller CancelController) *CancelNamespace {
	return &CancelNamespace{controller: controller}
}

func (c *CancelNamespace) Prefix() rune {
	return 'c'
}

func (c *CancelNamespace) Name() string {
	return "cancel"
}

func (c *CancelNamespace) Commands() []leader.Command {
	return []leader.Command{
		{Sequence: []rune{'c'}, Description: "Cancel running query", Handler: c.cancelQuery},
		{Sequence: []rune{'l'}, Description: "Clear results", Handler: c.clearResults},
	}
}

func (c *CancelNamespace) cancelQuery() error {
	if c.controller == nil {
		return errors.New("cancel controller not configured")
	}
	return c.controller.CancelQuery()
}

func (c *CancelNamespace) clearResults() error {
	if c.controller == nil {
		return errors.New("cancel controller not configured")
	}
	return c.controller.ClearResults()
}
