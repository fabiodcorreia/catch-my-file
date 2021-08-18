package peer

import (
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/fabiodcorreia/catch-my-file/pkg/layout"
)

func Test_PeerList_length(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	t.Run("initial length is 0", func(t *testing.T) {
		st := NewStore()
		pl := NewView(st)

		if pl.Length() != 0 {
			t.Errorf("Length expected = %v but got = %v", 0, pl.Length())
		}

	})

	t.Run("length increases when items are added to the store", func(t *testing.T) {
		st := NewStore()
		pl := NewView(st)

		st.Add(newPeer("peer", net.ParseIP("192.168.1.1"), 0, &net.TCPAddr{
			IP:   net.ParseIP("192.168.1.1"),
			Port: 8822,
		}))

		if pl.Length() != 1 {
			t.Errorf("Length expected = %v but got = %v", 1, pl.Length())
		}
	})
}

func Test_prepFileDialog(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	t.Run("its visible after show", func(t *testing.T) {
		w := a.NewWindow("test prep file dialog")
		w.Resize(fyne.NewSize(900, 600))
		d := prepFileDialog(w)
		d.Show()

		test.AssertImageMatches(t, "prep-dialog-open.png", w.Canvas().Capture())

	})

	t.Run("its not visible after tap cancel", func(t *testing.T) {
		w := a.NewWindow("test prep file dialog")
		w.Resize(fyne.NewSize(900, 600))
		d := prepFileDialog(w)
		d.Show()

		test.TapCanvas(w.Canvas(), fyne.NewPos(450, 350))
		test.AssertImageMatches(t, "prep-dialog-cancel.png", w.Canvas().Capture())

	})
}

func Test_prepRequest(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	path, _ := os.Executable()

	t.Run("request prepared and dialog shown", func(t *testing.T) {
		w := a.NewWindow("request prepared and dialog shown")
		w.Resize(fyne.NewSize(900, 600))

		var wg sync.WaitGroup
		wg.Add(1)
		go prepRequest(
			path,
			func(filePath, fileName, checksum, peerNames string, size int64, addr net.Addr) {
				if filePath != path {
					t.Errorf("prepRequest expected file path = %v but got =%v", path, fileName)
				}
				wg.Done()
			},
			"testfile",
			"peer-1",
			30233144,
			nil,
			w,
		)
		test.AssertImageMatches(t, "prep-request-open.png", w.Canvas().Capture())
		wg.Wait()
	})

	t.Run("request canceled and dialog closed", func(t *testing.T) {
		w := a.NewWindow("request canceled and dialog closed")
		w.Resize(fyne.NewSize(900, 600))
		var wg sync.WaitGroup
		wg.Add(1)
		go prepRequest(
			path,
			func(filePath, fileName, checksum, peerNames string, size int64, addr net.Addr) {
				wg.Done()
			},
			"testfile",
			"peer-1",
			30233144,
			nil,
			w,
		)
		test.TapCanvas(w.Canvas(), fyne.NewPos(450, 350))
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "prep-request-cancel.png", w.Canvas().Capture())
		wg.Wait()
	})

	t.Run("invalid file path", func(t *testing.T) {
		w := a.NewWindow("invalid file path")
		w.Resize(fyne.NewSize(900, 600))
		go prepRequest(
			".invalid",
			func(filePath, fileName, checksum, peerNames string, size int64, addr net.Addr) {},
			"testfile",
			"peer-1",
			30233144,
			nil,
			w,
		)
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "prep-request-invalid-file.png", w.Canvas().Capture())
	})

	t.Run("invalid file path press OK", func(t *testing.T) {
		w := a.NewWindow("invalid file path press OK")
		w.Resize(fyne.NewSize(900, 600))
		go prepRequest(
			".invalid",
			func(filePath, fileName, checksum, peerNames string, size int64, addr net.Addr) {},
			"testfile",
			"peer-1",
			30233144,
			nil,
			w,
		)
		time.Sleep(100 * time.Millisecond)
		test.TapCanvas(w.Canvas(), fyne.NewPos(450, 350))
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "prep-request-invalid-file-ok.png", w.Canvas().Capture())
	})

}

func TestPeerList_updateItem(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	t.Run("update item renders new peers", func(t *testing.T) {
		st := NewStore()
		pl := NewView(st)

		w := a.NewWindow("update item render new peers")
		w.SetContent(container.NewAppTabs(layout.NewPeersTab(pl)))
		w.Resize(fyne.NewSize(900, 600))

		st.Add(newPeer("peer-1", net.ParseIP("192.168.1.1"), 0, &net.TCPAddr{
			IP:   net.ParseIP("192.168.1.1"),
			Port: 8822,
		}))

		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "update-item-1.png", w.Canvas().Capture())

		st.Add(newPeer("peer-2", net.ParseIP("192.168.1.2"), 0, &net.TCPAddr{
			IP:   net.ParseIP("192.168.1.2"),
			Port: 8822,
		}))

		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "update-item-2.png", w.Canvas().Capture())
	})

	t.Run("update item send file", func(t *testing.T) {
		st := NewStore()
		pl := NewView(st)

		w := a.NewWindow("update item send file")
		w.SetContent(container.NewAppTabs(layout.NewPeersTab(pl)))
		w.Resize(fyne.NewSize(900, 600))
		pl.Parent = w

		st.Add(newPeer("peer-1", net.ParseIP("192.168.1.1"), 0, &net.TCPAddr{
			IP:   net.ParseIP("192.168.1.1"),
			Port: 8822,
		}))

		time.Sleep(100 * time.Millisecond)
		test.TapCanvas(w.Canvas(), fyne.NewPos(875, 60))
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "update-item-send.png", w.Canvas().Capture())
	})

	t.Run("update item send file cancel", func(t *testing.T) {
		st := NewStore()
		pl := NewView(st)

		w := a.NewWindow("update item send file cancel")
		w.SetContent(container.NewAppTabs(layout.NewPeersTab(pl)))
		w.Resize(fyne.NewSize(900, 600))
		pl.Parent = w

		st.Add(newPeer("peer-1", net.ParseIP("192.168.1.1"), 0, &net.TCPAddr{
			IP:   net.ParseIP("192.168.1.1"),
			Port: 8822,
		}))

		time.Sleep(100 * time.Millisecond)
		test.TapCanvas(w.Canvas(), fyne.NewPos(875, 60))
		time.Sleep(100 * time.Millisecond)
		test.TapCanvas(w.Canvas(), fyne.NewPos(600, 450))
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "update-item-send-cancel.png", w.Canvas().Capture())
	})
}
