package buffer

import (
	"bufio"
	"os"
	"slices"
)

type Buffer struct {
	Lines  []string
	Cursor *Cursor
}

// Constructors

func New() *Buffer {
	return &Buffer{
		Lines:  []string{""},
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
		Lines:  lines,
		Cursor: &Cursor{X: 0, Y: 0},
	}, scanner.Err()
}

// Editing

func (b *Buffer) Insert(ch rune) {
	line := b.Lines[b.Cursor.Y]

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
	b.Lines[b.Cursor.Y] = line
}

func (b *Buffer) Delete() {
	line := b.Lines[b.Cursor.Y]

	if b.Cursor.X == 0 {
		if b.Cursor.Y > 0 {
			above := b.Lines[b.Cursor.Y-1]

			b.Lines[b.Cursor.Y-1] += line
			b.Lines = slices.Delete(b.Lines, b.Cursor.Y, b.Cursor.Y+1)

			b.Cursor.X = len(above)
			b.Cursor.Y--
		}

		return
	}

	left := line[:b.Cursor.X-1]
	right := line[b.Cursor.X:]

	b.Lines[b.Cursor.Y] = left + right
	b.Cursor.X--
}

func (b *Buffer) NewLine() {
	line := b.Lines[b.Cursor.Y]

	left := line[:b.Cursor.X]
	right := line[b.Cursor.X:]

	b.Lines[b.Cursor.Y] = left
	b.Lines = slices.Insert(b.Lines, b.Cursor.Y+1, right)

	b.Cursor.X = 0
	b.Cursor.Y++
}
