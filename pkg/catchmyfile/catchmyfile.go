package catchmyfile

import (
	"context"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/fabiodcorreia/catch-my-file/pkg/clog"
	"github.com/fabiodcorreia/catch-my-file/pkg/layout"
	"github.com/fabiodcorreia/catch-my-file/pkg/network"
	"github.com/fabiodcorreia/catch-my-file/pkg/peer"
	"github.com/fabiodcorreia/catch-my-file/pkg/transfer"
	"github.com/fabiodcorreia/catch-my-file/pkg/worker"
)

type CatchMyFileApp struct {
	a     fyne.App
	w     fyne.Window
	ctx   context.Context
	port  int
	wPool worker.WorkerPool
}

// New will create a new instance of the appplication.
func New(port int) *CatchMyFileApp {
	c := &CatchMyFileApp{
		a:     app.NewWithID("github.fabiodcorreia.catch-my-file"),
		ctx:   context.Background(),
		port:  port,
		wPool: worker.NewPool(2),
	}
	return c
}

// initSetup will initialize the window, size, position
// and the OnClose action of the window.
func (c *CatchMyFileApp) initSetup() {
	ctx, cancel := context.WithCancel(c.ctx)
	c.ctx = ctx

	c.w = c.a.NewWindow("Catch My File")
	c.w.Resize(fyne.NewSize(900, 600))
	c.w.CenterOnScreen()
	c.w.SetMaster()

	c.w.SetOnClosed(func() {
		cancel()
	})
}

// handleError will log the error and show a popup on the ui.
func handleError(err error, parent fyne.Window) {
	clog.Error(err)
	dialog.ShowError(err, parent)
}

// Run will start all the required resources and open the
// application window.
func (c *CatchMyFileApp) Run() error {
	c.initSetup()

	pStore := peer.NewStore()
	pView := peer.NewView(pStore)
	pServer := peer.NewServer(network.Hostname(), c.port, pStore)

	tStore := transfer.NewStore()
	tView := transfer.NewView(tStore)
	tReceiver := transfer.NewReceiver(c.port, tStore)

	pView.TransferRequest = func(filePath, fileName, checksum, peerName string, size int64, addr net.Addr) {
		t := transfer.NewTransfer(fileName, checksum, peerName, size, addr, transfer.Upload)
		t.LocalFilePath = filePath
		i := tStore.Add(t)
		c.onTransferRequest(i, tStore, pStore)
	}

	pDone := make(chan interface{})
	if err := pServer.Run(c.ctx, pDone); err != nil {
		handleError(err, c.w)
		c.a.Quit()
	}

	rDone := make(chan interface{})
	if err := tReceiver.Run(c.ctx, rDone); err != nil {
		handleError(err, c.w)
		c.a.Quit()
	}

	c.w.SetContent(container.NewAppTabs(
		layout.NewPeersTab(pView),
		layout.NewTransferTab(tView),
	))

	c.wPool.Run(c.ctx)

	c.w.ShowAndRun()

	<-pDone        // Wait for Peers server to finish
	<-rDone        // Wait for Receiver server to finish
	c.wPool.Stop() // Wait for the WorkerPool to finish

	return nil
}

// onTransferRequest is the action that is executed everytime
// a new transfer is added by the user to be sent to a peer.
func (c *CatchMyFileApp) onTransferRequest(i int, tStore *transfer.TransferStore, pStore *peer.PeerStore) {
	t := tStore.Get(i)
	p := pStore.GetMe()

	err := c.wPool.AddTask(func(ctx context.Context) {
		clog.Info("Added transfer idx:%d to the worker", i)
		t.SenderName = p.Name
		conn, err := transfer.SendTransferReq(ctx, t)
		if err != nil {
			clog.Error(err)
			t.SetError(err)
			tStore.Update(i, t)
			return
		}

		defer conn.Close()

		if err = transfer.WaitConfirmation(ctx, i, conn, tStore); err != nil {
			switch err {
			case transfer.ErrRejected:
				t.Status = transfer.Rejected
			default:
				clog.Error(err)
				t.SetError(err)
			}
			tStore.Update(i, t)
			return
		}

		t.Status = transfer.Completed
		tStore.Update(i, t)
	})

	if err != nil {
		clog.Error(err)
		t.SetError(err)
		tStore.Update(i, t)
	}
}
