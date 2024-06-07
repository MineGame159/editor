package main

import (
	"bufio"
	"os"
	"slices"
)

type Cursor struct {
	X int
	Y int
}

func (c *Cursor) Move(dX int, dY int) {
	c.X += dX
	c.Y += dY
}

type Buffer struct {
	lines  []string
	cursor Cursor
}

func NewBuffer() *Buffer {
	return &Buffer{
		lines: []string{""},
	}
}

func (b *Buffer) LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	scanner := bufio.NewScanner(file)
	b.lines = nil

	for scanner.Scan() {
		line := scanner.Text()

		b.lines = append(b.lines, line)
	}

	b.cursor.Y = len(b.lines) - 1
	b.cursor.X = len(b.lines[b.cursor.Y])

	return scanner.Err()
}

func (b *Buffer) Insert(ch rune) {
	line := b.lines[b.cursor.Y]

	if b.cursor.X == 0 {
		line = string(ch) + line
	} else if b.cursor.X == len(line) {
		line = line + string(ch)
	} else {
		left := line[:b.cursor.X]
		right := line[b.cursor.X:]

		line = left + string(ch) + right
	}

	b.cursor.X++
	b.lines[b.cursor.Y] = line
}

func (b *Buffer) Delete() {
	line := b.lines[b.cursor.Y]

	if b.cursor.X == 0 {
		if b.cursor.Y > 0 {
			above := b.lines[b.cursor.Y-1]

			b.lines[b.cursor.Y-1] += line
			b.lines = slices.Delete(b.lines, b.cursor.Y, b.cursor.Y+1)

			b.cursor.X = len(above)
			b.cursor.Y--
		}

		return
	}

	left := line[:b.cursor.X-1]
	right := line[b.cursor.X:]

	b.lines[b.cursor.Y] = left + right
	b.cursor.X--
}

func (b *Buffer) NewLine() {
	line := b.lines[b.cursor.Y]

	left := line[:b.cursor.X]
	right := line[b.cursor.X:]

	b.lines[b.cursor.Y] = left
	b.lines = slices.Insert(b.lines, b.cursor.Y+1, right)

	b.cursor.X = 0
	b.cursor.Y++
}
