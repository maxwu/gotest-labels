# gotestlabels

GoTestLabels is a Go package that provides a way to select test cases by labels in testing function comments and it doesn't need to change the source code.

This package attempts to offload users from the golang `testing` package internal implementation details and the the ASK parsing process
steps and enable readers to filter test cases with labels. The go testing CLI provided the `-run` and `-list` flags to filter tests via regex against the test function names. However, sometimes it's not convenient to keep long function names as the convention or only selec
tests via the linear matching mechanism. For example, if the tests are cataloged in multiple dimensions like `TestAddByClickInDevForIntegration` and `TestAddByClickInProdForRegression`. This package is designed to provide a multiple dimension test
case selector with no actual coding maintenance efforts.

## Usage

### Filter the tests with labels

To use GoTestLabels, import the package and add labels to your test cases:

To install GoTestLabels, use `go get`:

```sh
go get github.com/maxwu/gotestlabels
```

Three ways to use the package to filter tests, one is just to import the `github.com/maxwu/gotestlabels/apply` package
in anonymous, which has an automatically init function to do the filtering. The other way is to explicitly import
the `github.com/maxwu/gotestlabels` package and call the `gotestlabels.MutateTestFuncsByLabels()` function in your
testing package's init function or TestMain function. If the parent package refers to a sub package underneath, adding
the invocation of `gotestlabels.MutateTestFuncsByLabels()` in `TestMain()` function is the required safe way.

Use the simple way, only one line of anonymous import is needed:

```go
_ "github.com/maxwu/gotestlabels/apply"
```

With the explicit way, one line of actual code is needed in the targeted testing package's init function:

```go
func init() {
    _ = gotestlabels.MutateTestFilterByLabels()
}
```

Or, if the parent package refers to a sub package underneath, here's the safe way:

```go
func TestMain(m *testing.M) {
    // The returned test case lists can be used to estimate the test costs or other tasks.
    _ = gotestlabels.MutateTestFilterByLabels()
    os.Exit(m.Run())
}
```

Users can refer to the [examples](examples) to see how to use the package with one line importing code or an explicit function call.

### Add labels to your test cases

To add labels to your test cases, add a comment to the test function in `@key=value` format. Double slash or slash start are both supported.

### Run Go Test with labels

Using env variable `TEST_LABELS` to specify the labels to run.

```sh
TEST_LABELS="group=demo" go test -v ./examples/simple/...
```

Or, if all the packages under test are equipted with the gotestlabels, the labels can be passed in CLI:

```sh
go test -v -labels="group=demo" ./examples/simple/...
```

### Compatibility

If there's no `TEST_LABELS` var or `-labels` flag passed in, the package will do nothing and go test runs normally.

If there are regex selectors like `-run` or `-list`, the labels will be applied after the regex selectors.

Due to go package loading mechanism, each involved package still needs to equip with the gotestlabels via one of the
provided three ways even the wildcard `your_package/...` is used in CLI.

### Examples

Here is an example of how to use GoTestLabels:

```go
package yourpackage

import (
    "testing"
    _ "github.com/maxwu/gotestlabels/apply"
)

// @group=demo
func TestExample(t *testing.T) {
    t.Log("Testing yourpackage.TestExample")
    // Your test code here
}
```

To run the above example, the CLI could be: 

```sh
TEST_LABELS="group=demo" go test -v ./yourpackage
```

or,

```sh
go test -v -labels="group=demo"./yourpackage
```

Given the tests in [examples](examples) folder, the below CLI runs tests with label `group=demo`:

```sh
❯ TEST_LABELS="group=demo" go test -v ./examples/simple
=== RUN   TestSimpleAlpha
    demo_test.go:13: Testing examples.simple.TestSimpleAlpha
--- PASS: TestSimpleAlpha (0.00s)
=== RUN   TestSimpleGamma
    demo_test.go:25: Testing examples.simple.TestSimpleGamma
--- PASS: TestSimpleGamma (0.00s)
PASS
ok  	gotestlabels/examples/simple	0.429s
```

To run the tests with label `env=dev`, the CLI could be:

```sh
❯ TEST_LABELS="env=dev" go test -v ./examples/simple -count=1
=== RUN   TestSimpleBeta
    demo_test.go:20: Testing examples.simple.TestSimpleBeta
--- PASS: TestSimpleBeta (0.00s)
PASS
ok  	gotestlabels/examples/simple	0.283s
```

Readers are kindly reminded to add `-count=1` to the CLI since there's only env var changes so the tests shall be forced
to rerun. Actually the test codes are rebuilt to new temporary binary file but the objective of this package is to offload
readers from golang `testing` package internal details.

### Debug

* To enable debug logs, set env var `GO_LOG=debug`.
* Add `-count=1` if only env vars are changed and the go test CLI still runs the same cases. It's due to go test internal
mechanism that the same built binary is used since there's no code change in between.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the Apache License 2.0. See the LICENSE file for details.
