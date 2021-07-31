package peer

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func Test_peerLayout_Layout(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	t.Run("invalid number of objects", func(t *testing.T) {
		l := &peerLayout{}

		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Layout expected panic = %v", err)
			}
		}()

		l.Layout([]fyne.CanvasObject{
			container.NewWithoutLayout(),
		}, fyne.NewSize(200, 200))
	})

	t.Run("valid number of objects", func(t *testing.T) {
		l := &peerLayout{}

		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Layout expected no panic = %v", err)
			}
		}()

		objects := []fyne.CanvasObject{
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
		}

		l.Layout(objects, fyne.NewSize(900, 600))

		if err0 := checkPosAndSize(objects[0], 540, 4); err0 != nil {
			t.Errorf("object 0: %v", err0)
		}

		if err1 := checkPosAndSize(objects[1], 151.2, 548); err1 != nil {
			t.Errorf("object 1: %v", err1)
		}

		if err2 := checkPosAndSize(objects[2], 40, 856); err2 != nil {
			t.Errorf("object 2: %v", err2)
		}

	})

}

func checkPosAndSize(obj fyne.CanvasObject, width, posX float32) error {
	if obj.Size().Width != width {
		return fmt.Errorf("expected width = %v but got width = %v", obj.Size().Width, width)
	}

	if obj.Position().X != posX {
		return fmt.Errorf("expected position x = %v but got position x = %v", obj.Position().X, posX)
	}
	return nil
}
