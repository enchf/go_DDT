package ddt

var defaultTestExecutor = func(_ []interface{}) bool { return true }
var defaultRowTransformer = func(row []interface{}) []interface{} { return row }
