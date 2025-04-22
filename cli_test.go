package gotest_labels

import (
	"testing"
)

func TestNewCliArgs(t *testing.T) {
	t.Setenv("TEST_LABELS", "group=demo")

	// Test NewCliArgs
	args := NewCliArgs()
	if args == nil {
		t.Fatalf("NewCliArgs() failed: args is nil")
		return
	}

	if len(args.labels) != 1 || args.labels != "group=demo" {
		t.Errorf("NewCliArgs() failed, expected group=demo but found %#v", args.labels)
		t.Fail()
	}
}
