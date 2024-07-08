package value

import (
	"encoding/binary"
	"math"
	"reflect"

	"github.com/phuslu/log"
)

const (
	// The number of bits in one block
	blockSizeBits = 32

	tagName   = "value"
	endMarker = "END"
)

type encoder struct {
	buf       []byte
	bufOffset int

	block uint32

	// cursor keeps track of how many bits have been written into the block.
	// Its value ranges between [0, 32]. A cursor value of 32 means the block is full.
	cursor int

	// root is the top level struct value that was passed to Encode.
	root reflect.Value
}

func (e *encoder) Encode(v any, sections []structSection) []byte {
	var totalSize int

	for _, s := range sections {
		totalSize += s.size
	}

	if totalSize == 0 {
		return nil
	}

	e.buf = make([]byte, totalSize*4)
	e.bufOffset = 0

	e.encodeRoot(reflect.Indirect(reflect.ValueOf(v)), sections)
	e.flush()

	return e.buf
}

// encodeRoot encodes the struct v with some additional logic to handle v as the root struct.
func (e *encoder) encodeRoot(v reflect.Value, sections []structSection) {
	if v.Kind() != reflect.Struct {
		panic("encode non-struct type " + v.Kind().String())
	}

	e.root = v

	for _, section := range sections {
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
	if e.cursor > 0 {
		binary.LittleEndian.PutUint32(e.buf[e.bufOffset:], e.block)

		e.block = 0
		e.bufOffset += 4 // Advance the buffer 4 bytes
		e.cursor = 0
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

		log.Panic().
			Int("count", 8-byteBits).
			Int("near", e.bufOffset).
			Str("type", e.root.Type().Name()).
			Msg("missing bit padding")
	}

	// Align to n bits
	if blockBits := e.cursor % n; blockBits != 0 {
		e.cursor += n - blockBits

		log.Panic().
			Int("count", n-blockBits).
			Int("near", e.bufOffset).
			Str("type", e.root.Type().Name()).
			Msg("missing byte padding")
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

func marshalValues(v any, onlyDirty bool, dirty *dirtyValues) ([]byte, []structSection) {
	var sections []structSection
	enc := &encoder{}

	if onlyDirty {
		sections = dirty.Sections()
	} else {
		sections = dirty.layout.sections
	}

	ret := enc.Encode(v, sections)
	dirty.Clear()

	return ret, sections
}
