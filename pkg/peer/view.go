package peer

import (
	"context"
	"image/color"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fabiodcorreia/catch-my-file/pkg/clog"
	"github.com/fabiodcorreia/catch-my-file/pkg/file"
)

// TransferRequest represents the callback that is executed when a new
// transfer is added to the queue to be transferred or waiting for confirmation.
type TransferRequest func(filePath, fileName, checksum, peerNames string, size int64, addr net.Addr)

// PeerList is an extended version of widget.List where is uses a store
// to hold the list items, has a callback to the ouside and has is own
// layout.
type PeerList struct {
	widget.List
	TransferRequest
	store  *PeerStore
	Parent fyne.Window
}

// NewView creates a new PeerList which is just an extended version
// of widget.List.
//
// Everytime the store content changes the view is refreshed.
func NewView(store *PeerStore) *PeerList {
	pl := &PeerList{
		store:  store,
		Parent: fyne.CurrentApp().Driver().AllWindows()[0],
	}

	pl.store.OnPeerStoreChange = func(i int) {
		pl.Refresh()
	}

	pl.List.Length = pl.length
	pl.List.CreateItem = pl.createItem
	pl.List.UpdateItem = pl.updateItem
	pl.ExtendBaseWidget(pl)

	return pl
}

// createItem creates a new template list item with the
// default widgets and custom layout.
func (pl *PeerList) createItem() fyne.CanvasObject {
	return container.New(
		&peerLayout{},
		widget.NewLabel(""), //Name
		widget.NewLabel(""), //Ip Address
		widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {}), //Send File
	)
}

// updateItem will be executed for each row of the list when it needs
// to be updated.
func (pl *PeerList) updateItem(i widget.ListItemID, item fyne.CanvasObject) {
	p := pl.store.Get(i)

	wName := item.(*fyne.Container).Objects[0].(*widget.Label)
	wAddress := item.(*fyne.Container).Objects[1].(*widget.Label)
	wSend := item.(*fyne.Container).Objects[2].(*widget.Button)

	if wName.Text == "" { // The peer information doesn't change
		wName.SetText(p.Name)
		wAddress.SetText(p.IPAddress.String())

		wSend.OnTapped = func() {
			dialog.ShowFileOpen(func(uc fyne.URIReadCloser, openErr error) {
				if openErr != nil || uc == nil {
					return
				}

				filePath := uc.URI().Path()
				name, size, err := file.Lookup(filePath)
				if err != nil {
					clog.Error(err)
					dialog.ShowError(err, pl.Parent)
					return
				}

				go prepRequest(filePath, pl.TransferRequest, name, p.Name, size, p.Address, pl.Parent)
			}, pl.Parent)
		}
	}
}

// length return the length of the List
func (pl *PeerList) length() int {
	return pl.store.Size()
}

// prepFileDialog will create and return a new dialog
// to inform the user that the file is getting prepared.
func prepFileDialog(parent fyne.Window) dialog.Dialog {
	bar := widget.NewProgressBarInfinite()
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(200, 0))
	wrapper := container.NewMax(rect, bar)

	return dialog.NewCustom("Preparing file to send", "Cancel", wrapper, parent)
}

// prepRequest will show a progress dialog while is generating the
// file checksum and send it to the worker to handle.
func prepRequest(path string, req TransferRequest, name, peer string, size int64, addr net.Addr, parent fyne.Window) {
	d := prepFileDialog(parent)
	d.Show()

	ctx, cancel := context.WithCancel(context.Background())

	d.SetOnClosed(func() {
		cancel()
	})

	f, oErr := file.Open(path, file.OPEN_READ)
	if oErr != nil {
		clog.Error(oErr)
		de := dialog.NewError(oErr, parent)
		de.SetOnClosed(func() {
			d.Hide()
		})
		de.Show()
		return
	}

	defer func() {
		if cErr := f.Close(); cErr != nil {
			clog.Error(cErr)
		}
	}()

	check, cErr := file.Checksum(ctx, f)
	if cErr != nil {
		if ctx.Err() == nil { //This means the user pressed cancel.
			dialog.ShowError(cErr, parent)
		}
		d.Hide()
		return
	}

	req(path, name, check, peer, size, addr)
	d.Hide()
}
