package gotest_labels

import (
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
		args := parseArgs([]string{"theBinDoesntMatter", "-labels", "group=cliFlag"})
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
