package ddt

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

type buildOperations struct {
	inputBuilder func(suite *SuiteBuilder) (*csv.Reader, error)
	rowBuilder   func(suite *SuiteBuilder, data []string) *TestCase
}

func setupOperations() buildOperations {
	return buildOperations{buildInput, buildRow}
}

// Build - Takes the current state of the suite builder and creates the test cases ready to be executed.
// Steps to follow (each one will return the proper error if present):
// - Verify input file.
// - Replace parameters with text/template.
// - Read CSV rows and transform them from []string to []interface{}.
// - Skip first row if headers are present.
// - Group the final row properly.
// - Return the TestCase groups.
// ---
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
		group := suite.determineGroup(testCase)

		if _, ok := cases[group]; !ok {
			cases[group] = make([]TestCase, 0)
		}

		cases[group] = append(cases[group], *testCase)
	}

	return cases, nil
}

func buildInput(suite *SuiteBuilder) (*csv.Reader, error) {
	return nil, errors.New("")
}

func buildRow(suite *SuiteBuilder, data []string) *TestCase {
	return nil
}

func (suite *SuiteBuilder) determineGroup(testCase *TestCase) string {
	group := ""

	if suite.groupColumn >= 0 && suite.groupColumn < len(testCase.data) {
		group = fmt.Sprintf("%v", testCase.data[suite.groupColumn])
	}

	return group
}
