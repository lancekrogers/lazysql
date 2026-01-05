package namespace

import (
	"errors"
	"testing"

	"github.com/lancekrogers/lazysql/internal/vim/buffer"
)

type fakeListHandler struct {
	calls   int
	buffers []*buffer.Buffer
	err     error
}

func (f *fakeListHandler) Show(buffers []*buffer.Buffer) error {
	f.calls++
	f.buffers = buffers
	return f.err
}

func TestBufferNamespaceCommands(t *testing.T) {
	manager := buffer.NewManager()
	navigator := buffer.NewNavigator(manager)
	list := &fakeListHandler{}
	ns := NewBufferNamespace(manager, navigator, list)

	cmds := ns.Commands()
	if len(cmds) == 0 {
		t.Fatal("expected commands")
	}

	found := map[string]bool{}
	for _, cmd := range cmds {
		found[string(cmd.Sequence)] = true
	}

	expect := []string{"n", "p", "b", "d", "D", "l", "N", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for _, key := range expect {
		if !found[key] {
			t.Fatalf("expected command %q", key)
		}
	}
}

func TestBufferNamespaceDelete(t *testing.T) {
	manager := buffer.NewManager()
	navigator := buffer.NewNavigator(manager)
	ns := NewBufferNamespace(manager, navigator, nil)

	buf := manager.Create("one")
	if err := manager.UpdateContent(buf.ID, "dirty"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	if err := ns.deleteBuffer(); err == nil {
		t.Fatal("expected error for dirty buffer")
	}

	if err := ns.forceDeleteBuffer(); err != nil {
		t.Fatalf("force delete: %v", err)
	}
}

func TestBufferNamespaceList(t *testing.T) {
	manager := buffer.NewManager()
	navigator := buffer.NewNavigator(manager)
	list := &fakeListHandler{}
	ns := NewBufferNamespace(manager, navigator, list)

	manager.Create("one")
	if err := ns.listBuffers(); err != nil {
		t.Fatalf("list buffers: %v", err)
	}
	if list.calls != 1 {
		t.Fatalf("expected list handler call, got %d", list.calls)
	}
	if len(list.buffers) != 1 {
		t.Fatalf("expected 1 buffer, got %d", len(list.buffers))
	}
}

func TestBufferNamespaceNavigation(t *testing.T) {
	manager := buffer.NewManager()
	first := manager.Create("one")
	second := manager.Create("two")
	navigator := buffer.NewNavigator(manager)
	ns := NewBufferNamespace(manager, navigator, nil)

	if err := manager.SetActive(first.ID); err != nil {
		t.Fatalf("set active: %v", err)
	}
	if err := ns.nextBuffer(); err != nil {
		t.Fatalf("next: %v", err)
	}
	if manager.ActiveID() != second.ID {
		t.Fatalf("expected second buffer, got %d", manager.ActiveID())
	}
	if err := ns.previousBuffer(); err != nil {
		t.Fatalf("previous: %v", err)
	}
	if manager.ActiveID() != first.ID {
		t.Fatalf("expected first buffer, got %d", manager.ActiveID())
	}

	_, _ = ns.navigator.Alternate()
	if err := ns.alternateBuffer(); err != nil {
		t.Fatalf("alternate: %v", err)
	}
}

func TestBufferNamespaceJumpAndNew(t *testing.T) {
	manager := buffer.NewManager()
	navigator := buffer.NewNavigator(manager)
	ns := NewBufferNamespace(manager, navigator, nil)

	manager.Create("one")
	manager.Create("two")

	var jumpTwo func() error
	for _, cmd := range ns.Commands() {
		if string(cmd.Sequence) == "2" {
			jumpTwo = cmd.Handler
		}
	}
	if jumpTwo == nil {
		t.Fatal("expected jump command for buffer 2")
	}
	if err := jumpTwo(); err != nil {
		t.Fatalf("jump to buffer 2: %v", err)
	}
	if manager.Active().Name != "two" {
		t.Fatalf("expected active buffer two, got %q", manager.Active().Name)
	}

	if err := ns.newBuffer(); err != nil {
		t.Fatalf("new buffer: %v", err)
	}
	if manager.Active() == nil || manager.Active().Name == "" {
		t.Fatal("expected new active buffer")
	}
}

func TestBufferNamespaceListNoHandler(t *testing.T) {
	manager := buffer.NewManager()
	navigator := buffer.NewNavigator(manager)
	ns := NewBufferNamespace(manager, navigator, nil)

	manager.Create("one")
	if err := ns.listBuffers(); err != nil {
		t.Fatalf("expected nil error with no list handler, got %v", err)
	}
}

func TestBufferNamespaceErrors(t *testing.T) {
	ns := NewBufferNamespace(nil, nil, nil)
	if err := ns.nextBuffer(); err == nil {
		t.Fatal("expected navigator error")
	}
	if err := ns.previousBuffer(); err == nil {
		t.Fatal("expected navigator error")
	}
	if err := ns.alternateBuffer(); err == nil {
		t.Fatal("expected navigator error")
	}
	if err := ns.deleteBuffer(); err == nil {
		t.Fatal("expected manager error")
	}
	if err := ns.forceDeleteBuffer(); err == nil {
		t.Fatal("expected manager error")
	}
	if err := ns.newBuffer(); err == nil {
		t.Fatal("expected manager error")
	}
	if err := ns.listBuffers(); err == nil {
		t.Fatal("expected manager error")
	}

	manager := buffer.NewManager()
	navigator := buffer.NewNavigator(manager)
	list := &fakeListHandler{err: errors.New("boom")}
	ns = NewBufferNamespace(manager, navigator, list)
	if err := ns.listBuffers(); err == nil {
		t.Fatal("expected list handler error")
	}
}
