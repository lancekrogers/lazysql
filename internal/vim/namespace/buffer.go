package namespace

import (
	"errors"
	"fmt"

	"github.com/lancekrogers/lazysql/internal/vim/buffer"
	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type BufferListHandler interface {
	Show(buffers []*buffer.Buffer) error
}

type BufferNamespace struct {
	manager     *buffer.Manager
	navigator   *buffer.Navigator
	listHandler BufferListHandler
}

func NewBufferNamespace(manager *buffer.Manager, navigator *buffer.Navigator, listHandler BufferListHandler) *BufferNamespace {
	return &BufferNamespace{
		manager:     manager,
		navigator:   navigator,
		listHandler: listHandler,
	}
}

func (b *BufferNamespace) Prefix() rune {
	return 'b'
}

func (b *BufferNamespace) Name() string {
	return "buffer"
}

func (b *BufferNamespace) Commands() []leader.Command {
	commands := []leader.Command{
		{Sequence: []rune{'n'}, Description: "Next buffer", Handler: b.nextBuffer},
		{Sequence: []rune{'p'}, Description: "Previous buffer", Handler: b.previousBuffer},
		{Sequence: []rune{'b'}, Description: "Alternate buffer", Handler: b.alternateBuffer},
		{Sequence: []rune{'d'}, Description: "Delete buffer", Handler: b.deleteBuffer},
		{Sequence: []rune{'D'}, Description: "Force delete buffer", Handler: b.forceDeleteBuffer},
		{Sequence: []rune{'l'}, Description: "List buffers", Handler: b.listBuffers},
		{Sequence: []rune{'N'}, Description: "New buffer", Handler: b.newBuffer},
	}

	for i := 1; i <= 9; i++ {
		index := i
		commands = append(commands, leader.Command{
			Sequence:    []rune{rune('0' + i)},
			Description: fmt.Sprintf("Jump to buffer %d", i),
			Handler: func() error {
				if b.navigator == nil {
					return errors.New("buffer navigator not configured")
				}
				return b.navigator.GoToIndex(index)
			},
		})
	}

	return commands
}

func (b *BufferNamespace) nextBuffer() error {
	if b.navigator == nil {
		return errors.New("buffer navigator not configured")
	}
	return b.navigator.Next()
}

func (b *BufferNamespace) previousBuffer() error {
	if b.navigator == nil {
		return errors.New("buffer navigator not configured")
	}
	return b.navigator.Previous()
}

func (b *BufferNamespace) alternateBuffer() error {
	if b.navigator == nil || b.manager == nil {
		return errors.New("buffer navigator not configured")
	}
	alternate, ok := b.navigator.Alternate()
	if !ok {
		return nil
	}
	return b.manager.SetActive(alternate)
}

func (b *BufferNamespace) deleteBuffer() error {
	if b.manager == nil {
		return errors.New("buffer manager not configured")
	}
	active := b.manager.Active()
	if active == nil {
		return errors.New("no active buffer")
	}
	if err := b.manager.Close(active.ID); err != nil {
		if errors.Is(err, buffer.ErrBufferDirty) {
			return errors.New("buffer has unsaved changes (use \\bD to force)")
		}
		return err
	}
	return nil
}

func (b *BufferNamespace) forceDeleteBuffer() error {
	if b.manager == nil {
		return errors.New("buffer manager not configured")
	}
	active := b.manager.Active()
	if active == nil {
		return errors.New("no active buffer")
	}
	return b.manager.CloseForce(active.ID)
}

func (b *BufferNamespace) newBuffer() error {
	if b.manager == nil {
		return errors.New("buffer manager not configured")
	}
	name := fmt.Sprintf("Untitled-%d", b.manager.Count()+1)
	buf := b.manager.Create(name)
	return b.manager.SetActive(buf.ID)
}

func (b *BufferNamespace) listBuffers() error {
	if b.manager == nil {
		return errors.New("buffer manager not configured")
	}
	if b.listHandler == nil {
		return nil
	}
	return b.listHandler.Show(b.manager.ListBuffers())
}
