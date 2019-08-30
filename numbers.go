package xdcodec

import (
	"encoding/binary"
	"math"
)

func (c *Codec) ReadInt8(n *int8) error {
	var un uint8
	err := c.ReadUint8(&un)
	if err == nil {
		*n = int8(un)
	}
	return err
}

func (c *Codec) WriteInt8(n int8) error {
	return c.WriteUint8(uint8(n))
}

func (c *Codec) ReadInt16(n *int16) error {
	var un uint16
	err := c.ReadUint16(&un)
	if err == nil {
		*n = int16(un)
	}
	return err
}

func (c *Codec) WriteInt16(n int16) error {
	return c.WriteUint16(uint16(n))
}

func (c *Codec) ReadInt32(n *int32) error {
	var un uint32
	err := c.ReadUint32(&un)
	if err == nil {
		*n = int32(un)
	}
	return err
}

func (c *Codec) WriteInt32(n int32) error {
	return c.WriteUint32(uint32(n))
}

func (c *Codec) ReadInt64(n *int64) error {
	var un uint64
	err := c.ReadUint64(&un)
	if err == nil {
		*n = int64(un)
	}
	return err

}

func (c *Codec) WriteInt64(n int64) error {
	return c.WriteUint64(uint64(n))
}

func (c *Codec) ReadUint8(n *uint8) error {
	if c.err == nil {
		var b byte
		if b, c.err = c.rw.ReadByte(); c.err == nil {
			*n = b
		}
	}
	return c.err
}

func (c *Codec) WriteUint8(n uint8) error {
	if c.err == nil {
		c.err = c.rw.WriteByte(n)
	}
	return c.err
}

func (c *Codec) ReadUint16(n *uint16) error {
	buf := c.buf[:2]
	_, err := c.Read(buf)
	if err == nil {
		*n = c.byteOrder.Uint16(buf)
	}
	return err
}

func (c *Codec) WriteUint16(n uint16) error {
	buf := c.buf[:2]
	c.byteOrder.PutUint16(buf, n)
	_, err := c.Write(buf)
	return err
}

func (c *Codec) ReadUint32(n *uint32) error {
	buf := c.buf[:4]
	_, err := c.Read(buf)
	if err == nil {
		*n = c.byteOrder.Uint32(buf)
	}
	return err
}

func (c *Codec) WriteUint32(n uint32) error {
	buf := c.buf[:4]
	c.byteOrder.PutUint32(buf, n)
	_, err := c.Write(buf)
	return err
}

func (c *Codec) ReadUint64(n *uint64) error {
	buf := c.buf[:8]
	_, err := c.Read(buf)
	if err == nil {
		*n = c.byteOrder.Uint64(buf)
	}
	return err
}

func (c *Codec) WriteUint64(n uint64) error {
	buf := c.buf[:8]
	c.byteOrder.PutUint64(buf, n)
	_, err := c.Write(buf)
	return err
}

func (c *Codec) ReadFloat32(f *float32) error {
	var bits uint32
	if err := c.ReadUint32(&bits); err != nil {
		return err
	}
	*f = math.Float32frombits(bits)
	return nil
}

func (c *Codec) WriteFloat32(f float32) error {
	bits := math.Float32bits(f)
	return c.WriteUint32(bits)
}

func (c *Codec) ReadFloat64(f *float64) error {
	var bits uint64
	if err := c.ReadUint64(&bits); err != nil {
		return err
	}
	*f = math.Float64frombits(bits)
	return nil
}

func (c *Codec) WriteFloat64(f float64) error {
	bits := math.Float64bits(f)
	return c.WriteUint64(bits)
}

func (c *Codec) ReadVarint(n *int64) error {
	if c.err == nil {
		*n, c.err = binary.ReadVarint(c.rw)
	}
	return c.err
}

func (c *Codec) WriteVarint(n int64) error {
	if c.err != nil {
		return c.err
	}
	buf := c.buf[:]
	offset := binary.PutVarint(buf, n)
	_, err := c.Write(buf[:offset])
	return err
}

func (c *Codec) ReadUvarint(n *uint64) error {
	if c.err == nil {
		*n, c.err = binary.ReadUvarint(c.rw)
	}
	return c.err
}

func (c *Codec) WriteUvarint(n uint64) error {
	if c.err != nil {
		return c.err
	}
	buf := c.buf[:]
	offset := binary.PutUvarint(buf, n)
	_, err := c.Write(buf[:offset])
	return err
}

func SizeOfVarint(n int64) int {
	un := uint64(n) << 1
	if n < 0 {
		un = ^un
	}
	return SizeOfUvarint(un)
}

func SizeOfUvarint(n uint64) int {
	i := 0
	for n >= 0x80 {
		n >>= 7
		i++
	}
	return i + 1
}
