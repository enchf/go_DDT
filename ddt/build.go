package ddt

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"path"
	"strings"
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

		if _, ok := cases[testCase.Group]; !ok {
			cases[testCase.Group] = make([]TestCase, 0)
		}

		cases[testCase.Group] = append(cases[testCase.Group], *testCase)
	}

	return cases, nil
}

func buildInput(suite *SuiteBuilder) (*csv.Reader, error) {
	tpl, err := template.New(path.Base(suite.inputFile)).ParseFiles(suite.inputFile)

	if err != nil {
		return nil, err
	}

	context := suite.variables
	var buf strings.Builder

	if context == nil {
		context = map[string]interface{}{}
	}

	err = tpl.Execute(&buf, context)

	if err != nil {
		return nil, err
	}

	stringReader := strings.NewReader(buf.String())
	csvReader := csv.NewReader(stringReader)

	return csvReader, nil
}

func buildRow(suite *SuiteBuilder, data []string) *TestCase {
	finalData := suite.transformer(data)
	return &TestCase{suite.determineGroup(finalData), finalData, suite.test}
}

func (suite *SuiteBuilder) determineGroup(finalData []interface{}) string {
	group := ""

	if suite.groupColumn >= 0 && suite.groupColumn < len(finalData) {
		group = fmt.Sprintf("%v", finalData[suite.groupColumn])
	}

	return group
}
