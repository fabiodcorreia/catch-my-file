package transfer

import (
	"fmt"
	"testing"
)

func Test_TransferStore_Get(t *testing.T) {
	t.Run("get with empty store", func(t *testing.T) {
		s := NewStore()
		tr := s.Get(0)

		if tr != nil {
			t.Errorf("get expected nil but got %v", tr)
		}
	})

	t.Run("get with invalid index", func(t *testing.T) {
		s := NewStore()
		tr := s.Get(-1)

		if tr != nil {
			t.Errorf("get expected nil but got %v", tr)
		}
	})

	t.Run("get with valid index", func(t *testing.T) {
		s := NewStore()
		i := s.Add(&Transfer{})
		tr := s.Get(i)

		if tr == nil {
			t.Errorf("get expected non nil but got %v", tr)
		}
	})
}

func Test_TransferStore_Add(t *testing.T) {
	t.Run("add nil transfer", func(t *testing.T) {
		s := NewStore()
		i := s.Add(nil)

		if i != -1 {
			t.Errorf("add expected nil but got %v", i)
		}
	})

	t.Run("add valid transfer", func(t *testing.T) {
		s := NewStore()
		i := s.Add(&Transfer{})

		if i != 0 {
			t.Errorf("add expected 0 but got %v", i)
		}
	})

	t.Run("add valid transfer with callback", func(t *testing.T) {
		s := NewStore()
		var changeEvent string
		s.OnStoreChange = func(i int) {
			changeEvent = fmt.Sprintf("Index Changed = %d", i)
		}
		_ = s.Add(&Transfer{})

		if changeEvent != "Index Changed = 0" {
			t.Errorf("add expected Index Changed = 0 but got %v", changeEvent)
		}
	})
}

func Test_TransferStore_Update(t *testing.T) {
	t.Run("update with empty store", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("update not expected panic = %v", err)
			}
		}()

		s := NewStore()
		s.Update(0, &Transfer{})
	})

	t.Run("update with invalid index", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("update not expected panic = %v", err)
			}
		}()

		s := NewStore()
		s.Update(-1, &Transfer{})
	})

	t.Run("update to nil transfer", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("update not expected panic = %v", err)
			}
		}()

		s := NewStore()
		i := s.Add(&Transfer{})
		s.Update(i, nil)
	})

	t.Run("update transfer to error state", func(t *testing.T) {
		s := NewStore()
		i := s.Add(&Transfer{})
		tr1 := s.Get(i)
		tr1.SetError(fmt.Errorf("my error"))
		s.Update(i, tr1)
		tr2 := s.Get(i)

		if tr2.Status != Error {
			t.Errorf("update expected status = Waiting but got = %v ", tr2.Status.String())
		}

		if tr2.Error().Error() != "my error" {
			t.Errorf("update expected error = my error but got = %v ", tr2.Error().Error())
		}
	})

	t.Run("update transfer to error state with callback ", func(t *testing.T) {
		s := NewStore()
		var changeEvent string
		s.OnStoreChange = func(i int) {
			changeEvent = fmt.Sprintf("Index Changed = %d", i)
		}
		i := s.Add(&Transfer{})
		tr1 := s.Get(i)
		tr1.SetError(fmt.Errorf("my error"))
		s.Update(i, tr1)

		if changeEvent != "Index Changed = 0" {
			t.Errorf("update expected Index Changed = 0 but got %v", changeEvent)
		}
	})

}

func Test_TransferStore_AddToWait(t *testing.T) {
	t.Run("add and wait nil transfer", func(t *testing.T) {
		s := NewStore()
		_, wait := s.AddToWait(nil)

		if wait != nil {
			t.Errorf("add to wait expected wait nil but got = %v", wait)
		}
	})

	t.Run("add to wait transfer to accept state and unlock", func(t *testing.T) {
		s := NewStore()
		i, wait := s.AddToWait(&Transfer{})

		go func() {
			tr := s.Get(i)
			tr.Status = Accepted
			s.Update(i, tr)
		}()

		<-wait

		tt := s.Get(i)
		if tt.Status != Accepted {
			t.Errorf("add to wait expected status Accepted but got = %v", tt.Status.String())
		}
	})
}

func Test_TransferStore_FollowProgress(t *testing.T) {
	t.Run("follow progress of invalid index", func(t *testing.T) {
		s := NewStore()
		progress := s.FollowProgress(-1)

		if progress != nil {
			t.Errorf("follow progress expected nil but got = %v", progress)
		}
	})

	t.Run("follow progress with empty store", func(t *testing.T) {
		s := NewStore()
		progress := s.FollowProgress(0)

		if progress != nil {
			t.Errorf("follow progress expected nil but got = %v", progress)
		}
	})

	t.Run("follow progress change to complete", func(t *testing.T) {
		s := NewStore()
		i := s.Add(&Transfer{})

		progress := s.FollowProgress(i)

		go func() {
			s.UpdateProgress(i, 10)
			tt := s.Get(i)
			tt.Status = Completed
			s.Update(i, tt)
		}()

		val := <-progress

		tr := s.Get(i)
		if tr.Status != Completed {
			t.Errorf("follow progress expected status Completed but got = %v", tr.Status.String())
		}

		if val != 10 {
			t.Errorf("ollow progress expected progress 10 but got = %v", val)
		}
	})
}

func Test_TransferStore_UpdateProgress(t *testing.T) {

}
