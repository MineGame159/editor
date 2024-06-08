package editor

import (
	"editor/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

func DrawStr(view views.View, x int, y int, str string, style tcell.Style) int {
	for _, ch := range str {
		if ch == '\t' {
			for i := 0; i < config.TabWidth; i++ {
				x = DrawChar(view, x, y, ' ', style)
			}
		} else {
			x = DrawChar(view, x, y, ch, style)
		}
	}

	return x
}

func DrawChar(view views.View, x int, y int, ch rune, style tcell.Style) int {
	view.SetContent(x, y, ch, nil, style)
	return x + 1
}
