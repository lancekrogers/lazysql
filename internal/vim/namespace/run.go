package namespace

import (
	"errors"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type RunController interface {
	RunCurrentQuery() error
	RunAllQueries() error
	RunSelectedLines() error
	RunWithParams() error
	ShowQueryHistory() error
}

type RunNamespace struct {
	controller RunController
}

func NewRunNamespace(controller RunController) *RunNamespace {
	return &RunNamespace{controller: controller}
}

func (r *RunNamespace) Prefix() rune {
	return 'r'
}

func (r *RunNamespace) Name() string {
	return "run"
}

func (r *RunNamespace) Commands() []leader.Command {
	return []leader.Command{
		{Sequence: []rune{'r'}, Description: "Run current query", Handler: r.runCurrentQuery},
		{Sequence: []rune{'a'}, Description: "Run all queries", Handler: r.runAllQueries},
		{Sequence: []rune{'l'}, Description: "Run selected lines", Handler: r.runSelectedLines},
		{Sequence: []rune{'p'}, Description: "Run with parameters", Handler: r.runWithParams},
		{Sequence: []rune{'h'}, Description: "Query history", Handler: r.showHistory},
	}
}

func (r *RunNamespace) runCurrentQuery() error {
	if r.controller == nil {
		return errors.New("run controller not configured")
	}
	return r.controller.RunCurrentQuery()
}

func (r *RunNamespace) runAllQueries() error {
	if r.controller == nil {
		return errors.New("run controller not configured")
	}
	return r.controller.RunAllQueries()
}

func (r *RunNamespace) runSelectedLines() error {
	if r.controller == nil {
		return errors.New("run controller not configured")
	}
	return r.controller.RunSelectedLines()
}

func (r *RunNamespace) runWithParams() error {
	if r.controller == nil {
		return errors.New("run controller not configured")
	}
	return r.controller.RunWithParams()
}

func (r *RunNamespace) showHistory() error {
	if r.controller == nil {
		return errors.New("run controller not configured")
	}
	return r.controller.ShowQueryHistory()
}
