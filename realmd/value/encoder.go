package value

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"sync"

	"github.com/phuslu/log"
)

type encoder struct {
	buf   bytes.Buffer
	block uint32

	// cursor keeps track of how many bits have been written into the block.
	// Its value ranges between [0, 32]. A cursor value of 32 means the block is full.
	cursor int

	// root is the top level struct value that was passed to Encode.
	root reflect.Value
}

func (e *encoder) Encode(v any, sections []int) []byte {
	e.buf.Reset()

	e.encodeRoot(reflect.Indirect(reflect.ValueOf(v)), sections)
	e.flush()

	return bytes.Clone(e.buf.Bytes())
}

// encodeRoot encodes the struct v with some additional logic to handle v as the root struct.
func (e *encoder) encodeRoot(v reflect.Value, sections []int) {
	if v.Kind() != reflect.Struct {
		panic("encode non-struct type " + v.Kind().String())
	}

	e.root = v
	info := getStructInfo(e.root)
	numField := v.NumField()

	for i := 0; i < numField; i++ {
		if i == info.endIndex {
			return
		}

		e.encode(v.Field(i))
	}
}

// encode writes v to the buffer as uint32 blocks.
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

	case reflect.Uint8:
		e.writeN(uint32(v.Uint()), 1)
	case reflect.Int8:
		e.writeN(uint32(v.Int()), 1)

	case reflect.Uint16:
		e.writeN(uint32(v.Uint()), 2)
	case reflect.Int16:
		e.writeN(uint32(v.Int()), 2)

	case reflect.Uint32:
		e.writeN(uint32(v.Uint()), 4)
	case reflect.Int32:
		e.writeN(uint32(v.Int()), 4)

	case reflect.Float32:
		e.writeN(math.Float32bits(float32(v.Float())), 4)

	case reflect.Uint64:
		e.writeN(uint32(v.Uint()), 4)
		e.writeN(uint32(v.Uint()>>32), 4)
	case reflect.Int64:
		e.writeN(uint32(v.Int()), 4)
		e.writeN(uint32(v.Int()>>32), 4)

	default:
		panic("unknown field " + v.Kind().String())
	}
}

// flush writes the block to the buffer if it's non-empty.
func (e *encoder) flush() {
	var data [4]byte

	if e.cursor > 0 {
		binary.LittleEndian.PutUint32(data[:], e.block)
		e.buf.Write(data[:])
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

		log.Warn().Func(
			func(entry *log.Entry) {
				entry.
					Int("count", 8-byteBits).
					Int("near", e.buf.Len()).
					Str("type", e.root.Type().Name()).
					Msg("missing bit padding")
			})
	}

	// Align to n bits
	if blockBits := e.cursor % n; blockBits != 0 {
		e.cursor += n - blockBits

		log.Warn().Func(
			func(entry *log.Entry) {
				entry.
					Int("count", n-blockBits).
					Int("near", e.buf.Len()).
					Str("type", e.root.Type().Name()).
					Msg("missing byte padding")
			})
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

// structSession represents a section or collection or fields in a values struct.
// Fields are added to sections until the section size is at least 4 bytes.
// Slices or structs cannot be split up, so it's possible for a section to only
// contain one field but be hundreds of bytes long.
type structSection struct {
	// fields is a list of field indexes in this section.
	fields []int

	// size is the number of 4 byte blocks in this section.
	size int
}

type structInfo struct {
	// endIndex is the index of the field that marks where encoding should stop.
	// This enables storing additional metadata inside the struct without it being encoded.
	// If the struct has no end field, endIndex is -1.
	endIndex int

	// TODO
	sections []structSection
}

var structInfoMap sync.Map // map[reflect.Type]*structInfo

func getStructInfo(v reflect.Value) *structInfo {
	if info, ok := structInfoMap.Load(v.Type()); ok {
		return info.(*structInfo)
	}

	var currentSection structSection

	t := v.Type()
	numField := t.NumField()
	info := &structInfo{endIndex: -1}
	bitSize := 0

	for i := 0; i < numField; i++ {
		f := t.Field(i)

		if f.Tag.Get(tagName) == endMarker {
			info.endIndex = i
			break
		}

		bitSize += dataSizeBits(f.Type)

		// Padding fields are not included in the section field list because they are never encoded.
		// However, their size should be taken into consideration.
		if f.Name == "_" {
			continue
		}

		currentSection.fields = append(currentSection.fields, i)

		if bitSize >= 32 {
			currentSection.size = bitSize / 8
			bitSize = 0
			info.sections = append(info.sections, currentSection)
			currentSection.fields = currentSection.fields[:0]
		}
	}

	if bitSize > 0 {
		currentSection.size = bitSize / 8
		info.sections = append(info.sections, currentSection)
	}

	structInfoMap.Store(t, info)

	return info
}

// dataSizeBits returns the number of bits needed to store t.
func dataSizeBits(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		return t.Len() * dataSizeBits(t.Elem())

	case reflect.Struct:
		size := 0
		numField := t.NumField()
		for i := 0; i < numField; i++ {
			size += dataSizeBits(t.Field(i).Type)
		}
		return size

	case reflect.Bool:
		return 1

	case reflect.Int8, reflect.Uint8:
		return 8

	case reflect.Int16, reflect.Uint16:
		return 16

	case reflect.Int32, reflect.Uint32:
		return 32

	case reflect.Int64, reflect.Uint64:
		return 64

	default:
		return -1
	}
}
