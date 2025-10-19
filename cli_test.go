package gotest_labels

import (
	"regexp"
	"slices"
	"testing"
)

func TestNewCliArgs(t *testing.T) {
	t.Run("New CliArgs with env var", func(t *testing.T) {
		t.Setenv("TEST_LABELS", "group=demo")

		args := NewCliArgs()
		if args == nil {
			t.Fatalf("NewCliArgs() failed: args shall not be nil")
			return
		}

		if args.labels != "group=demo" {
			t.Errorf("NewCliArgs() failed, expected group=demo but found %#v", args.labels)
			t.Fail()
		}

		if !args.labelsEnabled() {
			t.Errorf("NewCliArgs() failed, expected labelsEnabled to be true")
			t.Fail()
		}

		if args.labelsAST == nil {
			t.Errorf("NewCliArgs() failed, expected labelsAST to be not nil")
			t.Fail()
		}
	})

	t.Run("CLI flag shall overwrite env var", func(t *testing.T) {
		t.Setenv("TEST_LABELS", "group=demo")
		args := parseArgs([]string{"theBinDoesntMatter.test", "-labels", "group=cliFlag"})
		if args == nil {
			t.Fatalf("NewCliArgs() failed: args shall not be nil")
			return
		}

		if args.labels != "group=cliFlag" {
			t.Errorf("NewCliArgs() failed, expected group=cliFlag but found %#v", args.labels)
			t.Fail()
		}

		if !args.labelsEnabled() {
			t.Errorf("NewCliArgs() failed, expected labelsEnabled to be true")
			t.Fail()
		}

		if args.labelsAST == nil {
			t.Errorf("NewCliArgs() failed, expected labelsAST to be not nil")
			t.Fail()
		}
	})

	t.Run("CLI flag with equal sign", func(t *testing.T) {
		t.Setenv("TEST_LABELS", "group=demo")
		args := parseArgs([]string{"theBinDoesntMatter.test", `-labels=group=cliFlag`})
		if args == nil {
			t.Fatalf("NewCliArgs() failed: args shall not be nil")
			return
		}

		if args.labels != "group=cliFlag" {
			t.Errorf("NewCliArgs() failed, expected group=cliFlag but found %#v", args.labels)
			t.Fail()
		}

		if !args.labelsEnabled() {
			t.Errorf("NewCliArgs() failed, expected labelsEnabled to be true")
			t.Fail()
		}

		if args.labelsAST == nil {
			t.Errorf("NewCliArgs() failed, expected labelsAST to be not nil")
			t.Fail()
		}
	})

	t.Run("CLI flag with equal sign and quote signs", func(t *testing.T) {
		t.Setenv("TEST_LABELS", "group=demo")
		args := parseArgs([]string{"theBinDoesntMatter.test", `-labels="group=cliFlag"`})
		if args == nil {
			t.Fatalf("NewCliArgs() failed: args shall not be nil")
			return
		}

		if args.labels != "group=cliFlag" {
			t.Errorf("NewCliArgs() failed, expected group=cliFlag but found %#v", args.labels)
			t.Fail()
		}

		if !args.labelsEnabled() {
			t.Errorf("NewCliArgs() failed, expected labelsEnabled to be true")
			t.Fail()
		}

		if args.labelsAST == nil {
			t.Errorf("NewCliArgs() failed, expected labelsAST to be not nil")
			t.Fail()
		}
	})

	t.Run("CLI flag with equal sign and single quote signs", func(t *testing.T) {
		t.Setenv("TEST_LABELS", "group=demo")
		args := parseArgs([]string{"theBinDoesntMatter.test", `-labels='group=cliFlag'`})
		if args == nil {
			t.Fatalf("NewCliArgs() failed: args shall not be nil")
			return
		}

		if args.labels != "group=cliFlag" {
			t.Errorf("NewCliArgs() failed, expected group=cliFlag but found %#v", args.labels)
			t.Fail()
		}

		if !args.labelsEnabled() {
			t.Errorf("NewCliArgs() failed, expected labelsEnabled to be true")
			t.Fail()
		}

		if args.labelsAST == nil {
			t.Errorf("NewCliArgs() failed, expected labelsAST to be not nil")
			t.Fail()
		}
	})
}

func TestRemoveLabelFlagsFromArgsWithoutEqualSign(t *testing.T) {
	origArgs := []string{"-test.v", "-labels", "group=demo", "-test.run", "Alpha"}

	newArgs := removeLabelFlagsFromArgs(origArgs)

	if len(newArgs) != 3 {
		t.Errorf("Expected 3 args after mutation, got %#v", newArgs)
		t.Fail()
	}

	if !slices.Equal(newArgs, []string{"-test.v", "-test.run", "Alpha"}) {
		t.Errorf("Expected [-test.v -test.run Alpha], got %v", newArgs)
		t.Fail()
	}
}

func TestRemoveLabelFlagsFromArgsWithEqualSign(t *testing.T) {
	origArgs := []string{"-test.v", "-labels='group=demo'", "-test.run", "Alpha"}

	newArgs := removeLabelFlagsFromArgs(origArgs)

	if len(newArgs) != 3 {
		t.Errorf("Expected 3 args after mutation, got %#v", newArgs)
		t.Fail()
	}

	if !slices.Equal(newArgs, []string{"-test.v", "-test.run", "Alpha"}) {
		t.Errorf("Expected [-test.v -test.run Alpha], got %v", newArgs)
		t.Fail()
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		osArgs   []string
		expected *cliArgs
	}{
		{
			name:   "Test run pattern with -test.run flag",
			osArgs: []string{"program", "-test.run", "TestPattern"},
			expected: &cliArgs{
				runRegex: regexp.MustCompile("TestPattern"),
				listMode: false,
				labels:   "",
			},
		},
		{
			name:   "Test list mode with -test.list flag",
			osArgs: []string{"program", "-test.list", "ListPattern"},
			expected: &cliArgs{
				runRegex: regexp.MustCompile("ListPattern"),
				listMode: true,
				labels:   "",
			},
		},
		{
			name:   "Test labels with -labels flag",
			osArgs: []string{"program", "-labels", "env=prod"},
			expected: &cliArgs{
				runRegex: nil,
				listMode: false,
				labels:   "env=prod",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArgs(tt.osArgs)

			// Compare runRegex separately as it is a pointer
			if (result.runRegex == nil && tt.expected.runRegex != nil) || (result.runRegex != nil && tt.expected.runRegex == nil) {
				t.Errorf("runRegex mismatch: got %v, want %v", result.runRegex, tt.expected.runRegex)
			} else if result.runRegex != nil && tt.expected.runRegex != nil && result.runRegex.String() != tt.expected.runRegex.String() {
				t.Errorf("runRegex mismatch: got %v, want %v", result.runRegex.String(), tt.expected.runRegex.String())
			}

			// Compare other fields
			if result.listMode != tt.expected.listMode {
				t.Errorf("listMode mismatch: got %v, want %v", result.listMode, tt.expected.listMode)
			}
			if result.labels != tt.expected.labels {
				t.Errorf("labels mismatch: got %v, want %v", result.labels, tt.expected.labels)
			}
		})
	}
}
