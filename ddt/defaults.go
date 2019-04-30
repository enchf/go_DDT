package ddt

var defaultTestExecutor TestExecutor = func(_ []interface{}) bool { return true }

var defaultRowTransformer = func(row []string) []interface{} {
	finalRow := make([]interface{}, len(row))
	for i, val := range row {
		finalRow[i] = val
	}
	return finalRow
}
