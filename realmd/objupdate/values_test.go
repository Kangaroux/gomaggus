package objupdate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBuildersOnce(t *testing.T) {
	v := Values{}

	assert.Equal(t, v.Object(), v.Object())
	assert.Equal(t, v.Unit(), v.Unit())
}

func TestUpdateExistingField(t *testing.T) {
	v := Values{}
	v.Unit().Health(100)
	v.Unit().Health(150)

	assert.Equal(t, 1, len(v.fields))
	assert.Equal(t, []uint32{150}, v.fields[0].value)
}

func TestValuesBytes(t *testing.T) {
	t.Run("no values set", func(t *testing.T) {
		v := Values{}
		expected := []byte{0} // single byte for the mask size

		assert.Equal(t, expected, v.Bytes())
	})

	t.Run("multiple values", func(t *testing.T) {
		v := Values{}
		v.Unit().Health(123)
		v.Object().Type(ObjectTypePlayer)
		expected := []byte{
			1,          // mask size
			4, 0, 0, 1, // mask (bits 2, 24)
			16, 0, 0, 0, // 1 << ObjectTypePlayer
			123, 0, 0, 0,
		}

		assert.Equal(t, expected, v.Bytes())
	})
}
