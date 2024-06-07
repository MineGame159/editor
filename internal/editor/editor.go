package editor

import (
	"editor/internal/buffer"
	"github.com/gdamore/tcell/v2"
	"os"
)

type Editor struct {
	screen tcell.Screen
	buffer *buffer.Buffer
}

// Constructor

// New creates an editor screen with no buffer
func New() (*Editor, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	screen.EnableMouse(tcell.MouseButtonEvents | tcell.MouseDragEvents | tcell.MouseMotionEvents)
	screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)

	return &Editor{screen: screen}, nil
}

func (e *Editor) LoadBuffer(path string) error {
	var buf *buffer.Buffer

	file, err := os.Open(path)

	if err != nil {
		if os.IsNotExist(err) {
			buf = buffer.New()
		} else {
			return err
		}
	} else {
		buf, err = buffer.FromFile(file)
		if err != nil {
			return err
		}
	}

	err = file.Close()
	if err != nil {
		return err
	}

	e.buffer = buf
	return nil
}

// Rendering

func (e *Editor) Render() {
	e.screen.SetStyle(tcell.StyleDefault)
	e.screen.Clear()

	e.buffer.Render(e.screen)

	e.updateCursor()
	e.screen.Sync()
}

// Input

func (e *Editor) GetCursor() *buffer.Cursor {
	return e.buffer.Cursor
}

func (e *Editor) updateCursor() {
	e.screen.ShowCursor(e.buffer.GetCursorVisibleOffset(), e.GetCursor().Y)
}

func (e *Editor) MoveCursor(dX int, dY int) {
	visibleOffset := 0

	if dY != 0 {
		visibleOffset = e.buffer.GetCursorVisibleOffset()
	}

	e.GetCursor().Move(dX, dY)
	e.buffer.ClampCursor()

	if dY != 0 {
		e.GetCursor().X = e.buffer.GetCursorCharacterOffset(visibleOffset)
		e.buffer.ClampCursorX()
	}

	e.updateCursor()
	e.Render()
	e.screen.Show()
}

func (e *Editor) Run() error {
	for {
		event := e.screen.PollEvent()

		switch event := event.(type) {
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyRight:
				e.MoveCursor(1, 0)
			case tcell.KeyLeft:
				e.MoveCursor(-1, 0)

			case tcell.KeyUp:
				e.MoveCursor(0, -1)
			case tcell.KeyDown:
				e.MoveCursor(0, 1)

			case tcell.KeyRune:
				e.buffer.Insert(event.Rune())
				e.Render()

			case tcell.KeyTab:
				e.buffer.Insert('\t')
				e.Render()

			case tcell.KeyBackspace2:
				e.buffer.Delete()
				e.Render()

			case tcell.KeyEnter:
				e.buffer.NewLine()
				e.Render()

			case tcell.KeyEscape:
				e.screen.Fini()
				return nil
			}

		case *tcell.EventMouse:
			if event.Buttons()&tcell.ButtonPrimary != 0 {
				c := e.GetCursor()
				c.X, c.Y = event.Position()

				e.buffer.ClampCursorY()
				c.X = e.buffer.GetCursorCharacterOffset(c.X)

				e.buffer.ClampCursorX()

				e.updateCursor()
				e.screen.Show()
			}

		case *tcell.EventResize:
			e.Render()
		}
	}
}
