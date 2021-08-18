package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// MinSize will calculate the minium size that will feat each object.
//
// The Height is returned in separate to make a standard height across
// all the rows.
func MinSize(maxMinSizeHeight float32, objects []fyne.CanvasObject) (fyne.Size, float32) {
	var maxMinSizeWidth float32
	for _, child := range objects {
		if child.Visible() {
			maxMinSizeWidth += child.MinSize().Width
			maxMinSizeHeight = fyne.Max(child.MinSize().Height, maxMinSizeHeight)
		}
	}

	return fyne.NewSize(maxMinSizeWidth, maxMinSizeHeight), maxMinSizeHeight
}

// ResizeAndMove will adjust the size of an fyne.CanvasObject and will
// move it to the specificed X position. The Y position will always be
// 0 of the containers of the object.
func ResizeAndMove(obj fyne.CanvasObject, width, posX, height float32) {
	obj.Resize(fyne.NewSize(width, height))
	obj.Move(fyne.NewPos(posX, 0))
}

// NewTab creates a new tab icon for the peers.
func NewPeersTab(w fyne.Widget) *container.TabItem {
	return container.NewTabItemWithIcon("Peers", theme.ComputerIcon(), w)
}

// NewTab creates a new tab icon for the Transfers.
func NewTransferTab(w fyne.Widget) *container.TabItem {
	return container.NewTabItemWithIcon("Transfers", theme.StorageIcon(), w)
}
