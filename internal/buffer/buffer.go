package buffer

import (
	"bufio"
	"editor/internal/config"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	"slices"
)

type Buffer struct {
	lines  []string
	Cursor *Cursor
}

// Constructors

func New() *Buffer {
	return &Buffer{
		lines:  []string{""},
		Cursor: &Cursor{X: 0, Y: 0},
	}
}

func FromFile(file *os.File) (*Buffer, error) {
	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return &Buffer{
		lines:  lines,
		Cursor: &Cursor{X: 0, Y: 0},
	}, scanner.Err()
}

// Rendering

func (b *Buffer) Render(s tcell.Screen) {
	for y := 0; y < len(b.lines); y++ {
		isCursorLine := y == b.Cursor.Y

		numStyle := config.LineNumberStyle
		if isCursorLine {
			numStyle = config.LineStyle
		}

		x := b.drawStr(s, 1, y, fmt.Sprintf("%2d", y+1), numStyle) + 1
		// TODO: x = b.drawChar(s, x, y, tcell.RuneVLine, config.LineNumberStyle)
		b.drawStr(s, x, y, b.lines[y], config.LineStyle)
	}
}

func (b *Buffer) drawStr(s tcell.Screen, x int, y int, str string, style tcell.Style) int {
	for _, ch := range str {
		if ch == '\t' {
			for i := 0; i < config.TabWidth; i++ {
				x = b.drawChar(s, x, y, ' ', style)
			}
		} else {
			x = b.drawChar(s, x, y, ch, style)
		}
	}

	return x
}

func (b *Buffer) drawChar(s tcell.Screen, x int, y int, ch rune, style tcell.Style) int {
	s.SetContent(x, y, ch, nil, style)
	return x + 1
}

// Editing

func (b *Buffer) Insert(ch rune) {
	line := b.lines[b.Cursor.Y]

	if b.Cursor.X == 0 {
		line = string(ch) + line
	} else if b.Cursor.X == len(line) {
		line = line + string(ch)
	} else {
		left := line[:b.Cursor.X]
		right := line[b.Cursor.X:]

		line = left + string(ch) + right
	}

	b.Cursor.X++
	b.lines[b.Cursor.Y] = line
}

func (b *Buffer) Delete() {
	line := b.lines[b.Cursor.Y]

	if b.Cursor.X == 0 {
		if b.Cursor.Y > 0 {
			above := b.lines[b.Cursor.Y-1]

			b.lines[b.Cursor.Y-1] += line
			b.lines = slices.Delete(b.lines, b.Cursor.Y, b.Cursor.Y+1)

			b.Cursor.X = len(above)
			b.Cursor.Y--
		}

		return
	}

	left := line[:b.Cursor.X-1]
	right := line[b.Cursor.X:]

	b.lines[b.Cursor.Y] = left + right
	b.Cursor.X--
}

func (b *Buffer) NewLine() {
	line := b.lines[b.Cursor.Y]

	left := line[:b.Cursor.X]
	right := line[b.Cursor.X:]

	b.lines[b.Cursor.Y] = left
	b.lines = slices.Insert(b.lines, b.Cursor.Y+1, right)

	b.Cursor.X = 0
	b.Cursor.Y++
}
