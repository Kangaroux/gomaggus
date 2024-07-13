package values

import (
	"math"
	"testing"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/stretchr/testify/assert"
)

func TestObjectMarshal(t *testing.T) {
	var data []byte

	o := NewObjectData()

	data, _ = o.Marshal(true)
	assert.Equal(t, 0, len(data))

	data, _ = o.Marshal(false)
	assert.Equal(t, 20, len(data))

	o.SetGUID(0x1122334455667788)
	o.SetEntry(0x12345678)
	o.SetType(ObjectTypeObject, ObjectTypeUnit, ObjectTypePlayer)
	o.SetScaleX(math.Float32frombits(0xEFBEADDE))

	data, _ = o.Marshal(true)
	assert.Equal(t, internal.MustDecodeHex("88776655443322111900000078563412DEADBEEF"), data)
}
