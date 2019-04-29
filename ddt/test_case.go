package ddt

// TestCase - Abstraction for the individual test case execution.
type TestCase struct {
	data  []interface{}
	test  TestExecutor
	suite string
}
