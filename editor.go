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
	e.screen.ShowCursor(e.buffer.cursor.X+4, e.buffer.cursor.Y)
}

func (e *Editor) UpdateCursor() {
	c := &e.buffer.cursor

	// Clamp Y
	if c.Y < 0 {
		c.Y = 0
	} else if c.Y >= len(e.buffer.lines) {
		c.Y = len(e.buffer.lines) - 1
	}

	// Clamp X
	width := stringWidth(e.buffer.lines[c.Y])

	if c.X < 0 {
		c.X = 0
	} else if c.X > width {
		c.X = width
	}

	e.setCursor()
	e.screen.Show()
}

func (e *Editor) MoveCursor(dX int, dY int) {
	e.buffer.cursor.Move(dX, dY)
	e.UpdateCursor()
}

func stringWidth(str string) int {
	width := 0

	for _, ch := range str {
		if ch == '\t' {
			width += 4
		} else {
			width++
		}
	}

	return width
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

			case tcell.KeyEscape:
				e.screen.Fini()
				return
			}

		case *tcell.EventMouse:
			if event.Buttons()&tcell.ButtonPrimary != 0 {
				c := &e.buffer.cursor

				c.X, c.Y = event.Position()
				c.X -= 4

				e.UpdateCursor()
			}

		case *tcell.EventResize:
			e.Render()
		}
	}
}
