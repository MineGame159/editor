package editor

import (
	"editor/internal/buffer"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"os"
)

type Editor struct {
	views.BoxLayout

	screen tcell.Screen
	app    *views.Application

	firstDraw bool

	buffer     *buffer.Buffer
	bufferView *BufferView
}

// Constructor

// New creates an editor screen with no buffer
func New() (*Editor, error) {
	// Create tcell.Screen
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	// Create editor
	app := &views.Application{}

	editor := &Editor{screen: screen, app: app, firstDraw: true}
	editor.SetOrientation(views.Vertical)

	editor.app.SetScreen(screen)
	editor.app.SetRootWidget(editor)

	// Add widgets
	editor.bufferView = &BufferView{}
	editor.AddWidget(editor.bufferView, 1)

	return editor, nil
}

func (e *Editor) LoadBuffer(path string) error {
	var buf *buffer.Buffer

	file, err := os.Open(path)

	if err != nil {
		if os.IsNotExist(err) {
			buf = buffer.New()
		} else {
			return err
		}
	} else {
		buf, err = buffer.FromFile(file)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	e.buffer = buf
	e.bufferView.buffer = buf

	return nil
}

func (e *Editor) Run() error {
	return e.app.Run()
}

// Widget

func (e *Editor) Draw() {
	if e.firstDraw {
		e.screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)
		e.firstDraw = false
	}

	e.screen.ShowCursor(e.buffer.GetCursorVisibleOffset(), e.buffer.Cursor.Y)

	e.BoxLayout.Draw()
}

func (e *Editor) HandleEvent(event tcell.Event) bool {
	switch event := event.(type) {
	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyEscape:
			e.app.Quit()
		default:
			return e.BoxLayout.HandleEvent(event)
		}

		return true
	}

	return e.BoxLayout.HandleEvent(event)
}

// Input
