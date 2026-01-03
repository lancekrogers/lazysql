package buffer

import "time"

type BufferID int

type Buffer struct {
	ID         BufferID
	Name       string
	Content    string
	FilePath   string
	Dirty      bool
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func newBuffer(id BufferID, name string) *Buffer {
	now := time.Now()
	return &Buffer{
		ID:         id,
		Name:       name,
		CreatedAt:  now,
		ModifiedAt: now,
	}
}

func (b *Buffer) SetContent(content string) {
	b.Content = content
	b.MarkDirty()
}

func (b *Buffer) MarkDirty() {
	b.Dirty = true
	b.ModifiedAt = time.Now()
}

func (b *Buffer) MarkSaved() {
	b.Dirty = false
	b.ModifiedAt = time.Now()
}
