package xdcodec

import (
	"encoding/binary"
	"io"
)

type Codec struct {
	buf       [binary.MaxVarintLen64]byte
	keybuf    []string
	byteOrder binary.ByteOrder

	rw  ByteReadWriter
	err error
}

func New(byteOrder binary.ByteOrder, rw io.ReadWriter) *Codec {
	return &Codec{
		byteOrder: byteOrder,
		rw:        NewByteReadWriter(rw),
	}
}

func (c *Codec) RawReadWriter() io.ReadWriter {
	return c.rw
}

func (c *Codec) Reset(rw io.ReadWriter) {
	for i := range c.buf { // Optimized by compiler.
		c.buf[i] = 0
	}
	c.rw = NewByteReadWriter(rw)
	c.err = nil
}

func (c *Codec) Err() error {
	return c.err
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
