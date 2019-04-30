package ddt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	original := []string{"A", "B", "C"}
	data := []interface{}{"A", "B", "C"}
	assert.True(t, defaultTestExecutor(nil))
	assert.Equal(t, data, defaultRowTransformer(original))
}
