package frontend

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fabiodcorreia/catch-my-file/internal/store"
)

type transferList struct {
	widget.List
	items []*store.Transfer
}

func newTransferList() *transferList {
	p := &transferList{
		items: make([]*store.Transfer, 0, 5),
	}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.ExtendBaseWidget(p)

	return p
}

func (tl *transferList) Length() int {
	return len(tl.items)
}

func (tl *transferList) CreateItem() fyne.CanvasObject {
	return container.New(&transferLayout{},
		widget.NewLabel(""), // Name
		widget.NewLabel(""), // Sender
		widget.NewLabel(""), // Size
		container.NewHBox(
			widget.NewProgressBar(),
			widget.NewButton("", func() {}),  // Accept
			widget.NewButton("", func() {})), // Reject
	)
}

func onConfirm(t *store.Transfer, cProAction *fyne.Container) {
	saveDialog := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil || uc == nil {
			return
		}
		t.Accept <- uc.URI().Path()
		close(t.Accept)
		showHideButtons(true, cProAction)
	}, fyne.CurrentApp().Driver().AllWindows()[0])
	saveDialog.SetFileName(t.Name)
	saveDialog.Show()
}

func showHideButtons(pbVisible bool, cProAction *fyne.Container) {
	if pbVisible {
		cProAction.Objects[0].(*widget.ProgressBar).Show()
		cProAction.Objects[1].(*widget.Button).Hide()
		cProAction.Objects[2].(*widget.Button).Hide()
	} else {
		cProAction.Objects[0].(*widget.ProgressBar).Hide()
		cProAction.Objects[1].(*widget.Button).Show()
		cProAction.Objects[2].(*widget.Button).Show()
	}
}

func (tl *transferList) UpdateItem(i int, item fyne.CanvasObject) {
	wName := item.(*fyne.Container).Objects[0].(*widget.Label)
	wSource := item.(*fyne.Container).Objects[1].(*widget.Label)
	wSize := item.(*fyne.Container).Objects[2].(*widget.Label)
	cProAction := item.(*fyne.Container).Objects[3].(*fyne.Container)

	// If label size is not set it's the frist update of the item
	if wSize.Text == "" {
		wName.SetText(tl.items[i].Name)
		wSource.SetText(tl.items[i].SourceName)
		wSize.SetText(byteCountSI(tl.items[i].Size))

		showHideButtons(tl.items[i].IsToSend, cProAction)

		tl.items[i].OnProgressChange(func(progress float64) {
			cProAction.Objects[0].(*widget.ProgressBar).SetValue(progress)
		})

		if cProAction.Objects[1].Visible() {
			cProAction.Objects[1] = widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
				onConfirm(tl.items[i], cProAction)
			})

			cProAction.Objects[2] = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				close(tl.items[i].Accept)
				cProAction.Objects[1].(*widget.Button).Hide()
				cProAction.Objects[2].(*widget.Button).Hide()
			})
		}
	}
}

func (tl *transferList) RemoveItem(i int) {
	copy(tl.items[i:], tl.items[i+1:])
	tl.items[tl.Length()-1] = nil
	tl.items = tl.items[:tl.Length()-1]
	tl.Refresh()
}

func (tl *transferList) NewTransfer(t *store.Transfer) {
	tl.items = append(tl.items, t)
	//tl.Refresh()
	if !t.IsToSend {
		go dialog.ShowInformation("Transfer Request Received", t.Name, fyne.CurrentApp().Driver().AllWindows()[0])
	}
}

func byteCountSI(b float64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%f B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

var maxMinSizeHeight float32 // Keeping all instances of the list layout consistent in height

type transferLayout struct{}

// Layout is called to pack all child objects into a specified size.
// The objects is the list of UI objects inside the layout and the
// size is the size of the container.
func (l *transferLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	col1Width := size.Width * 0.60
	col2Width := size.Width * 0.15
	col3Width := size.Width * 0.10
	col4Width := size.Width - col1Width - col2Width - col3Width

	resizeAndMove(objects[0], col1Width, 0, 0)
	resizeAndMove(objects[1], col2Width, objects[0].Position().X, objects[0].Size().Width)
	resizeAndMove(objects[2], col3Width, objects[1].Position().X, objects[1].Size().Width)
	resizeAndMove(objects[3], col4Width, objects[2].Position().X, objects[2].Size().Width)

	cont := objects[3].(*fyne.Container)
	// ProgressBar is visible
	if cont.Objects[0].Visible() {
		cont.Objects[0].Resize(fyne.NewSize(col4Width, objects[3].Size().Height))
		return
	}

	// Buttons are visible
	cont.Objects[1].Resize(fyne.NewSize(col4Width/2, objects[3].Size().Height))
	cont.Objects[2].Resize(fyne.NewSize(col4Width/2, objects[3].Size().Height))
	cont.Objects[2].Move(fyne.NewPos(cont.Objects[1].Position().X+cont.Objects[1].Size().Width, 0))
}

func resizeAndMove(obj fyne.CanvasObject, width float32, prevPositionX float32, prevObjWidth float32) {
	obj.Resize(fyne.NewSize(width, obj.Size().Height))
	obj.Move(fyne.NewPos(prevPositionX+prevObjWidth, 0))
}

// MinSize finds the smallest size that satisfies all the child objects.
// Height will stay consistent between each each instance.
func (g *transferLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	maxMinSizeWidth := float32(0)
	for _, child := range objects {
		if child.Visible() {
			maxMinSizeWidth += child.MinSize().Width
			maxMinSizeHeight = fyne.Max(child.MinSize().Height, maxMinSizeHeight)
		}
	}

	return fyne.NewSize(maxMinSizeWidth, maxMinSizeHeight+theme.Padding())
}
