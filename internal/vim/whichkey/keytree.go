package whichkey

import "strings"

type KeyNode struct {
	Key         rune
	Description string
	Children    map[rune]*KeyNode
	Action      func()
	Parent      *KeyNode
}

func (n *KeyNode) IsLeaf() bool {
	return n != nil && n.Action != nil
}

type KeyTree struct {
	Root    *KeyNode
	Current *KeyNode
}

func NewKeyTree() *KeyTree {
	root := &KeyNode{
		Children: make(map[rune]*KeyNode),
	}
	return &KeyTree{
		Root:    root,
		Current: root,
	}
}

func (t *KeyTree) AddGroup(sequence []rune, description string) *KeyNode {
	node := t.ensurePath(sequence)
	node.Description = description
	return node
}

func (t *KeyTree) AddCommand(sequence []rune, description string, action func()) *KeyNode {
	node := t.ensurePath(sequence)
	node.Description = description
	node.Action = action
	return node
}

func (t *KeyTree) Enter(key rune) bool {
	if t.Current == nil {
		return false
	}
	child, ok := t.Current.Children[key]
	if !ok {
		return false
	}
	t.Current = child
	return true
}

func (t *KeyTree) Back() bool {
	if t.Current == nil || t.Current == t.Root {
		return false
	}
	t.Current = t.Current.Parent
	return true
}

func (t *KeyTree) Reset() {
	if t.Root == nil {
		t.Root = &KeyNode{Children: make(map[rune]*KeyNode)}
	}
	t.Current = t.Root
}

func (t *KeyTree) Path() []rune {
	if t.Current == nil {
		return nil
	}
	var path []rune
	node := t.Current
	for node != nil && node != t.Root {
		path = append([]rune{node.Key}, path...)
		node = node.Parent
	}
	return path
}

func (t *KeyTree) Breadcrumb(leader rune) string {
	parts := []string{string(leader)}
	for _, key := range t.Path() {
		parts = append(parts, formatKeyLabel(key))
	}
	return strings.Join(parts, " > ") + " > "
}

func (t *KeyTree) ensurePath(sequence []rune) *KeyNode {
	if t.Root == nil {
		t.Root = &KeyNode{Children: make(map[rune]*KeyNode)}
	}
	current := t.Root
	for _, key := range sequence {
		child := current.Children[key]
		if child == nil {
			child = &KeyNode{
				Key:      key,
				Parent:   current,
				Children: make(map[rune]*KeyNode),
			}
			current.Children[key] = child
		}
		current = child
	}
	return current
}
