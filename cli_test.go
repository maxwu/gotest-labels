package gotest_labels

import (
	"testing"
)

func TestNewCliArgs(t *testing.T) {
	t.Setenv("TEST_LABELS", "group=demo")

	// Test NewCliArgs
	args := NewCliArgs()
	if args == nil {
		t.Errorf("NewCliArgs() failed")
		t.Fail()
	}
	if len(args.labels) != 1 || args.labels["group"] != "demo" {
		t.Errorf("NewCliArgs() failed, expected map{group: demo} but found %#v", args.labels)
		t.Fail()
	}
}