package namespace

import (
	"errors"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type ExplainController interface {
	ExplainCurrentQuery() error
	ExplainAnalyzeCurrentQuery() error
}

type ExplainNamespace struct {
	controller ExplainController
}

func NewExplainNamespace(controller ExplainController) *ExplainNamespace {
	return &ExplainNamespace{controller: controller}
}

func (e *ExplainNamespace) Prefix() rune {
	return 'x'
}

func (e *ExplainNamespace) Name() string {
	return "explain"
}

func (e *ExplainNamespace) Commands() []leader.Command {
	return []leader.Command{
		{Sequence: []rune{'x'}, Description: "Explain query", Handler: e.explain},
		{Sequence: []rune{'X'}, Description: "Explain analyze", Handler: e.explainAnalyze},
	}
}

func (e *ExplainNamespace) explain() error {
	if e.controller == nil {
		return errors.New("explain controller not configured")
	}
	return e.controller.ExplainCurrentQuery()
}

func (e *ExplainNamespace) explainAnalyze() error {
	if e.controller == nil {
		return errors.New("explain controller not configured")
	}
	return e.controller.ExplainAnalyzeCurrentQuery()
}
