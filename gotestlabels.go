package gotestlabels

import (
	"log/slog"
	"os"
	"strings"

	slogenv "github.com/cbrewster/slog-env"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slogenv.NewHandler(slog.NewTextHandler(os.Stderr, nil), slogenv.WithDefaultLevel(slog.LevelWarn)))
}

// The actually exposed entrypoint to mutate the test functions by labels
// It can be called in the TestMain function of the test package.
// If the test command is running tests with wildcards for sub packages, either set the labels
// via the TEST_LABELS environment variable or set the labels in the TestMain function in every
// involved package
// The function returns the list of test functions that matched the labels as well. The result can be used to estimate
// the test costs or support the test operation/observability/report features.
func MutateTestFilterByLabels() []string {
	tests, listMode := getTestFuncsByLabels()
	funcs := []string{}
	for _, test := range tests {
		funcs = append(funcs, "^" + test + "$")
	}
	pattern := strings.Join(funcs, "|")
	logger.Debug("The pattern to run or list is", "pattern", pattern)

	if listMode {
		logger.Info("The tests to list are", "tests", tests)
		os.Args = append(os.Args, "-test.list", pattern)
	} else {
		logger.Info("The tests to execute are", "tests", tests)
		os.Args = append(os.Args, "-test.run", pattern)
	}

	logger.Debug("New os.Args", "args", os.Args)
	return tests
}

func getTestFuncsByLabels() ([]string, bool) {
    args := parseArgs()

	allPkgs, err := getPackages()
	if err != nil {
		logger.Error("Error resolving packages", "err", err)
		return nil, args.listMode
	}

	logger.Debug("Packages to parse", "pkgs", allPkgs)

    var allTestFuncs []string
    for _, pkg := range allPkgs {
		files := getTestFiles(pkg)
        funcs, err := FindTestFuncs(files, args.labels)
        if err != nil {
            logger.Error("Error parsing tests", "pkg", pkg, "err", err)
            return nil, args.listMode
        }
        allTestFuncs = append(allTestFuncs, funcs...)
    }

    matchedFuncs := filterTestFuncs(allTestFuncs, args.runRegex)
    return matchedFuncs, args.listMode
}
