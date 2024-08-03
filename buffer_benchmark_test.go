package seekbuffer

import (
	"io"
	"testing"
)

func Benchmark_grow(b *testing.B) {
	var bb SeekBuffer

	for i := 0; i < b.N; i++ {
		bb.grow(i * 100)
	}
}

func Benchmark_Write(b *testing.B) {
	var bb SeekBuffer

	for i := 0; i < b.N; i++ {
		bb.Write([]byte("Hello"))
	}
}

func Benchmark_Read(b *testing.B) {
	var bb SeekBuffer
	bb.Write([]byte("Hello"))

	for i := 0; i < b.N; i++ {
		buf := make([]byte, 5, 5)

		bb.Read(buf)
		bb.Seek(0, io.SeekStart)
	}
}
