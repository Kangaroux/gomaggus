package value

import "github.com/phuslu/log"

type dirtyValues struct {
	layout   *structLayout
	sections map[int]bool
}

func newDirtyValues(layout *structLayout) *dirtyValues {
	return &dirtyValues{
		layout:   layout,
		sections: make(map[int]bool),
	}
}

// Clear clears all the flagged sections.
func (dv *dirtyValues) Clear() {
	dv.sections = make(map[int]bool)
}

// Flag marks the field's section as dirty. Sections marked as dirty will be sent
// to clients on the next update.
func (dv *dirtyValues) Flag(fieldName string) {
	i, ok := dv.layout.nameToSection[fieldName]
	if !ok {
		log.Warn().Str("field", fieldName).Msg("unknown field name")
		return
	}

	dv.sections[i] = true
}

// Sections returns a list of dirty section indexes.
func (dv *dirtyValues) Sections() []int {
	sections := make([]int, 0, len(dv.sections))

	for i := range dv.sections {
		sections = append(sections, i)
	}

	return sections
}
