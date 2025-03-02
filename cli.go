package gotestlabels

import (
	"log/slog"
	"os"
	"regexp"
	"strings"

	slogenv "github.com/cbrewster/slog-env"
)

type TestLabels map[string]string

type cliArgs struct {
	runRegex *regexp.Regexp
	listMode bool
	jsonLog bool
	labels TestLabels
}

func NewCliArgs() *cliArgs {
	cliArgs := &cliArgs{
		labels: make(TestLabels),
	}
	testLabels := os.Getenv("TEST_LABELS")
	if testLabels != "" {
		pairs := strings.Split(testLabels, ",")
		for _, pair := range pairs {
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
	cliArgs := NewCliArgs()
    runPattern := ""

	args := os.Args[1:]

    for i := 0; i < len(args); i++ {
        arg := args[i]

        if arg == "-run" {
            if i+1 < len(args) {
                runPattern = args[i+1]
                i++
            }
            continue
        } else if strings.HasPrefix(arg, "-run=") {
            runPattern = strings.TrimPrefix(arg, "-run=")
            continue
        }

        if arg == "-list" {
            cliArgs.listMode = true
            continue
        }

        if arg == "-json" {
            logger = slog.New(slogenv.NewHandler(slog.NewJSONHandler(os.Stderr, nil)))
			cliArgs.jsonLog = true
            continue
        }

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

	removeLabelFlags()
	logger.Debug("Parsed CLI args", "args", *cliArgs)

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
