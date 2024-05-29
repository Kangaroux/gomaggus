package srpv2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPad(t *testing.T) {
	assert.Equal(t, []byte{}, pad(0, []byte{}))
	assert.Equal(t, []byte{0, 0, 0, 0}, pad(4, []byte{}))
	assert.Equal(t, []byte{0, 1, 2}, pad(3, []byte{1, 2}))
}

func TestReverse(t *testing.T) {
	assert.Equal(t, []byte{0, 1}, reverse([]byte{1, 0}))
	assert.Equal(t, []byte{0, 1, 2}, reverse([]byte{2, 1, 0}))
}
