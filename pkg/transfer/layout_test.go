package transfer

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func Test_transferListLayout_Layout(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	t.Run("invalid number of objects", func(t *testing.T) {
		l := &transferLayout{}

		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Layout expected panic = %v", err)
			}
		}()

		l.Layout([]fyne.CanvasObject{
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
		}, fyne.NewSize(200, 200))
	})

	t.Run("valid number of objects with progress bar", func(t *testing.T) {
		l := &transferLayout{}

		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Layout expected no panic = %v", err)
			}
		}()

		objects := []fyne.CanvasObject{
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(
				container.NewWithoutLayout(),
				container.NewWithoutLayout(),
				container.NewWithoutLayout(),
			),
		}

		l.Layout(objects, fyne.NewSize(900, 600))

		if err0 := checkPosAndSize(objects[0], 25, 4); err0 != nil {
			t.Errorf("object 0: %v", err0)
		}

		if err1 := checkPosAndSize(objects[1], 367.5, 33); err1 != nil {
			t.Errorf("object 1: %v", err1)
		}

		if err2 := checkPosAndSize(objects[2], 133.125, 404.5); err2 != nil {
			t.Errorf("object 2: %v", err2)
		}

		if err3 := checkPosAndSize(objects[3], 76.6875, 541.625); err3 != nil {
			t.Errorf("object 3: %v", err3)
		}

		if err4 := checkPosAndSize(objects[4], 107.030624, 622.3125); err4 != nil {
			t.Errorf("object 4: %v", err4)
		}

		if err5 := checkPosAndSize(objects[5], 170.65688, 725.34314); err5 != nil {
			t.Errorf("object 5: %v", err5)
		}

		actions := objects[5].(*fyne.Container)

		if err6 := checkPosAndSize(actions.Objects[0], 170.65688, 4); err6 != nil {
			t.Errorf("object 6: %v", err6)
		}

	})

	t.Run("valid number of objects without progress bar", func(t *testing.T) {
		l := &transferLayout{}

		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Layout expected no panic = %v", err)
			}
		}()

		pb := container.NewWithoutLayout()
		pb.Hide()

		objects := []fyne.CanvasObject{
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(),
			container.NewWithoutLayout(
				pb,
				container.NewWithoutLayout(),
				container.NewWithoutLayout(),
			),
		}

		l.Layout(objects, fyne.NewSize(900, 600))

		if err0 := checkPosAndSize(objects[0], 25, 4); err0 != nil {
			t.Errorf("object 0: %v", err0)
		}

		if err1 := checkPosAndSize(objects[1], 367.5, 33); err1 != nil {
			t.Errorf("object 1: %v", err1)
		}

		if err2 := checkPosAndSize(objects[2], 133.125, 404.5); err2 != nil {
			t.Errorf("object 2: %v", err2)
		}

		if err3 := checkPosAndSize(objects[3], 76.6875, 541.625); err3 != nil {
			t.Errorf("object 3: %v", err3)
		}

		if err4 := checkPosAndSize(objects[4], 107.030624, 622.3125); err4 != nil {
			t.Errorf("object 4: %v", err4)
		}

		if err5 := checkPosAndSize(objects[5], 170.65688, 725.34314); err5 != nil {
			t.Errorf("object 5: %v", err5)
		}

		actions := objects[5].(*fyne.Container)

		if err7 := checkPosAndSize(actions.Objects[1], 40, 4); err7 != nil {
			t.Errorf("object 7: %v", err7)
		}

		if err8 := checkPosAndSize(actions.Objects[2], 40, 48); err8 != nil {
			t.Errorf("object 8: %v", err8)
		}
	})
}

func checkPosAndSize(obj fyne.CanvasObject, width, posX float32) error {
	if obj.Size().Width != width {
		return fmt.Errorf("expected width = %v but got width = %v", width, obj.Size().Width)
	}

	if obj.Position().X != posX {
		return fmt.Errorf("expected position x = %v but got position x = %v", posX, obj.Position().X)
	}
	return nil
}
