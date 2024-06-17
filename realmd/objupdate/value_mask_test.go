package objupdate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueMask(t *testing.T) {
	t.Run("empty mask", func(t *testing.T) {
		assert.Empty(t, (&ValueMask{}).Mask())
	})
}
