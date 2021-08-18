package transfer

import (
	"fmt"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fabiodcorreia/catch-my-file/pkg/layout"
)

func Test_setItemDirection(t *testing.T) {
	t.Run("set upload icon", func(t *testing.T) {
		icon := widget.NewIcon(nil)
		setItemDirection(icon, Upload)

		if icon.Resource.Name() != theme.UploadIcon().Name() {
			t.Errorf("set item direction expected icon = %v but got =%v", theme.UploadIcon().Name(), icon.Resource.Name())
		}
	})

	t.Run("set download icon", func(t *testing.T) {
		icon := widget.NewIcon(nil)
		setItemDirection(icon, Download)

		if icon.Resource.Name() != theme.DownloadIcon().Name() {
			t.Errorf("set item direction expected icon = %v but got =%v", theme.DownloadIcon().Name(), icon.Resource.Name())
		}
	})

	t.Run("set download but its already set", func(t *testing.T) {
		icon := widget.NewIcon(theme.CancelIcon())
		setItemDirection(icon, Upload)

		if icon.Resource.Name() != theme.CancelIcon().Name() {
			t.Errorf("set item direction expected icon = %v but got =%v", theme.CancelIcon().Name(), icon.Resource.Name())
		}
	})

	t.Run("widget icon is nil", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("set item idrenction not expected panic = %v", err)
			}
		}()

		setItemDirection(nil, Upload)
	})
}

func Test_byteCountSI(t *testing.T) {
	t.Run("convert bytes", func(t *testing.T) {
		output := byteCountSI(10)
		if output != "10 B" {
			t.Errorf("byteCountSI expected = %v, got %v", "10 B", output)
		}
	})

	t.Run("convert kb", func(t *testing.T) {
		output := byteCountSI(1000)
		if output != "1.0 KB" {
			t.Errorf("byteCountSI expected = %v, got %v", "1.0 KB", output)
		}
	})

	t.Run("convert mb", func(t *testing.T) {
		output := byteCountSI(1000000)
		if output != "1.0 MB" {
			t.Errorf("byteCountSI expected = %v, got %v", "1.0 MB", output)
		}
	})

	t.Run("convert gb", func(t *testing.T) {
		output := byteCountSI(1000000000)
		if output != "1.0 GB" {
			t.Errorf("byteCountSI expected = %v, got %v", "1.0 GB", output)
		}
	})

	t.Run("convert tb", func(t *testing.T) {
		output := byteCountSI(1000000000000)
		if output != "1.0 TB" {
			t.Errorf("byteCountSI expected = %v, got %v", "1.0 TB", output)
		}
	})
}

