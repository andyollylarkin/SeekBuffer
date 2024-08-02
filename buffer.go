package seekbuffer

import (
	"errors"
	"fmt"
	"io"
)

type SeekBuffer struct {
	buff     []byte
	readPos  int64
	writePos int64
	closed   bool
}

func (b *SeekBuffer) grow(n int) {
	if cap(b.buff) > len(b.buff)+n {
		return
	}

	newBuff := make([]byte, int(b.writePos)+n, int(b.writePos)+n)

	copy(newBuff, b.buff)

	b.buff = newBuff
}

func (b *SeekBuffer) Bytes() []byte {
	return b.buff[:b.writePos]
}

func (b *SeekBuffer) Reset() {
	b.buff = nil
	b.writePos = 0
	b.readPos = 0
}

func (b *SeekBuffer) Len() int {
	if b.closed {
		return 0
	}

	return int(b.writePos)
}

func (b *SeekBuffer) Read(p []byte) (n int, err error) {
	if b.closed {
		return 0, fmt.Errorf("buffer is closed")
	}

	n = copy(p, b.buff[b.readPos:])

	b.readPos += int64(n)

	if b.readPos == b.writePos {
		return n, io.EOF
	}

	return n, nil
}

func (b *SeekBuffer) Close() error {
	b.buff = nil

	b.closed = true

	return nil
}

func (b *SeekBuffer) Seek(offset int64, whence int) (int64, error) {
	if b.closed {
		return 0, fmt.Errorf("buffer is closed")
	}

	newPos := b.readPos

	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = b.readPos + offset
	case io.SeekEnd:
		newPos = b.writePos - offset
	default:
		return 0, errors.New("invalid whence value")
	}

	if newPos < 0 || newPos > b.writePos {
		return 0, errors.New("seek position out of range")
	}

	b.readPos = newPos

	return b.readPos, nil
}

func (b *SeekBuffer) Write(p []byte) (n int, err error) {
	if b.closed {
		return 0, fmt.Errorf("buffer is closed")
	}

	if len(b.buff) == 0 {
		tmp := make([]byte, len(p), cap(p))
		copy(tmp, p)

		b.buff = tmp

		b.writePos = int64(len(p))

		return len(p), nil
	}

	b.grow(len(p))

	copy(b.buff[b.writePos:], p)

	b.writePos += int64(len(p))

	return len(p), nil
}
