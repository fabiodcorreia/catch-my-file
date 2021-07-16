package frontend

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type logTable struct {
	widget.Table
	items []error
}

func newLogTable() *logTable {
	lt := &logTable{
		items: make([]error, 0, 20),
	}
	lt.Table.Length = lt.Length
	lt.Table.CreateCell = lt.CreateCell
	lt.Table.UpdateCell = lt.UpdateCell
	lt.Table.SetColumnWidth(0, 160)
	lt.ExtendBaseWidget(lt)

	return lt
}

func (lt *logTable) Length() (int, int) {
	return len(lt.items), 2
}

func (lt *logTable) CreateCell() fyne.CanvasObject {
	return widget.NewLabel("Name")
}

func (lt *logTable) UpdateCell(id widget.TableCellID, item fyne.CanvasObject) {
	switch id.Col {
	case 0:
		item.(*widget.Label).SetText(time.Now().Format("2006-01-02 15:04:05"))
	case 1:
		item.(*widget.Label).SetText(lt.items[id.Row].Error())
	}
}

func (lt *logTable) NewLogRecord(t time.Time, err error) {
	if err != nil {
		lt.items = append(lt.items, err)
		lt.Refresh()
	}
}
