package main

import (
	"bufio"
	"os"
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
	return &Buffer{}
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
