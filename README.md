# go_DDT
Applying DichloroDiphenylTrichloroethane to code bugs using Data Driven Testing with Golang

# Goals and usage

DDT is a simple data-driven testing approach for executing testing:

* Reads input test cases from a CSV.
* Test cases can contain parameters in Mustache format: `{{THIS_IS_A_PARAMETER}}`.
* Parameters can be replaced prior to test execution.
* Each row can be transformed to change input data types or to lookup ID's in a database.

First, a DDTSuite builder is created. It accepts the following configurations:

* It is instantiated with the pointer to the CSV where to read the test cases.
* Test cases can be grouped by the values of certain column.
* Input file can contain headers.
* Values for the variables can be set passing a `map[string]interface{}` as argument.
* Test suite can have a global name.
* A row transformer can be set in the form of `func(row []string) []interface{}`.

```go
vars := map[string]interface{} { "VAR_A": true, "VAR_B": false }

transformer := func(row []string) []interface{} {
  return []interface{} { row[0], int(row[1]), Database.lookup(int(row[2])), row[3], int(row[4]) }
}

test := func(data []interface{}) bool {
  return data[2].FunctionToBeTested(data[1]) == data[4]
}

/**
 * Input file can be something like:
 * A,1,2,VAR_A,5
 * A,2,2,VAR_B,10
 * Z,20,2,"Complex value",100
 */
suite := ddt.NewDDTSuite("input.csv")       // Path to the input file.
            .GroupBy(0)                  // Group by column 0.
            .Headers(false)              // True/False if input file has headers.
            .Variables(vars)             // Variables values. If variable is not in the map it keeps the string value.
            .GlobalName("DDT Suite")     // Suite Global name.
            .RowTransformer(transformer) // Function to transform row values.
            .TestFunction(test)          // The test to be executed against all test cases.
            .Build()                     // Returns a map[string][]DDTCase.
                                         // Groups the test cases and each group have an array of test cases.
                                         // If no grouping is set, there is a single no-name group for all test cases.
```

Each DDTCase has the following properties:

* Data: `[]interface{}` input data.
* Test: `func(data []interface{}) bool` test execution.
* Suite: `string` the group this test case belongs to.

Each test case can be run using `testCase.Run()`.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/enchf/go_DDT/tags). 

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Authors

* **Ernesto Espinosa** - *Initial work* - [enchf](https://github.com/enchf)
