package xdcodec

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(binary.BigEndian, buf)

	l := List{
		TypedInt(0),
		TypedInt(math.MaxInt32),
		TypedInt(math.MaxInt64),
		TypedUint(0),
		TypedUint(math.MaxUint32),
		TypedUint(math.MaxUint64),
		TypedFloat(0),
		TypedFloat(math.MaxFloat32),
		TypedFloat(math.MaxFloat64),
		TypedBytes(nil), // FIXME(damnever: zero values.
		TypedBytes(RandString(128)),
		TypedBytes(RandString(1024)),
		TypedString(""), // FIXME(damnever: zero values.
		TypedString(RandString(128)),
		TypedString(RandString(1024)),
		TypedList{},
		TypedList{
			List: List{
				TypedInt(128),
				TypedString(RandString(2096)),
			},
		},
		TypedMap{},
		TypedMap{
			Map: Map{"a": TypedUint(123)},
		},
	}

	require.Nil(t, c.WriteList(l))
	var lx List
	require.Nil(t, c.ReadList(&lx))
	require.Equal(t, l, lx)
}

func TestMap(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(binary.BigEndian, buf)

	m := Map{
		RandString(1):  TypedInt(0),
		RandString(2):  TypedInt(math.MaxInt32),
		RandString(3):  TypedInt(math.MaxInt64),
		RandString(4):  TypedUint(0),
		RandString(5):  TypedUint(math.MaxUint32),
		RandString(6):  TypedUint(math.MaxUint64),
		RandString(7):  TypedFloat(0),
		RandString(8):  TypedFloat(math.MaxFloat32),
		RandString(9):  TypedFloat(math.MaxFloat64),
		RandString(10): TypedBytes(nil), // FIXME(damnever: zero values.
		RandString(11): TypedBytes(RandString(128)),
		RandString(12): TypedBytes(RandString(1024)),
		RandString(13): TypedString(""), // FIXME(damnever: zero values.
		RandString(14): TypedString(RandString(128)),
		RandString(15): TypedString(RandString(1024)),
		RandString(16): TypedList{},
		RandString(17): TypedList{
			List: List{
				TypedInt(math.MaxInt8),
				TypedString(RandString(2096)),
			},
		},
		RandString(18): TypedMap{},
		RandString(255): TypedMap{
			Map: Map{"a": TypedUint(math.MaxUint16)},
		},
	}

	require.Nil(t, c.WriteMap(m))
	var mx Map
	require.Nil(t, c.ReadMap(&mx))
	require.Equal(t, m, mx)
}
