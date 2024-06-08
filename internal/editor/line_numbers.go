package editor

import (
	"editor/internal/buffer"
	"editor/internal/config"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"math"
)

type LineNumbers struct {
	views.WidgetWatchers

	view views.View

	cursor *buffer.Cursor
	min    int
	max    int

	width int
}

func (l *LineNumbers) SetRange(buffer *buffer.Buffer) {
	l.cursor = buffer.Cursor
	l.min = 1
	l.max = len(buffer.Lines)

	l.width = 1 + (int(math.Log10(float64(l.max))) + 1) + 1
}

func (l *LineNumbers) Draw() {
	_, height := l.view.Size()
	y := 0

	for i := l.min; i <= min(l.max, height-1); i++ {
		isCursorLine := i-1 == l.cursor.Y

		numStyle := config.LineNumberStyle
		if isCursorLine {
			numStyle = config.LineStyle
		}

		DrawStr(l.view, 1, y, fmt.Sprintf("%2d", i), numStyle)
		y++
	}
}

func (l *LineNumbers) Resize() {
}

func (l *LineNumbers) HandleEvent(ev tcell.Event) bool {
	return false
}

func (l *LineNumbers) SetView(view views.View) {
	l.view = view
}

func (l *LineNumbers) Size() (int, int) {
	return l.width, 0
}