func Test_setItemLabels(t *testing.T) {
	t.Run("set labels from transfer", func(t *testing.T) {
		wStatus := widget.NewLabel("")
		wName := widget.NewLabel("")
		wSize := widget.NewLabel("")
		wSource := widget.NewLabel("")
		tt := &Transfer{
			Status:     Waiting,
			SenderName: "Peer 1",
			FileName:   "File.txt",
			FileSize:   1000,
		}

		setItemLabels(tt, wStatus, wName, wSize, wSource)

		if wStatus.Text != tt.Status.String() {
			t.Errorf("setItemLabels expected status = %v but got = %v", tt.Status.String(), wStatus.Text)
		}

		if wName.Text != tt.FileName {
			t.Errorf("setItemLabels expected name = %v but got = %v", tt.FileName, wName.Text)
		}

		if wSize.Text != "1.0 KB" {
			t.Errorf("setItemLabels expected size = %v but got = %v", "1.0 KB", wSize.Text)
		}

		if wSource.Text != tt.SenderName {
			t.Errorf("setItemLabels expected source = %v but got = %v", tt.SenderName, wSource.Text)
		}
	})

	t.Run("set labels from transfer and update status if changed", func(t *testing.T) {
		wStatus := widget.NewLabel("")
		wName := widget.NewLabel("")
		wSize := widget.NewLabel("")
		wSource := widget.NewLabel("")
		tt := &Transfer{
			Status:     Waiting,
			SenderName: "Peer 1",
			FileName:   "File.txt",
			FileSize:   1000,
		}

		setItemLabels(tt, wStatus, wName, wSize, wSource)

		if wStatus.Text != tt.Status.String() {
			t.Errorf("setItemLabels expected status = %v but got = %v", tt.Status.String(), wStatus.Text)
		}

		setItemLabels(tt, wStatus, wName, wSize, wSource)

		tt.Status = Completed

		setItemLabels(tt, wStatus, wName, wSize, wSource)

		if wStatus.Text != tt.Status.String() {
			t.Errorf("setItemLabels expected status = %v but got = %v", tt.Status.String(), wStatus.Text)
		}

		if wName.Text != tt.FileName {
			t.Errorf("setItemLabels expected name = %v but got = %v", tt.FileName, wName.Text)
		}

		if wSize.Text != "1.0 KB" {
			t.Errorf("setItemLabels expected size = %v but got = %v", "1.0 KB", wSize.Text)
		}

		if wSource.Text != tt.SenderName {
			t.Errorf("setItemLabels expected source = %v but got = %v", tt.SenderName, wSource.Text)
		}
	})

	t.Run("set labels from transfer and not updated if only labels change", func(t *testing.T) {
		wStatus := widget.NewLabel("")
		wName := widget.NewLabel("")
		wSize := widget.NewLabel("")
		wSource := widget.NewLabel("")
		tt := &Transfer{
			Status:     Waiting,
			SenderName: "Peer 1",
			FileName:   "File.txt",
			FileSize:   1000,
		}

		setItemLabels(tt, wStatus, wName, wSize, wSource)

		tt.FileName = "Super File"

		setItemLabels(tt, wStatus, wName, wSize, wSource)

		if wStatus.Text != tt.Status.String() {
			t.Errorf("setItemLabels expected status = %v but got = %v", tt.Status.String(), wStatus.Text)
		}

		if wName.Text != "File.txt" {
			t.Errorf("setItemLabels expected name = %v but got = %v", "File.txt", wName.Text)
		}

		if wSize.Text != "1.0 KB" {
			t.Errorf("setItemLabels expected size = %v but got = %v", "1.0 KB", wSize.Text)
		}

		if wSource.Text != tt.SenderName {
			t.Errorf("setItemLabels expected source = %v but got = %v", tt.SenderName, wSource.Text)
		}
	})
}

func Test_showHideActions(t *testing.T) {
	t.Run("waiting status", func(t *testing.T) {
		c := container.NewWithoutLayout()
		showHideActions(Waiting, c)
		if c.Hidden {
			t.Errorf("showHideActions expected container hidden but got visible")
		}
	})

	t.Run("accepted status", func(t *testing.T) {
		c := container.NewWithoutLayout()
		showHideActions(Accepted, c)
		if c.Hidden {
			t.Errorf("showHideActions expected container hidden but got visible")
		}
	})

	t.Run("rejected status", func(t *testing.T) {
		c := container.NewWithoutLayout()
		showHideActions(Rejected, c)
		if !c.Hidden {
			t.Errorf("showHideActions expected container visible but got hidden")
		}
	})

	t.Run("completed status", func(t *testing.T) {
		c := container.NewWithoutLayout()
		showHideActions(Completed, c)
		if !c.Hidden {
			t.Errorf("showHideActions expected container visible but got hidden")
		}
	})

	t.Run("error status", func(t *testing.T) {
		c := container.NewWithoutLayout()
		showHideActions(Error, c)
		if !c.Hidden {
			t.Errorf("showHideActions expected container visible but got hidden")
		}
	})
}

func Test_showHideAccRej(t *testing.T) {
	t.Run("show accept and reject buttons", func(t *testing.T) {
		c := container.NewWithoutLayout(
			container.NewWithoutLayout(),
			widget.NewButton("", func() {}),
			widget.NewButton("", func() {}),
		)

		showHideAccRej(true, c)

		if !c.Objects[1].Visible() {
			t.Errorf("showHideAccRej expected accept button visible but got %v", c.Objects[1].Visible())
		}

		if !c.Objects[2].Visible() {
			t.Errorf("showHideAccRej expected accept button visible but got %v", c.Objects[2].Visible())
		}
	})

	t.Run("hide accept and reject buttons", func(t *testing.T) {
		c := container.NewWithoutLayout(
			container.NewWithoutLayout(),
			widget.NewButton("", func() {}),
			widget.NewButton("", func() {}),
		)

		showHideAccRej(false, c)

		if c.Objects[1].Visible() {
			t.Errorf("showHideAccRej expected accept button hiden but got %v", c.Objects[1].Visible())
		}

		if c.Objects[2].Visible() {
			t.Errorf("showHideAccRej expected accept button hiden but got %v", c.Objects[2].Visible())
		}
	})
}

