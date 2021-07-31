package protocol

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_fillMessageField(t *testing.T) {

	t.Run("zero value content", func(t *testing.T) {
		input := ""
		output := make([]byte, 5)
		want := make([]byte, 5)

		if err := fillMessageField(input, output); err != nil {
			t.Errorf("fillMessageField not expected error = %v", err)
		}

		if !reflect.DeepEqual(output, want) {
			t.Errorf("fillMessageField expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("zero string content", func(t *testing.T) {
		input := "0"
		output := make([]byte, 5)
		want := []byte{48, 0, 0, 0, 0}

		if err := fillMessageField(input, output); err != nil {
			t.Errorf("fillMessageField not expected error = %v", err)
		}

		if !reflect.DeepEqual(output, want) {
			t.Errorf("fillMessageField expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("content length larger than field length", func(t *testing.T) {
		input := "000000"
		output := make([]byte, 5)

		if err := fillMessageField(input, output); err == nil {
			t.Errorf("fillMessageField expected error = %v", err)
		}
	})

	t.Run("content length smaller than field length", func(t *testing.T) {
		input := "000"
		output := make([]byte, 5)
		want := []byte{48, 48, 48, 0, 0}

		if err := fillMessageField(input, output); err != nil {
			t.Errorf("fillMessageField not expected error = %v", err)
		}

		if !reflect.DeepEqual(output, want) {
			t.Errorf("fillMessageField expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("content with special chars", func(t *testing.T) {
		intput := "Ts-?41ç$%"
		output := make([]byte, 12)
		want := []byte{84, 115, 45, 63, 52, 49, 195, 167, 36, 37, 0, 0}

		if err := fillMessageField(intput, output); err != nil {
			t.Errorf("fillMessageField not expected error = %v", err)
		}

		if !reflect.DeepEqual(output, want) {
			t.Errorf("fillMessageField expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("content with 3 special chars and field length of 3", func(t *testing.T) {
		intput := "ç±€"
		output := make([]byte, 3)

		if err := fillMessageField(intput, output); err == nil {
			t.Errorf("fillMessageField expected error = %v", err)
		}
	})

}

func Test_trimMessageField(t *testing.T) {
	t.Run("zero value content", func(t *testing.T) {
		input := []byte{0, 0, 0, 0, 0}
		want := ""

		output := trimMessageField(input)
		if !reflect.DeepEqual(output, want) {
			t.Errorf("trimMessageField expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("zero string content", func(t *testing.T) {
		input := []byte{48, 0, 0, 0, 0}
		want := "0"

		output := trimMessageField(input)
		if !reflect.DeepEqual(output, want) {
			t.Errorf("trimMessageField expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("content length smaller then field length", func(t *testing.T) {
		input := []byte{115, 97, 109, 112, 108, 101, 0, 0, 0, 0}
		want := "sample"

		output := trimMessageField(input)

		if !reflect.DeepEqual(output, want) {
			t.Errorf("trimMessageField expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("content with special chars", func(t *testing.T) {
		input := []byte{84, 115, 45, 63, 52, 49, 195, 167, 36, 37, 0, 0}
		want := "Ts-?41ç$%"

		output := trimMessageField(input)

		if !reflect.DeepEqual(output, want) {
			t.Errorf("trimMessageField expected output = %v but got output = %v", want, output)
		}
	})
}

func Test_WriteRequestMessage(t *testing.T) {
	t.Run("message completed and valid", func(t *testing.T) {
		input := RequestMessage{
			Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			FileName: "file-name.txt",
			Hostname: "my-hostname",
			FileSize: 99999,
		}
		output := &bytes.Buffer{}
		want := []byte{
			57, 102, 56, 54, 100, 48, 56, 49, 56, 56, 52, 99, 55, 100, 54, 53, 57, 97, 50, 102,
			101, 97, 97, 48, 99, 53, 53, 97, 100, 48, 49, 53, 97, 51, 98, 102, 52, 102, 49, 98,
			50, 98, 48, 98, 56, 50, 50, 99, 100, 49, 53, 100, 54, 99, 49, 53, 98, 48, 102, 48,
			48, 97, 48, 56, 57, 57, 57, 57, 57, 0, 0, 0, 0, 0, 0, 0, 0, 102, 105, 108, 101, 45,
			110, 97, 109, 101, 46, 116, 120, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 109, 121, 45, 104, 111, 115, 116, 110, 97, 109, 101,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}

		if err := WriteRequestMessage(input, output); err != nil {
			t.Errorf("WriteRequestMessage not expected error = %v", err)
		}
		if !reflect.DeepEqual(output.Bytes(), want) {
			t.Errorf("WriteRequestMessage expected output = %v but got output = %v", want, output.Bytes())
		}
	})

	t.Run("message with a file size bigger than 10GB", func(t *testing.T) {
		input := RequestMessage{
			Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			FileName: "file-name.txt",
			Hostname: "my-hostname",
			FileSize: 99999999999,
		}
		output := &bytes.Buffer{}
		want := []byte{
			57, 102, 56, 54, 100, 48, 56, 49, 56, 56, 52, 99, 55, 100, 54, 53, 57, 97, 50, 102,
			101, 97, 97, 48, 99, 53, 53, 97, 100, 48, 49, 53, 97, 51, 98, 102, 52, 102, 49, 98,
			50, 98, 48, 98, 56, 50, 50, 99, 100, 49, 53, 100, 54, 99, 49, 53, 98, 48, 102, 48,
			48, 97, 48, 56, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 0, 0, 102, 105, 108, 101,
			45, 110, 97, 109, 101, 46, 116, 120, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 109, 121, 45, 104, 111, 115, 116, 110, 97, 109,
			101, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}

		if err := WriteRequestMessage(input, output); err != nil {
			t.Errorf("WriteRequestMessage not expected error = %v", err)
		}
		if !reflect.DeepEqual(output.Bytes(), want) {
			t.Errorf("WriteRequestMessage expected output = %v but got output = %v", want, output.Bytes())
		}
	})

	t.Run("message with checksum field too long", func(t *testing.T) {
		input := RequestMessage{
			Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08123123123123123312312",
			FileName: "file-name.txt",
			Hostname: "my-hostname",
			FileSize: 99999999999,
		}
		output := &bytes.Buffer{}

		if err := WriteRequestMessage(input, output); err == nil {
			t.Errorf("WriteRequestMessage expected error = %v", err)
		}
	})

	t.Run("message with file size field too long", func(t *testing.T) {
		input := RequestMessage{
			Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			FileName: "file-name.txt",
			Hostname: "my-hostname",
			FileSize: 99999999999000,
		}
		output := &bytes.Buffer{}

		if err := WriteRequestMessage(input, output); err == nil {
			t.Errorf("WriteRequestMessage expected error = %v", err)
		}
	})

	t.Run("message with filename field too long", func(t *testing.T) {
		input := RequestMessage{
			Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			FileName: "file-naaaaaaaaaaammmmmmmmeeeeeeeeme-sooooooooooooooooooooo------looooooooooooooooooooonnnnnnnnnnnnggggggggggggggggggggggggggg.txt",
			Hostname: "my-hostname",
			FileSize: 99999999999,
		}
		output := &bytes.Buffer{}

		if err := WriteRequestMessage(input, output); err == nil {
			t.Errorf("WriteRequestMessage expected error = %v", err)
		}
	})

	t.Run("message with hosname field too long", func(t *testing.T) {
		input := RequestMessage{
			Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			FileName: "file-name.txt",
			Hostname: "my-hostname-that-is-tooooooo-long",
			FileSize: 99999999999,
		}
		output := &bytes.Buffer{}

		if err := WriteRequestMessage(input, output); err == nil {
			t.Errorf("WriteRequestMessage expected error = %v", err)
		}
	})

	t.Run("message empty", func(t *testing.T) {
		input := RequestMessage{}
		output := &bytes.Buffer{}
		want := []byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}

		if err := WriteRequestMessage(input, output); err != nil {
			t.Errorf("WriteRequestMessage not expected error = %v", err)
		}
		if !reflect.DeepEqual(output.Bytes(), want) {
			t.Errorf("WriteRequestMessage expected output = %v but got output = %v", want, output.Bytes())
		}
	})

	t.Run("output writer is nil", func(t *testing.T) {
		input := RequestMessage{}
		if err := WriteRequestMessage(input, nil); err == nil {
			t.Errorf("WriteRequestMessage expected error = %v", err)
		}
	})

}

func Test_ReadRequestMessage(t *testing.T) {
	t.Run("message completed and valid", func(t *testing.T) {
		input := bytes.NewBuffer([]byte{
			57, 102, 56, 54, 100, 48, 56, 49, 56, 56, 52, 99, 55, 100, 54, 53, 57, 97, 50, 102,
			101, 97, 97, 48, 99, 53, 53, 97, 100, 48, 49, 53, 97, 51, 98, 102, 52, 102, 49, 98,
			50, 98, 48, 98, 56, 50, 50, 99, 100, 49, 53, 100, 54, 99, 49, 53, 98, 48, 102, 48,
			48, 97, 48, 56, 50, 51, 49, 50, 51, 50, 49, 0, 0, 0, 0, 0, 0, 102, 105, 108, 101, 45,
			110, 97, 109, 101, 46, 116, 120, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 109, 121, 45, 104, 111, 115, 116, 110, 97, 109, 101,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		})
		var output RequestMessage
		want := RequestMessage{
			Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			FileName: "file-name.txt",
			Hostname: "my-hostname",
			FileSize: 2312321,
		}

		if err := ReadRequestMessage(&output, input); err != nil {
			t.Errorf("ReadRequestMessage not expected error = %v", err)
		}
		if !reflect.DeepEqual(output, want) {
			t.Errorf("ReadRequestMessage expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("message with invalid length", func(t *testing.T) {
		input := bytes.NewBuffer([]byte{
			57, 102, 56, 54, 100, 48, 56, 49, 56, 56, 52, 99, 55, 100, 54, 53, 57, 97, 50, 102,
		})
		var output RequestMessage

		if err := ReadRequestMessage(&output, input); err == nil {
			t.Errorf("ReadRequestMessage expected error = %v", err)
		}
	})

	t.Run("input reader is nil", func(t *testing.T) {
		var output RequestMessage

		if err := ReadRequestMessage(&output, nil); err == nil {
			t.Errorf("ReadRequestMessage expected error = %v", err)
		}
	})

	t.Run("request message pointer is nil", func(t *testing.T) {
		input := bytes.NewBuffer([]byte{
			57, 102, 56, 54, 100, 48, 56, 49, 56, 56, 52, 99, 55, 100, 54, 53, 57, 97, 50, 102,
			101, 97, 97, 48, 99, 53, 53, 97, 100, 48, 49, 53, 97, 51, 98, 102, 52, 102, 49, 98,
			50, 98, 48, 98, 56, 50, 50, 99, 100, 49, 53, 100, 54, 99, 49, 53, 98, 48, 102, 48,
			48, 97, 48, 56, 50, 51, 49, 50, 51, 50, 49, 0, 0, 0, 0, 0, 0, 102, 105, 108, 101, 45,
			110, 97, 109, 101, 46, 116, 120, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 109, 121, 45, 104, 111, 115, 116, 110, 97, 109, 101,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		})

		if err := ReadRequestMessage(nil, input); err == nil {
			t.Errorf("ReadRequestMessage expected error = %v", err)
		}
	})

	t.Run("request message with invalid file size", func(t *testing.T) {
		input := bytes.NewBuffer([]byte{
			57, 102, 56, 54, 100, 48, 56, 49, 56, 56, 52, 99, 55, 100, 54, 53, 57, 97, 50, 102,
			101, 97, 97, 48, 99, 53, 53, 97, 100, 48, 49, 53, 97, 51, 98, 102, 52, 102, 49, 98,
			50, 98, 48, 98, 56, 102, 102, 99, 100, 49, 53, 100, 54, 99, 49, 53, 98, 48, 102, 48,
			48, 97, 48, 56, 50, 102, 49, 50, 51, 50, 49, 0, 0, 0, 0, 0, 0, 102, 105, 108, 101, 45,
			110, 97, 109, 101, 46, 116, 120, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 109, 121, 45, 104, 111, 115, 116, 110, 97, 109, 101,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		})
		var output RequestMessage

		if err := ReadRequestMessage(&output, input); err == nil {
			t.Errorf("ReadRequestMessage expected error = %v", err)
		}
	})
}

func Test_WriteDecision(t *testing.T) {
	t.Run("send accept decision", func(t *testing.T) {
		output := &bytes.Buffer{}
		want := []byte{1}

		if err := WriteDecision(true, output); err != nil {
			t.Errorf("WriteDecision not expected error = %v", err)
		}
		if !reflect.DeepEqual(output.Bytes(), want) {
			t.Errorf("WriteDecision expected output = %v but got output = %v", want, output.Bytes())
		}
	})

	t.Run("send reject decision", func(t *testing.T) {
		output := &bytes.Buffer{}
		want := []byte{0}

		if err := WriteDecision(false, output); err != nil {
			t.Errorf("WriteDecision not expected error = %v", err)
		}
		if !reflect.DeepEqual(output.Bytes(), want) {
			t.Errorf("WriteDecision expected output = %v but got output = %v", want, output.Bytes())
		}
	})

	t.Run("output writer is nil", func(t *testing.T) {
		if err := WriteDecision(true, nil); err == nil {
			t.Errorf("WriteDecision expected error = %v", err)
		}
	})
}

func Test_ReadDecision(t *testing.T) {
	t.Run("send accept decision", func(t *testing.T) {
		input := bytes.NewBuffer([]byte{1})
		want := true

		output, err := ReadDecision(input)
		if err != nil {
			t.Errorf("ReadDecision not expected error = %v", err)
		}
		if output != want {
			t.Errorf("ReadDecision expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("send reject decision", func(t *testing.T) {
		input := bytes.NewBuffer([]byte{0})
		want := false

		output, err := ReadDecision(input)
		if err != nil {
			t.Errorf("ReadDecision not expected error = %v", err)
		}
		if output != want {
			t.Errorf("ReadDecision expected output = %v but got output = %v", want, output)
		}
	})

	t.Run("intput reader is nil", func(t *testing.T) {
		if _, err := ReadDecision(nil); err == nil {
			t.Errorf("ReadDecision expected error = %v", err)
		}
	})
}

func Benchmark_WriteRequestMessage(b *testing.B) {
	p := make([]byte, messageRequestLen)
	buffer := bytes.NewBuffer(p)
	m := RequestMessage{
		FileName: "file-name.txt",
		FileSize: 231231321,
		Hostname: "my-hostname",
		Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	}
	for n := 0; n < b.N; n++ {
		WriteRequestMessage(m, buffer)
	}
	//TODO check if we send the writer to the fill and write it directly is better
}

func Benchmark_ReadRequestMessage(b *testing.B) {
	m1 := RequestMessage{
		FileName: "file-name.txt",
		FileSize: 231231321,
		Hostname: "my-hostname",
		Checksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	}
	buf := bytes.NewBuffer(make([]byte, 0, messageRequestLen))
	WriteRequestMessage(m1, buf)
	for n := 0; n < b.N; n++ {
		var m2 RequestMessage
		ReadRequestMessage(&m2, buf)
	}
}

/*
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
Benchmark_writeRequestMessage-12    	 3859276	       307.6 ns/op	     771 B/op	       2 allocs/op
Benchmark_writeRequestMessage-12    	 7932200	       272.3 ns/op	     757 B/op	       2 allocs/op
Benchmark_writeRequestMessage-12    	 7391088	       295.3 ns/op	     793 B/op	       2 allocs/op
Benchmark_writeRequestMessage-12    	 7269085	       160.0 ns/op	     803 B/op	       2 allocs/op
Benchmark_writeRequestMessage-12    	 7672149	       161.8 ns/op	     774 B/op	       2 allocs/op
Benchmark_readRequestMessage-12     	 4604421	       262.8 ns/op	     320 B/op	       3 allocs/op
Benchmark_readRequestMessage-12     	 4550571	       258.4 ns/op	     320 B/op	       3 allocs/op
Benchmark_readRequestMessage-12     	 4647685	       257.6 ns/op	     320 B/op	       3 allocs/op
Benchmark_readRequestMessage-12     	 4630645	       257.9 ns/op	     320 B/op	       3 allocs/op
Benchmark_readRequestMessage-12     	 4638036	       258.0 ns/op	     320 B/op	       3 allocs/op
*/
