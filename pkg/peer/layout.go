package peer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/fabiodcorreia/catch-my-file/pkg/layout"
)

type peerLayout struct {
	maxMinSizeHeight float32
}

// Layout will calculate the size and position of each object in a row
// of the Peers List.
func (l *peerLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	col1Width := size.Width * 0.60
	col1X := theme.Padding()

	col2Width := (size.Width - col1Width) * 0.42
	col2X := col1X + col1Width + theme.Padding()

	col3Width := float32(40)
	col3X := size.Width - theme.Padding() - col3Width

	layout.ResizeAndMove(objects[0], col1Width, col1X, l.maxMinSizeHeight)
	layout.ResizeAndMove(objects[1], col2Width, col2X, l.maxMinSizeHeight)
	layout.ResizeAndMove(objects[2], col3Width, col3X, l.maxMinSizeHeight)
}

// MinSize will calculate the minimum size allowed that
func (l *peerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	size, height := layout.MinSize(l.maxMinSizeHeight, objects)
	l.maxMinSizeHeight = height
	return size
}
