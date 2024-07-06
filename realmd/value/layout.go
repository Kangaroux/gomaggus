package value

import (
	"reflect"
	"sync"
)

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

	// nameToSection maps a field name to the section index it belongs to.
	nameToSection map[string]int
}

var structLayoutMap sync.Map // map[reflect.Type]*structLayout

func getStructLayout(v reflect.Value) *structLayout {
	if info, ok := structLayoutMap.Load(v.Type()); ok {
		return info.(*structLayout)
	}

	var currentSection structSection

	t := v.Type()
	numField := t.NumField()
	info := &structLayout{
		nameToSection: make(map[string]int),
	}
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
			info.nameToSection[f.Name] = len(info.sections)
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
