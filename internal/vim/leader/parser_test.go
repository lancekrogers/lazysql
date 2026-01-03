package leader

import (
	"testing"

	"github.com/jorgerojas26/lazysql/internal/vim/whichkey"
)

func TestSequenceParserParse(t *testing.T) {
	tree := whichkey.NewKeyTree()
	called := false
	tree.AddGroup([]rune{'f'}, "find")
	tree.AddCommand([]rune{'f', 't'}, "table", func() {
		called = true
	})

	parser := NewSequenceParser(tree)
	result := parser.Parse('f')
	if result.Status != ParsePartial {
		t.Fatalf("expected partial status, got %v", result.Status)
	}
	if result.Node == nil || result.Node.Key != 'f' {
		t.Fatal("expected node for 'f'")
	}

	result = parser.Parse('t')
	if result.Status != ParseComplete {
		t.Fatalf("expected complete status, got %v", result.Status)
	}
	if result.Action == nil {
		t.Fatal("expected action for complete match")
	}
	result.Action()
	if !called {
		t.Fatal("expected action to be invoked")
	}
}

func TestSequenceParserInvalid(t *testing.T) {
	tree := whichkey.NewKeyTree()
	parser := NewSequenceParser(tree)
	result := parser.Parse('x')
	if result.Status != ParseInvalid {
		t.Fatalf("expected invalid status, got %v", result.Status)
	}
}

func TestSequenceParserReset(t *testing.T) {
	tree := whichkey.NewKeyTree()
	tree.AddGroup([]rune{'f'}, "find")
	parser := NewSequenceParser(tree)
	parser.Parse('f')

	parser.Reset()
	if len(parser.keys) != 0 {
		t.Fatal("expected keys to be cleared")
	}
	if parser.current != tree.Root {
		t.Fatal("expected current to be root after reset")
	}
}

func TestSequenceParserChildren(t *testing.T) {
	tree := whichkey.NewKeyTree()
	tree.AddGroup([]rune{'f'}, "find")
	parser := NewSequenceParser(tree)
	children := parser.CurrentChildren()
	if _, ok := children['f']; !ok {
		t.Fatal("expected child 'f' in current children")
	}
}

func TestSequenceParserNoTree(t *testing.T) {
	parser := NewSequenceParser(nil)
	result := parser.Parse('f')
	if result.Status != ParseInvalid {
		t.Fatalf("expected invalid status without tree, got %v", result.Status)
	}
}
