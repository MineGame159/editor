package editor

import (
	"editor/internal/buffer"
	"editor/internal/config"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type BufferView struct {
	buffer *buffer.Buffer

	view views.View
}

// Drawing

func (b *BufferView) Draw() {
	if b.buffer == nil {
		return
	}

	for y := 0; y < len(b.buffer.Lines); y++ {
		isCursorLine := y == b.buffer.Cursor.Y

		numStyle := config.LineNumberStyle
		if isCursorLine {
			numStyle = config.LineStyle
		}

		x := b.drawStr(1, y, fmt.Sprintf("%2d", y+1), numStyle) + 1
		// TODO: x = b.drawChar(s, x, y, tcell.RuneVLine, config.LineNumberStyle)
		b.drawStr(x, y, b.buffer.Lines[y], config.LineStyle)
	}
}

func (b *BufferView) drawStr(x int, y int, str string, style tcell.Style) int {
	for _, ch := range str {
		if ch == '\t' {
			for i := 0; i < config.TabWidth; i++ {
				x = b.drawChar(x, y, ' ', style)
			}
		} else {
			x = b.drawChar(x, y, ch, style)
		}
	}

	return x
}

func (b *BufferView) drawChar(x int, y int, ch rune, style tcell.Style) int {
	b.view.SetContent(x, y, ch, nil, style)
	return x + 1
}

// Events

func (b *BufferView) HandleEvent(ev tcell.Event) bool {
	if b.buffer == nil {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRight:
			b.MoveCursor(1, 0)
		case tcell.KeyLeft:
			b.MoveCursor(-1, 0)

		case tcell.KeyUp:
			b.MoveCursor(0, -1)
		case tcell.KeyDown:
			b.MoveCursor(0, 1)

		case tcell.KeyRune:
			b.buffer.Insert(ev.Rune())
		case tcell.KeyTab:
			b.buffer.Insert('\t')
		case tcell.KeyBackspace2:
			b.buffer.Delete()
		case tcell.KeyEnter:
			b.buffer.NewLine()

		default:
			return false
		}

		return true

	case *tcell.EventMouse:
		if ev.Buttons()&tcell.ButtonPrimary != 0 {
			c := b.GetCursor()
			c.X, c.Y = ev.Position()

			b.buffer.ClampCursorY()
			c.X = b.buffer.GetCursorCharacterOffset(c.X)

			b.buffer.ClampCursorX()

			return true
		}
	}

	return false
}

func (b *BufferView) GetCursor() *buffer.Cursor {
	return b.buffer.Cursor
}

func (b *BufferView) MoveCursor(dX int, dY int) {
	visibleOffset := 0

	if dY != 0 {
		visibleOffset = b.buffer.GetCursorVisibleOffset()
	}

	b.GetCursor().Move(dX, dY)
	b.buffer.ClampCursor()

	if dY != 0 {
		b.GetCursor().X = b.buffer.GetCursorCharacterOffset(visibleOffset)
		b.buffer.ClampCursorX()
	}
}

// Widget stuff

func (b *BufferView) SetView(view views.View) {
	b.view = view
}

func (b *BufferView) Resize() {
}

func (b *BufferView) Size() (int, int) {
	return 0, 0
}

func (b *BufferView) Watch(handler tcell.EventHandler) {
}

func (b *BufferView) Unwatch(handler tcell.EventHandler) {
}
