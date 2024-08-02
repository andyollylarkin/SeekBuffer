package seekbuffer

import (
	"fmt"
	"io"
	"testing"
)

func TestSeekBuffer_grow(t *testing.T) {
	var bb SeekBuffer

	bb.grow(10)

	if cap(bb.buff) != 10 {
		t.Error("Invalid cap")
	}

	// with data

	bb.Write([]byte{'1', '2', '3'})

	bb.grow(10)

	if string(bb.Bytes()) != string([]byte{'1', '2', '3'}) {
		t.Error("buffer not copied properly")
	}
}

func TestSeekBuffer_Write(t *testing.T) {
	var bb SeekBuffer

	bytesToWrite := []byte("Hello")

	n, err := bb.Write(bytesToWrite)

	if err != nil {
		t.Error(fmt.Errorf("cant write to buffer: %w", err).Error())
	}

	if n != len(bytesToWrite) {
		t.Errorf("written less len len of content: n=%d, len_content=%d", n, len(bytesToWrite))
	}

	if bb.writePos != int64(len(bytesToWrite)) {
		t.Errorf("Invalid write pos. Should be %d, actual %d ", len(bytesToWrite), bb.writePos)
	}

	bb.Write(bytesToWrite)

	if string(bb.buff) != string(append(bytesToWrite, bytesToWrite...)) {
		t.Errorf("invalid buffer content")
	}

	if bb.Len() != len(bytesToWrite)*2 {
		t.Errorf("Invalid content len. Should be %d. Actual %d", len(bytesToWrite)*2, bb.Len())
	}
}

func TestSeekBuffer_Read(t *testing.T) {
	var bb SeekBuffer

	bb.Write([]byte("Hello"))
	bb.Write([]byte(" World!"))

	content, err := io.ReadAll(&bb)
	if err != nil {
		t.Errorf("Error read: %s", err.Error())
	}

	if string(content) != "Hello World!" {
		t.Errorf("Error read content. Readed content not equal to original content")
	}
}

func TestSeekBuffer_Seek(t *testing.T) {
	var bb SeekBuffer

	bb.Write([]byte("Hello World"))

	_, err := bb.Seek(3, io.SeekCurrent)
	if err != nil {
		t.Errorf("Seek error: %s", err.Error())
	}

	if bb.readPos != 3 {
		t.Errorf("Seek error: position should be 3, actual - %d", bb.readPos)
	}

	_, err = bb.Seek(3, io.SeekEnd)
	if err != nil {
		t.Errorf("Seek error: %s", err.Error())
	}

	if bb.readPos != 8 {
		t.Errorf("Seek error: position should be 8, actual - %d", bb.readPos)
	}

	_, err = bb.Seek(1, io.SeekCurrent)
	if err != nil {
		t.Errorf("Seek error: %s", err.Error())
	}

	if bb.readPos != 9 {
		t.Errorf("Seek error: position should be 9, actual - %d", bb.readPos)
	}

	_, err = bb.Seek(10, io.SeekCurrent)
	if err == nil {
		t.Errorf("Seek should cause error, becase its out of range. Error is nil")
	}

	bb.Seek(1, io.SeekStart)

	content, err := io.ReadAll(&bb)
	if err != nil {
		t.Errorf("Error read seeked buffer: %s", err.Error())
	}

	if string(content) != "ello World" {
		t.Errorf("Error read seeked buffer. Invalid readed content")
	}
}
