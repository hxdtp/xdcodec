package xdcodec

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesArray(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(binary.BigEndian, buf)

	for i := 0; i < 2048; i++ {
		buf.Reset()
		s := RandString(i)
		{
			require.Nil(t, c.WriteSizedBytes(Atob(s)))
			var x []byte
			require.Nil(t, c.ReadSizedBytes(&x))
			require.Equal(t, s, Btoa(x))
		}
		{
			_, err := c.Write(Atob(s))
			require.Nil(t, err)
			x := make([]byte, len(s))
			_, err = c.Read(x)
			require.Nil(t, err)
			require.Equal(t, s, Btoa(x))
		}
		{
			require.Nil(t, c.WriteSizedString(s))
			var x string
			require.Nil(t, c.ReadSizedString(&x))
			require.Equal(t, s, x)
		}
	}
}