func Test_showHidePBar(t *testing.T) {
	t.Run("show progress bar", func(t *testing.T) {
		c := container.NewWithoutLayout(
			widget.NewProgressBar(),
			widget.NewButton("", func() {}),
			widget.NewButton("", func() {}),
		)

		showHidePBar(true, c)

		if !c.Objects[0].Visible() {
			t.Errorf("showHidePBar expected accept progress bar visible but got %v", c.Objects[0].Visible())
		}

	})

	t.Run("hide progress bar", func(t *testing.T) {
		c := container.NewWithoutLayout(
			widget.NewProgressBar(),
			widget.NewButton("", func() {}),
			widget.NewButton("", func() {}),
		)

		showHidePBar(false, c)

		if c.Objects[0].Visible() {
			t.Errorf("showHidePBar expected accept progress bar hiden but got %v", c.Objects[0].Visible())
		}
	})
}

func Test_TransferList_length(t *testing.T) {

	//
	//
	t.Skip()
	//
	//

	a := test.NewApp()
	defer a.Quit()

	t.Run("initial length is 0", func(t *testing.T) {
		st := NewStore()
		tl := NewView(st)

		if tl.Length() != 0 {
			t.Errorf("Length expected = %v but got = %v", 0, tl.Length())
		}

	})

	t.Run("length increases when items are added to the store", func(t *testing.T) {
		st := NewStore()
		tl := NewView(st)

		st.Add(&Transfer{})

		if tl.Length() != 1 {
			t.Errorf("Length expected = %v but got = %v", 1, tl.Length())
		}
	})
}

func TestTransferList_updateItem(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	t.Run("add transfer upload and download", func(t *testing.T) {
		st := NewStore()
		tl := NewView(st)

		w := a.NewWindow("update item render new transfers")
		w.SetContent(container.NewAppTabs(layout.NewTransferTab(tl)))
		w.Resize(fyne.NewSize(900, 600))

		st.Add(NewTransfer("uploadfile.txt", "123123", "peer-1", 1000000, nil, Upload))

		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "update-item-new-upload.png", w.Canvas().Capture())

		st.Add(NewTransfer("downloadfile.txt", "asdasd", "peer-1", 10000000, nil, Download))

		time.Sleep(500 * time.Millisecond)
		test.AssertImageMatches(t, "update-item-new-download.png", w.Canvas().Capture())

	})

	t.Run("add transfer upload on error state", func(t *testing.T) {
		st := NewStore()
		tl := NewView(st)

		w := a.NewWindow("update item render new transfers")
		w.SetContent(container.NewAppTabs(layout.NewTransferTab(tl)))
		w.Resize(fyne.NewSize(900, 600))
		tl.Parent = w

		tt := NewTransfer("uploadfile.txt", "123123", "peer-1", 1000000, nil, Upload)

		st.Add(tt)
		tt.SetError(fmt.Errorf("sample test error"))
		st.Update(0, tt)

		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "update-item-transfer-error.png", w.Canvas().Capture())
	})

	t.Run("reject transfer", func(t *testing.T) {
		st := NewStore()
		tl := NewView(st)

		w := a.NewWindow("update item render new transfers")
		w.SetContent(container.NewAppTabs(layout.NewTransferTab(tl)))
		w.Resize(fyne.NewSize(900, 600))
		tl.Parent = w

		st.Add(NewTransfer("downloadfile.txt", "123123", "peer-1", 1000000, nil, Download))

		time.Sleep(100 * time.Millisecond)
		test.TapCanvas(w.Canvas(), fyne.NewPos(845, 60))

		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "update-item-transfer-reject.png", w.Canvas().Capture())

		tt := st.Get(0)
		if tt.Status != Rejected {
			t.Errorf("updateItem expected status rejected but got = %v", tt.Status.String())
		}
	})

	t.Run("accept transfer", func(t *testing.T) {
		st := NewStore()
		tl := NewView(st)

		w := a.NewWindow("update item render new transfers")
		w.SetContent(container.NewAppTabs(layout.NewTransferTab(tl)))
		w.Resize(fyne.NewSize(900, 600))
		tl.Parent = w

		tt := NewTransfer("downloadfile.txt", "123123", "peer-1", 1000000, nil, Download)

		st.Add(tt)

		time.Sleep(100 * time.Millisecond)
		test.TapCanvas(w.Canvas(), fyne.NewPos(815, 60))

		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "update-item-transfer-accept.png", w.Canvas().Capture())

		test.TapCanvas(w.Canvas(), fyne.NewPos(600, 450))
		time.Sleep(100 * time.Millisecond)

	})
}
