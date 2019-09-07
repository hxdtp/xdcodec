package xdcodec

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
)

var (
	ErrExceedContainerCap = errors.New("xdcodec: containers can only contain 255 elements")
)

const (
	ContainerCapacity = 255

	typeInt    = 'i'
	typeUint   = 'u'
	typeFloat  = 'f'
	typeBytes  = 'b'
	typeString = 's'
	typeList   = 'l'
	typeMap    = 'm'
)

type (
	// Typed using type-length-value encoding.
	Typed interface {
		T() byte
	}
	List []Typed
	Map  map[string]Typed

	TypedInt    int64
	TypedUint   uint64
	TypedFloat  float64
	TypedBytes  []byte
	TypedString string
	TypedList   struct{ List }
	TypedMap    struct{ Map }
)

func (t TypedInt) T() byte    { return typeInt }
func (t TypedUint) T() byte   { return typeUint }
func (t TypedFloat) T() byte  { return typeFloat }
func (t TypedBytes) T() byte  { return typeBytes }
func (t TypedString) T() byte { return typeString }
func (t TypedList) T() byte   { return typeList }
func (t TypedMap) T() byte    { return typeMap }

func (c *Codec) ReadList(l *List) error {
	var b byte
	if b, c.err = c.rw.ReadByte(); c.err != nil {
		return c.err
	}
	nelem := uint8(b)
	if nelem == 0 {
		return nil
	}

	if int(nelem) > len(*l) {
		if c := cap(*l); c >= int(nelem) {
			*l = (*l)[:c]
		} else {
			*l = make(List, nelem, nelem)
		}
	}

	for i := uint8(0); i < nelem; i++ {
		e, err := c.ReadTyped()
		if err != nil {
			return err
		}
		(*l)[i] = e
	}
	return nil
}

func (c *Codec) WriteList(l List) error {
	n := len(l)
	if n > ContainerCapacity {
		return ErrExceedContainerCap
	}

	c.WriteUint8(uint8(n))
	for _, e := range l {
		if err := c.WriteTyped(e); err != nil {
			return err
		}
	}
	return nil
}

func (c *Codec) ReadMap(m *Map) error {
	var b byte
	if b, c.err = c.rw.ReadByte(); c.err != nil {
		return c.err
	}
	nkv := uint8(b)
	if nkv == 0 {
		return nil
	}
	if *m == nil {
		*m = Map{}
	}

	for i := uint8(0); i < nkv; i++ {
		var k string
		if err := c.ReadSizedString(&k); err != nil {
			return err
		}
		v, err := c.ReadTyped()
		if err != nil {
			return err
		}
		(*m)[k] = v
	}
	return nil
}

func (c *Codec) keys(n int) []string {
	if cap(c.keybuf) < n {
		c.keybuf = make([]string, n, n)
	}
	return c.keybuf[:n]
}

func (c *Codec) WriteMap(m Map) error {
	n := len(m)
	if n > ContainerCapacity {
		return ErrExceedContainerCap
	}
	keys := c.keys(n)[:0]
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	c.WriteUint8(uint8(n))
	for _, k := range keys {
		c.WriteSizedString(k)
		if err := c.WriteTyped(m[k]); err != nil {
			return err
		}
	}
	return nil
}

func (c *Codec) ReadTyped() (t Typed, err error) {
	if c.err != nil {
		err = c.err
		return
	}
	var typ byte
	if typ, c.err = c.rw.ReadByte(); c.err != nil {
		err = c.err
		return
	}

	switch typ {
	case typeInt:
		var n int64
		if err = c.ReadVarint(&n); err == nil {
			t = TypedInt(n)
		}
	case typeUint:
		var n uint64
		if err = c.ReadUvarint(&n); err == nil {
			t = TypedUint(n)
		}
	case typeFloat:
		var n uint64
		if err = c.ReadUvarint(&n); err == nil {
			t = TypedFloat(math.Float64frombits(n))
		}
	case typeBytes:
		var bs []byte
		if err = c.ReadSizedBytes(&bs); err == nil {
			t = TypedBytes(bs)
		}
	case typeString:
		var s string
		if err = c.ReadSizedString(&s); err == nil {
			t = TypedString(s)
		}
	case typeList:
		var l List
		if err = c.ReadList(&l); err == nil {
			t = TypedList{List: l}
		}
	case typeMap:
		var m Map
		if err = c.ReadMap(&m); err == nil {
			t = TypedMap{Map: m}
		}
	default:
		err = fmt.Errorf("unknown type: %v", reflect.TypeOf(t))
	}
	return
}

func (c *Codec) WriteTyped(t Typed) error {
	if c.err != nil {
		return c.err
	}
	switch x := t.(type) {
	case TypedInt:
		c.Write([]byte{typeInt})
		c.WriteVarint(int64(x))
	case TypedUint:
		c.Write([]byte{typeUint})
		c.WriteUvarint(uint64(x))
	case TypedFloat:
		c.Write([]byte{typeFloat})
		c.WriteUvarint(math.Float64bits(float64(x)))
	case TypedBytes:
		c.Write([]byte{typeBytes})
		c.WriteSizedBytes(x)
	case TypedString:
		// No copy, Ref:
		// https://stackoverflow.com/questions/32253768/does-convertion-between-alias-types-in-go-create-copies
		c.Write([]byte{typeString})
		c.WriteSizedString(string(x))
	case TypedList:
		c.Write([]byte{typeList})
		c.WriteList(x.List)
	case TypedMap:
		c.Write([]byte{typeMap})
		c.WriteMap(x.Map)
	default:
		return fmt.Errorf("unsupported type: %v", reflect.TypeOf(t))
	}
	return c.err
}

func SizeOfTyped(t Typed) int {
	switch x := t.(type) {
	case TypedInt:
		return 1 + SizeOfVarint(int64(x))
	case TypedUint:
		return 1 + SizeOfUvarint(uint64(x))
	case TypedFloat:
		return 1 + SizeOfUvarint(math.Float64bits(float64(x)))
	case TypedBytes:
		return 1 + SizeOfSizedBytes(x)
	case TypedString:
		return 1 + SizeOfSizedString(string(x))
	case TypedList:
		return 1 + x.List.Size()
	case TypedMap:
		return 1 + x.Map.Size()
	default:
		panic("unsupported type")
	}
}

func (l List) Size() int {
	n := SizeOfUvarint(uint64(len(l)))
	for _, e := range l {
		n += SizeOfTyped(e)
	}
	return n
}

func (l List) Reset() {
	for i := range l { // Optimized by compiler.
		l[i] = nil
	}
}

func (m Map) Size() int {
	n := SizeOfUvarint(uint64(len(m)))
	for k, v := range m {
		n += SizeOfSizedString(k)
		n += SizeOfTyped(v)
	}
	return n
}

func (m Map) Reset() {
	for k := range m {
		delete(m, k)
	}
}
