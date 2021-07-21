package transfer

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_fillMessageField(t *testing.T) {
	type args struct {
		content string
		buffer  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "empty field content",
			args: args{
				content: "",
				buffer:  make([]byte, 5),
			},
			want:    make([]byte, 5),
			wantErr: false,
		},
		{
			name: "number zero as string content",
			args: args{
				content: "0",
				buffer:  make([]byte, 5),
			},
			want:    []byte{48, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "field content longger than field len",
			args: args{
				content: "big field for the len allowed",
				buffer:  make([]byte, 5),
			},
			want:    make([]byte, 5),
			wantErr: true,
		},
		{
			name: "field content smaller than fiedl len",
			args: args{
				content: "sample",
				buffer:  make([]byte, 10),
			},
			want:    []byte{115, 97, 109, 112, 108, 101, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "field content with special chars and len equals to the number of chars",
			args: args{
				content: "Ts-?41ç$%",
				buffer:  make([]byte, 9),
			},
			want:    make([]byte, 9),
			wantErr: true,
		},
		{
			name: "field content with special chars",
			args: args{
				content: "Ts-?41ç$%",
				buffer:  make([]byte, 12),
			},
			want:    []byte{84, 115, 45, 63, 52, 49, 195, 167, 36, 37, 0, 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fillMessageField(tt.args.content, tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("fillMessageField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.buffer, tt.want) {
				t.Errorf("fillMessageField() = %v, want %v", tt.args.buffer, tt.want)
			}
		})
	}
}

func Test_trimMessageField(t *testing.T) {
	type args struct {
		field []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty field content",
			args: args{
				field: []byte{0, 0, 0, 0, 0},
			},
			want: "",
		},
		{
			name: "field content smaller than fiedl len",
			args: args{
				field: []byte{115, 97, 109, 112, 108, 101, 0, 0, 0, 0},
			},
			want: "sample",
		},
		{
			name: "field content with special chars",
			args: args{
				field: []byte{84, 115, 45, 63, 52, 49, 195, 167, 36, 37, 0, 0},
			},
			want: "Ts-?41ç$%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimMessageField(tt.args.field); got != tt.want {
				t.Errorf("trimMessageField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeRequestMessage(t *testing.T) {
	type args struct {
		m requestMessage
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "message completed and valid",
			args: args{
				m: requestMessage{
					Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
					FileName: "file-name.txt",
					Hostname: "my-hostname",
					FileSize: 2312321,
				},
			},
			wantW: string([]byte{
				57, 102, 56, 54, 100, 48, 56, 49, 56, 56, 52, 99, 55, 100, 54, 53, 57, 97, 50, 102,
				101, 97, 97, 48, 99, 53, 53, 97, 100, 48, 49, 53, 97, 51, 98, 102, 52, 102, 49, 98,
				50, 98, 48, 98, 56, 50, 50, 99, 100, 49, 53, 100, 54, 99, 49, 53, 98, 48, 102, 48,
				48, 97, 48, 56, 50, 51, 49, 50, 51, 50, 49, 0, 0, 0, 102, 105, 108, 101, 45, 110, 97,
				109, 101, 46, 116, 120, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 109, 121, 45, 104, 111, 115, 116, 110, 97, 109, 101, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			}),
			wantErr: false,
		}, {
			name: "message with invalid field len",
			args: args{
				m: requestMessage{
					Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
					FileName: "file-name.txt",
					Hostname: "my-hostname-that-is-tooooooo-long",
					FileSize: 2312321,
				},
			},
			wantErr: true,
		}, {
			name: "empty message",
			args: args{
				m: requestMessage{},
			},
			wantW: string([]byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := writeRequestMessage(tt.args.m, w); (err != nil) != tt.wantErr {
				t.Errorf("writeRequestMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeRequestMessage() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_readRequestMessage(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    requestMessage
		wantErr bool
	}{
		{
			name: "valid request message",
			args: args{
				r: bytes.NewBuffer([]byte{
					57, 102, 56, 54, 100, 48, 56, 49, 56, 56, 52, 99, 55, 100, 54, 53, 57, 97, 50, 102,
					101, 97, 97, 48, 99, 53, 53, 97, 100, 48, 49, 53, 97, 51, 98, 102, 52, 102, 49, 98,
					50, 98, 48, 98, 56, 50, 50, 99, 100, 49, 53, 100, 54, 99, 49, 53, 98, 48, 102, 48,
					48, 97, 48, 56, 50, 51, 49, 50, 51, 50, 49, 0, 0, 0, 102, 105, 108, 101, 45, 110, 97,
					109, 101, 46, 116, 120, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 109, 121, 45, 104, 111, 115, 116, 110, 97, 109, 101, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				}),
			},
			want: requestMessage{
				Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
				FileName: "file-name.txt",
				Hostname: "my-hostname",
				FileSize: 2312321,
			},
			wantErr: false,
		}, {
			name: "request message with invalid length",
			args: args{
				r: bytes.NewBuffer([]byte{
					57, 102, 56, 54, 100, 48, 56, 49, 56, 56, 52, 99, 55, 100, 54, 53, 57, 97, 50, 102,
					101, 97, 97, 48, 99, 53, 53, 97, 100, 48, 49, 53, 97, 51, 98, 102, 52, 102, 49, 98,
					50, 98, 48, 98, 56, 50, 50, 99, 100, 49, 53, 100, 54, 99, 49, 53, 98, 48, 102, 48,
					48, 97, 48, 56, 50, 51, 49, 50, 51, 50, 49, 0, 0, 0, 102, 105, 108, 101, 45, 110, 97,
					109, 101, 46, 116, 120, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 109, 121, 45, 104, 111, 115, 116, 110, 97, 109, 101, 0, 0,
					0, 0, 0, 0, 0, 0,
				}),
			},
			want:    requestMessage{},
			wantErr: true,
		}, {
			name: "invalid buffer",
			args: args{
				r: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m requestMessage
			if err := readRequestMessage(&m, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("writeRequestMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(m, tt.want) {
				t.Errorf("writeRequestMessage() = %v, want %v", m, tt.want)
			}
		})
	}
}

func Benchmark_writeRequestMessage(b *testing.B) {
	p := make([]byte, messageRequestLen)
	buffer := bytes.NewBuffer(p)
	m := requestMessage{
		FileName: "file-name.txt",
		FileSize: 231231321,
		Hostname: "my-hostname",
		Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	}
	for n := 0; n < b.N; n++ {
		writeRequestMessage(m, buffer)
	}
	//TODO check if we send the writer to the fill and write it directly is better
}

func Benchmark_readRequestMessage(b *testing.B) {
	m := requestMessage{
		FileName: "file-name.txt",
		FileSize: 231231321,
		Hostname: "my-hostname",
		Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	}
	buf := bytes.NewBuffer(make([]byte, messageRequestLen))
	writeRequestMessage(m, buf)
	for n := 0; n < b.N; n++ {
		var m requestMessage
		readRequestMessage(&m, buf)
	}
}

/*
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
Benchmark_writeRequestMessage-12    	 3176440	       358.2 ns/op	    1038 B/op	       5 allocs/op
Benchmark_writeRequestMessage-12    	 5904570	       336.2 ns/op	    1098 B/op	       5 allocs/op
Benchmark_writeRequestMessage-12    	 5854117	       476.4 ns/op	    1105 B/op	       5 allocs/op
Benchmark_writeRequestMessage-12    	 5284982	       210.0 ns/op	     726 B/op	       5 allocs/op
Benchmark_writeRequestMessage-12    	 5807980	       249.4 ns/op	    1112 B/op	       5 allocs/op
Benchmark_readRequestMessage-12     	 4979452	       241.8 ns/op	     288 B/op	       3 allocs/op
Benchmark_readRequestMessage-12     	 4942868	       238.8 ns/op	     288 B/op	       3 allocs/op
Benchmark_readRequestMessage-12     	 5021805	       238.6 ns/op	     288 B/op	       3 allocs/op
Benchmark_readRequestMessage-12     	 5008231	       239.0 ns/op	     288 B/op	       3 allocs/op
Benchmark_readRequestMessage-12     	 5013057	       238.8 ns/op	     288 B/op	       3 allocs/op

*/
