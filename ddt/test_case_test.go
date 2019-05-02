package ddt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type execution struct{ executed bool }

func TestRun(t *testing.T) {
	executionFlag := execution{false}
	test := func(data []interface{}) bool {
		data[0].(*execution).executed = true
		return true
	}
	data := []interface{}{&executionFlag}
	testCase := TestCase{"group", data, test}

	assert.True(t, testCase.Run())
	assert.True(t, executionFlag.executed)
}
