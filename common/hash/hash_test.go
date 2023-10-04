package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMD5(t *testing.T) {
	hash := HashMD5("123")
	assert.Equal(t, "202cb962ac59075b964b07152d234b70", hash)
}
