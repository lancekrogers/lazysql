package namespace

import (
	"errors"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type DescribeController interface {
	DescribeTable() error
	DescribeView() error
	DescribeIndexes() error
	DescribeSequences() error
	DescribeFunctions() error
	DescribeDatabase() error
	DescribeUnderCursor() error
}

type DescribeNamespace struct {
	controller DescribeController
}

func NewDescribeNamespace(controller DescribeController) *DescribeNamespace {
	return &DescribeNamespace{controller: controller}
}

func (d *DescribeNamespace) Prefix() rune {
	return 'd'
}

func (d *DescribeNamespace) Name() string {
	return "describe"
}

func (d *DescribeNamespace) Commands() []leader.Command {
	return []leader.Command{
		{Sequence: []rune{'t'}, Description: "Describe table", Handler: d.describeTable},
		{Sequence: []rune{'v'}, Description: "Describe view", Handler: d.describeView},
		{Sequence: []rune{'i'}, Description: "Describe indexes", Handler: d.describeIndexes},
		{Sequence: []rune{'s'}, Description: "Describe sequences", Handler: d.describeSequences},
		{Sequence: []rune{'f'}, Description: "Describe functions", Handler: d.describeFunctions},
		{Sequence: []rune{'d'}, Description: "Describe database", Handler: d.describeDatabase},
		{Sequence: []rune{'\n'}, Description: "Describe under cursor", Handler: d.describeUnderCursor},
	}
}

func (d *DescribeNamespace) describeTable() error {
	if d.controller == nil {
		return errors.New("describe controller not configured")
	}
	return d.controller.DescribeTable()
}

func (d *DescribeNamespace) describeView() error {
	if d.controller == nil {
		return errors.New("describe controller not configured")
	}
	return d.controller.DescribeView()
}

func (d *DescribeNamespace) describeIndexes() error {
	if d.controller == nil {
		return errors.New("describe controller not configured")
	}
	return d.controller.DescribeIndexes()
}

func (d *DescribeNamespace) describeSequences() error {
	if d.controller == nil {
		return errors.New("describe controller not configured")
	}
	return d.controller.DescribeSequences()
}

func (d *DescribeNamespace) describeFunctions() error {
	if d.controller == nil {
		return errors.New("describe controller not configured")
	}
	return d.controller.DescribeFunctions()
}

func (d *DescribeNamespace) describeDatabase() error {
	if d.controller == nil {
		return errors.New("describe controller not configured")
	}
	return d.controller.DescribeDatabase()
}

func (d *DescribeNamespace) describeUnderCursor() error {
	if d.controller == nil {
		return errors.New("describe controller not configured")
	}
	return d.controller.DescribeUnderCursor()
}
