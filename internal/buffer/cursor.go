package buffer

type Cursor struct {
	X int
	Y int
}

func (c *Cursor) Move(dX int, dY int) {
	c.X += dX
	c.Y += dY
}

func (b *Buffer) ClampCursor() {
	b.ClampCursorY()
	b.ClampCursorX()
}

func (b *Buffer) ClampCursorX() {
	width := len(b.Lines[b.Cursor.Y])

	if b.Cursor.X < 0 {
		b.Cursor.X = 0
	} else if b.Cursor.X > width {
		b.Cursor.X = width
	}
}

func (b *Buffer) ClampCursorY() {
	if b.Cursor.Y < 0 {
		b.Cursor.Y = 0
	} else if b.Cursor.Y >= len(b.Lines) {
		b.Cursor.Y = len(b.Lines) - 1
	}
}

func (b *Buffer) GetCursorVisibleOffset() int {
	offset := 4

	for i, ch := range b.Lines[b.Cursor.Y] {
		if i >= b.Cursor.X {
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

func (b *Buffer) GetCursorCharacterOffset(visible int) int {
	offset := visible - 4

	for i, ch := range b.Lines[b.Cursor.Y] {
		if i >= b.Cursor.X {
			break
		}

		if ch == '\t' {
			offset -= 3
		}
	}

	return offset
}
