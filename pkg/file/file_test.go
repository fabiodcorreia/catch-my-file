package file

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"testing"
)

func Test_Stream(t *testing.T) {
	t.Run("input reader is nil", func(t *testing.T) {
		_, err := Stream(context.Background(), nil, nil, nil)

		if err == nil {
			t.Errorf("Stream expected error = %v", err)
		}
	})

	t.Run("output writer is nil", func(t *testing.T) {
		input := bytes.NewBufferString("stream great content! 88392931 :)")

		_, err := Stream(context.Background(), input, nil, nil)

		if err == nil {
			t.Errorf("Stream expected error = %v", err)
		}
	})

	t.Run("output writer is nil", func(t *testing.T) {
		content := "stream great content! 88392931 :)"
		input := bytes.NewBufferString(content)
		output := &bytes.Buffer{}

		count, err := Stream(context.Background(), input, output, nil)

		if err != nil {
			t.Errorf("Stream not expected error = %v", err)
		}

		if len(content) != count {
			t.Errorf("Stream expected count = %v but got count = %v", len(content), count)
		}

		if !reflect.DeepEqual(output.String(), content) {
			t.Errorf("Stream expected output = %v but got output = %v", content, output.String())
		}
	})

	t.Run("output writer is nil", func(t *testing.T) {
		content := "stream great content! 88392931 :)"
		input := bytes.NewBufferString(content)
		output := &bytes.Buffer{}
		var progress int

		count, err := Stream(context.Background(), input, output, func(transferred int) {
			progress = transferred
		})

		if err != nil {
			t.Errorf("Stream not expected error = %v", err)
		}

		if len(content) != count {
			t.Errorf("Stream expected count = %v but got count = %v", len(content), count)
		}

		if count != progress {
			t.Errorf("Stream expected progress = %v but got progress = %v", count, progress)
		}
	})
}

func Test_Lookup(t *testing.T) {
	t.Run("empty file path", func(t *testing.T) {
		_, _, err := Lookup("")

		if err == nil {
			t.Errorf("Lookup expected error = %v", err)
		}
	})

	t.Run("file path doesn't exists", func(t *testing.T) {
		input := "/sample-file.txt"
		_, _, err := Lookup(input)

		if err == nil {
			t.Errorf("Lookup expected error = %v", err)
		}
	})

	t.Run("file path is a dir", func(t *testing.T) {
		input := os.TempDir()
		_, _, err := Lookup(input)

		if err == nil {
			t.Errorf("Lookup expected error = %v", err)
		}
	})

	t.Run("file path valid", func(t *testing.T) {
		input, _ := os.Executable()
		wantName := "file.test"
		outName, outSize, err := Lookup(input)

		if err != nil {
			t.Errorf("Lookup not expected error = %v", err)
		}
		if wantName != outName {
			t.Errorf("Lookup expected output name = %v but got output size = %v", wantName, outName)
		}
		if outSize < 1000 {
			t.Errorf("Lookup expected output size > 1000 but got output size = %v", outSize)
		}
	})
}

func Test_Checksum(t *testing.T) {
	t.Run("valid input reader", func(t *testing.T) {
		input := bytes.NewBufferString("sample text to hash")
		want := "b1668ccc2110b0ce6d103144e5eabc6b0cc59ec84cc58be543536c62f8d6fc00"

		output, err := Checksum(context.Background(), input)

		if err != nil {
			t.Errorf("Checksum not expected error = %v", err)
		}
		if !reflect.DeepEqual(output, want) {
			t.Errorf("Checksum expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("input reader is nil", func(t *testing.T) {
		_, err := Checksum(context.Background(), nil)

		if err == nil {
			t.Errorf("Checksum not expected error = %v", err)
		}
	})

	t.Run("valid input reader", func(t *testing.T) {
		input := bytes.NewBufferString("sample text to hash")

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := Checksum(ctx, input)

		if err == nil {
			t.Errorf("Checksum not expected error = %v", err)
		}
	})
}

func Benchmark_Checksum(b *testing.B) {
	ctx := context.Background()
	input := bytes.NewBufferString("sample text to hash")
	for n := 0; n < b.N; n++ {
		Checksum(ctx, input)
	}
}

/*
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
Benchmark_Checksum-12    	  220658	      5349 ns/op	   65817 B/op	       5 allocs/op
Benchmark_Checksum-12    	  224128	      5562 ns/op	   65817 B/op	       5 allocs/op
Benchmark_Checksum-12    	  227800	      5338 ns/op	   65817 B/op	       5 allocs/op
Benchmark_Checksum-12    	  227668	      5195 ns/op	   65817 B/op	       5 allocs/op
Benchmark_Checksum-12    	  229113	      5178 ns/op	   65817 B/op	       5 allocs/op
*/
