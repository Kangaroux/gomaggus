package value

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"sync"

	"github.com/phuslu/log"
)

const (
	// The number of bits in one block
	blockSizeBits = 32

	tagName   = "value"
	endMarker = "END"
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
	info := getStructLayout(e.root)

	for _, sectionIndex := range sections {
		section := info.sections[sectionIndex]

		for _, fieldIndex := range section.fields {
			e.encode(v.Field(fieldIndex))
		}
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
	if e.cursor+n > blockSizeBits {
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

	if e.cursor == blockSizeBits {
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

// structSection represents a section of blocks in a struct. Each section contains
// one or more blocks. Fields which are structs or slices cannot be split up and
// are considered a single section.
type structSection struct {
	// blockStart is the block within the struct this section starts at.
	blockStart int

	// fields is a list of field indexes in this section. Padding fields are not included.
	fields []int

	// size is the number of blocks in this section.
	size int
}

type structLayout struct {
	// TODO
	sections []structSection

	// size is the total number of blocks.
	size int
}

var structLayoutMap sync.Map // map[reflect.Type]*structLayout

func getStructLayout(v reflect.Value) *structLayout {
	if info, ok := structLayoutMap.Load(v.Type()); ok {
		return info.(*structLayout)
	}

	var currentSection structSection

	t := v.Type()
	numField := t.NumField()
	info := &structLayout{}
	bitSize := 0
	block := 0

	addBlock := func() {
		currentSection.blockStart = block
		currentSection.size = bitSize / blockSizeBits

		block += currentSection.size
		bitSize = 0

		info.sections = append(info.sections, currentSection)
		currentSection = structSection{}
	}

	for i := 0; i < numField; i++ {
		f := t.Field(i)

		if f.Tag.Get(tagName) == endMarker {
			break
		}

		size := dataSizeBits(f.Type)
		if size == -1 {
			panic("invalid field type " + f.Type.Kind().String())
		}

		bitSize += size

		// Padding fields are not included in the section list to avoid encoding them, however
		// their size is taken into account. Padding fields can also be their own section if
		// they are large enough, but the field list for the section will be empty.
		if f.Name != "_" {
			currentSection.fields = append(currentSection.fields, i)
		}

		if bitSize >= blockSizeBits {
			addBlock()
		}
	}

	if bitSize > 0 {
		addBlock()
	}

	info.size = block
	structLayoutMap.Store(t, info)

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

	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 32

	case reflect.Int64, reflect.Uint64:
		return 64

	default:
		return -1
	}
}
