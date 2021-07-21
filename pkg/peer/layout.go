package peer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type peerListLayout struct {
	maxMinSizeHeight float32
}

func (l *peerListLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	col1Width := size.Width * 0.60
	col1X := theme.Padding()

	col2Width := (size.Width - col1Width) * 0.42
	col2X := col1X + col1Width + theme.Padding()

	col3Width := float32(40)
	col3X := size.Width - theme.Padding() - col3Width

	l.resizeAndMove(objects[0], col1Width, col1X)
	l.resizeAndMove(objects[1], col2Width, col2X)
	l.resizeAndMove(objects[2], col3Width, col3X)
}

func (l *peerListLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var maxMinSizeWidth float32
	for _, child := range objects {
		if child.Visible() {
			maxMinSizeWidth += child.MinSize().Width
			l.maxMinSizeHeight = fyne.Max(child.MinSize().Height, l.maxMinSizeHeight)
		}
	}

	return fyne.NewSize(maxMinSizeWidth, l.maxMinSizeHeight)
}

func (l *peerListLayout) resizeAndMove(obj fyne.CanvasObject, width float32, posX float32) {
	obj.Resize(fyne.NewSize(width, l.maxMinSizeHeight))
	obj.Move(fyne.NewPos(posX, 0))
}
