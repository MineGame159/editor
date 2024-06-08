package editor

import (
	"editor/internal/buffer"
	"github.com/gdamore/tcell/v2/views"
)

type BufferView struct {
	views.BoxLayout

	numbers *LineNumbers
	area    *TextArea
}

func NewBufferView() *BufferView {
	view := &BufferView{}
	view.SetOrientation(views.Horizontal)

	view.numbers = &LineNumbers{}
	view.AddWidget(view.numbers, 0)

	view.area = &TextArea{}
	view.AddWidget(view.area, 1)

	return view
}

func (b *BufferView) GetBuffer() *buffer.Buffer {
	return b.area.buffer
}

func (b *BufferView) SetBuffer(buffer *buffer.Buffer) {
	b.numbers.SetRange(buffer)
	b.area.buffer = buffer
}
