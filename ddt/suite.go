package ddt

// RowTransformer - Signature for the row transformer function.
type RowTransformer func(row []interface{}) []interface{}

// TestExecutor - Signature for the test case executor.
type TestExecutor func(data []interface{}) bool

// VariablesFormat - Format in which the variable values should be set.
type VariablesFormat map[string]interface{}

// SuiteBuilder - Test suite main class.
type SuiteBuilder struct {
	Name        string
	inputFile   string
	groupColumn int
	headers     bool
	variables   VariablesFormat
	transformer RowTransformer
	test        TestExecutor
}

// NewSuiteBuilder - Creates a new SuiteBuilder setting up input data file.
func NewSuiteBuilder(inputFile string) *SuiteBuilder {
	builder := new(SuiteBuilder)
	builder.Name = ""
	builder.inputFile = inputFile
	builder.groupColumn = -1
	builder.headers = false
	builder.variables = nil
	builder.transformer = defaultRowTransformer
	builder.test = defaultTestExecutor

	return builder
}

// GroupBy - Sets up the column by which test cases should be grouped. If column value is invalid it will be ignored.
func (suite *SuiteBuilder) GroupBy(column int) *SuiteBuilder {
	suite.groupColumn = column
	return suite
}

// Headers - Sets up whether or not the input file have headers. False by default.
func (suite *SuiteBuilder) Headers(haveHeaders bool) *SuiteBuilder {
	suite.headers = haveHeaders
	return suite
}

// Variables - Sets up a map with the input variables values.
func (suite *SuiteBuilder) Variables(vars VariablesFormat) *SuiteBuilder {
	suite.variables = vars
	return suite
}

// GlobalName - Sets up the suite name.
func (suite *SuiteBuilder) GlobalName(name string) *SuiteBuilder {
	suite.Name = name
	return suite
}

// RowTransformer - Sets up a custom way to transform input values.
func (suite *SuiteBuilder) RowTransformer(transformer RowTransformer) *SuiteBuilder {
	suite.transformer = transformer
	return suite
}

// TestExecutor - Sets up the way in which data will be validated.
func (suite *SuiteBuilder) TestExecutor(executor TestExecutor) *SuiteBuilder {
	suite.test = executor
	return suite
}
