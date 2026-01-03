package leader

import "github.com/lancekrogers/lazysql/internal/vim/whichkey"

type ParseStatus uint8

const (
	ParsePartial ParseStatus = iota
	ParseComplete
	ParseInvalid
)

type ParseResult struct {
	Status   ParseStatus
	Node     *whichkey.KeyNode
	Action   func()
	Sequence []rune
}

type SequenceParser struct {
	tree    *whichkey.KeyTree
	current *whichkey.KeyNode
	keys    []rune
}

func NewSequenceParser(tree *whichkey.KeyTree) *SequenceParser {
	var current *whichkey.KeyNode
	if tree != nil {
		current = tree.Root
	}
	return &SequenceParser{
		tree:    tree,
		current: current,
		keys:    make([]rune, 0),
	}
}

func (p *SequenceParser) Parse(key rune) ParseResult {
	if p.tree == nil {
		return ParseResult{Status: ParseInvalid}
	}
	if p.current == nil {
		p.current = p.tree.Root
	}

	child, ok := p.current.Children[key]
	p.keys = append(p.keys, key)
	if !ok {
		return ParseResult{
			Status:   ParseInvalid,
			Sequence: append([]rune{}, p.keys...),
		}
	}

	p.current = child
	if child.Action != nil {
		return ParseResult{
			Status:   ParseComplete,
			Node:     child,
			Action:   child.Action,
			Sequence: append([]rune{}, p.keys...),
		}
	}

	return ParseResult{
		Status:   ParsePartial,
		Node:     child,
		Sequence: append([]rune{}, p.keys...),
	}
}

func (p *SequenceParser) Reset() {
	p.keys = p.keys[:0]
	if p.tree != nil {
		p.current = p.tree.Root
	}
}

func (p *SequenceParser) CurrentChildren() map[rune]*whichkey.KeyNode {
	if p.current == nil {
		return nil
	}
	return p.current.Children
}
