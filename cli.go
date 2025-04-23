package gotest_labels

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type cliArgs struct {
	runRegex *regexp.Regexp  // The regex pattern for -run or -list
	listMode bool  // Whether the -list flag is used
	labels   string // The labels filter from the -labels flag or TEST_LABELS env variable
	labelsAST Node // The parsed AST of the labels filter
}

func (c *cliArgs) labelsEnabled() bool {
	return c.labels != "" && c.labelsAST != nil
}

func (c *cliArgs) buildLabelsAST() {
	if c.labels == "" {
		c.labelsAST = nil
		return
	}
	ast, err := ParseLabelExp(c.labels)
	if err != nil {
		fmt.Println("Error parsing label expression:", err)
	} else {
		c.labelsAST = ast
	}
}

func NewCliArgs() *cliArgs {
	cliArgs := &cliArgs{
		labels: os.Getenv("TEST_LABELS"),
	}
	cliArgs.buildLabelsAST()
	return cliArgs
}

// Parse the os.Args for -run, -list, -json and the new added -labels flags in go test command.
func ParseOSArgs() *cliArgs {
	defer removeLabelFlags()

	return parseArgs(os.Args)
}

func parseArgs(osArgs []string) *cliArgs {
	cliArgs := NewCliArgs()
	runPattern := ""
	args := osArgs[1:]

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "-test.run" {
			if i+1 < len(args) {
				runPattern = args[i+1]
				i++
			}
			continue
		} else if strings.HasPrefix(arg, "-test.run=") {
			runPattern = strings.TrimPrefix(arg, "-test.run=")
			continue
		}

		if arg == "-test.list" {
			cliArgs.listMode = true
			if i+1 < len(args) {
				runPattern = args[i+1]
				i++
			}
			continue
		} else if strings.HasPrefix(arg, "-test.list=") {
			cliArgs.listMode = true
			runPattern = strings.TrimPrefix(arg, "-test.list=")
			continue
		}

		// -labels flag overwrites the values from TEST_LABELS env var
		if arg == "-labels" && i+1 < len(args) {
			filter := args[i+1]
			cliArgs.labels = filter
			i++
		}
	}

	if runPattern != "" {
		cliArgs.runRegex = regexp.MustCompile(runPattern)
	}

	cliArgs.buildLabelsAST()

	return cliArgs
}

// Remove the -labels flag from os.Args after parsing it. This flag isn't a std go test flag.
func removeLabelFlags() {
	os.Args = removeLabelFlagsFromArgs(os.Args)
}

func removeLabelFlagsFromArgs(args []string) []string {
	newArgs := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if args[i] == "-labels" {
			i++
			continue
		}
		newArgs = append(newArgs, args[i])
	}
	return newArgs
}
