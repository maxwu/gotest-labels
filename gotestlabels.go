package gotestlabels

import (
	"log"
	"os"
	"strings"
)

// The actually exposed entrypoint to mutate the test functions by labels
// It can be called in the TestMain function of the test package.
// If the test command is running tests with wildcards for sub packages, either set the labels
// via the TEST_LABELS environment variable or set the labels in the TestMain function in every
// involved package
// The function returns the list of test functions that matched the labels as well. The result can be used to estimate
// the test costs or support the test operation/observability/report features.
func MutateTestFilterByLabels() []string {
	args := parseArgs()
	tests, listMode := getTestFuncsByLabels(args)

	// If the labels are not enabled, return the original tests without mutating the os.Args.
	// The results are useful to estimate the test time and costs.
	if !args.labelsEnabled() {
		log.Printf("Labels are not enabled, running tests as normal and still collect the tests for estimation")
		return tests
	}

	// If the labels are enabled, mutate the os.Args to run the selected tests.
	pattern := "^" + strings.Join(tests, "|") + "$"
	if listMode {
		os.Args = append(os.Args, "-test.list", pattern)
	} else {
		os.Args = append(os.Args, "-test.run", pattern)
	}

	return tests
}

// The internal function to get the selected test functions by labels and whether it's in listing mode
func getTestFuncsByLabels(args *cliArgs) ([]string, bool) {

	allPkgs, err := getPackages()
	if err != nil {
		log.Printf("Error resolving packages: %#v", err)
		return nil, args.listMode
	}

	var allTestFuncs []string
	for _, pkg := range allPkgs {
		files := getTestFiles(pkg)
		funcs, err := FindTestFuncs(files, args.labels)
		if err != nil {
			log.Printf("Error parsing tests %s: : %#v", pkg.Name, err)
			return nil, args.listMode
		}
		allTestFuncs = append(allTestFuncs, funcs...)
	}

	matchedFuncs := filterTestFuncs(allTestFuncs, args.runRegex)
	return matchedFuncs, args.listMode
}
