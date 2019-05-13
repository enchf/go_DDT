package ddt

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var badInputBuilder = func(suite *SuiteBuilder) (*csv.Reader, error) { return nil, errors.New("Error") }

func TestDetermineGroup(t *testing.T) {
	suite := NewSuiteBuilder("fake.file")
	data := []interface{}{"A", "B"}

	assert.Equal(t, "", suite.determineGroup(data))

	suite.GroupBy(1)
	assert.Equal(t, "B", suite.determineGroup(data))
}

func TestBuildRow(t *testing.T) {
	suite := NewSuiteBuilder("fake.file")
	data := []string{"A", "B"}
	expectedData := []interface{}{"A", "B"}

	testCase := buildRow(suite, data)

	assert.Equal(t, "", testCase.Group)
	assert.Equal(t, expectedData, testCase.data)

	assert.True(t, testCase.Run())
}

func TestBuildInputErrors(t *testing.T) {
	suite := NewSuiteBuilder("unexisting.file")

	reader, err := buildInput(suite)
	assert.Nil(t, reader)
	assert.NotNil(t, err)

	tmp, err := createTempFile()
	if err != nil {
		assert.Fail(t, err.Error())
	}

	defer os.Remove(tmp.Name())

	if err := writeFile(tmp, "Bad {{.template"); err != nil {
		assert.Fail(t, err.Error())
	}

	vars := map[string]interface{}{"template": "value", "notused": "never"}

	suite.inputFile = tmp.Name()
	suite.Variables(vars)

	reader, err = buildInput(suite)
	assert.Nil(t, reader)
	assert.NotNil(t, err)

	if err := tmp.Close(); err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestBuildInputVariablesError(t *testing.T) {
	tmp, err := createTempFile()
	if err != nil {
		assert.Fail(t, err.Error())
	}

	defer os.Remove(tmp.Name())

	suite := NewSuiteBuilder(tmp.Name()).Variables(map[string]interface{}{"A": nil})

	if err := writeFile(tmp, "{{len .A}}"); err != nil {
		assert.Fail(t, err.Error())
	}

	reader, err := buildInput(suite)
	assert.Nil(t, reader)
	assert.NotNil(t, err)
}

func TestBuildInput(t *testing.T) {
	tmp, err := createTempFile()
	if err != nil {
		assert.Fail(t, err.Error())
	}

	defer os.Remove(tmp.Name())
	suite := NewSuiteBuilder(tmp.Name())

	if err := writeFile(tmp, "{{.A}},{{.B}}"); err != nil {
		assert.Fail(t, err.Error())
	}

	reader, err := buildInput(suite)

	assert.NotNil(t, reader)
	assert.Nil(t, err)

	rows, err := getRows(reader)
	assert.Equal(t, 1, len(rows))
	assert.Equal(t, []string{"", ""}, rows[0])

	vars := map[string]interface{}{"A": "value", "B": "other", "C": "notused"}
	suite.Variables(vars)

	reader, err = buildInput(suite)

	assert.NotNil(t, reader)
	assert.Nil(t, err)

	rows, err = getRows(reader)
	assert.Equal(t, 1, len(rows))
	assert.Equal(t, []string{"value", "other"}, rows[0])

	if err := tmp.Close(); err != nil {
		assert.Fail(t, err.Error())
	}
}

/**
func (suite *SuiteBuilder) Build() (map[string][]TestCase, error) {
	reader, err := suite.operations.inputBuilder(suite)

	if err != nil {
		return nil, err
	}

	cases := make(map[string][]TestCase)

	for {
		row, err := reader.Read()

		if err != nil {
			if err != io.EOF {
				return nil, err
			}

			break
		}

		testCase := suite.operations.rowBuilder(suite, row)

		if _, ok := cases[testCase.Group]; !ok {
			cases[testCase.Group] = make([]TestCase, 0)
		}

		cases[testCase.Group] = append(cases[testCase.Group], *testCase)
	}

	return cases, nil
}
*/

func TestBuildWithGroup(t *testing.T) {
	suite := NewSuiteBuilder("fake.file").GroupBy(0)
	suite.operations = buildOperations{goodInputBuilder("1,2\n1,3\n2,4"), buildRow}

	expected := map[string][][]interface{}{"1": [][]interface{}{{"1", "2"}, {"1", "3"}}, "2": [][]interface{}{{"2", "4"}}}

	cases, err := suite.Build()
	assert.Nil(t, err)

	for key, values := range expected {
		caseValues, ok := cases[key]

		assert.True(t, ok)
		assert.Equal(t, len(values), len(caseValues))

		for i := 0; i < len(values); i++ {
			assert.Equal(t, values[i], caseValues[i].data)
		}
	}
}

func TestBuildWithoutGroup(t *testing.T) {
	suite := NewSuiteBuilder("fake.file")
	suite.operations = buildOperations{goodInputBuilder("1,2\n1,3\n2,4"), buildRow}

	cases, err := suite.Build()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(cases))

	caseValues := make([][]interface{}, 3)
	uniqueCase := cases[""]

	assert.NotNil(t, uniqueCase)

	for i := 0; i < len(uniqueCase); i++ {
		caseValues[i] = uniqueCase[i].data
	}

	expected := [][]interface{}{[]interface{}{"1", "2"}, []interface{}{"1", "3"}, []interface{}{"2", "4"}}
	assert.Equal(t, expected, caseValues)
}

func TestBuildErrors(t *testing.T) {
	suite := NewSuiteBuilder("fake.file")
	suite.operations = buildOperations{badInputBuilder, buildRow}

	cases, err := suite.Build()
	assert.Nil(t, cases)
	assert.NotNil(t, err)

	suite.operations.inputBuilder = goodInputBuilder("1\",B")
	cases, err = suite.Build()
	assert.Nil(t, cases)
	assert.NotNil(t, err)
}

// Helper functions

func goodInputBuilder(content string) func(suite *SuiteBuilder) (*csv.Reader, error) {
	return func(suite *SuiteBuilder) (*csv.Reader, error) {
		return csv.NewReader(strings.NewReader(content)), nil
	}
}

func createTempFile() (*os.File, error) {
	return ioutil.TempFile("", "build_test")
}

func writeFile(file *os.File, content string) error {
	_, err := file.Write([]byte(content))
	return err
}

func getRows(reader *csv.Reader) ([][]string, error) {
	rows := make([][]string, 0)

	for {
		row, err := reader.Read()

		if err != nil {
			if err != io.EOF {
				return nil, err
			}

			break
		}

		rows = append(rows, row)
	}

	return rows, nil
}
