package xdcodec

import "io"

var (
	NopReader = nopReader{}
	NopWriter = nopWriter{}
)

type nopReader struct{}

func (r nopReader) Read([]byte) (int, error) { return 0, nil }

type nopWriter struct{}

func (w nopWriter) Write([]byte) (int, error) { return 0, nil }

type WithReadWriter struct {
	io.Writer
	io.Reader
}

// ByteReadWriter implements ByteReader and ByteWriter.
type ByteReadWriter struct {
	p [1]byte
	io.ReadWriter
}

// NewByteReadWriter creates a new ByteReadWriter.
func NewByteReadWriter(rw io.ReadWriter) ByteReadWriter {
	return ByteReadWriter{ReadWriter: rw}
}

// ReadByte reads a byte.
func (brw ByteReadWriter) ReadByte() (b byte, err error) {
	p := brw.p[:]
	if _, err = brw.ReadWriter.Read(p); err != nil {
		return
	}
	b = p[0]
	return
}

// WriteByte writes a byte.
func (brw ByteReadWriter) WriteByte(b byte) error {
	_, err := brw.ReadWriter.Write([]byte{b})
	return err
}
