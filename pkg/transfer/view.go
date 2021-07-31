package transfer

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TransferList is an extended version of widget.List where is uses a store
// to hold the list items.
type TransferList struct {
	widget.List
	store  *TransferStore
	Parent fyne.Window
}

// NewView creates a new TransferList which is just an extended version
// of widget.List.
//
// Everytime the store content changes the view is refreshed.
func NewView(store *TransferStore) *TransferList {
	tl := &TransferList{
		store:  store,
		Parent: fyne.CurrentApp().Driver().AllWindows()[0],
	}

	store.OnStoreChange = func(i int) {
		tl.Refresh()
	}

	tl.List.Length = tl.length
	tl.List.CreateItem = tl.createItem
	tl.List.UpdateItem = tl.updateItem
	tl.ExtendBaseWidget(tl)

	return tl
}

// createItem creates a new template list item with the
// default widgets and custom layout.
func (tl *TransferList) createItem() fyne.CanvasObject {
	return container.New(
		&transferLayout{},
		widget.NewIcon(nil), // Upload/Download
		widget.NewLabel(""), // Name
		widget.NewLabel(""), // Sender
		widget.NewLabel(""), // Size
		widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{
			Bold: true,
		}), // Status
		container.NewHBox(
			widget.NewProgressBar(),
			widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {}), // Accept
			widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {}),  // Reject
		),
	)
}

// updateItem will be executed for each row of the list when it needs
// to be updated.
func (tl *TransferList) updateItem(i widget.ListItemID, item fyne.CanvasObject) {
	t := tl.store.Get(i)

	wDirection := item.(*fyne.Container).Objects[0].(*widget.Icon)
	wName := item.(*fyne.Container).Objects[1].(*widget.Label)
	wSource := item.(*fyne.Container).Objects[2].(*widget.Label)
	wSize := item.(*fyne.Container).Objects[3].(*widget.Label)
	wStatus := item.(*fyne.Container).Objects[4].(*widget.Label)
	cActions := item.(*fyne.Container).Objects[5].(*fyne.Container)

	// Set the text labels and status
	setItemLabels(t, wStatus, wName, wSize, wSource)

	// If the status is a final status hide the actions container
	showHideActions(t.Status, cActions)

	// These values will not change so they only need to be initiazlied once.
	if wDirection.Resource == nil {

		followProgress(i, cActions.Objects[0].(*widget.ProgressBar), tl.store)

		// Set the transfer direction icon based on the transfer direction
		setItemDirection(wDirection, t.Direction)

		if t.Direction == Upload {
			showHidePBar(true, cActions)
			showHideAccRej(false, cActions)
		} else {
			showHidePBar(false, cActions)
			showHideAccRej(true, cActions)
		}

		if cActions.Objects[1].Visible() {
			cActions.Objects[1].(*widget.Button).OnTapped = func() {
				saveDialog := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
					if err != nil || uc == nil {
						return // if error or uc is null user cancel dialog
					}

					t.LocalFilePath = uc.URI().Path()
					t.Status = Accepted
					tl.store.Update(i, t)

					showHidePBar(true, cActions)
					showHideAccRej(false, cActions)

				}, tl.Parent)
				saveDialog.SetFileName(t.FileName)
				saveDialog.Show()
			}
			cActions.Objects[2].(*widget.Button).OnTapped = func() {
				t.Status = Rejected
				tl.store.Update(i, t)
			}
		}
	}

	if t.Status == Error && t.Error() != nil {
		dialog.ShowError(t.Error(), tl.Parent)

		//! Realy don't like this
		go func() {
			t.err = nil
			tl.store.Update(i, t)
		}()
	}
}

// length return the length of the List
func (tl *TransferList) length() int {
	return tl.store.Size()
}

// followProgress will start a new goroutine to follow the progress
// of a transfer after it got accepted.
func followProgress(i int, pbar *widget.ProgressBar, store *TransferStore) {
	go func() {
		for {
			val, ok := <-store.FollowProgress(i)
			if !ok {
				break
			}
			pbar.SetValue(val)
		}
	}()
}

// showHidePBar will show or hide the progress bar.
func showHidePBar(show bool, cProAction *fyne.Container) {
	if show {
		cProAction.Objects[0].(*widget.ProgressBar).Show()
	} else {
		cProAction.Objects[0].(*widget.ProgressBar).Hide()
	}
}

// showHideAccRej will show or hide the buttons accept and reject.
func showHideAccRej(show bool, cProAction *fyne.Container) {
	if show {
		cProAction.Objects[1].(*widget.Button).Show()
		cProAction.Objects[2].(*widget.Button).Show()
	} else {
		cProAction.Objects[1].(*widget.Button).Hide()
		cProAction.Objects[2].(*widget.Button).Hide()
	}
}

// showHideActions will hide the progress bar and buttons if the status
// is a terminal status or show them otherwise.
func showHideActions(status Status, container *fyne.Container) {
	if status.IsFinal() {
		container.Hide()
	} else {
		container.Show()
	}
}

// setItemLabels will setup the labels displayed on each row.
func setItemLabels(t *Transfer, wStatus, wName, wSize, wSource *widget.Label) {
	// Only update the status if it changes.
	if wStatus.Text != t.Status.String() {
		wStatus.SetText(t.Status.String())
	}

	// Since the labels will not change this will be executed only one time.
	if wName.Wrapping != fyne.TextTruncate {
		wName.Wrapping = fyne.TextTruncate
		wSource.Wrapping = fyne.TextTruncate
		wStatus.Wrapping = fyne.TextTruncate
		wSize.Wrapping = fyne.TextTruncate

		wName.SetText(t.FileName)
		wSize.SetText(byteCountSI(t.FileSize))
		wSource.SetText(t.SenderName)
	}
}

// setItemDirection will set the icon upload or download depending on
// the direction of the transfer.
func setItemDirection(wDirection *widget.Icon, direction Direction) {
	if wDirection != nil && wDirection.Resource == nil {
		switch direction {
		case Upload:
			wDirection.SetResource(theme.UploadIcon())
		case Download:
			wDirection.SetResource(theme.DownloadIcon())
		}
	}
}

// bytCountSI will convert a number of bytes into SI unit string.
func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
