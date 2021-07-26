package peer

import (
	"fmt"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fabiodcorreia/catch-my-file/pkg2/transfer"
)

type PeerList struct {
	widget.List
	store *PeerStore
	//addTrans func(t *transfer.Transfer) int
}

//func NewListView(store *PeerStore, addTransfer func(t *transfer.Transfer) int) *PeerList {
func newListView(store *PeerStore) *PeerList {
	pl := &PeerList{
		store: store,
		//addTrans: addTransfer,
	}
	pl.store.AddOnChangeListener(func(i int) {
		pl.Refresh()
	})

	pl.List.Length = pl.length
	pl.List.CreateItem = pl.createItem
	pl.List.UpdateItem = pl.updateItem
	pl.ExtendBaseWidget(pl)

	return pl
}

func (pl *PeerList) createItem() fyne.CanvasObject {
	return container.New(&peerListLayout{},
		widget.NewLabel(``),             //Name
		widget.NewLabel(``),             //Ip Address
		widget.NewButton(``, func() {}), //Send File
	)
}

func (pl *PeerList) updateItem(i widget.ListItemID, item fyne.CanvasObject) {
	p := pl.store.Get(i)

	item.(*fyne.Container).Objects[0].(*widget.Label).SetText(p.Name)
	item.(*fyne.Container).Objects[1].(*widget.Label).SetText(p.Address.String())
	item.(*fyne.Container).Objects[2] = widget.NewButtonWithIcon(``, theme.MailSendIcon(), func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}

			name, size, err := transfer.LookupFile(uc.URI().Path())
			if err != nil {
				dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				return
			}

			check, err := transfer.ChecksumFile(uc.URI().Path())
			if err != nil {
				dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				return
			}

			addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%v:%d", p.Address, p.Port))
			if err != nil {
				dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				return
			}

			t := transfer.NewUpTransfer(name, check, "hostname", size, addr)
			t.LocalFilePath = uc.URI().Path()

			//pl.addTrans(t)

		}, fyne.CurrentApp().Driver().AllWindows()[0])
	})
	item.Refresh()
}

func (pl *PeerList) length() int {
	return pl.store.Size()
}
