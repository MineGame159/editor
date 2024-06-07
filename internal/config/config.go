package config

import "github.com/gdamore/tcell/v2"

var (
	TabWidth        = 4
	LineStyle       = tcell.StyleDefault.Foreground(tcell.ColorWhite)
	LineNumberStyle = tcell.StyleDefault.Foreground(tcell.ColorDimGray)
)
