package value

import (
	"bytes"
	"encoding/binary"
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
