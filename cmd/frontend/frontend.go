package frontend

import (
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/fabiodcorreia/catch-my-file/cmd/internal/backend"
)

const (
	prefPort     = "port"
	prefHostname = "hostname"
)

type Frontend struct {
	a fyne.App
	w fyne.Window
}

func New() *Frontend {
	f := &Frontend{
		a: app.NewWithID("github.fabiodcorreia.catch-my-file"),
	}
	f.w = f.a.NewWindow("Catch My File")
	f.w.Resize(fyne.NewSize(880, 600))
	return f
}

func (f *Frontend) Run() error {
	pl := newPeersList()
	rl := newTransferList()
	sl := newTransferList()
	lt := newLogTable()
	f.w.SetContent(container.NewAppTabs(
		container.NewTabItemWithIcon("Receiving", theme.DownloadIcon(), rl),
		container.NewTabItemWithIcon("Sending", theme.UploadIcon(), sl),
		container.NewTabItemWithIcon("Peers", theme.ComputerIcon(), pl),
		container.NewTabItemWithIcon("Log", theme.ErrorIcon(), lt),
	))

	hn, err := os.Hostname()
	if err != nil {
		return err
	}

	f.a.Preferences().SetString(prefHostname, strings.ReplaceAll(hn, ".", "-")) //! ?
	f.a.Preferences().SetInt(prefPort, 8820)

	e := backend.NewEngine(f.a.Preferences().String(prefHostname), f.a.Preferences().Int(prefPort))

	go func() {
		for {
			select {
			case p := <-e.DiscoverPeers():
				pl.NewPeer(p)
			case t := <-e.ReceiveTransferNotification():
				rl.NewTransfer(t)
			case t := <-pl.SendTransfer:
				sl.NewTransfer(t)
			case rds := <-e.DiscoverServerError():
				lt.NewLogRecord(time.Now(), rds)
			case rdc := <-e.DiscoverClientError():
				lt.NewLogRecord(time.Now(), rdc)
			case rts := <-e.TransferServerError():
				lt.NewLogRecord(time.Now(), rts)
			}
		}
	}()

	f.w.CenterOnScreen()
	f.w.SetMaster()
	f.w.SetOnClosed(func() {
		e.Shutdown()
	})

	err = e.Start()
	if err != nil {
		f.a.SendNotification(fyne.NewNotification("Catch My File - Fail to Start", err.Error()))
		return err
	}

	f.w.ShowAndRun()
	return nil
}
