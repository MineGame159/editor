package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	screen tcell.Screen
	buffer *Buffer
}

func NewEditor() *Editor {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	err = screen.Init()
	if err != nil {
		panic(err)
	}

	screen.EnableMouse(tcell.MouseButtonEvents | tcell.MouseDragEvents | tcell.MouseMotionEvents)
	screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)

	e := &Editor{
		screen: screen,
		buffer: NewBuffer(),
	}

	e.Render()
	e.setCursor()

	return e
}

func (e *Editor) Render() {
	e.screen.SetStyle(tcell.StyleDefault)
	e.screen.Clear()

	lineNumberStyle := tcell.StyleDefault.Foreground(tcell.ColorDimGray)

	_, height := e.screen.Size()

	for y := 0; y < height; y++ {
		x := e.write(1, y, fmt.Sprintf("%2d", y+1), lineNumberStyle) + 1

		if y < len(e.buffer.lines) {
			e.write(x, y, e.buffer.lines[y], tcell.StyleDefault)
		}
	}

	e.setCursor()
	e.screen.Sync()
}

func (e *Editor) write(x int, y int, str string, style tcell.Style) int {
	for _, ch := range str {
		if ch == '\t' {
			x += 4
		} else {
			e.screen.SetContent(x, y, ch, nil, style)
			x++
		}
	}

	return x
}

func (e *Editor) setCursor() {
	e.screen.ShowCursor(e.getCursorVisibleOffset(), e.buffer.cursor.Y)
}

func (e *Editor) ClampCursorY() {
	c := &e.buffer.cursor

	if c.Y < 0 {
		c.Y = 0
	} else if c.Y >= len(e.buffer.lines) {
		c.Y = len(e.buffer.lines) - 1
	}
}

func (e *Editor) ClampCursorX() {
	c := &e.buffer.cursor
	width := len(e.buffer.lines[c.Y])

	if c.X < 0 {
		c.X = 0
	} else if c.X > width {
		c.X = width
	}
}

func (e *Editor) MoveCursor(dX int, dY int) {
	visibleOffset := 0

	if dY != 0 {
		visibleOffset = e.getCursorVisibleOffset()
	}

	e.buffer.cursor.Move(dX, dY)
	e.ClampCursorY()
	e.ClampCursorX()

	if dY != 0 {
		e.buffer.cursor.X = e.getCursorCharacterOffset(visibleOffset)
		e.ClampCursorX()
	}

	e.setCursor()
	e.screen.Show()
}

func (e *Editor) getCursorVisibleOffset() int {
	offset := 4

	for i, ch := range e.buffer.lines[e.buffer.cursor.Y] {
		if i >= e.buffer.cursor.X {
			break
		}

		if ch == '\t' {
			offset += 4
		} else {
			offset++
		}
	}

	return offset
}

func (e *Editor) getCursorCharacterOffset(visible int) int {
	offset := visible - 4

	for i, ch := range e.buffer.lines[e.buffer.cursor.Y] {
		if i >= e.buffer.cursor.X {
			break
		}

		if ch == '\t' {
			offset -= 3
		}
	}

	return offset
}

func (e *Editor) Run() {
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
				return
			}

		case *tcell.EventMouse:
			if event.Buttons()&tcell.ButtonPrimary != 0 {
				c := &e.buffer.cursor
				c.X, c.Y = event.Position()

				e.ClampCursorY()
				c.X = e.getCursorCharacterOffset(c.X)

				e.ClampCursorX()

				e.setCursor()
				e.screen.Show()
			}

		case *tcell.EventResize:
			e.Render()
		}
	}
}
