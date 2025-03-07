package gotestlabels

import (
	"os"
	"testing"
)

func TestMutateTestFilterByLabels(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"theBinDoesntMatter", "-test.v", "-labels", "group=demo", "-test.list=."}
	origDefaultPkg := defaultPkg
	defer func() { defaultPkg = origDefaultPkg }()

	defaultPkg = "./examples/simple"

	tests := MutateTestFilterByLabels()

	t.Logf("Filtered tests: %v\n", tests)

	if len(tests) != 2 {
		t.Errorf("Expected 2 tests, got %v", len(tests))
	}

	if tests[0] != "TestSimpleAlpha" {
		t.Errorf("Expected TestSimpleAlpha, got %v", tests[0])
	}

	if tests[1] != "TestSimpleGamma" {
		t.Errorf("Expected TestSimpleGamma, got %v", tests[1])
	}
}
