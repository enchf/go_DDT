package ddt

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilderMethods(t *testing.T) {
	inputFile := "fake.file"
	client := NewSuiteBuilder(inputFile)

	assert.Equal(t, inputFile, client.inputFile)
	assert.False(t, client.headers)
	assert.Nil(t, client.variables)
	assert.Equal(t, "", client.Name)
	assert.Equal(t, -1, client.groupColumn)

	name := "Suite Test"
	vars := map[string]interface{}{"A": true, "B": false}

	transformer := func(row []string) []interface{} {
		val1, _ := strconv.Atoi(row[1])
		val2, _ := strconv.Atoi(row[2])
		val4, _ := strconv.Atoi(row[4])

		return []interface{}{row[0], val1, val2, row[3], val4}
	}

	test := func(data []interface{}) bool {
		// Trivial as it is not executed.
		return data[2] == data[4]
	}

	client.GroupBy(0).Headers(true).Variables(vars)
	client.GlobalName(name).RowTransformer(transformer).TestExecutor(test)

	assert.Equal(t, inputFile, client.inputFile)
	assert.True(t, client.headers)
	assert.Equal(t, vars, client.variables)
	assert.Equal(t, name, client.Name)
	assert.Equal(t, 0, client.groupColumn)
}
