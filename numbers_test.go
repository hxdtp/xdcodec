package xdcodec

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFixedInt(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(binary.BigEndian, buf)
	{
		for _, n := range []int8{0, 64, math.MaxInt8} {
			buf.Reset()
			require.Nil(t, c.WriteInt8(n))
			var x int8
			require.Nil(t, c.ReadInt8(&x))
			require.Equal(t, n, x)
		}
	}
	{
		for _, n := range []int16{0, 64, math.MaxInt8, math.MaxInt16} {
			buf.Reset()
			require.Nil(t, c.WriteInt16(n))
			var x int16
			require.Nil(t, c.ReadInt16(&x))
			require.Equal(t, n, x)
		}
	}
	{
		for _, n := range []int32{0, 64, math.MaxInt8, math.MaxInt16, math.MaxInt32} {
			buf.Reset()
			require.Nil(t, c.WriteInt32(n))
			var x int32
			require.Nil(t, c.ReadInt32(&x))
			require.Equal(t, n, x)
		}
	}
	{
		for _, n := range []int64{0, 64, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64} {
			buf.Reset()
			require.Nil(t, c.WriteInt64(n))
			var x int64
			require.Nil(t, c.ReadInt64(&x))
			require.Equal(t, n, x)
		}
	}

	{ // Uint without reset.
		for _, n := range []uint8{0, 64, math.MaxInt8, math.MaxUint8} {
			require.Nil(t, c.WriteUint8(n))
			var x uint8
			require.Nil(t, c.ReadUint8(&x))
			require.Equal(t, n, x)
		}
	}
	{
		for _, n := range []uint16{0, 64, math.MaxInt8, math.MaxUint8, math.MaxInt16, math.MaxUint16} {
			require.Nil(t, c.WriteUint16(n))
			var x uint16
			require.Nil(t, c.ReadUint16(&x))
			require.Equal(t, n, x)
		}
	}
	{
		for _, n := range []uint32{0, 64, math.MaxInt8, math.MaxUint8,
			math.MaxInt16, math.MaxUint16, math.MaxInt32, math.MaxUint32} {
			require.Nil(t, c.WriteUint32(n))
			var x uint32
			require.Nil(t, c.ReadUint32(&x))
			require.Equal(t, n, x)
		}
	}
	{
		for _, n := range []uint64{0, 64, math.MaxInt8, math.MaxUint8,
			math.MaxInt16, math.MaxUint16, math.MaxInt32, math.MaxUint32,
			math.MaxInt64, math.MaxUint64} {
			require.Nil(t, c.WriteUint64(n))
			var x uint64
			require.Nil(t, c.ReadUint64(&x))
			require.Equal(t, n, x)
		}
	}
}

func TestVarInt(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(binary.BigEndian, buf)
	dataset := []uint64{
		0, 64,
		math.MaxInt8, math.MaxUint8,
		math.MaxInt16, math.MaxUint16,
		math.MaxInt32, math.MaxUint32,
		math.MaxInt64,
	}
	for _, v := range dataset {
		n := int64(v)
		buf.Reset()
		require.Nil(t, c.WriteVarint(n))
		var x int64
		require.Nil(t, c.ReadVarint(&x))
		require.Equal(t, n, x)
	}

	for _, un := range append(dataset, math.MaxUint64) {
		require.Nil(t, c.WriteUvarint(un))
		var x uint64
		require.Nil(t, c.ReadUvarint(&x))
		require.Equal(t, un, x)
	}
}
