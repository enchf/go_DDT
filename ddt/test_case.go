package ddt

// TestCase - Abstraction for the individual test case execution.
type TestCase struct {
	Group string
	data  []interface{}
	test  TestExecutor
}

// Run - Execute the test case itself.
func (testCase *TestCase) Run() bool {
	return testCase.test(testCase.data)
}
