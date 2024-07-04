package value

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"sync"
	"unsafe"
)

type blockField struct {
	index int
	rtype reflect.Type
}

type fieldBlock struct {
	index  int
	fields []blockField
}

// struct value
// inspect the size of each field
// group the fields (by index) into 4 byte blocks
// use the indexes with t.FieldByIndex and write them
func StructFieldBlocks(value any) []fieldBlock {
	var blocks []fieldBlock
	var currentBlock fieldBlock

	rv := reflect.Indirect(reflect.ValueOf(value))
	t := rv.Type()
	numField := rv.NumField()
	pos := 0

	for i := 0; i < numField; i++ {
		f := t.Field(i)

		if f.Tag.Get(tagName) == endMarker {
			break
		}

		if f.Name != "_" {
			blockIndex := pos / 4
			if blockIndex != currentBlock.index && len(currentBlock.fields) > 0 {
				blocks = append(blocks, currentBlock)
				currentBlock = fieldBlock{index: blockIndex}
			}

			currentBlock.fields = append(currentBlock.fields, blockField{index: i, rtype: f.Type})
		}

		pos += dataSize(rv.Field(i))
	}

	if len(currentBlock.fields) > 0 {
		blocks = append(blocks, currentBlock)
	}

	return blocks
}

func Encode(value any, fieldBlocks []fieldBlock, encodeBlocks []int) []byte {
	var buf bytes.Buffer
	var blockBuf bytes.Buffer

	rv := reflect.Indirect(reflect.ValueOf(value))
	j := 0

	for _, i := range encodeBlocks {
		var block *fieldBlock

		for ; j < len(fieldBlocks); j++ {
			if fieldBlocks[j].index == i {
				block = &fieldBlocks[j]
				break
			}
		}

		if block == nil {
			panic("encode block out of range")
		}

		for _, b := range block.fields {
			f := rv.Field(b.index)
			binary.Write(&blockBuf, binary.BigEndian, f.Interface())
		}

		sd := (*uint32)(unsafe.Pointer(unsafe.SliceData(blockBuf.Bytes())))
		data := unsafe.Slice(sd, blockBuf.Len()/4)
		binary.Write(&buf, binary.LittleEndian, data)
	}

	return buf.Bytes()
}

var structSize sync.Map

// dataSize returns the number of bytes the actual data represented by v occupies in memory.
// For compound structures, it sums the sizes of the elements. Thus, for instance, for a slice
// it returns the length of the slice times the element size and does not count the memory
// occupied by the header. If the type of v is not acceptable, dataSize returns -1.
func dataSize(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Slice:
		if s := sizeof(v.Type().Elem()); s >= 0 {
			return s * v.Len()
		}
		return -1

	case reflect.Struct:
		t := v.Type()
		if size, ok := structSize.Load(t); ok {
			return size.(int)
		}
		size := sizeof(t)
		structSize.Store(t, size)
		return size

	default:
		return sizeof(v.Type())
	}
}

// sizeof returns the size >= 0 of variables for the given type or -1 if the type is not acceptable.
func sizeof(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		if s := sizeof(t.Elem()); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := sizeof(t.Field(i).Type)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(t.Size())
	}

	return -1
}

type encoder struct {
	buf   bytes.Buffer
	block uint32

	// cursor keeps track of how many bits have been written into the block.
	// Its value ranges between [0, 32]. A cursor value of 32 means the block is full.
	cursor int
}

func (e *encoder) Encode(v any) []byte {
	e.buf.Reset()

	e.encode(reflect.Indirect(reflect.ValueOf(v)))
	e.flush()

	return bytes.Clone(e.buf.Bytes())
}

func (e *encoder) encode(v reflect.Value) {
	switch v.Kind() {
	case reflect.Struct:
		numField := v.NumField()
		for i := 0; i < numField; i++ {
			e.encode(v.Field(i))
		}

	case reflect.Array:
		length := v.Len()
		for i := 0; i < length; i++ {
			e.encode(v.Index(i))
		}

	case reflect.Bool:
		e.writeBit(v.Bool())

	case reflect.Uint8, reflect.Int8:
		e.writeN(uint32(v.Uint()), 1)

	case reflect.Uint16, reflect.Int16:
		e.writeN(uint32(v.Uint()), 2)

	case reflect.Uint32, reflect.Int32:
		e.writeN(uint32(v.Uint()), 4)

	case reflect.Float32:
		e.writeN(math.Float32bits(float32(v.Float())), 4)

	default:
		panic("unknown field")
	}
}

// flush writes the block to the buffer if it's non-empty.
func (e *encoder) flush() {
	if e.cursor > 0 {
		binary.Write(&e.buf, binary.LittleEndian, e.block)
		e.cursor = 0
		e.block = 0
	}
}

// align makes room for n bytes inside the block and aligns the cursor to be a multiple of n.
// The block is automatically flushed if it can't fit n bytes.
// align panics if n is not 1, 2, or 4.
func (e *encoder) align(n int) {
	switch n {
	case 1, 2, 4:
		// Convert bytes to bits
		n *= 8
	default:
		panic("n must be 1, 2, or 4")
	}

	// Start a new byte
	if byteBits := e.cursor % 8; byteBits != 0 {
		e.cursor += 8 - byteBits
	}

	// Align to n bits
	if blockBits := e.cursor % n; blockBits != 0 {
		e.cursor += n - blockBits
	}

	// Block can't fit n bits
	if e.cursor+n > 32 {
		e.flush()
	}
}

// writeN interprets val as an n-byte value and writes it to the block.
// writeN panics if n is not 1, 2, or 4.
func (e *encoder) writeN(val uint32, n int) {
	e.align(n)
	e.block |= val << uint32(e.cursor)
	e.cursor += n * 8
}

func (e *encoder) writeBit(b bool) {
	var v uint32

	if e.cursor == 32 {
		e.flush()
	}

	if b {
		v = 1
	} else {
		v = 0
	}

	e.block |= v << uint32(e.cursor)
	e.cursor++
}
