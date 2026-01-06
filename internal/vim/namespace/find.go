package namespace

import (
	"errors"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type FindController interface {
	FindTable() error
	FindColumn() error
	FindFunction() error
	FindView() error
	FindSchema() error
	FindInQueries() error
	FindRecent() error
}

type FindNamespace struct {
	controller FindController
}

func NewFindNamespace(controller FindController) *FindNamespace {
	return &FindNamespace{controller: controller}
}

func (f *FindNamespace) Prefix() rune {
	return 'f'
}

func (f *FindNamespace) Name() string {
	return "find"
}

func (f *FindNamespace) Commands() []leader.Command {
	return []leader.Command{
		{Sequence: []rune{'t'}, Description: "Find table", Handler: f.findTable},
		{Sequence: []rune{'c'}, Description: "Find column", Handler: f.findColumn},
		{Sequence: []rune{'f'}, Description: "Find function", Handler: f.findFunction},
		{Sequence: []rune{'v'}, Description: "Find view", Handler: f.findView},
		{Sequence: []rune{'s'}, Description: "Find schema", Handler: f.findSchema},
		{Sequence: []rune{'q'}, Description: "Find in queries", Handler: f.findInQueries},
		{Sequence: []rune{'r'}, Description: "Find recent", Handler: f.findRecent},
	}
}

func (f *FindNamespace) findTable() error {
	if f.controller == nil {
		return errors.New("find controller not configured")
	}
	return f.controller.FindTable()
}

func (f *FindNamespace) findColumn() error {
	if f.controller == nil {
		return errors.New("find controller not configured")
	}
	return f.controller.FindColumn()
}

func (f *FindNamespace) findFunction() error {
	if f.controller == nil {
		return errors.New("find controller not configured")
	}
	return f.controller.FindFunction()
}

func (f *FindNamespace) findView() error {
	if f.controller == nil {
		return errors.New("find controller not configured")
	}
	return f.controller.FindView()
}

func (f *FindNamespace) findSchema() error {
	if f.controller == nil {
		return errors.New("find controller not configured")
	}
	return f.controller.FindSchema()
}

func (f *FindNamespace) findInQueries() error {
	if f.controller == nil {
		return errors.New("find controller not configured")
	}
	return f.controller.FindInQueries()
}

func (f *FindNamespace) findRecent() error {
	if f.controller == nil {
		return errors.New("find controller not configured")
	}
	return f.controller.FindRecent()
}
