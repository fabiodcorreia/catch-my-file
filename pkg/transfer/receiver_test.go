package transfer

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/fabiodcorreia/catch-my-file/pkg/network"
	"github.com/fabiodcorreia/catch-my-file/pkg/transfer/internal/protocol"
)

func Test_verifyTransfer(t *testing.T) {
	t.Run("transfer and file valid and matching", func(t *testing.T) {
		input := bytes.NewBufferString("File Super Secret Content :)")
		ctx := context.Background()
		tr := NewTransfer("File.txt", "92343059e81e2c7b0a589c7f2a7583cec6024a48fcdad981080dccbaa8ec61c1", "", 1000, nil, Upload)
		err := verifyTransfer(ctx, tr, 1000, input)

		if err != nil {
			t.Errorf("verify transfer not expected error but got %v", err)
		}
	})

	t.Run("size don't match", func(t *testing.T) {
		input := bytes.NewBufferString("File Super Secret Content :)")
		ctx := context.Background()
		tr := NewTransfer("File.txt", "92343059e81e2c7b0a589c7f2a7583cec6024a48fcdad981080dccbaa8ec61c1", "", 1000, nil, Upload)
		err := verifyTransfer(ctx, tr, 400, input)

		if err == nil {
			t.Errorf("verify transfer expected error but got %v", err)
		}
	})

	t.Run("size don't match", func(t *testing.T) {
		input := bytes.NewBufferString("File Super Secret Content :)")
		ctx := context.Background()
		tr := NewTransfer("File.txt", "92343059e81e2c7b0a589c7f2a7583cec6024a48fcdad981080dccbaa8ec61c1", "", 1000, nil, Upload)
		err := verifyTransfer(ctx, tr, 400, input)

		if err == nil {
			t.Errorf("verify transfer expected error but got %v", err)
		}
	})

	t.Run("checksum don't match", func(t *testing.T) {
		input := bytes.NewBufferString("File Super Secret Content :)")
		ctx := context.Background()
		tr := NewTransfer("File.txt", "not correct checksum", "", 1000, nil, Upload)
		err := verifyTransfer(ctx, tr, 1000, input)

		if err == nil {
			t.Errorf("verify transfer expected error but got %v", err)
		}
	})

	t.Run("transfer is nil", func(t *testing.T) {
		input := bytes.NewBufferString("File Super Secret Content :)")
		ctx := context.Background()
		err := verifyTransfer(ctx, nil, 400, input)

		if err == nil {
			t.Errorf("verify transfer expected error but got %v", err)
		}
	})

	t.Run("input reader is nil", func(t *testing.T) {
		ctx := context.Background()
		tr := NewTransfer("File.txt", "not correct checksum", "", 1000, nil, Upload)
		err := verifyTransfer(ctx, tr, 400, nil)

		if err == nil {
			t.Errorf("verify transfer expected error but got %v", err)
		}
	})

	t.Run("cancel context", func(t *testing.T) {
		input := bytes.NewBufferString("File Super Secret Content :)")
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		tr := NewTransfer("File.txt", "92343059e81e2c7b0a589c7f2a7583cec6024a48fcdad981080dccbaa8ec61c1", "", 1000, nil, Upload)
		err := verifyTransfer(ctx, tr, 1000, input)

		if err == nil {
			t.Errorf("verify transfer expected error but got %v", err)
		}
	})
}

func Test_watchdog(t *testing.T) {
	t.Run("cancel context", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		listener, _ := net.Listen(network.Type, "localhost:9933")
		done := make(chan interface{})

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()
		go func() {
			watchdog(ctx, listener)
			close(done)
		}()

		timeout := time.After(100 * time.Millisecond)

		select {
		case <-timeout:
			t.Errorf("watchdog not interrupted")
		case <-done:
		}
	})

	t.Run("not cancel context", func(t *testing.T) {
		ctx := context.Background()
		listener, _ := net.Listen(network.Type, "localhost:9933")
		done := make(chan interface{})

		go func() {
			watchdog(ctx, listener)
		}()

		timeout := time.After(100 * time.Millisecond)

		select {
		case <-timeout:
		case <-done:
			t.Errorf("watchdog interrupted")
		}
	})
}

func Test_reqDecisionAndWait(t *testing.T) {
	rm := protocol.RequestMessage{
		FileName: "file-1",
		FileSize: 1000,
		Hostname: "peer-1",
		Checksum: "check123",
	}

	t.Run("send request message and wait for decision accepted", func(t *testing.T) {
		ctx := context.Background()
		inOut := bytes.NewBuffer(make([]byte, 0))
		store := NewStore()

		protocol.WriteRequestMessage(rm, inOut)

		go func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Store=%p and MU=%p\n", store, &store.mu) //For some reason the race detector stops after this log
			tt := store.Get(0)
			tt.Status = Accepted
			store.Update(0, tt)
		}()

		_, err := reqDecisionAndWait(ctx, store, inOut, nil)

		if err != nil {
			t.Errorf("request decision and wait not excepted error but got = %v", err)
		}
	})

	t.Run("send request message and wait for decision rejected", func(t *testing.T) {
		ctx := context.Background()
		inOut := bytes.NewBuffer(make([]byte, 0))
		store := NewStore()

		protocol.WriteRequestMessage(rm, inOut)

		go func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Store=%p and MU=%p\n", store, &store.mu) //For some reason the race detector stops after this log
			tt := store.Get(0)
			tt.Status = Rejected
			store.Update(0, tt)
		}()

		_, err := reqDecisionAndWait(ctx, store, inOut, nil)

		if err != nil {
			t.Errorf("request decision and wait not excepted error but got = %v", err)
		}
	})

	t.Run("send request message and wait unitl cancelation", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		inOut := bytes.NewBuffer(make([]byte, 0))
		store := NewStore()

		protocol.WriteRequestMessage(rm, inOut)

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		_, err := reqDecisionAndWait(ctx, store, inOut, nil)

		if err == nil {
			t.Errorf("request decision and wait not excepted error but got = %v", err)
		}
	})

	t.Run("store is nil", func(t *testing.T) {
		ctx := context.Background()
		_, err := reqDecisionAndWait(ctx, nil, nil, nil)

		if err == nil {
			t.Errorf("request decision and wait excepted error but got = %v", err)
		}
	})
}
