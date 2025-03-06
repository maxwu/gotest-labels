package gotestlabels

import (
	"os"
	"regexp"
	"strings"
)

type TestLabels map[string]string

type cliArgs struct {
	runRegex *regexp.Regexp
	listMode bool
	labels   TestLabels
}

func (c *cliArgs) labelsEnabled() bool {
	return len(c.labels) > 0
}

func NewCliArgs() *cliArgs {
	cliArgs := &cliArgs{
		labels: make(TestLabels),
	}
	testLabels := os.Getenv("TEST_LABELS")
	if testLabels != "" {
		pairs := strings.SplitSeq(testLabels, ",")
		for pair := range pairs {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 {
				cliArgs.labels[parts[0]] = parts[1]
			}
		}
	}
	return cliArgs
}

// Parse the go test CLI -run, -list, -json and the new added -labels flags.
func parseArgs() *cliArgs {
	defer removeLabelFlags()

	cliArgs := NewCliArgs()
	runPattern := ""
	args := os.Args[1:]

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

		// -labels flag is merged up with values from TEST_LABELS env var
		if arg == "-labels" && i+1 < len(args) {
			pair := args[i+1]
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 {
				cliArgs.labels[parts[0]] = parts[1]
			}
			i++
		}
	}

	if runPattern != "" {
		cliArgs.runRegex = regexp.MustCompile(runPattern)
	}

	return cliArgs
}

// Remove the -labels flag from os.Args after parsing it. This flag isn't a std go test flag.
func removeLabelFlags() {
	var newArgs []string
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-labels" {
			i++ // 跳过值
			continue
		}
		newArgs = append(newArgs, os.Args[i])
	}
	os.Args = newArgs
}
