package editor

import (
	"editor/internal/buffer"
	"editor/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type TextArea struct {
	views.WidgetWatchers

	view views.View

	buffer *buffer.Buffer
}

// Drawing

func (t *TextArea) Draw() {
	if t.buffer == nil {
		return
	}

	for y := 0; y < len(t.buffer.Lines); y++ {
		DrawStr(t.view, 0, y, t.buffer.Lines[y], config.LineStyle)
	}
}

// Events

func (t *TextArea) HandleEvent(ev tcell.Event) bool {
	if t.buffer == nil {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRight:
			t.MoveCursor(1, 0)
		case tcell.KeyLeft:
			t.MoveCursor(-1, 0)

		case tcell.KeyUp:
			t.MoveCursor(0, -1)
		case tcell.KeyDown:
			t.MoveCursor(0, 1)

		case tcell.KeyRune:
			t.buffer.Insert(ev.Rune())
		case tcell.KeyTab:
			t.buffer.Insert('\t')
		case tcell.KeyBackspace2:
			t.buffer.Delete()
		case tcell.KeyEnter:
			t.buffer.NewLine()

		default:
			return false
		}

		return true

	case *tcell.EventMouse:
		if ev.Buttons()&tcell.ButtonPrimary != 0 {
			c := t.GetCursor()
			c.X, c.Y = ev.Position()

			t.buffer.ClampCursorY()
			c.X = t.buffer.GetCursorCharacterOffset(c.X)

			t.buffer.ClampCursorX()

			return true
		}
	}

	return false
}

func (t *TextArea) GetCursor() *buffer.Cursor {
	return t.buffer.Cursor
}

func (t *TextArea) MoveCursor(dX int, dY int) {
	visibleOffset := 0

	if dY != 0 {
		visibleOffset = t.buffer.GetCursorVisibleOffset()
	}

	t.GetCursor().Move(dX, dY)
	t.buffer.ClampCursor()

	if dY != 0 {
		t.GetCursor().X = t.buffer.GetCursorCharacterOffset(visibleOffset)
		t.buffer.ClampCursorX()
	}
}

// Widget stuff

func (t *TextArea) SetView(view views.View) {
	t.view = view
}

func (t *TextArea) Resize() {
}

func (t *TextArea) Size() (int, int) {
	return 0, 0
}
