package xdcodec

import (
	"io"
	"unsafe"
)

func (c *Codec) Read(p []byte) (int, error) {
	var n int
	if c.err == nil {
		n, c.err = io.ReadFull(c.rw, p)
	}
	return n, c.err
}

func (c *Codec) Write(b []byte) (int, error) {
	var n int
	if c.err == nil {
		n, c.err = c.rw.Write(b)
	}
	return n, c.err
}

func (c *Codec) ReadSizedBytes(p *[]byte) error {
	var n uint64
	if err := c.ReadUvarint(&n); err != nil {
		return err
	}
	if int(n) > len(*p) {
		if c := cap(*p); c >= int(n) {
			*p = (*p)[:c]
		} else {
			*p = make([]byte, n, n)
		}
	}
	_, err := c.Read((*p)[:n])
	return err
}

func (c *Codec) WriteSizedBytes(b []byte) error {
	if err := c.WriteUvarint(uint64(len(b))); err != nil {
		return err
	}
	_, err := c.Write(b)
	return err
}

func (c *Codec) ReadSizedString(s *string) error {
	var p []byte
	err := c.ReadSizedBytes(&p)
	if err == nil {
		*s = Btoa(p)
	}
	return err
}

func (c *Codec) WriteSizedString(s string) error {
	return c.WriteSizedBytes(Atob(s))
}

// Atob converts string into bytes without copy.
func Atob(s string) []byte {
	if s == "" {
		return []byte{}
	}
	return *(*[]byte)(unsafe.Pointer(&s))
}

// Btoa converts bytes into string without copy.
func Btoa(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}

func SizeOfSizedBytes(b []byte) int {
	return SizeOfUvarint(uint64(len(b))) + len(b)
}

func SizeOfSizedString(s string) int {
	return SizeOfSizedBytes(Atob(s))
}
