package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizeImage(t *testing.T) {
	var v string

	v = sizeImage("image.pn?v=1", 22)
	assert.Equal(t, "image.pn?v=1&s=22", v)

	v = sizeImage("image.pn", 22)
	assert.Equal(t, "image.pn?s=22", v)
}
