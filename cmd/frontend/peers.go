package frontend

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/fabiodcorreia/catch-my-file/internal/store"
	"github.com/fabiodcorreia/catch-my-file/internal/transfer"
)

type peersList struct {
	widget.List
	items        []store.Peer
	SendTransfer chan *store.Transfer
}

func newPeersList() *peersList {
	p := &peersList{
		items:        make([]store.Peer, 0, 1),
		SendTransfer: make(chan *store.Transfer, 1),
	}

	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.ExtendBaseWidget(p)

	return p
}

func (pl *peersList) Length() int {
	return len(pl.items)
}

func (pl *peersList) CreateItem() fyne.CanvasObject {
	return container.NewAdaptiveGrid(
		4,
		widget.NewLabel(""), //Name
		widget.NewLabel(""), //Ip Address
		widget.NewLabel(""), //Port
		widget.NewButton("", func() {}),
	)
}

func (pl *peersList) UpdateItem(i int, item fyne.CanvasObject) {
	item.(*fyne.Container).Objects[0].(*widget.Label).SetText(pl.items[i].Name)
	item.(*fyne.Container).Objects[1].(*widget.Label).SetText(pl.items[i].Address.String())
	item.(*fyne.Container).Objects[2].(*widget.Label).SetText(fmt.Sprintf("%d", pl.items[i].Port))
	item.(*fyne.Container).Objects[3] = widget.NewButton("Send File", func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}

			ft, err := transfer.NewFileTransfer(uc.URI().Path())
			if err != nil {
				return
			}

			rect := canvas.NewRectangle(color.Transparent)
			rect.SetMinSize(fyne.NewSize(200, 0))

			d := dialog.NewCustom(
				"Send File Request",
				fmt.Sprintf("Preparing %s to send to %s", ft.FileName, pl.items[i].Name),
				container.NewMax(rect, widget.NewProgressBarInfinite()),
				fyne.CurrentApp().Driver().AllWindows()[0],
			)

			tc := transfer.NewClient(pl.items[i].Address, pl.items[i].Port, ft)

			tt := store.NewTransfer("", ft.FileName, float64(ft.FileSize), pl.items[i].Name, nil)
			tt.IsToSend = true

			d.Show()
			err = tc.SendRequest()
			d.Hide()

			if err != nil {
				return
			}

			go tc.WaitSendOrStop(tt)
			pl.SendTransfer <- tt
		}, fyne.CurrentApp().Driver().AllWindows()[0])
	})
	item.Refresh()
}
func (pl *peersList) RemoveItem(i int) {
	copy(pl.items[i:], pl.items[i+1:])
	//pl.items[pl.Length()-1] = nil
	pl.items = pl.items[:pl.Length()-1]
	pl.Refresh()
}

func (pl *peersList) NewPeer(p store.Peer) {
	pl.items = append(pl.items, p)
	pl.Refresh()
}
