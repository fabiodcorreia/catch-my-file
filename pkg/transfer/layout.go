package transfer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/fabiodcorreia/catch-my-file/pkg/layout"
)

type transferLayout struct {
	maxMinSizeHeight float32
}

// Layout will calculate the size and position of each object in a row
// of the Transfers List.
func (l *transferLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	col1Width := float32(25)
	col1X := theme.Padding()

	col2Width := (size.Width - col1Width) * 0.42
	col2X := col1X + col1Width + theme.Padding()

	col3Width := (size.Width - col2Width) * 0.25
	col3X := col2X + col2Width + theme.Padding()

	col4Width := (size.Width - col3Width) * 0.10
	col4X := col3X + col3Width + theme.Padding()

	col5Width := (size.Width - col4Width) * 0.13
	col5X := col4X + col4Width + theme.Padding()

	col6Width := (size.Width - col5X - col5Width)
	col6X := col5X + col5Width - theme.Padding()

	layout.ResizeAndMove(objects[0], col1Width, col1X, l.maxMinSizeHeight)
	layout.ResizeAndMove(objects[1], col2Width, col2X, l.maxMinSizeHeight)
	layout.ResizeAndMove(objects[2], col3Width, col3X, l.maxMinSizeHeight)
	layout.ResizeAndMove(objects[3], col4Width, col4X, l.maxMinSizeHeight)
	layout.ResizeAndMove(objects[4], col5Width, col5X, l.maxMinSizeHeight)
	layout.ResizeAndMove(objects[5], col6Width, col6X, l.maxMinSizeHeight)

	col6 := objects[5].(*fyne.Container)
	if col6.Objects[0].Visible() {
		layout.ResizeAndMove(col6.Objects[0], col6Width, theme.Padding(), l.maxMinSizeHeight)
		return
	}

	col62Width := float32(40)
	col62X := theme.Padding()
	layout.ResizeAndMove(col6.Objects[1], col62Width, col62X, l.maxMinSizeHeight)

	col63Width := col62Width
	col63X := col62X + col62Width + theme.Padding()
	layout.ResizeAndMove(col6.Objects[2], col63Width, col63X, l.maxMinSizeHeight)
}

// MinSize will calculate the minimum size allowed that.
func (l *transferLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	size, height := layout.MinSize(l.maxMinSizeHeight, objects)
	l.maxMinSizeHeight = height
	return size
}
