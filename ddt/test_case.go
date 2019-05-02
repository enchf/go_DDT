package ddt

// TestCase - Abstraction for the individual test case execution.
type TestCase struct {
	data  []interface{}
	test  TestExecutor
	group string
}

// Run - Execute the test case itself.
func (testCase *TestCase) Run() bool {
	return testCase.test(testCase.data)
}
